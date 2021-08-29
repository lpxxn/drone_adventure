package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type A struct {
	Name string `column:"email"`
}

func TestSlice1(t *testing.T) {
	bbb(&A{})
}

func aaa(v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	} else {
		panic("Input param is not a slice")
	}

	sl := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		sl = sl.Elem()
	}

	st := sl.Type()
	fmt.Printf("Slice Type %s:\n", st)

	sliceType := st.Elem()
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	fmt.Printf("Slice Elem Type %v:\n", sliceType)

	for i := 0; i < 5; i++ {
		newitem := reflect.New(sliceType)
		newitem.Elem().FieldByName("Name").SetString(fmt.Sprintf("Grzes %d", i))

		s := newitem.Elem()
		for i := 0; i < s.NumField(); i++ {
			col := s.Type().Field(i).Tag.Get("column")
			fmt.Println(col, s.Field(i).Addr().Interface())
		}

		sl.Set(reflect.Append(sl, newitem))
	}
}

func bbb(v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("Input param is not a struct")
	}

	models := reflect.New(reflect.SliceOf(reflect.TypeOf(v))).Interface()

	aaa(models)

	fmt.Println(models)

	b, err := json.Marshal(models)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
