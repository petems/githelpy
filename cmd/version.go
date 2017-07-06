package cmd

import (
	"github.com/petems/githelpy/githelpy"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Run: func(cmd *cobra.Command, args []string) {
		info(githelpy.GetVersion())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
