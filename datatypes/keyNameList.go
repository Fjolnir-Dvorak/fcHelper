package datatypes

type KeyNameList struct {
	list list
}

type list struct {
	key  []string
	name []string
}

type KeyName struct {
	Key  string
	Name string
}

func findFirst(slice []string, item string) int {
	for index, s := range slice {
		if s == item {
			return index
		}
	}
	return -1
}

func InitEmptyKeyNameList() KeyNameList {
	return KeyNameList{list{make([]string, 0), make([]string, 0)}}
}

func CreateKeyNameList(keyNames []KeyName) KeyNameList {
	keys, names := generateLists(keyNames)
	return KeyNameList{list{keys, names}}
}

func generateLists(keyNames []KeyName) ([]string, []string) {
	keys := make([]string, len(keyNames))
	names := make([]string, len(keyNames))
	for index, keyName := range keyNames {
		keys[index] = keyName.Key
		names[index] = keyName.Name
	}
	return keys, names
}

func (list *KeyNameList) GetKey(name string) string {
	index := findFirst(list.list.name, name)
	if index < 0 {
		return ""
	}
	return list.list.key[index]
}
func (list *KeyNameList) GetName(key string) string {
	index := findFirst(list.list.key, key)
	if index < 0 {
		return ""
	}
	return list.list.name[index]

}

func (list *KeyNameList) Iterate() []KeyName {
	out := make([]KeyName, len(list.list.key))
	for index, key := range list.list.key {
		name := list.list.name[index]
		out[index] = KeyName{Key: key, Name: name}
	}
	return out
}

func (list *KeyNameList) Append(keyMap []KeyName) {
	newKeys, newNames := generateLists(keyMap)

	list.list.key = append(list.list.key, newKeys...)
	list.list.name = append(list.list.name, newNames...)
}

func (list *KeyNameList) Merge(nameList KeyNameList) {
	list.list.key = append(list.list.key, nameList.list.key...)
	list.list.name = append(list.list.name, nameList.list.name...)
}
