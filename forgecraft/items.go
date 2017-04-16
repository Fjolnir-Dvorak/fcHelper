package forgecraft

import (
	"fmt"
	"github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/beevik/etree"
)

type ItemEntries struct {
	Entries []ItemEntry
}

type ItemEntry struct {
	ItemID          string
	Key             string
	Name            string
	Plural          string
	Type            string
	Object          string
	Sprite          string
	Category        string
	ResearchReqKey  []string
	ScanReqKey      []string
	ItemAction      string
	Hidden          string
	MaxDurability   string
	ActionParameter string
	DecomposeValue  string
}

type ItemEntryDef struct {
	ItemID          string
	Key             string
	Name            string
	Plural          string
	Type            string
	Object          string
	Sprite          string
	Category        string
	ResearchReqKey  string
	ScanReqKey      string
	ItemAction      string
	Hidden          string
	MaxDurability   string
	ActionParameter string
	DecomposeValue  string
}

var (
	ItemsDef = ItemEntryDef{
		ItemID:          "ItemID",
		Key:             "Key",
		Name:            "Name",
		Plural:          "Plural",
		Type:            "Type",
		Object:          "Object",
		Sprite:          "Sprite",
		Category:        "Category",
		ResearchReqKey:  "ResearchRequirements",
		ScanReqKey:      "ScanRequirements",
		ItemAction:      "ItemAction",
		Hidden:          "Hidden",
		MaxDurability:   "MaxDurability",
		ActionParameter: "ActionParameter",
		DecomposeValue:  "DecomposeValue",
	}
	ItemsFilename = "Items.xml"
	Parent = "Entries"
)

func parseItem(element *etree.Element) ItemEntry {
	children := element.ChildElements()
	var item ItemEntry
	for _, child := range children {
		switch child.Tag {
		case ItemsDef.ItemID:
			item.ItemID = child.Text()
			break
		case ItemsDef.Key:
			item.Key = child.Text()
			break
		case ItemsDef.Name:
			item.Name = child.Text()
			break
		case ItemsDef.Plural:
			item.Plural = child.Text()
			break
		case ItemsDef.Type:
			item.Type = child.Text()
			break
		case ItemsDef.Object:
			item.Object = child.Text()
			break
		case ItemsDef.Sprite:
			item.Sprite = child.Text()
			break
		case ItemsDef.Category:
			item.Category = child.Text()
			break
		case ItemsDef.ResearchReqKey:
			item.ResearchReqKey = getChildText(child)
			break
		case ItemsDef.ScanReqKey:
			item.ScanReqKey = getChildText(child)
			break
		case ItemsDef.ItemAction:
			item.ItemAction = child.Text()
			break
		case ItemsDef.Hidden:
			item.Hidden = child.Text()
			break
		case ItemsDef.MaxDurability:
			item.MaxDurability = child.Text()
			break
		case ItemsDef.ActionParameter:
			item.ActionParameter = child.Text()
			break
		case ItemsDef.DecomposeValue:
			item.DecomposeValue = child.Text()
			break
		default:
			fmt.Println(child.Tag)
		}
	}
	return item
}

func getChildText(element *etree.Element) []string {
	children := element.ChildElements()
	strings := make([]string, len(children))
	for index, child := range children {
		strings[index] = child.Text()
	}
	return strings
}

func ParseItemsXML(root *etree.Element) ItemEntries {
	children := root.ChildElements()
	items := make([]ItemEntry, len(children))
	for index, child := range children {
		items[index] = parseItem(child)
	}
	return ItemEntries{items}
}

func ParseItemsXMLFile(filename string) (ItemEntries, *etree.Document) {
	doc := etree.NewDocument()
	doc.ReadFromFile(filename)
	root := doc.Root()
	return ParseItemsXML(root), doc
}

func (items *ItemEntries) CreateKeyMap() datatypes.KeyNameList {

	keyNames := make([]datatypes.KeyName, len(items.Entries))
	for index, item := range items.Entries {
		keyName := datatypes.KeyName{Key: item.Key, Name: item.Name}
		keyNames[index] = keyName
	}

	return datatypes.CreateKeyNameList(keyNames)
}
