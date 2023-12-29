package main

import (
	"encoding/json"
	"fmt"

	"github.com/egibs/deepwalk/pkg/deepsearch"
	"github.com/egibs/deepwalk/pkg/deepwalk"
)

func main() {
	exampleJson := []byte(`{
		"a": {
			"b": {
				"c": [
					{
						"d": "foo"
					},
					{
						"d": "bar"
					},
					{
						"d": "baz"
					}
				],
				"e": [
					[
						[
							{
								"f": "foo"
							}
						]
					]
				],
				"g": [[[[[[[[]]]]]]]]
			}
		}
	}`)

	var object map[string]interface{}
	err := json.Unmarshal(exampleJson, &object)
	if err != nil {
		fmt.Println(err)
	}
	value, err := deepwalk.DeepWalk(object, []string{"a", "b", "e", "f"}, "default", "all")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	value2, err := deepsearch.DeepSearch(object, "d", "default", "all")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Value should be 'foo': %v\n", value)
	fmt.Printf("Value2 should be '[foo bar baz]': %v\n", value2)
}
