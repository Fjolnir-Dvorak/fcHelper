package cmd

import "github.com/spf13/cobra"

// extractCmd represents the extract command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compares two language xmls and extracts all different (obsolete/ keys.",
	Long: `Compares two different files or two different directory-trees with the same
	 structure. It will generate a file for each comparison showing missing keys (keys
	 in second file) with '+++' in front and the surplus (keys in first file) with '---'
	 in front. If the comparison file will be applied to the base file it will have the
	 same key structure as the second one.`,
	Run: doCompare,
}

func init() {
	RootCmd.AddCommand(compareCmd)
	compareCmd.Flags().StringVarP(&gameDir, "gameDir", "g", "", "Directory containing FortressCraft Evolved.")
	compareCmd.Flags().StringVarP(&destination, "destination", "d", out, "Destination Directory to create the parsed files.")
	compareCmd.Flags().StringVarP(&language, "language", "l", "", "Used language shortkey.")
	compareCmd.Flags().BoolVarP(&createTemplates, "createTemplates", "t", false, "Wether to generate xml-templates.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func doCompare(cmd *cobra.Command, args []string) {

}
