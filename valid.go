package main

import (
	"github.com/lufia/go-validator"
)

func ParseEmailAddress(s string) (EmailAddress, error) {
	return EmailAddress(s), nil
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var (
	validateString50       = validator.Length[string](1, 50)
	validateOptionString50 = validator.Length[string](0, 50)
)

func main() {
}
