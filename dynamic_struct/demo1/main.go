package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "Height",
			Type: reflect.TypeOf(float64(0)),
			Tag:  `json:"height"`,
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"age"`,
		},
	})

	newPtrValue := reflect.New(typ)
	v := newPtrValue.Elem()
	v.Field(0).SetFloat(0.4)
	v.Field(1).SetInt(2)
	s := v.Addr().Interface()

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

	fmt.Printf("value: %+v\n", s)
	fmt.Printf("json:  %s", w.Bytes())

	r := bytes.NewReader([]byte(`{"height":1.5,"age":10}`))
	if err := json.NewDecoder(r).Decode(s); err != nil {
		panic(err)
	}
	fmt.Printf("value: %+v\n", s)

	sv := reflect.MakeSlice(reflect.SliceOf(v.Type()), 0, 0)
	fmt.Println(sv.Type())

	svPtr := reflect.MakeSlice(reflect.SliceOf(newPtrValue.Type()), 0, 0)
	fmt.Println(svPtr.Type())
	/*
		[]struct { Height float64 "json:\"height\""; Age int "json:\"age\"" }
		[]*struct { Height float64 "json:\"height\""; Age int "json:\"age\"" }
	*/

	newSvPtr := reflect.New(reflect.SliceOf(newPtrValue.Type()))
	fmt.Println(newSvPtr.Type())
	newSvPtrV := newSvPtr.Interface()
	fmt.Println(newSvPtrV)

	newPtrValue = reflect.New(typ)
	v = newPtrValue.Elem()
	v.Field(0).SetFloat(1.1)
	v.Field(1).SetInt(234)
	fmt.Printf("%#v \n", newPtrValue.Interface())
	// 注意
	sv = reflect.Append(sv, v)
	fmt.Printf("%#v \n", sv.Interface())
	loopSlice(sv.Interface())

	v.Field(0).SetFloat(2.3)
	v.Field(1).SetInt(980)

	// 注意
	newSvPtr = reflect.Append(newSvPtr.Elem(), newPtrValue)
	fmt.Printf("%#v \n", newSvPtr)
	fmt.Printf("type := %#v \n", newSvPtr.Type())
	fmt.Printf("%#v \n", newSvPtr.Interface())


	loopSlice(newSvPtr.Interface())
}

func loopSlice(t interface{}) {
	sL := reflect.ValueOf(t)

	for i := 0; i < sL.Len(); i++ {
		fmt.Println("item: ", sL.Index(i))
	}
}
