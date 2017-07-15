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
	. "github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/Fjolnir-Dvorak/fcHelper/util"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (

	// Keywords to extract
	title     = "Title"
	header    = "Header"
	paragraph = "Paragraph"
	left      = "Left"
	right     = "Right"

	// xml structure
	optNode    = "Text"
	contNode   = "Pages"
	filePrefix = "handbook-"

	// Replacekey design
	surround       = "{}"
	surroundAmount = 2
)

var (
	keywords     = [...]string{title, header, paragraph, left, right}
	Out          string
	NoTemplate   bool
	NoExtraction bool
	Lang         string
)

func init() {
	NoExtraction = false
}

func DoExtract(cmd *cobra.Command, args []string) {
	if NoExtraction && NoTemplate {
		return
	}

	if GameDir != "" {
		// The user specified another installation directory
		gameDir = GameDir
	} else {
		// Default Steam installation directory
		gameDir = SteamGameDir
	}

	// Ensure the template base directory is exsisting
	templateDir := Out
	translationDir := Out
	if NoTemplate == false {
		templateDir = filepath.Join(Out, structures.GitTemplateHandbookDir)
	}
	if NoExtraction == false {
		languageDir := structures.GitLangDirPrefix
		if Lang != "en" {
			languageDir = languageDir + structures.GitLangDirSeparator + Lang
		}
		translationDir = filepath.Join(Out, structures.GitResDir, languageDir)
	}
	_ = os.MkdirAll(templateDir, os.ModePerm)
	_ = os.MkdirAll(translationDir, os.ModePerm)
	createExtract(structures.Av, structures.Av_key, templateDir, translationDir)
	createExtract(structures.Co, structures.Co_key, templateDir, translationDir)
	createExtract(structures.Cr, structures.Cr_key, templateDir, translationDir)
	createExtract(structures.Ma, structures.Ma_key, templateDir, translationDir)
	createExtract(structures.Su, structures.Su_key, templateDir, translationDir)

}

func createExtract(name, namecode, templDir, transDir string) {
	handbookDir := filepath.Join(gameDir, structures.Handbook, name)
	dir, _ := ioutil.ReadDir(handbookDir)

	destTempl := filepath.Join(templDir, name)
	if NoTemplate == false {
		os.Mkdir(destTempl, os.ModePerm)
	}

	var keys []KeyValue
	for _, fileInf := range dir {
		if fileInf.IsDir() {
			continue
		}
		basefile := filepath.Join(handbookDir, fileInf.Name())
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

		if NoTemplate == false {
			out, _ := doc.WriteToString()
			util.WriteRawStringToFile(filepath.Join(destTempl, fileInf.Name()), out)
			//_ = doc.WriteToFile(filepath.Join(destTempl, fileInf.Name()))
		}
	}
	if NoExtraction == false {
		var filename string
		filename = filepath.Join(transDir, filePrefix+strings.Replace(name, " ", "_", -1)+".xml")
		util.CreateXLFkv(keys, filename)
	}

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
	code = " index .Data \"" + code + "\" "
	surrounded := surroundLeft + code + surroundRight
	return surrounded
}
