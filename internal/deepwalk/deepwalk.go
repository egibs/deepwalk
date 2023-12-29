package deepwalk

import (
	"reflect"

	util "github.com/egibs/deepwalk/internal/util"
)

// DeepWalk traverses a map object and returns the value of the specified key
// DeepWalk also supports traversing slices and structs
// If the key is not found, the default value is returned
// If the key is found, all values are returned as a slice by default
// If the returnVal argument is set to "first", the first value is returned
// If the returnVal argument is set to "last", the last value is returned
func DeepWalk(
	obj interface{},
	keys []string,
	defaultVal string,
	returnVal string,
) (interface{}, error) {
	// Return the default value if the object is empty or if the keys or return value are invalid
	if util.IsEmpty(obj) || !util.ValidKeys(keys) || !util.ValidReturnVal(returnVal) {
		return defaultVal, nil
	}

	// Return the object if there are no keys to traverse
	if len(keys) == 0 {
		return obj, nil
	}

	currentKey := keys[0]
	var foundList []interface{}

	r := reflect.ValueOf(obj)
	if r.Kind() == reflect.Struct {
		return deepWalkStruct(r, currentKey, keys, defaultVal, returnVal)
	}

	switch object := obj.(type) {
	case map[string]interface{}:
		return deepWalkMap(object, currentKey, keys, defaultVal, returnVal)
	case []interface{}:
		return deepWalkSlice(object, keys, defaultVal, returnVal, &foundList)
	}

	return util.HandleReturnVal(&foundList, defaultVal, returnVal)
}

// handleStruct handles the case where the object is a struct
func deepWalkStruct(
	r reflect.Value,
	currentKey string,
	keys []string,
	defaultVal string,
	returnVal string,
) (interface{}, error) {
	f := r.FieldByName(currentKey)
	if !f.IsValid() {
		return defaultVal, nil
	}
	return DeepWalk(f.Interface(), keys[1:], defaultVal, returnVal)
}

// handleMap handles the case where the object is a map
func deepWalkMap(
	object map[string]interface{},
	currentKey string,
	keys []string,
	defaultVal string,
	returnVal string,
) (interface{}, error) {
	nextObj, ok := object[currentKey]
	if !ok {
		return defaultVal, nil
	}
	return DeepWalk(nextObj, keys[1:], defaultVal, returnVal)
}

// handleSlice handles the case where the object is a slice
func deepWalkSlice(
	object []interface{},
	keys []string,
	defaultVal string,
	returnVal string,
	foundList *[]interface{},
) (interface{}, error) {
	for _, item := range object {
		val, err := DeepWalk(item, keys, defaultVal, returnVal)
		if err != nil {
			return defaultVal, err
		}
		if val != defaultVal {
			*foundList = append(*foundList, val)
		}
	}

	return util.HandleReturnVal(foundList, defaultVal, returnVal)
}
