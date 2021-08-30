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
	// []struct { Height float64 "json:\"height\""; Age int "json:\"age\"" }
	fmt.Println(sv.Type())

	svPtr := reflect.MakeSlice(reflect.SliceOf(newPtrValue.Type()), 0, 0)
	// []*struct { Height float64 "json:\"height\""; Age int "json:\"age\"" }
	fmt.Println(svPtr.Type())

	newSvPtr := reflect.New(reflect.SliceOf(newPtrValue.Type()))
	// *[]*struct { Height float64 "json:\"height\""; Age int "json:\"age\"" }
	fmt.Println("newSvPtr: ", newSvPtr.Type())
	newSvPtrV := newSvPtr.Interface()
	body := []byte(`[{"height":1.5,"age":10}, {"height":8.2,"age":765}]`)
	if err := json.Unmarshal(body, newSvPtrV); err != nil {
		panic(err)
	}
	loopSlice(newSvPtr.Elem().Interface())
	fmt.Println(newSvPtr)

	fmt.Println(`-------`)

	newPtrValue = reflect.New(typ)
	v = newPtrValue.Elem()
	v.Field(0).SetFloat(1.1)
	v.Field(1).SetInt(234)
	fmt.Printf("%#v \n", newPtrValue.Interface())
	// 注意
	svInterface := sv.Interface()
	if err := json.Unmarshal(body, &svInterface); err != nil {
		panic(err)
	}
	// 会变成map
	loopSlice(svInterface)
	fmt.Println("begin loop sv 1 ")
	// 没有数据
	loopSlice(sv.Interface())
	fmt.Println("end loop sv 1 ")
	// 这样，创建一个pointer
	// Create a pointer to a slice value and set it to the slice
	x := reflect.New(sv.Type())
	x.Elem().Set(sv)
	xInterface := x.Interface()

	if err := json.Unmarshal(body, xInterface); err != nil {
		panic(err)
	}
	loopSlice(x.Elem().Interface())

	//sv = reflect.AppendSlice(sv, reflect.ValueOf(svInterface))
	//sv.Set(reflect.ValueOf(svInterface))
	sv = reflect.Append(sv, v)
	fmt.Printf("sv type := %#v \n", sv.Type())

	fmt.Printf("%#v \n", sv.Interface())
	loopSlice(sv.Interface())

	fmt.Println("end sv ---")
	v.Field(0).SetFloat(2.3)
	v.Field(1).SetInt(980)


	// 注意
	newSvPtr = reflect.Append(newSvPtr.Elem(), newPtrValue)
	fmt.Printf("%#v \n", newSvPtr)
	fmt.Printf("type := %#v \n", newSvPtr.Type())
	fmt.Printf("%#v \n", newSvPtr.Interface())
	newSvPtrInterface := newSvPtr.Interface()

	loopSlice(newSvPtrInterface)
}

func loopSlice(t interface{}) {
	sL := reflect.ValueOf(t)
	fmt.Println("loopSlice type: ", sL.Type())
	for i := 0; i < sL.Len(); i++ {
		fmt.Println("item: ", sL.Index(i))
	}
}
