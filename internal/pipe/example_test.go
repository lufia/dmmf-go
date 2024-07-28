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

func require[T ~string](v T) (T, error) {
	if len(v) == 0 {
		return "", fmt.Errorf("zero length")
	}
	return v, nil
}

func ExampleValue() {
	v := pipe.Value("hello world").Catch(require)(tee)(strings.ToUpper)
	a := pipe.From(v, pipe.Errorable(strings.Fields))
	fmt.Println(a.Value())
	// Output:
	// hello world
	// [HELLO WORLD]
}
