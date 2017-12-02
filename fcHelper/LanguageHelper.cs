using System;
using System.Collections.Generic;

namespace fcHelper
{
    class LanguageHelper
    {
        private readonly Dictionary<string, string> codeToEnglish = new Dictionary<string, string>();
        private static readonly LanguageHelper inst = new LanguageHelper();

        private LanguageHelper()
        {
            codeToEnglish.Add("zh", "chinese");
            codeToEnglish.Add("zh-rCN", "chinese_simple");
            codeToEnglish.Add("cs", "czech");
            codeToEnglish.Add("nl", "dutch");
            codeToEnglish.Add("en", "english");
            codeToEnglish.Add("fi", "finnish");
            codeToEnglish.Add("fr", "french");
            codeToEnglish.Add("de", "german");
            codeToEnglish.Add("ko", "korean");
            codeToEnglish.Add("nn", "norwegian_nynorsk");
            codeToEnglish.Add("pl", "polish");
            codeToEnglish.Add("pt-rBR", "portuguese_brazil");
            codeToEnglish.Add("ru", "russian");
            codeToEnglish.Add("es", "spanish");
            codeToEnglish.Add("sv", "swedish");
            codeToEnglish.Add("tr", "turkish");
            codeToEnglish.Add("uk", "ukrainian");
        }

        public static string getLanguageName(string code)
        {
            try
            {
                return inst.codeToEnglish[code];
            }
            catch (Exception)
            {
                return null;
            }
        }
    }
}