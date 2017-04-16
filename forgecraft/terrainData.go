package forgecraft

import (
	"github.com/Fjolnir-Dvorak/fcHelper/datatypes"
	"github.com/beevik/etree"
)

type TerrainData struct {
	Doc     *etree.Document
	Entries []TerrainDataEntry
}

type TerrainDataEntry struct {
	ItemID              string
	Key                 string
	Name                string
	Plural              string
	Type                string
	Object              string
	Sprite              string
	Category            string
	ResearchReqKey      []string
	ResearchValue       string
	ScanReqKey          []string
	ItemAction          string
	Hidden              string
	MaxDurability       string
	ActionParameter     string
	DecomposeValue      string
	CubeType            string
	LayerType           string
	TopTexture          string
	SideTexture         string
	BottomTexture       string
	GUITexture          string
	IsSolid             string
	IsTransparent       string
	IsHollow            string
	IsGlass             string
	IsPassable          string
	IsGarbage           string
	IsHidden            string
	IsReinforced        string
	IsPaintable         string
	IsColorised         string
	IsMultiBlockMachine string
	HasObject           string
	HasEntity           string
	AudioWalkType       string
	AudioBuildType      string
	AudioDestroyType    string
	Tags                Tags
	Fuel                []string
	Description         string
	IconName            string
	Hardness            string
	DefaultValue        string
	Value               string
	Values              Values
	ModStartValue       string
	PickReplacement     string
	Stages              Stages
	OreType             string
	MaxStack            string
}
type Tags struct {
	Tags []string
}
type Values struct {
	ValueEntries []ValueEntry
}
type ValueEntry struct {
	Value       string
	Key         string
	Name        string
	IconName    string
	Description string
}
type Stages struct {
	Stages []Stage
}
type Stage struct {
	RangeMinimum  string
	RangeMaximum  string
	TopTexture    string
	SideTexture   string
	BottomTexture string
	GUITexture    string
}

type TerrainDataEntryDef struct {
	ItemID              string
	Key                 string
	Name                string
	Plural              string
	Type                string
	Object              string
	Sprite              string
	Category            string
	ResearchReqKey      string
	ResearchValue       string
	ScanReqKey          string
	ItemAction          string
	Hidden              string
	MaxDurability       string
	ActionParameter     string
	DecomposeValue      string
	CubeType            string
	LayerType           string
	TopTexture          string
	SideTexture         string
	BottomTexture       string
	GUITexture          string
	IsSolid             string
	IsTransparent       string
	IsHollow            string
	IsGlass             string
	IsPassable          string
	IsGarbage           string
	IsHidden            string
	IsReinforced        string
	IsPaintable         string
	IsColorised         string
	IsMultiBlockMachine string
	HasObject           string
	HasEntity           string
	AudioWalkType       string
	AudioBuildType      string
	AudioDestroyType    string
	Tags                string
	Fuel                string
	Description         string
	IconName            string
	Hardness            string
	DefaultValue        string
	Value               string
	Values              string
	ModStartValue       string
	PickReplacement     string
	Stages              string
	OreType             string
	MaxStack            string
}

var (
	TerrainDataDef = TerrainDataEntryDef{
		ItemID:              "ItemID",
		Key:                 "Key",
		Name:                "Name",
		Plural:              "Plural",
		Type:                "Type",
		Object:              "Object",
		Sprite:              "Sprite",
		Category:            "Category",
		ResearchReqKey:      "ResearchRequirements",
		ResearchValue:       "ResearchValue",
		ScanReqKey:          "ScanRequirements",
		ItemAction:          "ItemAction",
		Hidden:              "Hidden",
		MaxDurability:       "MaxDurability",
		ActionParameter:     "ActionParameter",
		DecomposeValue:      "DecomposeValue",
		CubeType:            "CubeType",
		LayerType:           "LayerType",
		TopTexture:          "TopTexture",
		SideTexture:         "SideTexture",
		BottomTexture:       "BottomTexture",
		GUITexture:          "GUITexture",
		IsSolid:             "isSolid",
		IsTransparent:       "isTransparent",
		IsHollow:            "isHollow",
		IsGlass:             "isGlass",
		IsPassable:          "isPassable",
		IsGarbage:           "isGarbage",
		IsHidden:            "isHidden",
		IsReinforced:        "isReinforced",
		IsPaintable:         "isPaintable",
		IsColorised:         "isColorised",
		IsMultiBlockMachine: "isMultiBlockMachine",
		HasObject:           "hasObject",
		HasEntity:           "hasEntity",
		AudioWalkType:       "AudioWalkType",
		AudioBuildType:      "AudioBuildType",
		AudioDestroyType:    "AudioDestroyType",
		Tags:                "tags",
		Fuel:                "Fuel",
		Description:         "Description",
		IconName:            "IconName",
		Hardness:            "Hardness",
		DefaultValue:        "DefaultValue",
		Value:               "Value",
		Values:              "Values",
		ModStartValue:       "ModStartValue",
		PickReplacement:     "PickReplacement",
		Stages:              "Stages",
		OreType:             "OreType",
		MaxStack:            "MaxStack",
	}
	TerrainDataFilename = "TerrainData.xml"
)

func parseTerrainData(element *etree.Element) TerrainDataEntry {
	children := element.ChildElements()
	var item TerrainDataEntry
	for _, child := range children {
		switch child.Tag {
		case TerrainDataDef.ItemID:
			item.ItemID = child.Text()
			break
		case TerrainDataDef.Key:
			item.Key = child.Text()
			break
		case TerrainDataDef.Name:
			item.Name = child.Text()
			break
		case TerrainDataDef.Plural:
			item.Plural = child.Text()
			break
		case TerrainDataDef.Type:
			item.Type = child.Text()
			break
		case TerrainDataDef.Object:
			item.Object = child.Text()
			break
		case TerrainDataDef.Sprite:
			item.Sprite = child.Text()
			break
		case TerrainDataDef.Category:
			item.Category = child.Text()
			break
		case TerrainDataDef.ResearchReqKey:
			item.ResearchReqKey = getChildText(child)
			break
		case TerrainDataDef.ResearchValue:
			item.ResearchValue = child.Text()
			break
		case TerrainDataDef.ScanReqKey:
			item.ScanReqKey = getChildText(child)
			break
		case TerrainDataDef.ItemAction:
			item.ItemAction = child.Text()
			break
		case TerrainDataDef.Hidden:
			item.Hidden = child.Text()
			break
		case TerrainDataDef.MaxDurability:
			item.MaxDurability = child.Text()
			break
		case TerrainDataDef.ActionParameter:
			item.ActionParameter = child.Text()
			break
		case TerrainDataDef.DecomposeValue:
			item.DecomposeValue = child.Text()
			break
		case TerrainDataDef.CubeType:
			item.CubeType = child.Text()
			break
		case TerrainDataDef.LayerType:
			item.LayerType = child.Text()
			break
		case TerrainDataDef.TopTexture:
			item.TopTexture = child.Text()
			break
		case TerrainDataDef.SideTexture:
			item.SideTexture = child.Text()
			break
		case TerrainDataDef.BottomTexture:
			item.BottomTexture = child.Text()
			break
		case TerrainDataDef.GUITexture:
			item.GUITexture = child.Text()
			break
		case TerrainDataDef.IsSolid:
			item.IsSolid = child.Text()
			break
		case TerrainDataDef.IsTransparent:
			item.IsTransparent = child.Text()
			break
		case TerrainDataDef.IsHollow:
			item.IsHollow = child.Text()
			break
		case TerrainDataDef.IsGlass:
			item.IsGlass = child.Text()
			break
		case TerrainDataDef.IsPassable:
			item.IsPassable = child.Text()
			break
		case TerrainDataDef.IsGarbage:
			item.IsGarbage = child.Text()
			break
		case TerrainDataDef.IsHidden:
			item.IsHidden = child.Text()
			break
		case TerrainDataDef.IsReinforced:
			item.IsReinforced = child.Text()
			break
		case TerrainDataDef.IsPaintable:
			item.IsPaintable = child.Text()
			break
		case TerrainDataDef.IsColorised:
			item.IsColorised = child.Text()
			break
		case TerrainDataDef.IsMultiBlockMachine:
			item.IsMultiBlockMachine = child.Text()
			break
		case TerrainDataDef.HasObject:
			item.HasObject = child.Text()
			break
		case TerrainDataDef.HasEntity:
			item.HasEntity = child.Text()
			break
		case TerrainDataDef.AudioWalkType:
			item.AudioWalkType = child.Text()
			break
		case TerrainDataDef.AudioBuildType:
			item.AudioBuildType = child.Text()
			break
		case TerrainDataDef.AudioDestroyType:
			item.AudioDestroyType = child.Text()
			break
		case TerrainDataDef.Tags:
			//item.Tags = getChildText(child)  // TODO
			break
		case TerrainDataDef.Fuel:
			item.Fuel = getChildText(child)
			break
		case TerrainDataDef.Description:
			item.Description = child.Text()
			break
		case TerrainDataDef.IconName:
			item.IconName = child.Text()
			break
		case TerrainDataDef.Hardness:
			item.Hardness = child.Text()
			break
		case TerrainDataDef.DefaultValue:
			item.DefaultValue = child.Text()
			break
		case TerrainDataDef.Value:
			item.Value = child.Text()
			break
		case TerrainDataDef.Values:
			//item.Values = getChildText(child)  // TODO
			break
		case TerrainDataDef.ModStartValue:
			item.ModStartValue = child.Text()
			break
		case TerrainDataDef.PickReplacement:
			item.PickReplacement = child.Text()
			break
		case TerrainDataDef.Stages:
			//item.Stages = getChildText(child)  // TODO
			break
		case TerrainDataDef.OreType:
			item.OreType = child.Text()
			break
		case TerrainDataDef.MaxStack:
			item.MaxStack = child.Text()
			break
		default:
			utilHandleUnknownElement(child.Tag)
		}
	}
	return item
}

func parseTerrainDataXML(root *etree.Element) []TerrainDataEntry {
	children := root.ChildElements()
	items := make([]TerrainDataEntry, len(children))
	for index, child := range children {
		items[index] = parseTerrainData(child)
	}
	return items
}

func ParseTerrainDataXMLFile(filename string) TerrainData {
	doc := etree.NewDocument()
	doc.ReadFromFile(filename)
	items := TerrainData{Entries: parseTerrainDataXML(doc.Root()), Doc: doc}
	return items
}

func (items *TerrainData) CreateKeyMap() datatypes.KeyNameList {

	keyNames := make([]datatypes.KeyName, len(items.Entries))
	for index, item := range items.Entries {
		keyName := datatypes.KeyName{Key: item.Key, Name: item.Name}
		keyNames[index] = keyName
	}

	return datatypes.CreateKeyNameList(keyNames)
}
