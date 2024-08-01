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
	return f
}

func addr[T any](f func() T) uintptr {
	return reflect.ValueOf(f).Pointer()
}

//go:noinline
func nested[T any](v T) func() T {
	return do(v)
}

func main() {
	fmt.Printf("do[int] = 0x%x\n", reflect.ValueOf(do[int]).Pointer())
	f1 := do(10)
	f2 := do(20)
	f3 := do("aaa")
	f4 := do("aaa")
	fmt.Printf("%v 0x%x 0x%x\n", f1(), reflect.ValueOf(f1).Pointer(), addr(f1))
	fmt.Printf("%v 0x%x 0x%x\n", f2(), reflect.ValueOf(f2).Pointer(), addr(f2))
	fmt.Printf("%v 0x%x 0x%x\n", f3(), reflect.ValueOf(f3).Pointer(), addr(f3))
	fmt.Printf("%v 0x%x 0x%x\n", f4(), reflect.ValueOf(f4).Pointer(), addr(f4))

	g1 := nested(1)
	g2 := nested(2)
	fmt.Printf("%v 0x%x 0x%x\n", g1(), reflect.ValueOf(g1).Pointer(), addr(g1))
	fmt.Printf("%v 0x%x 0x%x\n", g2(), reflect.ValueOf(g2).Pointer(), addr(g2))
}
