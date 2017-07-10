package game

import (
	"bytes"
	"fmt"
	"github.com/Fjolnir-Dvorak/fcHelper/util"
	"github.com/beevik/etree"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	cmd2 "github.com/Fjolnir-Dvorak/fcHelper/cmd"
	"github.com/hashicorp/hcl/json/token"
)

const (
	gitBranch = "weblate-master"
	gitName = "Fortress-Craft-Evolved-Translation"
)

var (
	TransProject string
	DeployDest   string
	TestLang     string
	UseCommit    string
	gitCloneConfig = &git.CloneOptions{
		URL: "https://github.com/zebra1993/Fortress-Craft-Evolved-Translation.git",
		Progress: os.Stdout,
	}
	gitPullConfig = &git.PullOptions{
		ReferenceName: gitBranch,
	}
)

type mapp struct {
	Data map[string]string
}

func init() {
}

func handleFatalError() {

}

func DoCreate(cmd *cobra.Command, args []string) {
	// 1. Test if git repository is initialized
	// 2. Test if git repository is clean to update
	// 3. Update git repository
	// 3.5 Checkout commit if wished.
	// 4. Pull registry key for FortressCraft Evolved to get the installation directory
	// 5. Fill the templates with content and deploy all the stuff.

	var gitPath = cmd2.Environ.DataLocal()
	repo, err := git.PlainOpen(filepath.Join(gitPath, gitName))
	if err == git.ErrRepositoryNotExists {
		repo, err = git.PlainClone(gitPath, false, gitCloneConfig)
		if err != nil {
			handleFatalError()
		}
	}

	err = repo.Pull(gitPullConfig)
	if err != nil {
		handleFatalError()
	}

	if !validInputCreate() {
		return
	}

	// Read template files into a map.

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
				_, data := util.ParseXLFMap(file)
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
