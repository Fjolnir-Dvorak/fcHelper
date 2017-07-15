package xml

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
)

var (
	fileToClear string
)

func init() {
}

func DoClearXML(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("ERROR: No input file given.")
		return
	}
	for _, file := range args {
		doc := etree.NewDocument()
		if err := doc.ReadFromFile(file); err != nil {
			return
		}
		recursiveClear(doc.Root())
		doc.WriteToFile(fileToClear)
	}
}

func recursiveClear(elem *etree.Element) {
	elem.SetText("")
	for _, child := range elem.ChildElements() {
		recursiveClear(child)
	}
}
