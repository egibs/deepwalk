package deepwalk

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// DeepWalk traverses a map object and returns the value of the specified key
// If the key is not found, the default value is returned
// If the key is found, all values are returned as a slice by default
// If the returnVal argument is set to "first", the first value is returned
// If the returnVal argument is set to "last", the last value is returned
func DeepWalk(obj interface{}, keys []string, defaultVal string, returnVal string) (interface{}, error) {
	if isEmpty(obj) {
		return defaultVal, nil
	}

	if len(keys) == 0 {
		return obj, nil
	}

	currentKey := keys[0]
	found := orderedmap.New[string, struct{}]()

	switch object := obj.(type) {
	case map[string]interface{}:
		nextObj, ok := object[currentKey]
		if !ok {
			return defaultVal, nil
		}
		return DeepWalk(nextObj, keys[1:], defaultVal, returnVal)
	case []interface{}:
		for _, item := range object {
			value, err := DeepWalk(item, keys, defaultVal, returnVal)
			if err != nil {
				return defaultVal, err
			}
			if value != nil {
				switch subObject := value.(type) {
				case []string:
					for _, subItem := range subObject {
						found.Set(subItem, struct{}{})
					}
				case string:
					found.Set(value.(string), struct{}{})
				}
			}
		}
	}

	var foundList []string
	for kv := found.Oldest(); kv != nil; kv = kv.Next() {
		foundList = append(foundList, kv.Key)
	}

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
