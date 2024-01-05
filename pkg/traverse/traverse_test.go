package traverse

import (
	"reflect"
	"testing"

	util "github.com/egibs/deepwalk/internal/util"
)

func TestTraverse(t *testing.T) {
	type args struct {
		obj        interface{}
		searchKey  string
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
			name: "Test case 1 - key found in map",
			args: args{
				obj: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				searchKey:  "key1",
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "Test case 2 - key not found in map",
			args: args{
				obj: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				searchKey:  "key3",
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "Test case 3 - key found in nested map",
			args: args{
				obj: map[string]interface{}{
					"key1": "value1",
					"key2": map[string]interface{}{
						"key3": "value3",
					},
				},
				searchKey:  "key3",
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "value3",
			wantErr: false,
		},
		{
			name: "Test case 4 - key not found in nested map",
			args: args{
				obj: map[string]interface{}{
					"key1": "value1",
					"key2": map[string]interface{}{
						"key3": "value3",
					},
				},
				searchKey:  "key4",
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "Test case 5 - key found in struct",
			args: args{
				obj: &struct {
					Key1 string
					Key2 string
				}{
					Key1: "value1",
					Key2: "value2",
				},
				searchKey:  "Key1",
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "Test case 6 - key not found in struct",
			args: args{
				obj: &struct {
					Key1 string
					Key2 string
				}{
					Key1: "value1",
					Key2: "value2",
				},
				searchKey:  "Key3",
				defaultVal: "default",
				returnVal:  "first",
			},
			want:    "default",
			wantErr: false,
		},
		{
			name: "Test case 7 - duplicate key found in map",
			args: args{
				obj: map[string]interface{}{
					"key1": "value1",
					"key2": map[string]interface{}{
						"key1": "value2",
					},
				},
				searchKey:  "key1",
				defaultVal: "default",
				returnVal:  "all",
			},
			want: []interface{}{"value1", "value2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var returnValueType ReturnControl
			switch tt.args.returnVal {
			case "first":
				returnValueType = First
			case "last":
				returnValueType = Last
			case "all":
				returnValueType = All
			}
			got, err := Traverse(tt.args.obj, tt.args.searchKey, tt.args.defaultVal, returnValueType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Traverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Traverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTraverse(b *testing.B) {
	mapObj := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"key3": "value3",
			"key4": []interface{}{
				"item1",
				"item2",
				map[string]interface{}{
					"key5": "value5",
				},
			},
		},
	}

	var obj interface{} = mapObj
	var searchKey string = "key5"
	var defaultVal string = "default"
	var returnValueType ReturnControl = First

	for i := 0; i < b.N; i++ {
		have, err := Traverse(obj, searchKey, defaultVal, returnValueType)
		expected := "value5"
		if err != nil {
			b.Errorf("Traverse() error = %v", err)
		}
		if !reflect.DeepEqual(have, expected) {
			b.Errorf("Traverse() = %v, want %v", have, expected)
		}
	}
}

func BenchmarkTraverseSuccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dataMap, keys, expected, err := util.SucessCases(nil, util.MaxDepth)
		if err != nil {
			b.Errorf("SucessCases() error = %v", err)
		}
		var data interface{} = dataMap
		var defaultVal string = "default"
		var returnValueType ReturnControl = All
		have, err := Traverse(data, keys[len(keys)-1], defaultVal, returnValueType)
		if err != nil {
			b.Errorf("Traverse() error = %v", err)
		}
		if !reflect.DeepEqual(have, expected) {
			b.Errorf("Traverse() = %v, want %v", have, expected)
		}
	}
}

func BenchmarkTraverseDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dataMap, _, err := util.DefaultCases(0, util.MaxDepth)
		if err != nil {
			b.Errorf("defaultCases() error = %v", err)
		}
		var data interface{} = dataMap
		var searchKey string = ""
		var defaultVal string = "default"
		var returnValueType ReturnControl = All
		have, err := Traverse(data, searchKey, defaultVal, returnValueType)
		if err := err; err != nil {
			b.Errorf("Traverse() error = %v", err)
		}
		if !reflect.DeepEqual(have, "default") {
			b.Errorf("Traverse() = %v, want %v", have, "default")
		}
	}
}
