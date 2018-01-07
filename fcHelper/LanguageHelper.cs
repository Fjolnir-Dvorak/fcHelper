using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Globalization;

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
            codeToEnglish.Add("it", "italian");
        }

        public static string getLanguageName(string code)
        {
            if (code == "zh-rCN")
            {
                code = "zh-Hans";
            } else if (code == "zh-rTW")
            {
                code = "zh-Hant";
            }

            code = code.Replace("-r", "-");

            Debug.WriteLine(code);

            try
            {
                return CultureInfo.GetCultureInfoByIetfLanguageTag(code).EnglishName;
            }
            catch (Exception)
            {
                return null;
            }

            //try
            //{
            //    Debug.WriteLine("I am here: " + code);
            //    return inst.codeToEnglish[code];
            //}
            //catch (Exception)
            //{
            //    Debug.WriteLine("And I am there: " + code);
            //    try
            //    {
            //        return CultureInfo.GetCultureInfoByIetfLanguageTag(code).EnglishName;
            //    }
            //    catch (Exception)
            //    {
            //        Debug.WriteLine("Oh, I am everywhere: " + code);
            //        return null;
            //    }
            //}
        }
    }
}