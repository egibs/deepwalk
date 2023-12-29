package deepwalk

import (
	"reflect"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// DeepSearch traverses a data structure and returns the value of the specified key
// DeepSearch does not need to know the path to the specified key
// Like DeepWalk, DeepSearch will return either the first, last, or all encountered values
// associated with the provided `searchKey`
func DeepSearch(
	obj interface{},
	searchKey string,
	defaultVal string,
	returnVal string,
) (interface{}, error) {
	// Return the default value if the object is empty or if the keys or return value are invalid
	if IsEmpty(obj) || !ValidKeys([]string{searchKey}) || !ValidReturnVal(returnVal) {
		return defaultVal, nil
	}

	// Return the object if there are no keys to traverse
	if len(searchKey) == 0 {
		return obj, nil
	}

	found := orderedmap.New[string, struct{}]()
	var foundList []interface{}
	for kv := found.Oldest(); kv != nil; kv = kv.Next() {
		foundList = append(foundList, kv.Key)
	}

	search(obj, searchKey, found, &foundList)

	return HandleReturnVal(&foundList, defaultVal, returnVal)
}

// search traverses a data structure and returns the value of the specified key
func search(obj interface{}, searchKey string, found *orderedmap.OrderedMap[string, struct{}], foundList *[]interface{}) {
	switch object := obj.(type) {
	case map[string]interface{}:
		deepSearchMap(object, searchKey, found, foundList)
	case []interface{}:
		deepSearchSlice(object, searchKey, found, foundList)
	default:
		deepSearchStruct(object, searchKey, found, foundList)
	}
}

// deepSearchMap handles the case where the object is a map
func deepSearchMap(
	obj map[string]interface{},
	searchKey string,
	found *orderedmap.OrderedMap[string, struct{}],
	foundList *[]interface{},
) {
	for key, value := range obj {
		if key == searchKey {
			found.Set(value.(string), struct{}{})
			*foundList = append(*foundList, value)
		}
		search(value, searchKey, found, foundList)
	}
}

// deepSearchSlice handles the case where the object is a slice
func deepSearchSlice(
	obj []interface{},
	searchKey string,
	found *orderedmap.OrderedMap[string, struct{}],
	foundList *[]interface{},
) {
	for _, item := range obj {
		search(item, searchKey, found, foundList)
	}
}

// deepSearchStruct handles the case where the object is a struct
func deepSearchStruct(
	obj interface{},
	searchKey string,
	found *orderedmap.OrderedMap[string, struct{}],
	foundList *[]interface{},
) {
	if r := reflect.ValueOf(obj); r.Kind() == reflect.Struct {
		for i := 0; i < r.NumField(); i++ {
			f := r.Field(i)
			if r.Type().Field(i).Name == searchKey {
				found.Set(f.Interface().(string), struct{}{})
				*foundList = append(*foundList, f.Interface())
			} else {
				search(f.Interface(), searchKey, found, foundList)
			}
		}
	}
}
