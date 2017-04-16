// Copyright © 2017 Raphael Tiersch <tiersch.raphael@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/Fjolnir-Dvorak/fcHelper/forgecraft"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	srcItems    string
	srcHandbook string
)

const (
	Name = "Name"
	Key  = "Key"
)

// repairKeysCmd represents the repairKeys command
var repairKeysCmd = &cobra.Command{
	Use:   "repairKeys",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: doRepair,
}

func init() {
	RootCmd.AddCommand(repairKeysCmd)
	repairKeysCmd.Flags().StringVarP(&srcItems, "itemFiles", "i", "", "Source of all itemfiles.")
	repairKeysCmd.Flags().StringVarP(&srcHandbook, "handbook", "s", "", "Source of handbook to repair.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repairKeysCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repairKeysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func doRepair(cmd *cobra.Command, args []string) {
	files, _ := ioutil.ReadDir(srcItems)
	keylist := datatypes.InitEmptyKeyNameList()
	for _, file := range files {
		switch file.Name() {
		case forgecraft.ItemsFilename:
			itemsParsed, _ := forgecraft.ParseItemsXMLFile(srcItems)
			keylist.Merge(itemsParsed.CreateKeyMap())
			break
		default:
			fmt.Printf("!!! NOT SUPPORTED FILETYPE: %s", file.Name())
		}
	}

	//for _, keyvalue := range keylist.Iterate() {
	//	fmt.Printf("%s:%s\n", keyvalue.Key, keyvalue.Name)
	//}

	readDir(srcHandbook, keylist)
}

func readDir(directory string, keylist datatypes.KeyNameList) {
	files, _ := ioutil.ReadDir(directory)

	for _, fileInf := range files {
		basefile := filepath.Join(directory, fileInf.Name())
		if fileInf.IsDir() {
			readDir(basefile, keylist)
			continue
		}
		doc := etree.NewDocument()
		//fmt.Printf("Reading: %s\n", basefile)
		if err := doc.ReadFromFile(basefile); err != nil {
			panic(err)
		}
		root := doc.Root()
		var name, key string
		for _, child := range root.ChildElements() {
			if child.Tag == Key {
				name = child.Text()
				key = keylist.GetKey(name)
				if key == "" {
					fmt.Printf("Key not Found for name: %s ||| File: %s\n", name, basefile)
					continue
				}
				child.Tag = key
				child.SetText(Name)
				_ = doc.WriteToFile(basefile)
				os.Rename(basefile, filepath.Join(directory, key+".xml"))
				break
			}
		}
	}
}
