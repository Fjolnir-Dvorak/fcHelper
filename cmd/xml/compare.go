package xml

import (
	"fmt"
	. "github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/Fjolnir-Dvorak/fcHelper/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	compOld string
	compNew string
	compOut string
)

func init() {
}

func DoCompare(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("ERROR: No input file given.")
		return
	}
	if len(args) != 2 {
		fmt.Println("ERROR: Please specify two files for comparison")
		return
	}
	compOld = args[0]
	compNew = args[1]

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
