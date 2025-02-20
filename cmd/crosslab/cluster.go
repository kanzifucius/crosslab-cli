package crosslab

import (
	"context"
	"fmt"

	"github.com/kanzifucius/crosslab/pkg/config"
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
	createCmd.Flags().StringVarP(&kindConfigFile, "config", "c", ".crosslab/kind-config.yaml", "Path to the Kind cluster configuration file")
	createCmd.Flags().StringVarP(&clusterConfig, "provider-config", "p", ".crosslab/config/crosslab-config.yaml", "Path to the provider configuration file")
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

		// Check if configuration files exist
		if err := config.CheckConfigFile(kindConfigFile); err != nil {
			return err
		}

		if err := config.CheckConfigFile(clusterConfig); err != nil {
			return err
		}

		// Check if cluster exists
		kindManager := kind.NewManager()
		exists, err := kindManager.ClusterExists(clusterName)
		if err != nil {
			return fmt.Errorf("failed to check cluster existence: %v", err)
		}

		if exists {
			if !forceCreate {
				return fmt.Errorf("cluster '%s' already exists. Use --force flag to recreate it", clusterName)
			}

			fmt.Printf("Deleting existing cluster '%s'...\n", clusterName)
			if err := kindManager.DeleteCluster(clusterName); err != nil {
				return fmt.Errorf("failed to delete existing cluster: %v", err)
			}
		}

		// Create Kind cluster
		fmt.Printf("Creating Kind cluster '%s'...\n", clusterName)
		if err := kindManager.CreateCluster(kindConfigFile, clusterName); err != nil {
			return fmt.Errorf("failed to create Kind cluster: %v", err)
		}
		fmt.Printf("Kind cluster '%s' created successfully!\n", clusterName)

		// Initialize provider manager
		manager, err := provider.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create provider manager: %v", err)
		}

		// Load provider configuration
		providerConfig, err := config.LoadConfig(clusterConfig)
		if err != nil {
			return fmt.Errorf("failed to load provider configuration: %v", err)
		}

		// Validate configuration
		if err := providerConfig.Validate(); err != nil {
			return fmt.Errorf("invalid provider configuration: %v", err)
		}

		// install crossplane helm chart
		fmt.Println("\nInstalling Crossplane Helm chart...")
		if err := InstallCrossplane(ctx, manager); err != nil {
			return err
		}

		// Install AWS family provider
		fmt.Println("\nInstalling AWS provider...")
		if err := InstallClusterProvider(ctx, manager, providerConfig.AWS.Family); err != nil {
			return err
		}

		// Install AWS service providers
		fmt.Println("\nInstalling AWS service providers...")
		for _, p := range providerConfig.AWS.Services {
			if err := InstallClusterProvider(ctx, manager, p); err != nil {
				return err
			}
		}

		// Install other providers
		fmt.Println("\nInstalling other providers...")
		for _, p := range providerConfig.OtherProviders {
			if err := InstallClusterProvider(ctx, manager, p); err != nil {
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

func InstallClusterProvider(ctx context.Context, manager provider.Manager, provider config.Provider) error {
	fmt.Printf("Installing provider %s...\n", provider.Name)
	if err := manager.Install(ctx, provider, forceProviders); err != nil {
		return fmt.Errorf("failed to install provider %s: %v", provider.Name, err)
	}

	fmt.Printf("Waiting for provider %s to become healthy...\n", provider.Name)
	if err := manager.WaitForHealth(ctx, provider.Name); err != nil {
		return fmt.Errorf("failed to wait for provider %s health: %v", provider.Name, err)
	}

	return nil
}

// installCrossplane installs the Crossplane Helm chart
func InstallCrossplane(ctx context.Context, manager provider.Manager) error {
	fmt.Println("Installing Crossplane...")
	if err := manager.InstallCrossplane(ctx); err != nil {
		return fmt.Errorf("failed to install Crossplane: %v", err)
	}

	fmt.Println("Waiting for Crossplane to become healthy...")
	if err := manager.WaitForCrossplaneHealth(ctx); err != nil {
		return fmt.Errorf("failed to wait for Crossplane health: %v", err)
	}

	fmt.Println("Crossplane is healthy ✓")
	return nil
}
