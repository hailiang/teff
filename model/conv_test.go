package model

import (
	"fmt"
	"reflect"
	"testing"
)

/*
TODO:
1. mismatch type for struct field
2. ignore setting unexported field
3. reading unexported field
4. type S []S
*/

func TestModel(t *testing.T) {
	for i, testcase := range []struct {
		v interface{}
		l List
	}{
		{nil, nil},

		{1, List{{Value: 1}}},
		{pi(1), List{{Value: 1}}},
		{"a", List{{Value: "a"}}},
		{ps("a"), List{{Value: "a"}}},

		{
			[]int{},
			List{},
		},
		{
			[]string{"a"},
			List{{Value: "a"}},
		},
		{
			[]int{1, 2},
			List{{Value: 1}, {Value: 2}},
		},
		{
			[][]int{{1, 2}, {3}},
			List{
				{List: List{{Value: 1}, {Value: 2}}},
				{List: List{{Value: 3}}},
			},
		},

		{
			[]*int{pi(1), pi(2)},
			List{{Value: 1}, {Value: 2}},
		},
		{
			func() []*int {
				i := pi(3)
				return []*int{i, i}
			}(),
			List{{RefID: RefID("1"), Value: 3}, {Value: RefID("1")}},
		},
		{
			struct{}{},
			List{},
		},
		{
			struct {
				I int
				S string
			}{1, "a"},
			List{
				{Value: Identifier("I"), List: List{{Value: 1}}},
				{Value: Identifier("S"), List: List{{Value: "a"}}},
			},
		},
		{
			func() struct {
				S1 *string
				S2 *string
			} {
				s := "a"
				return struct {
					S1 *string
					S2 *string
				}{&s, &s}
			}(),
			List{
				{Value: Identifier("S1"), List: List{{RefID: "1", Value: "a"}}},
				{Value: Identifier("S2"), List: List{{Value: RefID("1")}}},
			},
		},
		{
			func() *struct { // return pointer so that S1 is addressable and can be correctly referenced.
				S1 string
				S2 *string
				//S3 **string
			} {
				s := struct {
					S1 string
					S2 *string
					//S3 **string
				}{S1: "a"}
				s.S2 = &s.S1
				//s.S3 = &s.S2
				return &s
			}(),
			List{
				{Value: Identifier("S1"), List: List{{RefID: "1", Value: "a"}}},
				{Value: Identifier("S2"), List: List{{Value: RefID("1")}}},
				//{Value: Identifier("S3"), List: List{{Value: RefID("1")}}},
			},
		},

		//{
		//	func() struct {
		//		S3 ***string
		//		S2 **string
		//		S1 *string
		//	} {
		//		s := "a"
		//		ps := &s
		//		pps := &ps
		//		ppps := &pps
		//		v := struct {
		//			S3 ***string
		//			S2 **string
		//			S1 *string
		//		}{ppps, pps, ps}
		//		return v
		//	}(),
		//	List{
		//		{RefID: "1", Value: IdentValue{"S1", "a"}},
		//		{Value: IdentValue{"S2", RefID("1")}},
		//	},
		//},
	} {
		//if i != 14 {
		//	continue
		//}
		{
			list, err := New(testcase.v)
			if err != nil {
				t.Fatalf("testcase %d: New: %v", i, err)
			}
			if !reflect.DeepEqual(list, testcase.l) {
				t.Fatalf("testcase %d: New: mismatch, expect \n%v\ngot\n%v", i, testcase.l, list)
			}
		}
		{
			v := newValueOf(testcase.v)
			if err := Fill(testcase.l, v); err != nil {
				t.Fatalf("testcase %d: Fill: %v", i, err)
			}
			list, err := New(v)
			if err != nil {
				t.Fatalf("testcase %d: New: %v", i, err)
			}
			if !reflect.DeepEqual(list, testcase.l) {
				t.Fatalf("testcase %d: Fill: mismatch, expect \n%v\ngot\n%v", i, testcase.l, list)
			}
		}
	}
}

func newValueOf(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return reflect.New(reflect.TypeOf(v)).Interface()
}

var p = fmt.Println

func pi(i int) *int {
	return &i
}

func ps(s string) *string {
	return &s
}

func (n *Node) String() string {
	return fmt.Sprintf("%#v", *n)
}
