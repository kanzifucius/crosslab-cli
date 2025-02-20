package crosslab

import (
	"context"
	"fmt"

	"github.com/kanzifucius/crosslab/pkg/kind"
	"github.com/kanzifucius/crosslab/pkg/provider"

	"github.com/spf13/cobra"
)

var (
	kindConfigFile string
	clusterConfig  string
	clusterName    string
	forceProviders bool
	forceCreate    bool
)

func init() {
	RootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(createCmd)

	// Add flags to create command
	createCmd.Flags().StringVarP(&kindConfigFile, "config", "c", "", "Path to the Kind cluster configuration file")
	createCmd.Flags().StringVarP(&clusterConfig, "provider-config", "p", "", "Path to the provider configuration file")
	createCmd.Flags().StringVarP(&clusterName, "name", "n", "kind", "Name of the Kind cluster")
	createCmd.Flags().BoolVarP(&forceProviders, "force-providers", "f", false, "Force reinstall providers if they exist")
	createCmd.Flags().BoolVar(&forceCreate, "force", false, "Force recreation of cluster if it exists")
	createCmd.MarkFlagRequired("config")
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage Kind clusters",
	Long:  `Create and manage Kind clusters with integrated provider installation`,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Kind cluster",
	Long:  `Create a new Kind cluster and install required providers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Check if cluster exists
		kindClient := kind.NewClient()
		exists, err := kindClient.ClusterExists(clusterName)
		if err != nil {
			return fmt.Errorf("failed to check cluster existence: %v", err)
		}

		if exists {
			if !forceCreate {
				return fmt.Errorf("cluster '%s' already exists. Use --force flag to recreate it", clusterName)
			}

			fmt.Printf("Deleting existing cluster '%s'...\n", clusterName)
			if err := kindClient.DeleteCluster(clusterName); err != nil {
				return fmt.Errorf("failed to delete existing cluster: %v", err)
			}
		}

		// Create Kind cluster
		fmt.Printf("Creating Kind cluster '%s'...\n", clusterName)
		if err := kindClient.CreateCluster(kindConfigFile, clusterName); err != nil {
			return fmt.Errorf("failed to create Kind cluster: %v", err)
		}
		fmt.Printf("Kind cluster '%s' created successfully!\n", clusterName)

		// Initialize provider manager
		manager, err := provider.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create provider manager: %v", err)
		}

		// Load provider configuration
		config, err := provider.LoadConfig(clusterConfig)
		if err != nil {
			return fmt.Errorf("failed to load provider configuration: %v", err)
		}

		// Validate configuration
		if err := config.Validate(); err != nil {
			return fmt.Errorf("invalid provider configuration: %v", err)
		}

		// install crossplane helm chart
		fmt.Println("\nInstalling Crossplane Helm chart...")
		if err := installCrossplane(ctx, manager); err != nil {
			return err
		}

		// Install AWS family provider
		fmt.Println("\nInstalling AWS provider...")
		if err := installClusterProvider(ctx, manager, config.AWS.Family); err != nil {
			return err
		}

		// Install AWS service providers
		fmt.Println("\nInstalling AWS service providers...")
		for _, p := range config.AWS.Services {
			if err := installClusterProvider(ctx, manager, p); err != nil {
				return err
			}
		}

		// Install other providers
		fmt.Println("\nInstalling other providers...")
		for _, p := range config.OtherProviders {
			if err := installClusterProvider(ctx, manager, p); err != nil {
				return err
			}
		}

		fmt.Println("\nCluster setup completed successfully!")

		// list providers
		fmt.Println("\nListing installed providers...")
		providers, err := manager.List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list providers: %v", err)
		}
		fmt.Println(providers)
		return nil
	},
}

func installClusterProvider(ctx context.Context, manager *provider.Manager, p provider.Provider) error {
	fmt.Printf("Installing %s...\n", p.Name)

	if err := manager.Install(ctx, p, forceProviders); err != nil {
		return fmt.Errorf("failed to install provider %s: %v", p.Name, err)
	}

	fmt.Printf("Waiting for %s to become healthy...\n", p.Name)
	if err := manager.WaitForHealth(ctx, p.Name); err != nil {
		return fmt.Errorf("failed while waiting for provider %s: %v", p.Name, err)
	}

	fmt.Printf("Provider %s is healthy ✓\n", p.Name)
	return nil
}

func installCrossplane(ctx context.Context, manager *provider.Manager) error {
	if err := manager.InstallCrossplane(ctx); err != nil {
		return fmt.Errorf("failed to install Crossplane: %v", err)
	}

	fmt.Println("Waiting for Crossplane to become healthy...")
	if err := manager.WaitForCrossplaneHealth(ctx); err != nil {
		return fmt.Errorf("failed while waiting for Crossplane: %v", err)
	}

	fmt.Println("Crossplane is healthy ✓")
	return nil
}
