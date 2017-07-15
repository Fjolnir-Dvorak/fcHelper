package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Fjolnir-Dvorak/fcHelper/cmd/game"
)

var (
	RootCMD = &cobra.Command{
		Use:   "fcHelper",
		Short: "A Utility program for FortressCraft Evolved",
		Long: `Used to to ease the burden of xmls and translation.
			Also used to validate wrong xml-tag usages. For more
			Information view the help of the commands.`,
	}
	GameCMD = &cobra.Command{
		Use:   "game",
		Short: "Utilities modifying the game data.",
		Long: `This package interacts directly with the game files
			or wil generate and deploy them. To archive that it looks
			for the location of the game files inside the windows registry.`,
	}
	GenerateCMD = &cobra.Command{
		Use:   "generate",
		Short: "Manipulates or generates game files.",
		Long:  `Please specify a subcommand for gamefile manipulation.`,
	}
	ConfigCMD = &cobra.Command{
		Use:   "config",
		Short: "Gives you the ability to change, import and export the config.",
	}
	XmlCMD = &cobra.Command{
		Use:   "xml",
		Short: "XML modifying utilities.",
		Long:  "XML utilities. These are not depending on FortressCraft.",
	}
	GameCreateTranslationCMD = &cobra.Command{
		Use:   "createTranslation",
		Short: "Creates valid game files from the translation data.",
		Long: `Injects the translated Keys back into the handbook files. changes
				the base nodes from the master translation files so they are valid, too.`,
		Run: game.DoCreate,
	}
	GameExtractCMD = &cobra.Command{
		Use:   "extract",
		Short: "Extracts strings to localize..",
		Long: `Extracts strings out of the handbook which needs to be localized.
			Those Keys will be stored in an Android-xml language files.
			It also is able to generate createTemplates-files which will be used to reinject
			the translated keys back into the xml.`,
	}
	GameRepairCMD = &cobra.Command{
		Use:   "repairKeys",
		Short: "Repaires missused name-tags in the handbook",
		Long: `The handbook contains both 'Key' and 'Value' tags as primary key.
			This will sort through supported Datafiles for key-name mappings and
			will replace all used 'Name' tags with their corresponding 'Key' tag.`,
	}
	GenerateCompletionCMD = &cobra.Command{
		Use:   "completion",
		Short: "Generate bash completion",
		Long: `Generates the bash completion. The system location to save these
			into is '/etc/bash_completion.d/'. If you made a user specific
			folder for completions you could also save this script in there.`,
	}
	GenerateManCMD = &cobra.Command{
		Use:   "man",
		Short: "Generates a manual for this program",
		Long:  `Man which can be used as manpage.`,
	}
	GenerateMarkdownCMD = &cobra.Command{
		Use:   "markdown",
		Short: "Generates a markdown manual for this program",
		Long:  `Man written in markdown.`,
	}
	GenerateYamlCMD = &cobra.Command{
		Use:   "yaml",
		Short: "Generates a yaml structure for this program",
		Long:  `Structure file.`,
	}
	XmlClearCMD = &cobra.Command{
		Use:   "clear",
		Short: "Deletes everything except the tags itself. No Comments, no Content.",
		Long: `If you do not know what you are doing exactly do not use this functionality.
			It is able to corrupt your whole game and to destroy your savegames. Do not ask
			me how, but I am shure there is a possibility for that.`,
	}
	XmlCompareCMD = &cobra.Command{
		Use:   "compare",
		Short: "Compares two language xmls and extracts all different keys.",
		Long: `Compares two different files or two different directory-trees with the same
	 		structure. It will generate a file for each comparison showing missing keys (keys
	 		in second file) with '+++' in front and the surplus (keys in first file) with '---'
	 		in front. If the comparison file will be applied to the base file it will have the
	 		same key structure as the second one.`,
	}
	XmlDuplicatesCMD = &cobra.Command{
		Use:   "duplicates",
		Short: "Searches through a File in interactive modeand removes key-duplicates",
		Long: `Searches through a single file and auto-deletes all duplicate keys where the
			value is identical. On different values the user will be asked which value should
			be deleted.`,
	}
	ConfigExportCMD = &cobra.Command{
		Use:   "export",
		Short: "Exports the config",
		Long:  "Exports the current config. There are different export modes.",
	}
	ConfigImportCMD = &cobra.Command{
		Use:   "import",
		Short: "Imports the config",
		Long:  "Imports a config file",
	}
)

func initCMD() {
	RootCMD.AddCommand(GameCMD)
	RootCMD.AddCommand(GenerateCMD)
	RootCMD.AddCommand(ConfigCMD)
	RootCMD.AddCommand(XmlCMD)
	RootCMD.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+defaultConf+")")

	GameCMD.AddCommand(GameCreateTranslationCMD)
	GameCMD.AddCommand(GameExtractCMD)
	GameCMD.AddCommand(GameRepairCMD)
	GameCMD.Flags().StringVarP(&game.GameDir, "gameDir", "g", "", "Optional. Define this if you do not want to use your Steam directory.")
	GameCMD.Flags().BoolVarP(&game.NoConfig, "noConfig", "n", false, "Deactivates the parsing of the config file.")

	GameCreateTranslationCMD.Flags().StringVarP(&game.TransProject, "transProject", "s", "", "Optional. This is the path of the git translation repo. Deactivates 'useCommit' and does not pull updates from remote.")
	GameCreateTranslationCMD.Flags().StringVarP(&game.DeployDest, "deployDest", "d", "", "Optional. Specify if you do not want to deploy into gameDir")
	GameCreateTranslationCMD.Flags().StringVarP(&game.TestLang, "TestLang", "l", "de", "Defines which language should be used for TestLang. Language-Code.")
	GameCreateTranslationCMD.Flags().StringVarP(&game.UseCommit, "useCommit", "c", "", "Specifies a commit which should be checked out and used.")

	GameExtractCMD.Flags().StringVarP(&game.Out, "out", "o", "fcOut", "Output Directory. Will be created if not available.")
	GameExtractCMD.Flags().StringVarP(&game.Lang, "lang", "l", "en", "Language of the extracted content.")
	GameExtractCMD.Flags().BoolVarP(&game.NoTemplate, "noTemplate", "t", false, "Deactivates the template generation.")

	GameRepairCMD.Flags().BoolVarP(&game.NameToKey, "nameToKey", "", false, "Changes Name-Tags from items to Key-Tags.")
	GameRepairCMD.Flags().BoolVarP(&game.NoSpace, "noSpace", "", false, "Replaces Spaces with underscored.")

	GenerateCMD.AddCommand(GenerateCompletionCMD)
	GenerateCMD.AddCommand(GenerateManCMD)
	GenerateCMD.AddCommand(GenerateMarkdownCMD)
	GenerateCMD.AddCommand(GenerateYamlCMD)

	// TODO define output key

	ConfigCMD.AddCommand(ConfigExportCMD)
	ConfigCMD.AddCommand(ConfigImportCMD)

	ConfigExportCMD.Flags().BoolVarP(&config.Full, "full", "f", false, "Exports the config containing all default keys.")
	ConfigExportCMD.Flags().BoolVarP(&config.Default, "default", "d", false, "Creates a config file with the programmed defaults.")
	ConfigExportCMD.Flags().BoolVarP(&config.Overwrite, "overwrite", "w", false, "Overwrites the programm config with the generated config file.")
	ConfigExportCMD.Flags().StringVarP(&config.Out, "out", "o", "", "Specifies the output channel. Default is STDOut.")

	ConfigImportCMD.Flags().BoolVarP(&config.IgnoreDef, "ignoreDef", "i", false, "Does not impoort keys specified with the default value.")

}
