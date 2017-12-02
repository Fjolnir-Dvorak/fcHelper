using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using System.Xml.Linq;
using DotLiquid;

namespace fcHelper
{
    internal static class XmlParser
    {
        // ReadHandbookMasterFile
        public static ConcurrentDictionary<string, object> Read(FileInfo file)
        {
            var dic = new ConcurrentDictionary<string, object>();
            XDocument doc = XDocument.Load(file.FullName);
            foreach (XElement el in doc.Root.Elements())
            {
                var tag = el.Attribute("name").Value.Replace(".", "_");
                var translation = el.Value;
                var escapedTranslation = System.Security.SecurityElement.Escape(translation);
                dic.TryAdd(tag, escapedTranslation);

            }
            return dic;
        }
        // ReadHandbookTemplatesIntoBuffer
        public static ConcurrentDictionary<string, Template> ReadHandbookTemplates(DirectoryInfo templateDir)
        {
            ConcurrentDictionary<string, Template> templates = new ConcurrentDictionary<string, Template>();
            foreach (FileInfo fileInfo in templateDir.GetFiles("*.xml"))
            {
                String text;
                try
                {
                    using (StreamReader sr = new StreamReader(fileInfo.FullName))
                    {
                        text = sr.ReadToEnd();
                    }
                    text = text.Replace("{{ index .Data \"", "{{ ")
                        .Replace("\" }}", " }}")
                        .Replace(".", "_");
                    //text = text.Replace("{{ index .Data \"", "{{ this.dictionary[\"")
                    //    .Replace("\" }}", "\"] }}");
                }
                catch (Exception e)
                {
                    continue;
                }
                Template template = Template.Parse(text);
                templates.TryAdd(fileInfo.Name, template);
            }
            return templates;
        }
        // ParseTemplates
        public static ConcurrentDictionary<string, string> ParseTemplates(ConcurrentDictionary<string, Template> templates, ConcurrentDictionary<string, object> dic)
        {
            ConcurrentDictionary<string, string> parsed = new ConcurrentDictionary<string, string>();
            var simpldDic = dic.ToDictionary(kvp => kvp.Key, kvp => (object) kvp.Value);
            foreach (var fileTemplate in templates)
            {
                //var parsedSingle = fileTemplate.Value.Render(Hash.FromAnonymousObject(new { @this = new DictionaryWrapper(simpldDic) }));
                var parsedSingle = fileTemplate.Value.Render(Hash.FromDictionary(simpldDic));
                parsed.TryAdd(fileTemplate.Key, parsedSingle);
            }
            return parsed;
        }
        // CorrectMasterFileLayout
        public static string correctMasterTree(FileInfo file)
        {
            XDocument doc = XDocument.Load(file.FullName);
            doc.Root.Name = "resource";
            return doc.ToString();
        }

        public static ConcurrentDictionary<string, ConcurrentDictionary<string, Template>> ReadAllHandbooks(DirectoryInfo templatePath)
        {
            var toReturn = new ConcurrentDictionary<string, ConcurrentDictionary<string, Template>>();
            var allHandbooks = templatePath.GetDirectories();
            Parallel.ForEach(allHandbooks, singleHandbookDir =>
            {
                var name = singleHandbookDir.Name;
                toReturn.TryAdd(name, ReadHandbookTemplates(singleHandbookDir));
            });
            return toReturn;
        }
    }

    internal class DictionaryWrapper
    {
        public Dictionary<string, object> dictionary { get; set; }

        public DictionaryWrapper(Dictionary<string, object> dic)
        {
            this.dictionary = dic;
        }
    }
}