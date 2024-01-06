# deepwalk

[![Go Reference](https://pkg.go.dev/badge/github.com/egibs/deepwalk.svg)](https://pkg.go.dev/github.com/egibs/deepwalk)

## Overview

`deepwalk` is a Golang implementation of the code documented [here](https://egibs.xyz/posts/technical/deep_walk/) which was originally written in Python.

The goal of `deepwalk` is to traverse an arbitrarily-nested map and retrieve the value associated with a given key.

This project was mostly done to hack on some Go after spending awhile way from it. That said, this package is still useful aside from the "see if it works" perspective.

Usage examples can be found in `pkg/traverse/traverse_test.go`, but because code is only as good as its documentation, examples will be added (and added to) below.

## Installation

To install `deepwalk`, run the following:
```sh
go install github.com/egibs/deepwalk/v2@latest
```

Import `traverse` like so:
```go
import (
    "github.com/egibs/deepwalk/pkg/v2/traverse"
)
```

## Examples

The original use-case for this in Python was to traverse Python dictionaries. While Golang has different names and conventions for this, the approach remains the same.

### CLI

[Cobra](https://github.com/spf13/cobra) is used to provide CLI support for `deepwalk`.

To get started, build the `deepwalk` binary with `make build`.

Then, run `deepwalk` to see the available commands:

```sh
❯ ./deepwalk -h
Usage:
  deepwalk [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  traverse    Traverse an object for the specified key and return the value associated with that key

Flags:
  -h, --help   help for deepwalk

Use "deepwalk [command] --help" for more information about a command.
```

To use `Traverse`, use the `traverse` command:
```sh
./deepwalk traverse -h
traverse utilizes the Traverse function which requires a ReturnControl value to determine what to return.
		First will return the first value found, Last will return the last value found, and All will return all values found.

Usage:
  deepwalk traverse --object <filename or JSON string> --search-key <search key> --default-value <default value> --return-value <return value> [flags]

Flags:
      --default-value string   Default value to return if search fails (default "NO_VALUE")
  -h, --help                   help for traverse
      --object string          Object to search
      --return-value string    Value to return if search succeeds (default "all")
      --search-key string      Key to search for
```

Examples of using `traverse`:
`search`:
```sh
❯ ./deepwalk traverse  --object complex.json --search-key name
[
  "Tamra Bennett",
  "Alana Hoover",
  "Ewing Williamson",
  "Webster Serrano",
  "Lea Bryant",
  "Sylvia Parks",
  "Hubbard Delgado",
  "Townsend Calderon",
  "Knapp Patton",
  "Barr Floyd",
  "Haynes Osborn",
  "Rebecca Walters",
  "Muriel Lindsay",
  "Osborne Reid",
  "Lois Chaney",
  "Contreras Wolfe",
  "Goodwin Christensen",
  "Rosa Luna",
  "Tabitha Moreno",
  "Oneil Carlson"
]

❯ ./deepwalk traverse --object '{"a": {"b": "c"}}' --search-key b
"c"
```

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
value, err := deepwalk.Traverse(object, "key", "<NO_VALUE>", "all")
if err != nil {
    fmt.Println(err)
}
fmt.Println(value)
```
which will result in `value` being printed.

To add some color to the aruguments used when calling `Traverse` --
* `object` is the starting map to traverse
* `key` is the target key to search for in the structure
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
var obj interface{} = exampleMap
values, err := deepwalk.DeepWalk(obj, "very_nested_key", "<NO_VALUE>", "all")
if err != nil {
    fmt.Println(err)
}
fmt.Println(value)
```

### Structs

Structs can also be traversed with `Traverse`:
```go
type TestStruct struct {
    Field1       string
    Field2       int
    NestedStruct struct {
        NestedStruct2 struct {
            NestedField1 string
            NestedField2 int
        }
        NestedField1 string
        NestedField2 int
    }
}

testStruct := TestStruct{
    Field1: "test",
    Field2: 123,
    NestedStruct: struct {
        NestedStruct2 struct {
            NestedField1 string
            NestedField2 int
        }
        NestedField1 string
        NestedField2 int
    }{
        NestedStruct2: struct {
            NestedField1 string
            NestedField2 int
        }{
            NestedField1: "nested",
            NestedField2: 456,
        },
        NestedField1: "nested",
        NestedField2: 456,
    },
}

values, err := deepwalk.DeepWalk(testStruct, "NestedField1", "<NO_VALUE>", "all")
if err != nil {
    fmt.Println(err)
}
fmt.Println(value)
```

## Testing
Several categories of tests are included:
1. 7 test cases covering a variety of uses
2. Randomly-generated test cases for successful path-based traversal
3. Randomly-generated test cases for traversals that return a default value

Categories three and four run one-thousand iterations of each variant and each distinct case uses a randomly-generated map structure, list of keys, and desired value (in the case of the path-based traversal tests)

Run all of the included tests by running `make test`:
```sh
❯ make test
go test ./... -v
go: downloading github.com/wk8/go-ordered-map/v2 v2.1.8
go: downloading github.com/spf13/cobra v1.8.0
go: downloading github.com/spf13/pflag v1.0.5
go: downloading gopkg.in/yaml.v3 v3.0.1
go: downloading github.com/bahlo/generic-list-go v0.2.0
go: downloading github.com/buger/jsonparser v1.1.1
go: downloading github.com/mailru/easyjson v0.7.7
?       github.com/egibs/deepwalk/v2    [no test files]
?       github.com/egibs/deepwalk/v2/cmd        [no test files]
?       github.com/egibs/deepwalk/v2/internal/util      [no test files]
=== RUN   TestTraverse
=== RUN   TestTraverse/Test_case_1_-_key_found_in_map
=== RUN   TestTraverse/Test_case_2_-_key_not_found_in_map
=== RUN   TestTraverse/Test_case_3_-_key_found_in_nested_map
=== RUN   TestTraverse/Test_case_4_-_key_not_found_in_nested_map
=== RUN   TestTraverse/Test_case_5_-_key_found_in_struct
=== RUN   TestTraverse/Test_case_6_-_key_not_found_in_struct
=== RUN   TestTraverse/Test_case_7_-_duplicate_key_found_in_map
--- PASS: TestTraverse (0.00s)
    --- PASS: TestTraverse/Test_case_1_-_key_found_in_map (0.00s)
    --- PASS: TestTraverse/Test_case_2_-_key_not_found_in_map (0.00s)
    --- PASS: TestTraverse/Test_case_3_-_key_found_in_nested_map (0.00s)
    --- PASS: TestTraverse/Test_case_4_-_key_not_found_in_nested_map (0.00s)
    --- PASS: TestTraverse/Test_case_5_-_key_found_in_struct (0.00s)
    --- PASS: TestTraverse/Test_case_6_-_key_not_found_in_struct (0.00s)
    --- PASS: TestTraverse/Test_case_7_-_duplicate_key_found_in_map (0.00s)
PASS
ok      github.com/egibs/deepwalk/v2/pkg/traverse       0.002s
```

## Benchmarks
Run the included benchmarks by running `make bench`:
```sh
❯ make bench
go test ./... -bench=. -benchmem
?       github.com/egibs/deepwalk/v2    [no test files]
?       github.com/egibs/deepwalk/v2/cmd        [no test files]
?       github.com/egibs/deepwalk/v2/internal/util      [no test files]
goos: linux
goarch: amd64
pkg: github.com/egibs/deepwalk/v2/pkg/traverse
cpu: 12th Gen Intel(R) Core(TM) i9-12900K
BenchmarkTraverse-24             1000000              1020 ns/op             764 B/op         28 allocs/op
BenchmarkTraverseSuccess-24         7668            156473 ns/op           27399 B/op       1525 allocs/op
BenchmarkTraverseDefault-24          706           1739258 ns/op          297695 B/op      16732 allocs/op
PASS
ok      github.com/egibs/deepwalk/v2/pkg/traverse       3.652s
```

## Miscellaneous

### Manually cutting a release

To cut a release for this package, do the following:
- Update the `VERSION` file
  - For now this is aesthetic but offers an easy way to reference the version outside of GitHub
- Update `CHANGELOG.md` with a new entry matching the new version
- Commit and push any changes to `main` (either directly or via  PR)
- Run `git tag vX.Y.Z`
- Run `git push origin --tags`
- Create a new release using the new tag

## Acknowledgements

[wk8](https://github.com/wk8) for their extremely handy [go-ordered-map](https://pkg.go.dev/github.com/wk8/go-ordered-map/v2@v2.1.8) package
  - It is easy to take `collections.OrderedDict` for granted in Python and this was extremely easy to implement
