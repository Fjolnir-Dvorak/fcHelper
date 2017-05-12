package util

var (
	LangCodes = make(map[string]string)
)

func init() {
	LangCodes["de"] = "german"
	LangCodes["fi"] = "finnish"
	LangCodes["fr"] = "french"
	LangCodes["nl"] = "dutch"
	LangCodes["ru"] = "russian"
	LangCodes["tr"] = "turkish"
	LangCodes["uk"] = "ukrainian"
}
