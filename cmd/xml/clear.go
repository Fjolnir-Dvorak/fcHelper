package xml

import (
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var clearCMD = &cobra.Command{
	Use:   "clearXML",
	Short: "Deletes everything except the tags itself. No Comments, no Content.",
	Long: `If you do not know what you are doing exactly do not use this functionality.
	It is able to corrupt your whole game and to destroy your savegames. Do not ask
	me how, but I am shure there is a possibility for that.`,
	Run: doClearXML,
}

var (
	fileToClear string
)

func init() {
	XmlCmd.AddCommand(clearCMD)
	clearCMD.Flags().StringVarP(&fileToClear, "file", "f", "", "File to clear.")
}

func doClearXML(cmd *cobra.Command, args []string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(fileToClear); err != nil {
		return
	}
	recursiveClear(doc.Root())
	doc.WriteToFile(fileToClear)
}
func recursiveClear(elem *etree.Element) {
	elem.SetText("")
	for _, child := range elem.ChildElements() {
		recursiveClear(child)
	}
}
