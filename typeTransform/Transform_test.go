package typeTransform

import "testing"
import (
	"github.com/bronze1man/kmg/test"
	"reflect"
)

func TestManager(ot *testing.T) {
	t := test.NewTestTools(ot)
	Int := 0
	ArrMapStringInt := []map[string]int{}
	type T1 struct {
		A int
		B string
	}
	ArrStruct := []T1{}
	testCaseTable := []struct {
		in  interface{}
		out interface{}
		exp interface{}
	}{
		{1, &Int, 1},
		{int64(1), &Int, 1},
		{
			[]map[string]string{
				{
					"a": "1",
				},
				{
					"b": "1",
				},
			},
			&ArrMapStringInt,
			[]map[string]int{
				{
					"a": 1,
				},
				{
					"b": 1,
				},
			},
		},
		{
			[]map[string]string{
				{
					"A": "1",
					"B": "abc",
					"C": "abd",
				},
				{
					"A": "",
					"C": "abd",
				},
			},
			&ArrStruct,
			[]T1{
				{
					A: 1,
					B: "abc",
				},
				{
					A: 0,
					B: "",
				},
			},
		},
	}
	for _, testCase := range testCaseTable {
		err := Transform(testCase.in, testCase.out)
		t.Equal(err, nil)
		t.Equal(reflect.ValueOf(testCase.out).Elem().Interface(), testCase.exp)
	}
}
