package game

import (
	"bytes"
	"fmt"
	cmd2 "github.com/Fjolnir-Dvorak/fcHelper/cmd"
	"github.com/Fjolnir-Dvorak/fcHelper/cmd/game/structures"
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
)

var (
	TransProject   string
	DeployDest     string
	TestLang       string
	UseCommit      string
	gitCloneConfig = &git.CloneOptions{
		URL:      "https://github.com/zebra1993/Fortress-Craft-Evolved-Translation.git",
		Progress: os.Stdout,
	}
	gitPullConfig = &git.PullOptions{
		ReferenceName: structures.GitBranch,
		RemoteName:    structures.GitRemote,
	}
	gitCeckoutConfig = &git.CheckoutOptions{
		Branch: structures.GitBranch,
	}
)

type mapp struct {
	Data map[string]string
}

func init() {
}

func handleFatalError() {
	// TODO IMPLEMENT ME!
}

func UpdateGit(path string, name string) *git.Repository {
	var gitPath = filepath.Join(path, name)
	repo, err := git.PlainOpen(gitPath)
	if err == git.ErrRepositoryNotExists {
		repo, err = git.PlainClone(gitPath, false, gitCloneConfig)
		if err != nil {
			handleFatalError()
		}
	}
	worktree, err := repo.Worktree()
	worktree.Checkout(gitCeckoutConfig)
	err = repo.Pull(gitPullConfig)
	if err != nil {
		handleFatalError()
	}
	return repo
}

func ValidateProject(project string) bool {
	fileInfo, err := os.Stat(project)
	if err != nil {
		handleFatalError()
	}
	if !fileInfo.IsDir() {
		handleFatalError()
	}
	return true
	// TODO IMPLEMENT ME!
}

func DoCreate(cmd *cobra.Command, args []string) {
	// 1. Test if git repository is initialized
	// 2. Test if git repository is clean to update
	// 3. Update git repository
	// 3.5 Checkout commit if wished.
	// 4. Pull registry key for FortressCraft Evolved to get the installation directory
	// 5. Fill the templates with content and deploy all the stuff.
	var repoDir string
	var valid bool
	if TransProject == "" {
		_ = UpdateGit(cmd2.Environ.DataLocal(), structures.GitName)
		repoDir = filepath.Join(cmd2.Environ.DataLocal(), structures.GitName)
		valid = ValidateProject(repoDir)
	} else {
		valid = ValidateProject(TransProject)
		repoDir = TransProject
	}
	if !valid {
		handleFatalError()
	}
	translationDir := filepath.Join(repoDir, structures.GitResDir)
	templateDir := filepath.Join(repoDir, structures.GitTemplateDir)

	var usedGameDir string

	if DeployDest != "" {
		// Tag deployDest was specified. This has the highest priority
		usedGameDir = DeployDest
	} else if GameDir != "" {
		// The user specified another installation directory
		usedGameDir = GameDir
	} else {
		// Default Steam installation directory
		usedGameDir = SteamGameDir
	}
	handbookBase := filepath.Join(usedGameDir, structures.Handbook)
	langBase := filepath.Join(usedGameDir, structures.Lang)

	// Read template files into a map.
	var translationLanguages, _ = ioutil.ReadDir(translationDir)
	for _, langInfo := range translationLanguages {
		langDirName := langInfo.Name()
		langDir := filepath.Join(translationDir, langDirName)

		splitted := strings.Split(langDirName, structures.GitLangDirSeparator)
		var langCode string

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
				currentHandbook := strings.Split(
					strings.Split(langFileInfo.Name(), "-")[1],
					".")[0]
				handbook := filepath.Join(handbookBase, currentHandbook)
				templates := filepath.Join(templateDir, "Handbook", currentHandbook)
				createHandbookFiles(file, handbook, langCode, templates)
			} else if langFileInfo.Name() == "master.xml" {
				createMasterFile(file, langCode, langBase)
			}
		}
	}
}

func createHandbookFiles(file, handbook, language, templates string) {
	//Handbook
	_, data := util.ParseXLFMap(file)
	temp := mapp{Data: data}

	var outputDir string
	if language == "english" {
		outputDir = handbook
	} else {
		outputDir = filepath.Join(handbook, language)
	}

	fmt.Printf("        Creating Directory if not existent: %s\n", outputDir)
	os.MkdirAll(outputDir, os.ModePerm)
	toParse, _ := ioutil.ReadDir(templates)
	for _, tempFile := range toParse {
		templateBase := filepath.Join(templates, tempFile.Name())
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
}

func createMasterFile(file, language, out string) {
	fmt.Println("Changing the first two nodes of " + file)
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(file); err != nil {
		panic(err)
	}
	newDoc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	newTag := newDoc.CreateElement("languages")
	for _, parent := range doc.ChildElements() {
		parent.Tag = language
		newTag.AddChild(parent)
	}
	var filename string
	if language == "english" {
		filename = "master_language_data.xml"
	} else {
		filename = "language_data_" + language + ".xml"
	}
	os.MkdirAll(out, os.ModePerm)
	outputFile := filepath.Join(out, filename)
	newDoc.WriteToFile(outputFile)
}
