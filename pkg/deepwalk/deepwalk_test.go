package deepwalk

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	util "github.com/egibs/deepwalk/internal/util"
)

const (
	maxRuns = 1000
)

var exampleJSON = []byte(`{
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

func TestDeepwalkMinimalJSON(t *testing.T) {
	var object map[string]interface{}
	err := json.Unmarshal(exampleJSON, &object)
	if err != nil {
		fmt.Println(err)
	}
	want := "foo"
	values, err := DeepWalk(object, []string{"a", "b", "e", "f"}, "<NO_VALUE>", "all")
	if err != nil {
		fmt.Println(err)
	}
	if values != want {
		t.Errorf("DeepWalk() = %v, want %v", values, want)
	}
}

func TestDeepwalkMinimalMap(t *testing.T) {
	want := "very_nested_value"
	values, err := DeepWalk(exampleMap, []string{"key", "inner_key", "very_nested_key"}, "<NO_VALUE>", "all")
	if err != nil {
		fmt.Println(err)
	}
	if values != want {
		t.Errorf("DeepWalk() = %v, want %v", values, want)
	}
}

func BenchmarkDeepwalkMinimalJSON(b *testing.B) {
	keys := map[int][]string{
		1: {"a", "b", "e", "f"},
		2: {"a", "b", "c", "d"},
		3: {"a", "b", "g"},
	}
	want := map[int]interface{}{
		1: "foo",
		2: []interface{}{"foo", "bar", "baz"},
		3: "<NO_VALUE>",
	}
	for i := 0; i < b.N; i++ {
		var object map[string]interface{}
		err := json.Unmarshal(exampleJSON, &object)
		if err != nil {
			b.Errorf("DeepWalk() = %v, want %v", nil, want)
		}
		for i, key := range keys {
			want, exists := want[i]
			if !exists {
				b.Errorf("DeepWalk() = %v, want %v", nil, want)
			}
			values, err := DeepWalk(object, key, "<NO_VALUE>", "all")
			if err != nil {
				b.Errorf("DeepWalk() = %v, want %v", values, want)
			}
			if !reflect.DeepEqual(values, want) {
				b.Errorf("DeepWalk() = %v, want %v", values, want)
			}
		}
	}
}

func TestDeepWalk(t *testing.T) {
	type args struct {
		obj        interface{}
		keys       []string
		defaultVal string
		returnVal  string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test case 1 - empty object",
			args: args{
				obj:        nil,
				keys:       []string{"key1", "key2"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "Test case 2 - empty keys",
			args: args{
				obj:        map[string]interface{}{"key1": "value1", "key2": "value2"},
				keys:       []string{},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    map[string]interface{}{"key1": "value1", "key2": "value2"},
			wantErr: false,
		},
		{
			name: "Test case 3 - key not found",
			args: args{
				obj:        map[string]interface{}{"key1": "value1", "key2": "value2"},
				keys:       []string{"key3"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "Test case 4 - nested key not found",
			args: args{
				obj:        map[string]interface{}{"key1": map[string]interface{}{"key2": "value2"}},
				keys:       []string{"key1", "key3"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "Test case 5 - nested key found",
			args: args{
				obj:        map[string]interface{}{"key1": map[string]interface{}{"key2": "value2"}},
				keys:       []string{"key1", "key2"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "value2",
			wantErr: false,
		},
		{
			name: "Test case 6 - array of strings",
			args: args{
				obj:        []interface{}{"value1", "value2", "value3"},
				keys:       []string{},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    []interface{}{"value1", "value2", "value3"},
			wantErr: false,
		},
		{
			name: "Test case 7 - array of maps",
			args: args{
				obj: []interface{}{
					map[string]interface{}{"key1": "value1", "key2": "value2"},
					map[string]interface{}{"key1": "value3", "key2": "value4"},
				},
				keys:       []string{"key1"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "Test case 8 - array of maps with no matching key",
			args: args{
				obj: []interface{}{
					map[string]interface{}{"key1": "value1", "key2": "value2"},
					map[string]interface{}{"key3": "value3", "key4": "value4"},
				},
				keys:       []string{"key5"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "Test case 9 - array of maps with multiple matching keys",
			args: args{
				obj: []interface{}{
					map[string]interface{}{"key1": "value1", "key2": "value2"},
					map[string]interface{}{"key1": "value3", "key2": "value4"},
				},
				keys:       []string{"key1"},
				defaultVal: "default",
				returnVal:  "all",
			},
			want:    []interface{}{"value1", "value3"},
			wantErr: false,
		},
		{
			name: "Test case 10 - array of maps with multiple matching keys, return last",
			args: args{
				obj: []interface{}{
					map[string]interface{}{"key1": "value1", "key2": "value2"},
					map[string]interface{}{"key1": "value3", "key2": "value4"},
				},
				keys:       []string{"key1"},
				defaultVal: "default",
				returnVal:  "last",
			},
			want:    "value3",
			wantErr: false,
		},
		{
			name: "Test case 11 - array of maps with multiple matching keys, return first",
			args: args{
				obj: []interface{}{
					map[string]interface{}{"key1": "value1", "key2": "value2"},
					map[string]interface{}{"key1": "value3", "key2": "value4"},
				},
				keys:       []string{"key1"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "Test case 12 - array of maps with multiple matching keys, return default",
			args: args{
				obj: []interface{}{
					map[string]interface{}{"key1": "value1", "key2": "value2"},
					map[string]interface{}{"key1": "value3", "key2": "value4"},
				},
				keys:       []string{"key5"},
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeepWalk(tt.args.obj, tt.args.keys, tt.args.defaultVal, tt.args.returnVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeepWalk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepWalk() = %v, want %v", got, tt.want)
			}
		})
	}
}

// defaultCases generate a list of non-existent keys that will be checked against the generated data
func defaultCases(depth, maxDepth int) (interface{}, []string, error) {
	nonexistentKeys := make([]string, rand.Intn(maxDepth)+1)
	for i := range nonexistentKeys {
		for j := 0; j < rand.Intn(100)+1; j++ {
			randomValue, err := util.RandomString(16)
			if err != nil {
				return nil, nil, err
			}
			nonexistentKeys[i] = randomValue
		}
	}
	data, _, err := util.DefaultCases(depth, maxDepth)
	if err != nil {
		return nil, nil, err
	}
	return data, nonexistentKeys, nil
}

func TestDeepwalkRandomSuccess(t *testing.T) {
	for i := 0; i < maxRuns; i++ {
		data, keys, expected, err := util.SucessCases(nil, util.MaxDepth)
		if err != nil {
			t.Errorf("SucessCases() error = %v", err)
		}
		got, err := DeepWalk(data, keys, "default", "all")
		if err != nil {
			t.Errorf("DeepWalk() error = %v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("DeepWalk() = %v, want %v", got, expected)
		}
	}
}

func TestDeepwalkRandomDefault(t *testing.T) {
	for i := 0; i < maxRuns; i++ {
		data, nonexistentKeys, err := defaultCases(0, util.MaxDepth)
		if err != nil {
			t.Errorf("defaultCases() error = %v", err)
		}
		got, err := DeepWalk(data, nonexistentKeys, "default", "all")
		if err := err; err != nil {
			t.Errorf("DeepWalk() error = %v", err)
		}
		if !reflect.DeepEqual(got, "default") {
			t.Errorf("DeepWalk() = %v, want %v", got, "default")
		}
	}
}

func TestDeepWalkWithStruct(t *testing.T) {
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

	tests := []struct {
		name       string
		keys       []string
		defaultVal interface{}
		returnVal  string
		want       interface{}
	}{
		{
			name:       "Test Field1",
			keys:       []string{"Field1"},
			defaultVal: "default",
			returnVal:  "first",
			want:       "test",
		},
		{
			name:       "Test Field2",
			keys:       []string{"Field2"},
			defaultVal: "default",
			returnVal:  "first",
			want:       123,
		},
		{
			name:       "Test NestedField1",
			keys:       []string{"NestedStruct", "NestedField1"},
			defaultVal: "default",
			returnVal:  "first",
			want:       "nested",
		},
		{
			name:       "Test NestedField2",
			keys:       []string{"NestedStruct", "NestedField2"},
			defaultVal: "default",
			returnVal:  "first",
			want:       456,
		},
		{
			name:       "Test Second-level NestedField1",
			keys:       []string{"NestedStruct", "NestedStruct2", "NestedField1"},
			defaultVal: "default",
			returnVal:  "first",
			want:       "nested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultVal, ok := tt.defaultVal.(string)
			if !ok {
				t.Errorf("defaultVal is not a string")
				return
			}
			got, err := DeepWalk(testStruct, tt.keys, defaultVal, tt.returnVal)
			if err != nil {
				t.Errorf("DeepWalk() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("DeepWalk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDeepwalkSuccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, keys, expected, err := util.SucessCases(nil, util.MaxDepth)
		if err != nil {
			b.Errorf("SucessCases() error = %v", err)
		}
		have, err := DeepWalk(data, keys, "default", "all")
		if err != nil {
			b.Errorf("DeepWalk() error = %v", err)
		}
		if !reflect.DeepEqual(have, expected) {
			b.Errorf("DeepWalk() = %v, want %v", have, expected)
		}
	}
}

func BenchmarkDeepwalkDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, nonexistentKeys, err := defaultCases(0, util.MaxDepth)
		if err != nil {
			b.Errorf("defaultCases() error = %v", err)
		}
		have, err := DeepWalk(data, nonexistentKeys, "default", "all")
		if err := err; err != nil {
			b.Errorf("DeepWalk() error = %v", err)
		}
		if !reflect.DeepEqual(have, "default") {
			b.Errorf("DeepWalk() = %v, want %v", have, "default")
		}
	}
}
