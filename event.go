package main

import (
	"iter"

	"github.com/lufia/dmmf-go/billing"
)

type PlacedOrderEvent interface {
	privatePlacedOrderEvent()
}

type OrderPlaced PricedOrder

func (*OrderPlaced) privatePlacedOrderEvent() {}

type OrderAcknowledgmentSent struct {
	OrderID      OrderID
	EmailAddress EmailAddress
}

func (*OrderAcknowledgmentSent) privatePlacedOrderEvent() {}

type BillableOrderPlaced struct {
	OrderID        OrderID
	BillingAddress *Address
	AmountToBill   billing.Amount
}

func (*BillableOrderPlaced) privatePlacedOrderEvent() {}

// Events returns a iterator. If event is nil, means ack was not send.
func Events(order *PricedOrder, seq iter.Seq[*OrderAcknowledgmentSent]) iter.Seq[PlacedOrderEvent] {
	return func(yield func(PlacedOrderEvent) bool) {
		yield((*OrderPlaced)(order))
		for e := range seq {
			yield(e)
		}
		// TODO: createBillingEvent ch9.4.2 (p170)
	}
}
