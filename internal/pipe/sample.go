//go:build ignore

package main

import (
	"fmt"
	"reflect"
)

// go:noinline
func do[T any](v T) func() T {
	var f func() T

	f = func() T {
		return v
	}
	//fmt.Println("addr:", addr(f))
	return f
}

func addr[T any](f func() T) uintptr {
	return reflect.ValueOf(f).Pointer()
}

func main() {
	f1 := do(10)
	f2 := do(20)
	f3 := do("aaa")
	f4 := do("aaa")
	fmt.Println(f1(), reflect.ValueOf(f1).Pointer(), addr(f1))
	fmt.Println(f2(), reflect.ValueOf(f2).Pointer(), addr(f2))
	fmt.Println(f3(), reflect.ValueOf(f3).Pointer(), addr(f3))
	fmt.Println(f4(), reflect.ValueOf(f4).Pointer(), addr(f4))
}
