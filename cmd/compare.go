package cmd

import (
	"fmt"
	. "github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/Fjolnir-Dvorak/fcHelper/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

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

var (
	compOld string
	compNew string
	compOut string
)

func init() {
	RootCmd.AddCommand(compareCmd)
	compareCmd.Flags().StringVarP(&compOld, "base", "b", "", "Old xml version.")
	compareCmd.Flags().StringVarP(&compNew, "new", "n", "", "New xml version.")
	compareCmd.Flags().StringVarP(&compOut, "out", "o", out, "Output directory")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func doCompare(cmd *cobra.Command, args []string) {
	// Load both files
	_, old := util.ParseXLFKeyName(compOld)
	_, new := util.ParseXLFKeyName(compNew)
	old = old.Mergesort()
	new = new.Mergesort()
	oldIt := old.Iterate()
	newIt := new.Iterate()
	oldExtra := []KeyName{}
	newExtra := []KeyName{}

	countOld, countNew := 0, 0

	fmt.Println("Start comparing")
	for countOld < len(oldIt) || countNew < len(newIt) {
		//fmt.Printf("countOld: %d, countNew: %d\n", countOld, countNew)
		if countOld >= len(oldIt) {
			// Old is empty. Only new keys are left
			fmt.Printf("Adding keys to newExtra: %d\n", len(newIt[countNew:]))
			newExtra = append(newExtra, newIt[countNew:]...)
			break
		} else if countNew >= len(newIt) {
			// New is empty. Only old keys are left
			fmt.Printf("Adding keys to oldExtra: %d\n", len(oldIt[countOld:]))
			oldExtra = append(oldExtra, oldIt[countOld:]...)
			break
		} else if oldIt[countOld].Key == newIt[countNew].Key {
			countOld++
			countNew++
		} else if oldIt[countOld].Key < newIt[countNew].Key {
			fmt.Printf("Adding key to oldExtra: %s\n", oldIt[countOld].Key)
			// old has a key new hasn't
			oldExtra = append(oldExtra, oldIt[countOld])
			countOld++
		} else {
			fmt.Printf("Adding key to newExtra: %s\n", newIt[countNew].Key)
			// new has a key old hasn't
			newExtra = append(newExtra, newIt[countNew])
			countNew++
		}
	}
	// parse xml and collect keys
	os.MkdirAll(compOut, os.ModePerm)
	oldFilename := filepath.Join(compOut, "baseExtra.xml")
	newFilename := filepath.Join(compOut, "newExtra.xml")
	util.CreateXLFkn(oldExtra, oldFilename)
	util.CreateXLFkn(newExtra, newFilename)
}
