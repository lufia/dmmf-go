package pipe_test

import (
	"fmt"
	"net/url"

	"github.com/lufia/dmmf-go/internal/pipe"
)

func WithPath(s string) func(u *url.URL) *url.URL {
	return func(u *url.URL) *url.URL {
		u.Path = s
		return u
	}
}

func WithParam(k, v string) func(u *url.URL) *url.URL {
	return func(u *url.URL) *url.URL {
		q := u.Query()
		q.Set(k, v)
		u.RawQuery = q.Encode()
		return u
	}
}

func ExampleValue_url() {
	p := pipe.From(pipe.Value("https://example.com"), url.Parse).
		Then(WithPath("/query"))(WithParam("key", "value"))
	fmt.Println(p.Value().String())
	// Output: https://example.com/query?key=value
}
