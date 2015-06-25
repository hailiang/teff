package model

import (
	"fmt"
	"reflect"
)

func New(v interface{}) (*Node, error) {
	if v == nil {
		return nil, nil
	}
	return newMaker().objectToNode(reflect.ValueOf(v))
}

func (n *Node) Fill(v interface{}) error {
	if v == nil {
		return nil
	}
	return newFiller().nodeToObject(n, reflect.ValueOf(v))
}

func (m *maker) objectToNode(v reflect.Value) (*Node, error) {
	var c C
	var err error
	switch v.Type().Kind() {
	case reflect.Int, reflect.String:
		c, err = m.objectToValue(v)
	case reflect.Slice, reflect.Array:
		c, err = m.arrayToArray(v)
	case reflect.Ptr:
		return m.ptrToNode(v)
	default:
		err = fmt.Errorf("maker.objectToNode: unsupported type: %v", v.Type())
	}
	if err != nil {
		return nil, err
	}
	return &Node{C: c}, nil
}

func (f *filler) nodeToObject(node *Node, v reflect.Value) error {
	switch v.Type().Kind() {
	case reflect.Int, reflect.String:
		if value, ok := node.C.(Value); ok {
			return f.valueToObject(value, v)
		}
	case reflect.Slice, reflect.Array:
		if array, ok := node.C.(Array); ok {
			return f.arrayToArray(array, v)
		}
	case reflect.Ptr:
		return f.nodeToPtr(node, v)
	}
	return fmt.Errorf("filler.nodeToObject: unsupported type: %v", v.Type())
}

func (m *maker) arrayToArray(v reflect.Value) (Array, error) {
	a := make(Array, v.Len())
	for i := 0; i < v.Len(); i++ {
		node, err := m.objectToNode(v.Index(i))
		if err != nil {
			return nil, err
		}
		a[i] = node
	}
	return a, nil
}

func (f *filler) arrayToArray(a Array, v reflect.Value) error {
	for i, n := range a {
		v.Set(reflect.Append(v, reflect.New(v.Type().Elem()).Elem()))
		elem := v.Index(i)
		if err := f.nodeToObject(n, elem); err != nil {
			return err
		}
	}
	return nil
}

func (m *maker) objectToValue(v reflect.Value) (Value, error) {
	switch v.Type().Kind() {
	case reflect.Int:
		return Value{int(v.Int())}, nil
	case reflect.String:
		return Value{v.String()}, nil
	}
	return Value{}, fmt.Errorf("maker.objectToValue: unsupported type: %v", v.Type())
}

func (f *filler) valueToObject(value Value, v reflect.Value) error {
	switch v.Type().Kind() {
	case reflect.Int, reflect.String:
		v.Set(reflect.ValueOf(value.V))
		return nil
	case reflect.Ptr:
		return f.valueToPtr(value, v)
	}
	return fmt.Errorf("filler.valueToObject: unsupported type: %v", v.Type())
}

func allocIndirect(v reflect.Value) reflect.Value {
	alloc(v)
	return reflect.Indirect(v)
}

func alloc(v reflect.Value) reflect.Value {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return v
}
