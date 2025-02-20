package crosslab

import (
	"github.com/kanzifucius/crosslab/pkg/config"
	"github.com/spf13/cobra"
)

var (
	outputDir string
)

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".crosslab", "Output directory for configuration files (default: .crosslab in current directory)")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration files",
	Long: `Create Kind and Crosslab configuration files in the specified directory.
By default, files will be created in the .crosslab directory in your current working directory.
This command will create:
- kind-config.yaml: Kind cluster configuration
- config/providers.yaml: Crossplane provider configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		initializer := config.NewInitializer(outputDir)
		return initializer.Initialize()
	},
}
