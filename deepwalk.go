package deepwalk

import (
	"reflect"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
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
	if isEmpty(obj) || !validKeys(keys) || !validReturnVal(returnVal) {
		return defaultVal, nil
	}

	// Return the object if there are no keys to traverse
	if len(keys) == 0 {
		return obj, nil
	}

	currentKey := keys[0]
	found := orderedmap.New[string, struct{}]()
	var foundList []string

	r := reflect.ValueOf(obj)
	if r.Kind() == reflect.Struct {
		return handleStruct(r, currentKey, keys, defaultVal, returnVal)
	}

	switch object := obj.(type) {
	case map[string]interface{}:
		return handleMap(object, currentKey, keys, defaultVal, returnVal)
	case []interface{}:
		return handleSlice(object, keys, defaultVal, returnVal, found, foundList)
	}

	for kv := found.Oldest(); kv != nil; kv = kv.Next() {
		foundList = append(foundList, kv.Key)
	}

	return handleReturnVal(found, foundList, defaultVal, returnVal)
}

// isEmpty checks if the specified object is empty
func isEmpty(subObj interface{}) bool {
	switch subObject := subObj.(type) {
	case string:
		return false
	case []interface{}:
		for _, nextObj := range subObject {
			if !isEmpty(nextObj) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// validKeys checks if the specified keys are valid
func validKeys(keys []string) bool {
	for _, key := range keys {
		if strings.TrimSpace(key) == "" {
			return false
		}
	}
	return true
}

// validReturnVal checks if the specified return value is valid
func validReturnVal(returnVal string) bool {
	switch returnVal {
	case "first", "last", "all":
		return true
	default:
		return false
	}
}

// handleStruct handles the case where the object is a struct
func handleStruct(
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
func handleMap(
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
func handleSlice(
	object []interface{},
	keys []string,
	defaultVal string,
	returnVal string,
	found *orderedmap.OrderedMap[string, struct{}],
	foundList []string,
) (interface{}, error) {
	for _, item := range object {
		val, err := DeepWalk(item, keys, defaultVal, returnVal)
		if err != nil {
			return defaultVal, err
		}
		if val != defaultVal {
			found.Set(val.(string), struct{}{})
			foundList = append(foundList, val.(string))
		}
	}

	return handleReturnVal(found, foundList, defaultVal, returnVal)
}

// handleReturnVal handles the the appropriate value to return based on the returnVal argument
func handleReturnVal(
	found *orderedmap.OrderedMap[string, struct{}],
	foundList []string,
	defaultVal string,
	returnVal string,
) (interface{}, error) {
	if len(foundList) == 0 {
		return defaultVal, nil
	}

	switch returnVal {
	case "first":
		return foundList[0], nil
	case "last":
		return foundList[len(foundList)-1], nil
	case "all":
		if len(foundList) == 1 {
			return foundList[0], nil
		}
		return foundList, nil
	default:
		return defaultVal, nil
	}
}
