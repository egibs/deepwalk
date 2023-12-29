package deepsearch

import (
	"reflect"

	util "github.com/egibs/deepwalk/internal/util"
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
	if util.IsEmpty(obj) || searchKey == "" || !util.ValidReturnVal(returnVal) {
		return defaultVal, nil
	}

	var foundList []interface{}

	search(obj, searchKey, &foundList)

	return util.HandleReturnVal(&foundList, defaultVal, returnVal)
}

// search traverses a data structure and returns the value of the specified key
func search(
	obj interface{},
	searchKey string,
	foundList *[]interface{},
) {
	switch v := obj.(type) {
	case map[string]interface{}:
		deepSearchMap(v, searchKey, foundList)
	case []interface{}:
		for _, item := range v {
			search(item, searchKey, foundList)
		}
	case interface{}:
		deepSearchStruct(v, searchKey, foundList)
	}
}

// deepSearchMap handles the case where the object is a map
func deepSearchMap(
	obj map[string]interface{},
	searchKey string,
	foundList *[]interface{},
) {
	for key, value := range obj {
		if key == searchKey {
			*foundList = append(*foundList, value)
		}
		// Recursively search in nested maps or slices
		switch v := value.(type) {
		case map[string]interface{}:
			search(v, searchKey, foundList)
		case []interface{}:
			for _, item := range v {
				search(item, searchKey, foundList)
			}
		}
	}
}

// deepSearchSlice handles the case where the object is a slice
func deepSearchSlice(
	obj []interface{},
	searchKey string,
	foundList *[]interface{},
) {
	for _, item := range obj {
		search(item, searchKey, foundList)
	}
}

// deepSearchStruct handles the case where the object is a struct
func deepSearchStruct(
	obj interface{},
	searchKey string,
	foundList *[]interface{},
) {
	if r := reflect.ValueOf(obj); r.Kind() == reflect.Struct {
		for i := 0; i < r.NumField(); i++ {
			f := r.Field(i)
			if r.Type().Field(i).Name == searchKey {
				*foundList = append(*foundList, f.Interface())
			}
			// Recursively search in nested maps, slices, or structs
			switch v := f.Interface().(type) {
			case map[string]interface{}, []interface{}:
				search(v, searchKey, foundList)
			case interface{}:
				if reflect.TypeOf(v).Kind() == reflect.Struct {
					search(v, searchKey, foundList)
				}
			}
		}
	}
}
