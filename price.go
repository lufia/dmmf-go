package main

import (
	"github.com/lufia/dmmf-go/billing"
)

type Price float64

func (p Price) Value() float64 {
	return float64(p)
}

func (p Price) Mul(n float64) Price {
	return Price(float64(p) * n)
}

type GetProductPrice func(code ProductCode) (Price, error)

type PricedOrder struct {
	OrderID         OrderID
	CustomerInfo    *CustomerInfo
	ShippingAddress *Address
	Lines           []*PricedOrderLine
	AmountToBill    billing.Amount
}

type PricedOrderLine struct {
	ValidatedOrderLine
	LinePrice Price
}

func PriceOrder(order *ValidatedOrder, getProductPrice GetProductPrice) (*PricedOrder, error) {
	prices := make([]Price, len(order.Lines))
	lines := make([]*PricedOrderLine, len(order.Lines))
	for i, l := range order.Lines {
		p, err := toPricedOrderLine(l, getProductPrice)
		if err != nil {
			return nil, err
		}
		prices[i] = p.LinePrice
		lines[i] = p
	}
	return &PricedOrder{
		OrderID:         order.OrderID,
		CustomerInfo:    order.CustomerInfo,
		ShippingAddress: order.ShippingAddress,
		Lines:           lines,
		AmountToBill:    billing.Sum(prices),
	}, nil
}

func toPricedOrderLine(line *ValidatedOrderLine, getProductPrice GetProductPrice) (*PricedOrderLine, error) {
	n := line.Quantity.Value()
	price, err := getProductPrice(line.ProductCode)
	if err != nil {
		return nil, err
	}
	return &PricedOrderLine{
		ValidatedOrderLine: *line,
		LinePrice:          price.Mul(n),
	}, nil
}
