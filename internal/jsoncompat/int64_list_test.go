package jsoncompat

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestInt64List(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  []int64
	}{
		{`["101",102,"-3"]`, []int64{101, 102, -3}},
		{`["9007199254740993",9223372036854775807,"-9223372036854775808"]`, []int64{9007199254740993, 9223372036854775807, -9223372036854775808}},
		{`[]`, []int64{}},
		{` null `, nil},
	} {
		t.Run(tc.input, func(t *testing.T) {
			ids := Int64List{7}
			if err := json.Unmarshal([]byte(tc.input), &ids); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual([]int64(ids), tc.want) {
				t.Fatalf("got %#v want %#v", ids, tc.want)
			}
		})
	}
}

func TestInt64ListRejectsInvalidElements(t *testing.T) {
	for _, input := range []string{`[1,"9223372036854775808"]`, `[1,-9223372036854775809]`, `[1,1.5]`, `[1,true]`, `[1,null]`, `[1,""]`, `[1,"bad"]`, `[1,{}]`} {
		t.Run(input, func(t *testing.T) {
			ids := Int64List{7}
			err := json.Unmarshal([]byte(input), &ids)
			if err == nil || !strings.Contains(err.Error(), "[1]") {
				t.Fatalf("want element index in error, got %v", err)
			}
			if !reflect.DeepEqual(ids, Int64List{7}) {
				t.Fatalf("failed decode modified destination: %v", ids)
			}
		})
	}
	for _, input := range []string{`{}`, `"123"`, `123`, `[1,`} {
		var ids Int64List
		if err := json.Unmarshal([]byte(input), &ids); err == nil {
			t.Fatalf("accepted %s", input)
		}
	}
}
