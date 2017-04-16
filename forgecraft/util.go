package forgecraft

import (
	"fmt"
	"github.com/beevik/etree"
)

var KeyContainer = make([]string, 0)

func getChildText(element *etree.Element) []string {
	children := element.ChildElements()
	strings := make([]string, len(children))
	for index, child := range children {
		strings[index] = child.Text()
	}
	return strings
}

func utilIsInKeyContainer(element string) bool {
	if contains(KeyContainer, element) {
		return true
	} else {
		KeyContainer = append(KeyContainer, element)
		return false
	}
}

func utilResetKeyContainer() {
	KeyContainer = make([]string, 0)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func utilHandleUnknownElement(element string) {
	if !utilIsInKeyContainer(element) {
		fmt.Println(element)
	}
}
