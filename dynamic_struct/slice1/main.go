package main

import (
	"fmt"
	"reflect"
)

type My struct {
	Name string
	Id   int
}

func main() {
	my := &My{}
	myType := reflect.TypeOf(my)
	slice := reflect.MakeSlice(reflect.SliceOf(myType), 10, 10)
	p := slice.Interface().([]*My)
	fmt.Printf("%T %d %d ", p, len(p), cap(p))
	fmt.Println(slice.Type())
}
