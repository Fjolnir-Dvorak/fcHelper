package util

import (
	"html"
	"io"
	"os"
	"strings"
)

func WriteStringToFile(filepath, s string) error {
	s = html.UnescapeString(s)
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
