package forgecraft

import (
	"github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/beevik/etree"
)

type Items struct {
	Doc     *etree.Document
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
	ItemsParent   = "ItemEntry"
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
			utilHandleUnknownElement(child.Tag)
		}
	}
	return item
}

func parseItemsXML(root *etree.Element) []ItemEntry {
	children := root.ChildElements()
	items := make([]ItemEntry, len(children))
	for index, child := range children {
		items[index] = parseItem(child)
	}
	return items
}

func ParseItemsXMLFile(filename string) Items {
	doc := etree.NewDocument()
	doc.ReadFromFile(filename)
	items := Items{Entries: parseItemsXML(doc.Root()), Doc: doc}
	return items
}

func (items *Items) CreateKeyMap() datatypes.KeyNameList {

	keyNames := make([]datatypes.KeyName, len(items.Entries))
	for index, item := range items.Entries {
		keyName := datatypes.KeyName{Key: item.Key, Name: item.Name}
		keyNames[index] = keyName
	}

	return datatypes.CreateKeyNameList(keyNames)
}
