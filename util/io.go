package util

import (
	"bytes"
	"fmt"
	. "github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/beevik/etree"
	"html"
	"io"
	"os"
	"regexp"
	"strings"
)

func WriteRawStringToFile(filepath, s string) error {
	re := regexp.MustCompile("(&quot;)|(&#34;)")
	cleaned := re.ReplaceAll([]byte(s), []byte("\""))
	re = regexp.MustCompile("\n\n")
	cleaned = re.ReplaceAll([]byte(cleaned), []byte("\n"))
	//s = html.UnescapeString(s)
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, bytes.NewReader(cleaned))
	if err != nil {
		return err
	}

	return nil
}

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}

// Parse xml Language File into a KeyNameList
func ParseXLFKeyName(filename string) (error, KeyNameList) {
	fmt.Println("Reading Keys from " + filename)
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filename); err != nil {
		return err, InitEmptyKeyNameList()
	}
	var values []KeyName
	//Bypass the resources tag.
	for _, parent := range doc.ChildElements() {
		values = make([]KeyName, len(parent.ChildElements()))
		for index, child := range parent.ChildElements() {
			key := child.SelectAttrValue("name", "")
			value := child.Text()
			value = html.EscapeString(value)
			values[index] = KeyName{Key: key, Name: value}
		}
	}
	return nil, CreateKeyNameList(values)
}

// Parse xml Language File into a Map
func ParseXLFMap(filename string) (error, map[string]string) {
	fmt.Println("Reading Keys from " + filename)
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filename); err != nil {
		return err, nil
	}
	values := make(map[string]string)
	//Bypass the resources tag.
	for _, parent := range doc.ChildElements() {
		for _, child := range parent.ChildElements() {
			key := child.SelectAttrValue("name", "")
			value := child.Text()
			value = html.EscapeString(value)
			values[key] = value
		}
	}
	return nil, values
}

// Create xml Language File
func CreateXLFknl(values KeyNameList, filename string) {
	CreateXLFkn(values.Iterate(), filename)
}

// Create xml Language File
func CreateXLFkn(values []KeyName, filename string) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	resources := doc.CreateElement("resources")
	for _, value := range values {
		key := resources.CreateElement("string")
		key.CreateAttr("name", value.Key)
		key.SetText(value.Name)
	}
	doc.WriteToFile(filename)
}

// Create xml Language File
func CreateXLFkv(values []KeyValue, filename string) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	resources := doc.CreateElement("resources")
	for _, value := range values {
		key := resources.CreateElement("string")
		key.CreateAttr("name", value.Key)
		key.SetText(value.Value)
	}
	doc.WriteToFile(filename)
}
