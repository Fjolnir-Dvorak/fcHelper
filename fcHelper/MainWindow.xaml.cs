using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using DotLiquid;
using LibGit2Sharp;
using Microsoft.Win32;
using Microsoft.WindowsAPICodePack.Dialogs;
using MouseEventArgs = System.Windows.Input.MouseEventArgs;
using System.Globalization;
using System.Runtime.CompilerServices;
using System.Threading;
using System.Windows.Threading;

namespace fcHelper
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        private String installPath = null;
        private String gitPath = null;
        private bool closeOnButton = false;
        private const String GITNAME = "Fortress-Craft-Evolved-Translation";

        //public ObservableCollection<LogLanguageEntry> TLoggingEntries = new ObservableCollection<LogLanguageEntry>();
        public ObservableCollection<LogLanguageEntry> TLoggingEntries = LogLanguageEntry.GenTestCollection();

        public MainWindow()
        {
            InitializeComponent();
        }

        private void XAutodetectPath(object sender, RoutedEventArgs e)
        {
            String gamepath = (string) Registry.LocalMachine
                                  .OpenSubKey(@"SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 254200\")
                                  ?.GetValue("InstallLocation", "") ?? "";
            this.xInstallPath.Text = gamepath;
            this.installPath = gamepath;
            this.validateForNextButton();
        }

        private void XSelectPath(object sender, RoutedEventArgs e)
        {
            var cofd = new CommonOpenFileDialog
            {
                IsFolderPicker = true,
                Title = "Installation Directory of FortressCraft Evolved"
            };
            var result = cofd.ShowDialog();
            if (result == CommonFileDialogResult.Ok)
            {
                this.xInstallPath.Text = cofd.FileName;
                this.installPath = cofd.FileName;
            }

            this.validateForNextButton();
        }
        private void XSelectGitPath(object sender, RoutedEventArgs e)
        {
            var cofd = new CommonOpenFileDialog
            {
                IsFolderPicker = true,
                Title = "Installation Directory of FortressCraft Evolved"
            };
            var result = cofd.ShowDialog();
            if (result == CommonFileDialogResult.Ok)
            {
                this.xGitPath.Text = cofd.FileName;
                this.gitPath = cofd.FileName;
            }

            this.validateForNextButton();
        }

        private void XMouseEnter(object sender, MouseEventArgs e)
        {
            // if (this.installPath != null && this.gitPath != null)
            //     this.xNext.Visibility = Visibility.Visible;
        }

        private void XMouseLeave(object sender, MouseEventArgs e)
        {
            // this.xNext.Visibility = Visibility.Hidden;
        }

        private void DoAllTheProcessing(bool doEnglish, string defaultLang)
        {
            DirectoryInfo languagePath = new DirectoryInfo(Path.Combine(this.gitPath, GITNAME, "res"));
            var templatePath = new DirectoryInfo(Path.Combine(this.gitPath, GITNAME, "templates", "Handbook"));
            DirectoryInfo gameHandbookPath = new DirectoryInfo(Path.Combine(this.installPath, "64", "Default", "Handbook"));
            DirectoryInfo gameMasterLangPath = new DirectoryInfo(Path.Combine(this.installPath, "64", "Default", "Lang"));
            string masterLangName = "language_data_{0}.xml";

            ConcurrentDictionary<string, ConcurrentDictionary<string, Template>> templateCollection =
                XmlParser.ReadAllHandbooks(templatePath);

            var templateCount = templateCollection.Count;

            var languageList = languagePath.GetDirectories();
            //foreach (var singleLangDir in languageList)
            Parallel.ForEach(languageList, singleLangDir =>
            //ForEach(languageList, singleLangDir =>
            {
                var name = singleLangDir.Name;
                if (name.StartsWith("values-") && name.Length >= 8)
                {
                    var langCode = name.Substring(7);
                    var langName = LanguageHelper.getLanguageName(langCode);
                    if (langName == null)
                    {
                        Debug.WriteLine("ERROR: Language does not exist: " + langCode);
                        Application.Current.Shutdown();
                        return;
                    }

                    bool isDefault = langCode.Equals(defaultLang);

                    var logLangName = (isDefault) ? "testLang" : langName;
                    LogLanguageEntry logEntry = new LogLanguageEntry(logLangName, 0, templateCount + 1, 0, "parsing handbooks");
                    updateTLoggingEntries(logEntry);

                    //foreach (var singleHandbook in templateCollection)
                    //Parallel.ForEach(templateCollection, singleHandbook =>
                    ForEach(templateCollection, singleHandbook =>
                    {
                        doingHandbookStuff(langName, singleHandbook.Key, singleHandbook.Value, singleLangDir, gameHandbookPath, isDefault);
                        logEntry.increaseProcess();
                    });
                    if (doEnglish)
                    {
                        logEntry.Message = "parsing Masterfile";
                        doingMasterfileRepairation(langName, singleLangDir, gameMasterLangPath, masterLangName, isDefault);
                    } else
                    {
                        logEntry.Message = "skipped Masterfile";
                    }
                    logEntry.increaseProcess();
                }
                else if (name.Equals("values"))
                {
                    if (!doEnglish)
                    {
                        return;
                    }
                    LogLanguageEntry logEntry = new LogLanguageEntry("English", 0, templateCount + 1, 0, "parsing handbooks");
                    updateTLoggingEntries(logEntry);

                    Parallel.ForEach(templateCollection, singleHandbook =>
                    //ForEach(templateCollection, singleHandbook =>
                    {
                        doingHandbookStuff(null, singleHandbook.Key, singleHandbook.Value, singleLangDir, gameHandbookPath, false);
                        logEntry.increaseProcess();
                    });
                    logEntry.Message = "parsing Masterfile";
                    doingMasterfileRepairation(null, singleLangDir, gameMasterLangPath, null, false);
                    logEntry.increaseProcess();
                }
                else
                {
                    // Invalid directory...
                }
            });
            Application.Current.Dispatcher.Invoke(new Action(() =>
            {
                this.closeOnButton = true;
                this.xNext.Content = "Close";
                this.xNext.IsEnabled = true;
            }));
        }

        private void XNext_OnClick(object sender, RoutedEventArgs e)
        {
            if (this.closeOnButton == true)
            {
                Application.Current.Shutdown();
                return;
            }

            bool doEnglish = xGenEnglish.IsChecked ?? false;
            // Get the default language:
            string defaultLang = (string)this.xDefaultLanguageSelect.SelectedItem;
            ThreadPool.QueueUserWorkItem(O => DoAllTheProcessing(doEnglish, defaultLang));
            this.xNext.IsEnabled = false;
        }

        private static void ForEach<T>(IEnumerable<T> items, Action<T> action)
        {
            foreach (var item in items)
            {
                action(item);
            }
        }

        [MethodImpl(MethodImplOptions.Synchronized)]
        private void updateTLoggingEntries(LogLanguageEntry log)
        {
            Application.Current.Dispatcher.Invoke(new Action(() =>
            {
                TLoggingEntries.Add(log);
                xLoggingBox.ItemsSource = TLoggingEntries;
                xLoggingBox.SelectedIndex = xLoggingBox.Items.Count - 1;
                xLoggingBox.ScrollIntoView(xLoggingBox.SelectedItem);
                xLoggingBox.UpdateLayout();
            }));
        }

        private static void doingMasterfileRepairation(string langName, DirectoryInfo singleLangDir, DirectoryInfo gameMasterLangPath, string masterLangName, bool isDefault)
        {
            langName = langName ?? "english";
            var tmpFilePath = singleLangDir.FullName;
            var tmpFileName = "master.xml";
            Debug.WriteLine(langName + "::" + tmpFileName);
            var fileInfo = new FileInfo(Path.Combine(tmpFilePath, tmpFileName));
            if (!fileInfo.Exists)
            {
                Debug.WriteLine("::WARN:: " + langName + " has no MasterFile.");
                return;
            }
            Directory.CreateDirectory(gameMasterLangPath.FullName);
            var fileContent = XmlParser.correctMasterTree(fileInfo);
            if (langName.Equals("english"))
            {
                File.WriteAllText(Path.Combine(gameMasterLangPath.FullName,
                    "master_language_data.xml"), fileContent);
            }
            else
            {
                File.WriteAllText(Path.Combine(gameMasterLangPath.FullName,
                    String.Format(masterLangName, langName)), fileContent);
                if (isDefault)
                {
                    File.WriteAllText(Path.Combine(gameMasterLangPath.FullName,
                        String.Format(masterLangName, "testlang")), fileContent);
                }
            }
        }

        private static void doingHandbookStuff(string langName, string handbook, ConcurrentDictionary<string, Template> templates, DirectoryInfo singleLangDir, DirectoryInfo gameHandbookPath, bool isDefault)
        {
            langName = langName ?? "english";
            var tmpFilePath = singleLangDir.FullName;
            var tmpFileName = "handbook-" + handbook + ".xml";
            Debug.WriteLine(langName + "::" + tmpFileName);
            var fileInfo = new FileInfo(Path.Combine(tmpFilePath, tmpFileName));
            if (!fileInfo.Exists)
            {
                Debug.WriteLine("::WARN:: " + langName + " has no " + handbook + ".");
                return;
            }
            var masterfile = fileInfo;
            string goalLanguageDir;
            if (langName.Equals("english"))
            {
                goalLanguageDir = gameHandbookPath.CreateSubdirectory(handbook).FullName;
            }
            else
            {
                goalLanguageDir = gameHandbookPath.CreateSubdirectory(Path.Combine(handbook, langName)).FullName;
            }
            var dict = XmlParser.Read(masterfile);
            //<Filename, Content>
            ConcurrentDictionary<string, string> parsedCollection = XmlParser.ParseTemplates(templates, dict);

            if (isDefault)
            {
                var defaultLanguageDir =
                    gameHandbookPath.CreateSubdirectory(Path.Combine(handbook, "testlang")).FullName;
                foreach (KeyValuePair<string, string> singleParsed in parsedCollection)
                {
                    var filename = singleParsed.Key;
                    var content = singleParsed.Value;
                    File.WriteAllText(Path.Combine(goalLanguageDir, filename), content, Encoding.UTF8);
                    File.WriteAllText(Path.Combine(defaultLanguageDir, filename), content, Encoding.UTF8);
                }
            }
            else
            {
                foreach (KeyValuePair<string, string> singleParsed in parsedCollection)
                {
                    var filename = singleParsed.Key;
                    var content = singleParsed.Value;
                    File.WriteAllText(Path.Combine(goalLanguageDir, filename), content, Encoding.UTF8);
                }
            }

        }

        private void XUseLocalAppData(object sender, RoutedEventArgs e)
        {
            var gitPath = Path.Combine(Environment.GetFolderPath(System.Environment.SpecialFolder.LocalApplicationData),
                "FjolnirDvorak",
                "fcHelper");
            this.xGitPath.Text = gitPath;
            this.gitUpdateOrInit(gitPath);

            var languagePath = Path.Combine(gitPath, GITNAME, "res");
            if (Directory.Exists(languagePath))
            {
                Debug.WriteLine("Directory exists.");
                List<String> langCodes = new List<String>();
                var subDirs = new DirectoryInfo(languagePath).GetDirectories();
                Debug.WriteLine("Found " + subDirs.Length + " potential languages");
                this.xDefaultLanguageSelect.Items.Clear();
                foreach (var subDir in subDirs)
                {
                    var name = subDir.Name;
                    Debug.WriteLine("Validating directory: " + name);
                    if (name.StartsWith("values-") && name.Length >= 8)
                    {
                        var langCode = name.Substring(7);
                        Debug.WriteLine("And I am adding an element to the dropdown menu: " + langCode);
                        langCodes.Add(langCode);
                        this.xDefaultLanguageSelect.Items.Add(langCode);
                    }
                }
                this.xDefaultLanguageSelect.UpdateLayout();
                this.gitPath = gitPath;
            }
            this.validateForNextButton();
        }

        private void gitUpdateOrInit(string gitPath)
        {
            var path = Path.Combine(gitPath, GITNAME);
            if (!Directory.Exists(path))
            {
                LibGit2Sharp.CloneOptions co = new CloneOptions();
                co.BranchName = "weblate-master";
                LibGit2Sharp.Repository.Clone("https://github.com/zebra1993/Fortress-Craft-Evolved-Translation.git",
                    path);
            }
            
            using (var repo = new Repository(path))
            {
                if (false)
                {
                    if (repo.Head != repo.Branches["origin/weblate-master"])
                    {
                        var localBranch = repo.Branches["origin/weblate-master"];
                        if (localBranch == null)
                        {
                            foreach (var repoBranch in repo.Branches)
                            {
                                Debug.WriteLine(repoBranch.IsRemote + " | " + repoBranch.FriendlyName + " | " + repoBranch.RemoteName);
                            }
                            Application.Current.Shutdown();
                            return;
                        }

                        LibGit2Sharp.CheckoutOptions cm = new CheckoutOptions();
                        cm.CheckoutModifiers = CheckoutModifiers.Force;
                        LibGit2Sharp.Commands.Checkout(repo, localBranch, cm);
                    }
                    LibGit2Sharp.PullOptions options = new LibGit2Sharp.PullOptions();
                    options.FetchOptions = new FetchOptions();
                    options.MergeOptions = new MergeOptions();
                    options.MergeOptions.FastForwardStrategy = FastForwardStrategy.FastForwardOnly;
                    options.MergeOptions.FileConflictStrategy = CheckoutFileConflictStrategy.Theirs;
                    options.MergeOptions.MergeFileFavor = MergeFileFavor.Theirs;
                    LibGit2Sharp.Commands.Pull(repo, new Signature("fcHelper", "fcHelper@demo.com", new DateTimeOffset(DateTime.Now)), options);
                }
                FetchOptions fOptions = new FetchOptions();
                string logMessage = "";
                foreach (Remote remote in repo.Network.Remotes)
                {
                    IEnumerable<string> refSpecs = remote.FetchRefSpecs.Select(x => x.Specification);
                    Commands.Fetch(repo, remote.Name, refSpecs, null, logMessage);
                }
                Debug.WriteLine(logMessage);
                repo.Reset(ResetMode.Hard, repo.Branches["origin/weblate-master"].Tip);
            }
            
        }

        private void validateForNextButton()
        {
            if (this.installPath != null && this.gitPath != null && this.xDefaultLanguageSelect.SelectedItem != null)
                this.xNext.Visibility = Visibility.Visible;
        }

        private void XDefaultLanguageSelect_OnSelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            this.validateForNextButton();
        }

        private void XDeleteAppData(object sender, RoutedEventArgs e)
        {
            var pathParent = Path.Combine(Environment.GetFolderPath(System.Environment.SpecialFolder.LocalApplicationData),
                "FjolnirDvorak");
            var path = Path.Combine(pathParent, "fcHelper");
            LogLanguageEntry logEntry = new LogLanguageEntry("Config", 0, 1, 0, "deleting config directory");
            updateTLoggingEntries(logEntry);
            DeleteDirectory(path, true);
            DeleteDirectory(pathParent, false);
            logEntry.increaseProcess();
            //Directory.Delete(path, true);
            //Directory.Delete(pathParent, false);
            return;
        }

        public void DeleteDirectory(string path, bool recursive)
        {
            var currentDir = new DirectoryInfo(path);
            if (!currentDir.Exists)
            {
                return;
            }
            if (recursive)
            {
                var subfolders = Directory.GetDirectories(path);
                foreach (var s in subfolders)
                {
                    DeleteDirectory(s, recursive);
                }
            }
            var files = Directory.GetFiles(path);
            foreach (var f in files)
            {
                try
                {
                    var attr = File.GetAttributes(f);
                    if ((attr & FileAttributes.ReadOnly) == FileAttributes.ReadOnly)
                    {
                        File.SetAttributes(f, attr ^ FileAttributes.ReadOnly);
                    }
                    File.Delete(f);
                }
                catch (IOException)
                {
                    //IOErrorOnDelete = true;
                }
            }

            // At this point, all the files and sub-folders have been deleted.
            // So we delete the empty folder using the OOTB Directory.Delete method.
            Directory.Delete(path);
        }
    }
}
