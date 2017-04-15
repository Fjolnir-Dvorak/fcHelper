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

	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type KeyValue struct {
	Key   string
	Value string
}

var gameDir string
var destination string
var template bool
var language string

const handbook = "64/Default/Handbook"

const av = "AvailableResearch"
const co = "CompletedResearch"
const cr = "Creative"
const ma = "Materials"
const su = "Survival"

const av_key = "AV"
const co_key = "CO"
const cr_key = "CR"
const ma_key = "MA"
const su_key = "SU"

const title = "Title"
const header = "Header"
const paragraph = "Paragraph"
const left = "Left"
const right = "Right"

var keywords = [...]string{title, header, paragraph, left, right}

const surround = "{}"
const surroundAmount = 2
const out = "fc_out"
const templates = "templates"

const optNode = "Text"
const contNode = "Pages"
const filePrefix = "handbook-"

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: doExtract,
}

func doExtract(cmd *cobra.Command, args []string) {
	destDir := destination
	if template {
		destDir = filepath.Join(destination, templates, "Handbook")
	}
	_ = os.MkdirAll(destDir, os.ModePerm)
	createExcract(av, av_key, destDir)
	createExcract(co, co_key, destDir)
	createExcract(cr, cr_key, destDir)
	createExcract(ma, ma_key, destDir)
	createExcract(su, su_key, destDir)

}

func createExcract(name, namecode, temp string) {
	gd := filepath.Join(gameDir, "/", handbook, name)
	dir, _ := ioutil.ReadDir(gd)

	dest := filepath.Join(temp, name)
	if template {
		os.Mkdir(dest, os.ModePerm)
	}

	var keys []KeyValue
	for _, fileInf := range dir {
		if fileInf.IsDir() {
			continue
		}
		basefile := filepath.Join(gd, fileInf.Name())
		basename := strings.TrimSuffix(fileInf.Name(), filepath.Ext(fileInf.Name()))
		basename = strings.Replace(basename, " ", "_", -1)
		doc := etree.NewDocument()
		fmt.Printf("Reading: %s\n", basefile)
		if err := doc.ReadFromFile(basefile); err != nil {
			panic(err)
		}
		root := doc.Root()
		for _, child := range root.ChildElements() {
			codebase := namecode + "." + basename
			if child.Tag == title {
				code := codebase + "." + title
				text := child.Text()
				keys = append(keys, KeyValue{code, text})
				child.SetText(surroundCode(code))
			} else if child.Tag == contNode {
				// Iterate through the whole tree
				iterator := 0
				_, keys = recursiveIteration(child, iterator, keys, codebase)

			}
		}

		if template {
			_ = doc.WriteToFile(filepath.Join(dest, fileInf.Name()))
		}
	}

	var filename string
	if language == "" {
		directory := filepath.Join(destination, "res", "values")
		os.MkdirAll(directory, os.ModePerm)
		filename = filepath.Join(directory, filePrefix+strings.Replace(name, " ", "_", -1)+".xml")
	} else {
		directory := filepath.Join(destination, "res", "values"+"-"+language)
		os.MkdirAll(directory, os.ModePerm)
		filename = filepath.Join(directory, filePrefix+strings.Replace(name, " ", "_", -1)+".xml")
	}
	createXML(keys, filename)
}
func createXML(values []KeyValue, filename string) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	resources := doc.CreateElement("resources")
	for _, value := range values {
		key := resources.CreateElement("string")
		key.CreateAttr("name", value.Key)
		key.SetText(value.Value)
	}
	doc.WriteToFile(filename)
}

func recursiveIteration(element *etree.Element, iterator int, values []KeyValue, codebase string) (int, []KeyValue) {
	for _, currentKey := range keywords {
		if element.Tag == currentKey {
			code := fmt.Sprintf("%s.%02d.%s", codebase, iterator, currentKey)
			var text string
			if hasChild(element) {
				children := element.ChildElements()
				text = children[0].Text()
				children[0].SetText(surroundCode(code))
			} else {
				text = element.Text()
				element.SetText(surroundCode(code))
			}
			values = append(values, KeyValue{code, text})
			iterator++

			return iterator, values
		}
	}
	for _, child := range element.ChildElements() {
		iterator, values = recursiveIteration(child, iterator, values, codebase)
	}
	return iterator, values
}

func hasChild(element *etree.Element) bool {
	length := len(element.ChildElements())
	if length == 0 {
		return false
	}
	return true
}

func surroundCode(code string) string {
	surroundLeft := strings.Repeat(string(surround[0]), surroundAmount)
	surroundRight := strings.Repeat(string(surround[1]), surroundAmount)
	surrounded := surroundLeft + code + surroundRight
	return surrounded
}

func init() {
	RootCmd.AddCommand(extractCmd)
	extractCmd.Flags().StringVarP(&gameDir, "gameDir", "g", "", "Directory containing FortressCraft Evolved.")
	extractCmd.Flags().StringVarP(&destination, "destination", "d", out, "Destination Directory to create the parsed files.")
	extractCmd.Flags().StringVarP(&language, "language", "l", "", "Used language shortkey.")
	extractCmd.Flags().BoolVarP(&template, "template", "t", false, "Wether to generate xml-templates.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
