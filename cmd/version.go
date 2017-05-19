package cmd

import "github.com/spf13/cobra"

// extractCmd represents the extract command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version number and the latest commit.",
	Run:   doVersion,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func doVersion(cmd *cobra.Command, args []string) {
}
