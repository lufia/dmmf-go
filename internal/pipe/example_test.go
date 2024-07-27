package pipe_test

import (
	"fmt"
	"strings"

	"github.com/lufia/dmmf-go/internal/pipe"
)

func tee[T any](v T) T {
	fmt.Println(v)
	return v
}

func ExampleValue() {
	v := pipe.Value("hello world")(tee)(strings.ToUpper)
	a := pipe.From(v, strings.Fields)
	fmt.Println(a.Value())
	// Output:
	// hello world
	// [HELLO WORLD]
}
