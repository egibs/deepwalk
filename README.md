# maptraverse

## Overview

`deepwalk` is a Golang implementation of the code documented [here](https://egibs.xyz/posts/technical/deep_walk/) which was originally written in Python.

The goal of `deepwalk` is to traverse an arbitrarily-nested map and retrieve the value associated with a given key. This key can be a singular key or a slice of keys with each key representing a deeper level of the map to traverse.

This project was mostly done to hack on some Go after spending awhile way from it; however, I also wanted to see how much more performant Go would be with this application. 

That said, this package is still useful aside from the "see if it works" perspective. It would be easy to specify a desired value and search the entire map, but I wanted to remain faithful to the original implementation.

Usage examples can be found in `deepwalk_test.go`, but because code is only as good as its documentation, I'll outline some examples below.

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

Given a simple JSON object like this:
```go
exampleJson = []byte(`{
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
--- PASS: TestDeepwalkRandomSuccess (0.26s)
=== RUN   TestDeepwalkRandomDefault
--- PASS: TestDeepwalkRandomDefault (3.20s)
PASS
ok  	egibs/deepwalk	3.554s
```

## Benchmarks
Run the three included benchmarks by running `make bench`:
```sh
make bench                                                                                                (deepwalk) 19:36:58
go test -bench=.
goos: darwin
goarch: arm64
pkg: egibs/deepwalk
BenchmarkDeepwalkMinimalJSON-10    	  186184	      6425 ns/op
BenchmarkDeepwalkSuccess-10        	    4609	    245696 ns/op
BenchmarkDeepwalkDefault-10        	     354	   3326961 ns/op
PASS
ok  	egibs/deepwalk	7.637s
```

## Acknowledgements

[wk8](https://github.com/wk8) for their extremely handy `go-ordered-map` [package](https://pkg.go.dev/github.com/wk8/go-ordered-map/v2@v2.1.8)
  - I take `collections.OrderedDict` for granted in Python and this was extremely easy to implement

## TODO

For now, I really just wanted to publish my first Go package since I'd never done that. Eventually I may add a `GoReleaser` Workflow to this repository
I do want to make a version that locates a desired value, if present, in a map.