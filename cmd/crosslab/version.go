package crosslab

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the current version of the Crosslocal CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Crosslocal CLI version %s\n", version)
	},
}
