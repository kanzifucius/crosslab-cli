package crosslab

import (
	"context"
	"fmt"

	"github.com/kanzifucius/crosslab/pkg/provider"

	"github.com/spf13/cobra"
)

var (
	providerPackage    string
	providerVersion    string
	forceReinstall     bool
	providerConfigFile string
)

func init() {
	RootCmd.AddCommand(providerCmd)
	providerCmd.AddCommand(installProviderCmd)
	providerCmd.AddCommand(listProvidersCmd)
	providerCmd.AddCommand(installAllCmd)

	// Add flags to install command
	installProviderCmd.Flags().StringVarP(&providerPackage, "package", "p", "", "Provider package (e.g., xpkg.upbound.io/upbound/provider-aws)")
	installProviderCmd.Flags().StringVarP(&providerVersion, "version", "v", "", "Provider version")
	installProviderCmd.Flags().StringVarP(&clusterName, "name", "n", "", "Provider name")
	installProviderCmd.Flags().BoolVarP(&forceReinstall, "force", "f", false, "Force reinstall if provider exists")
	installProviderCmd.MarkFlagRequired("package")
	installProviderCmd.MarkFlagRequired("version")
	installProviderCmd.MarkFlagRequired("name")

	// Add flags to install-all command
	installAllCmd.Flags().BoolVarP(&forceReinstall, "force", "f", false, "Force reinstall if providers exist")
	installAllCmd.Flags().StringVarP(&providerConfigFile, "config", "c", "", "Path to provider configuration file")
}

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage Crossplane providers",
	Long:  `Install, list, and manage Crossplane providers`,
}

var installProviderCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a Crossplane provider",
	Long:  `Install a Crossplane provider with the specified package and version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		manager, err := provider.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create provider manager: %v", err)
		}

		p := provider.Provider{
			Name:    clusterName,
			Package: providerPackage,
			Version: providerVersion,
		}

		fmt.Printf("Installing provider '%s' from package %s:%s...\n", p.Name, p.Package, p.Version)

		if err := manager.Install(ctx, p, forceReinstall); err != nil {
			return fmt.Errorf("failed to install provider: %v", err)
		}

		fmt.Printf("Waiting for provider '%s' to become healthy...\n", p.Name)
		if err := manager.WaitForHealth(ctx, p.Name); err != nil {
			return fmt.Errorf("failed while waiting for provider: %v", err)
		}

		fmt.Printf("Provider '%s' installed and healthy!\n", p.Name)
		return nil
	},
}

var listProvidersCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Crossplane providers",
	Long:  `List all installed Crossplane providers and their status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		manager, err := provider.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create provider manager: %v", err)
		}

		providers, err := manager.List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list providers: %v", err)
		}

		if len(providers) == 0 {
			fmt.Println("No Crossplane providers installed")
			return nil
		}

		fmt.Println("Installed Crossplane providers:")
		for _, p := range providers {
			fmt.Printf("- %s (%s)\n", p.Name, p.Package)
		}

		return nil
	},
}

var installAllCmd = &cobra.Command{
	Use:   "install-all",
	Short: "Install all required Crossplane providers",
	Long:  `Install all required Crossplane providers for the project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Load provider configuration
		config, err := provider.LoadConfig(providerConfigFile)
		if err != nil {
			return fmt.Errorf("failed to load provider configuration: %v", err)
		}

		// Validate configuration
		if err := config.Validate(); err != nil {
			return fmt.Errorf("invalid provider configuration: %v", err)
		}

		manager, err := provider.NewManager()
		if err != nil {
			return fmt.Errorf("failed to create provider manager: %v", err)
		}

		// Install AWS family provider
		fmt.Println("Installing AWS provider...")
		if err := installAndWait(ctx, manager, config.AWS.Family); err != nil {
			return err
		}

		// Install AWS service providers
		fmt.Println("\nInstalling AWS service providers...")
		for _, p := range config.AWS.Services {
			if err := installAndWait(ctx, manager, p); err != nil {
				return err
			}
		}

		// Install other providers
		fmt.Println("\nInstalling other providers...")
		for _, p := range config.OtherProviders {
			if err := installAndWait(ctx, manager, p); err != nil {
				return err
			}
		}

		fmt.Println("\nAll providers installed successfully!")
		return nil
	},
}

func installAndWait(ctx context.Context, manager *provider.Manager, p provider.Provider) error {
	fmt.Printf("Installing %s...\n", p.Name)

	if err := manager.Install(ctx, p, forceReinstall); err != nil {
		return fmt.Errorf("failed to install provider %s: %v", p.Name, err)
	}

	fmt.Printf("Waiting for %s to become healthy...\n", p.Name)
	if err := manager.WaitForHealth(ctx, p.Name); err != nil {
		return fmt.Errorf("failed while waiting for provider %s: %v", p.Name, err)
	}

	fmt.Printf("Provider %s is healthy ✓\n", p.Name)
	return nil
}
