# deepwalk

## Overview

`deepwalk` is a Golang implementation of the code documented [here](https://egibs.xyz/posts/technical/deep_walk/) which was originally written in Python.

The goal of `deepwalk` is to traverse an arbitrarily-nested map and retrieve the value associated with a given key. This key can be a singular key or a slice of keys with each key representing a deeper level of the map to traverse.

This project was mostly done to hack on some Go after spending awhile way from it. That said, this package is still useful aside from the "see if it works" perspective. It would be easy to specify a desired value and search the entire map, but I wanted to remain faithful to the original Python implementation.

Usage examples can be found in `deepwalk_test.go`, but because code is only as good as its documentation, examples will be added (and added to) below.

## Installation

To install `deepwalk`, run the following:
```sh
go get -u github.com/egibs/deepwalk
```

Import `deepwalk` like so:
```go
import (
    "github.com/egibs/deepwalk"
)
```

## Examples

The original use-case for this in Python was to traverse Python dictionaries. While Golang has different names and conventions for this, the approach remains the same.

### JSON
Given a simple JSON object like this:
```go
var exampleJSON = []byte(`{
    "key": [
        {
            "inner_key": "value"
        }
    ]
}`)
```

The value of `inner_key` can be retrieved like so:
```go
var object map[string]interface{}
err := json.Unmarshal(exampleJSON, &object)
if err != nil {
    fmt.Println(err)
}
value, err := DeepWalk(object, []string{"key", "inner_key"}, "<NO_VALUE>", "all")
if err != nil {
    fmt.Println(err)
}
fmt.Println(value)
```
which will result in `value` being printed.

To add some color to the aruguments used when calling `DeepWalk` --
* `object` is the starting map to traverse
* `[]string{...}` is the slice of keys, the last of which is the key with the desired value
* `<NO_VALUE>` is the default value to return if the desired value is not found
* `all` specifies the values to return, if found
  * `first` and `last` are also available which will return the first and last found values, respectively

### Maps

Maps can be used directly as well (this example is taken from the tests) --
```go
var exampleMap = map[string]interface{}{
	"key": map[string]interface{}{
		"inner_key": []interface{}{
			[]interface{}{
				[]interface{}{
					map[string]interface{}{
						"very_nested_key": "very_nested_value",
					},
				},
			},
		},
	},
}
values, err := DeepWalk(exampleMap, []string{"key", "inner_key", "very_nested_key"}, "<NO_VALUE>", "all")
if err != nil {
    fmt.Println(err)
}
fmt.Println(value)
```

## Testing
Several categories of tests are included:
1. A manual JSON object
2. 14 map-based tests
3. Randomly-generated test cases for successful path-based traversal
4. Randomly-generated test cases for traversals that return a default value

Categories three and four run one-thousand iterations of each variant and each distinct case uses a randomly-generated map structure, list of keys, and desired value (in the case of the path-based traversal tests)

Run all of the included tests by running `make test`:
```sh
go test ./... -v
=== RUN   TestDeepwalkMinimalJSON
--- PASS: TestDeepwalkMinimalJSON (0.00s)
=== RUN   TestDeepwalkMinimalMap
--- PASS: TestDeepwalkMinimalMap (0.00s)
=== RUN   TestDeepWalk
=== RUN   TestDeepWalk/Test_case_1_-_empty_object
=== RUN   TestDeepWalk/Test_case_2_-_empty_keys
=== RUN   TestDeepWalk/Test_case_3_-_key_not_found
=== RUN   TestDeepWalk/Test_case_4_-_nested_key_not_found
=== RUN   TestDeepWalk/Test_case_5_-_nested_key_found
=== RUN   TestDeepWalk/Test_case_6_-_array_of_strings
=== RUN   TestDeepWalk/Test_case_7_-_array_of_maps
=== RUN   TestDeepWalk/Test_case_8_-_array_of_maps_with_no_matching_key
=== RUN   TestDeepWalk/Test_case_9_-_array_of_maps_with_multiple_matching_keys
=== RUN   TestDeepWalk/Test_case_10_-_array_of_maps_with_multiple_matching_keys,_return_last
=== RUN   TestDeepWalk/Test_case_11_-_array_of_maps_with_multiple_matching_keys,_return_first
=== RUN   TestDeepWalk/Test_case_12_-_array_of_maps_with_multiple_matching_keys,_return_default
=== RUN   TestDeepWalk/Test_case_13_-_array_of_maps_with_multiple_matching_keys,_return_default
=== RUN   TestDeepWalk/Test_case_14_-_array_of_maps_with_multiple_matching_keys,_return_default
--- PASS: TestDeepWalk (0.00s)
    --- PASS: TestDeepWalk/Test_case_1_-_empty_object (0.00s)
    --- PASS: TestDeepWalk/Test_case_2_-_empty_keys (0.00s)
    --- PASS: TestDeepWalk/Test_case_3_-_key_not_found (0.00s)
    --- PASS: TestDeepWalk/Test_case_4_-_nested_key_not_found (0.00s)
    --- PASS: TestDeepWalk/Test_case_5_-_nested_key_found (0.00s)
    --- PASS: TestDeepWalk/Test_case_6_-_array_of_strings (0.00s)
    --- PASS: TestDeepWalk/Test_case_7_-_array_of_maps (0.00s)
    --- PASS: TestDeepWalk/Test_case_8_-_array_of_maps_with_no_matching_key (0.00s)
    --- PASS: TestDeepWalk/Test_case_9_-_array_of_maps_with_multiple_matching_keys (0.00s)
    --- PASS: TestDeepWalk/Test_case_10_-_array_of_maps_with_multiple_matching_keys,_return_last (0.00s)
    --- PASS: TestDeepWalk/Test_case_11_-_array_of_maps_with_multiple_matching_keys,_return_first (0.00s)
    --- PASS: TestDeepWalk/Test_case_12_-_array_of_maps_with_multiple_matching_keys,_return_default (0.00s)
    --- PASS: TestDeepWalk/Test_case_13_-_array_of_maps_with_multiple_matching_keys,_return_default (0.00s)
    --- PASS: TestDeepWalk/Test_case_14_-_array_of_maps_with_multiple_matching_keys,_return_default (0.00s)
=== RUN   TestDeepwalkRandomSuccess
--- PASS: TestDeepwalkRandomSuccess (0.25s)
=== RUN   TestDeepwalkRandomDefault
--- PASS: TestDeepwalkRandomDefault (3.22s)
PASS
ok  	github.com/egibs/deepwalk	3.600s
```

## Benchmarks
Run the included benchmarks by running `make bench`:
```sh
make bench                                                                                                (deepwalk) 20:35:43
go test -bench=.
goos: darwin
goarch: arm64
pkg: github.com/egibs/deepwalk
BenchmarkDeepwalkMinimalJSON-10    	  188026	      6118 ns/op
BenchmarkDeepwalkSuccess-10        	    5114	    240216 ns/op
BenchmarkDeepwalkDefault-10        	     378	   3192538 ns/op
PASS
ok  	github.com/egibs/deepwalk	7.537s
```

## Acknowledgements

[wk8](https://github.com/wk8) for their extremely handy [go-ordered-map](https://pkg.go.dev/github.com/wk8/go-ordered-map/v2@v2.1.8) package
  - It is easy to take `collections.OrderedDict` for granted in Python and this was extremely easy to implement

## TODO

- For now, the functionality of the package is the primary goal, but eventually more robust examples, functionality, and documentation are on the table
- Adding a `GoReleaser` Workflow and tweaking the package to account for that is potentially worth doing as well
