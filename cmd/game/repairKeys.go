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

package game

import (
	"fmt"

	"github.com/Fjolnir-Dvorak/fcHelper/cmd/game/structures"
	"github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/Fjolnir-Dvorak/fcHelper/forgecraft"
	"github.com/Fjolnir-Dvorak/fcHelper/util"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	srcItems     string
	srcHandbook  string
	noName       = false
	spaceReplace = false

	NoSpace   bool
	NameToKey bool
)

const (
	Name = "Name"
	Key  = "Key"
)

func init() {
}

func DoRepair(cmd *cobra.Command, args []string) {
	if GameDir != "" {
		// The user specified another installation directory
		gameDir = GameDir
	} else {
		// Default Steam installation directory
		gameDir = SteamGameDir
	}
	noName = !NameToKey
	spaceReplace = !NoSpace

	srcItems = filepath.Join(gameDir, structures.Data)
	srcHandbook = filepath.Join(gameDir, structures.Handbook)

	files, _ := ioutil.ReadDir(srcItems)
	keylist := datatypes.InitEmptyKeyNameList()
	for _, file := range files {
		baseFile := filepath.Join(srcItems, file.Name())
		fmt.Printf("Reading dictionary file: %s\n", baseFile)
		switch file.Name() {
		case forgecraft.ItemsFilename:
			items := forgecraft.ParseItemsXMLFile(baseFile)
			keylist.ConCat(items.CreateKeyMap())
			break
		case forgecraft.TerrainDataFilename:
			items := forgecraft.ParseTerrainDataXMLFile(baseFile)
			keylist.ConCat(items.CreateKeyMap())
			break
		default:
			fmt.Printf("!!! NOT SUPPORTED FILETYPE: %s\n", baseFile)
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
			if child.Tag == Name {
				name = child.Text()
				key = keylist.GetKey(name)
				if key == "" {
					fmt.Printf("Key not Found for name: %s ||| File: %s\n", name, basefile)
					continue
				}
				//fmt.Printf("Key {{%s}} Found for name: %s ||| File: %s\n", key, name, basefile)
				child.Tag = Key
				child.SetText(key)
				out, _ := doc.WriteToString()
				util.WriteRawStringToFile(basefile, out)
				if noName {
					if spaceReplace {
						sanitizedName := strings.Replace(fileInf.Name(), " ", "_", -1)
						os.Rename(basefile, filepath.Join(directory, sanitizedName))
					}
				} else {
					os.Rename(basefile, filepath.Join(directory, key+".xml"))
				}
				//os.Rename(basefile, filepath.Join(directory, "_"+key+"_"+fileInf.Name()))
				break
			}
		}
	}
}
