package main

import (
	"github.com/lufia/go-validator"

	"github.com/lufia/dmmf-go/internal/pipe"
)

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

func PlaceOrder(order *UnvalidatedOrder) {
	// TODO: set dependencies
	var (
		validateOrderConfig    ValidateOrderConfig
		priceOrderConfig       PriceOrderConfig
		acknowledgeOrderConfig AcknowledgeOrderConfig
	)
	p1 := pipe.Value(order)
	p2 := pipe.From(p1, validateOrderConfig.ValidateOrder)
	p3 := pipe.From(p2, priceOrderConfig.PriceOrder)
	p4 := pipe.From(p3, acknowledgeOrderConfig.AcknowledgeOrder)
	// TODO: events
	_ = p4.Value()
}

func main() {
}
