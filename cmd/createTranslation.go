package cmd

import (
	"bytes"
	"fmt"
	"github.com/Fjolnir-Dvorak/fcHelper/util"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	translationDir string
	templateDir    string
	createOut      string
)

type mapp struct {
	Data map[string]string
}

// createCmd represents the createTranslation command
var createCmd = &cobra.Command{
	Use:   "createTranslation",
	Short: "Creates valid game files from the translation data.",
	Long: `Injects the translated Keys back into the handbook files. changes
the base nodes from the master translation files so they are valid, too.`,
	Run: doCreate,
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&translationDir, "translationFiles", "g", filepath.Join(out, "res"), "Directory containing translated Files (calles 'res').")
	createCmd.Flags().StringVarP(&templateDir, "templates", "d", filepath.Join(out, "templates"), "Directory containing the handbook templates.")
	createCmd.Flags().StringVarP(&createOut, "destination", "l", filepath.Join(out, "translated"), "Destination directory.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func doCreate(cmd *cobra.Command, args []string) {
	if !validInputCreate() {
		return
	}

	// Read createTemplates files into a map.

	var translationLanguages, _ = ioutil.ReadDir(translationDir)

	for _, langInfo := range translationLanguages {
		langDirName := langInfo.Name()
		langDir := filepath.Join(translationDir, langDirName)

		splitted := strings.Split(langDirName, "-")
		langCode := ""

		if len(splitted) == 2 {
			langCode = splitted[1]
		}
		if langCode == "" {
			langCode = "english"
		} else {
			if val, ok := util.LangCodes[langCode]; ok {
				langCode = val
			}
		}

		langFiles, _ := ioutil.ReadDir(langDir)
		for _, langFileInfo := range langFiles {
			// TODO Make this more recursive
			file := filepath.Join(langDir, langFileInfo.Name())
			if strings.HasPrefix(langFileInfo.Name(), "handbook") {
				data := createMap(file)
				temp := mapp{Data: data}

				handbookDirName := strings.Split(strings.Split(langFileInfo.Name(), "-")[1], ".")[0]

				var outputDir string
				if langCode == "english" {
					outputDir = filepath.Join(createOut, "Handbook", handbookDirName)
				} else {
					outputDir = filepath.Join(createOut, "Handbook", handbookDirName, langCode)
				}
				handbookDir := filepath.Join(templateDir, "Handbook", handbookDirName)

				fmt.Printf("        Creating Directory if not existent: %s\n", outputDir)
				os.MkdirAll(outputDir, os.ModePerm)
				toParse, _ := ioutil.ReadDir(handbookDir)
				for _, tempFile := range toParse {
					templateBase := filepath.Join(handbookDir, tempFile.Name())
					fmt.Printf("- Parsing %s\n", templateBase)

					t, _ := template.ParseFiles(templateBase)
					var b bytes.Buffer
					_ = t.Execute(&b, temp)

					re := regexp.MustCompile("<Key>(.*?)</Key>")
					match := re.FindStringSubmatch(b.String())
					filename := tempFile.Name()
					if len(match) >= 1 {
						filename = match[1] + ".xml"
					}

					outputFile := filepath.Join(outputDir, filename)
					util.WriteStringToFile(outputFile, b.String())
				}
			} else if langFileInfo.Name() == "master.xml" {
				fmt.Println("Changing the first two nodes of " + file)
				doc := etree.NewDocument()
				if err := doc.ReadFromFile(file); err != nil {
					panic(err)
				}
				newDoc := etree.NewDocument()
				doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
				newTag := newDoc.CreateElement("languages")
				for _, parent := range doc.ChildElements() {
					parent.Tag = langCode
					newTag.AddChild(parent)
				}
				outputDir := filepath.Join(createOut, "master", langCode)
				os.MkdirAll(outputDir, os.ModePerm)
				outputFile := filepath.Join(outputDir, "master.xml")
				newDoc.WriteToFile(outputFile)
			}
		}
	}
}

func readTemplates() {

}
func createMap(filename string) map[string]string {
	fmt.Println("Reading Keys from " + filename)
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filename); err != nil {
		panic(err)
	}
	values := make(map[string]string)
	for _, parent := range doc.ChildElements() {
		for _, child := range parent.ChildElements() {
			key := child.SelectAttrValue("name", "")
			value := child.Text()
			values[key] = value
		}
	}
	return values
}

func validInputCreate() bool {
	fileInf, err := os.Stat(translationDir)
	if err != nil {
		return false
	}
	if !fileInf.IsDir() {
		return false
	}
	if fileInf.Name() != "res" {
		return false
	}

	fileInf, err = os.Stat(templateDir)
	if err != nil {
		return false
	}
	if !fileInf.IsDir() {
		return false
	}
	if fileInf.Name() != templates {
		return false
	}
	return true
}
