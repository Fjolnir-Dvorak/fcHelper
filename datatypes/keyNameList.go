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

func createKeyNameList(keys, names []string) KeyNameList {
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

func (list *KeyNameList) ConCat(nameList KeyNameList) {
	list.list.key = append(list.list.key, nameList.list.key...)
	list.list.name = append(list.list.name, nameList.list.name...)
}

func (list *KeyNameList) Mergesort() KeyNameList {
	keys, names := mergesort(list.list.key, list.list.name)
	return createKeyNameList(keys, names)
}

func mergesort(key, name []string) ([]string, []string) {
	if len(key) < 2 {
		return key, name
	}
	mid := len(key) / 2
	key1, name1 := mergesort(key[:mid], name[:mid])
	key2, name2 := mergesort(key[mid:], name[mid:])
	return merge(key1, name1, key2, name2)
}
func merge(key1, name1, key2, name2 []string) ([]string, []string) {
	if len(key1) == 0 {
		return key2, name2
	}
	if len(key2) == 0 {
		return key1, name1
	}
	size := len(key1) + len(key2)
	newKey := make([]string, size, size)
	newName := make([]string, size, size)

	count1, count2 := 0, 0
	for i := 0; i < size; i++ {
		if len(key1) <= count1 {
			newKey[i] = key2[count2]
			newName[i] = name2[count2]
			count2++
		} else if len(key2) <= count2 {
			newKey[i] = key1[count1]
			newName[i] = name1[count1]
			count1++
		} else if key1[count1] <= key2[count2] {
			newKey[i] = key1[count1]
			newName[i] = name1[count1]
			count1++
		} else {
			newKey[i] = key2[count2]
			newName[i] = name2[count2]
			count2++
		}
	}
	return newKey, newName
}
