package util

import (
	"bytes"
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
