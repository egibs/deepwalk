package traverse

import (
	"fmt"
	"reflect"
	"sort"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type ReturnControl int

const (
	First ReturnControl = iota
	Last
	All
)

func Traverse(v interface{}, targetKey string, defaultValue string, control ReturnControl) (interface{}, error) {
	found := orderedmap.New[interface{}, struct{}]()
	var processQueue func(reflect.Value)

	processQueue = func(rv reflect.Value) {
		// Handle pointers and interface{} by obtaining the value they point to
		for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				return
			}
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Map:
			for _, key := range rv.MapKeys() {
				val := rv.MapIndex(key)
				if fmt.Sprintf("%v", key.Interface()) == targetKey {
					found.Set(val.Interface(), struct{}{})
				}
				processQueue(val)
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				processQueue(rv.Index(i))
			}
		case reflect.Struct:
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				if rv.Type().Field(i).Name == targetKey {
					found.Set(field.Interface(), struct{}{})
				}
				processQueue(field)
			}
		}
	}

	// Convert the interface{} to reflect.Value before processing
	initialValue := reflect.ValueOf(v)
	processQueue(initialValue)

	var results []interface{}
	for kv := found.Oldest(); kv != nil; kv = kv.Next() {
		results = append(results, kv.Key)
	}
	// sort results to ensure deterministic output
	sort.Slice(results, func(i, j int) bool {
		return fmt.Sprintf("%v", results[i]) < fmt.Sprintf("%v", results[j])
	})

	switch control {
	case First:
		if len(results) > 0 {
			return results[0], nil
		}
	case Last:
		if len(results) > 0 {
			return results[len(results)-1], nil
		}
	case All:
		if len(results) > 0 {
			if len(results) == 1 {
				return results[0], nil
			}
			return results, nil
		}
	}
	return defaultValue, nil
}
