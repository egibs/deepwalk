package deepsearch

import (
	"reflect"
	"testing"

	util "github.com/egibs/deepwalk/internal/util"
)

func TestDeepSearch(t *testing.T) {
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
			name: "Test case 7 - key found in struct",
			args: args{
				obj: struct {
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
			name: "Test case 8 - key not found in struct",
			args: args{
				obj: struct {
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
			name: "Test case 9 - duplicate key found in map",
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
			got, err := DeepSearch(tt.args.obj, tt.args.searchKey, tt.args.defaultVal, tt.args.returnVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeepSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDeepSearch(b *testing.B) {
	obj := map[string]interface{}{
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

	for i := 0; i < b.N; i++ {
		have, err := DeepSearch(obj, "key5", "default", "first")
		expected := "value5"
		if err != nil {
			b.Errorf("DeepSearch() error = %v", err)
		}
		if !reflect.DeepEqual(have, expected) {
			b.Errorf("DeepSearch() = %v, want %v", have, expected)
		}
	}
}

func BenchmarkDeepsearchSuccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, keys, expected, err := util.SucessCases(nil, util.MaxDepth)
		if err != nil {
			b.Errorf("SucessCases() error = %v", err)
		}
		have, err := DeepSearch(data, keys[len(keys)-1], "default", "all")
		if err != nil {
			b.Errorf("DeepSearch() error = %v", err)
		}
		if !reflect.DeepEqual(have, expected) {
			b.Errorf("DeepSearch() = %v, want %v", have, expected)
		}
	}
}

func BenchmarkDeepsearchDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, _, err := util.DefaultCases(0, util.MaxDepth)
		if err != nil {
			b.Errorf("defaultCases() error = %v", err)
		}
		have, err := DeepSearch(data, "", "default", "all")
		if err := err; err != nil {
			b.Errorf("DeepSearch() error = %v", err)
		}
		if !reflect.DeepEqual(have, "default") {
			b.Errorf("DeepSearch() = %v, want %v", have, "default")
		}
	}
}
