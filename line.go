package main

import (
	"errors"
)

type UnvalidatedOrderLine struct {
	OrderLineID string
	ProductCode string
	Quantity    float64
}

type ValidatedOrderLine struct {
	OrderLineID OrderLineID
	ProductCode ProductCode
	Quantity    OrderQuantity
}

type OrderLineID string

func ParseOrderLineID(s string) (OrderLineID, error) {
	return OrderLineID(s), nil
}

// ProductCode はプロダクトコードを抽象化した型です。以下の2つが存在します。
//
//   - [Widget]
//   - [Gizmo]
type ProductCode interface {
	privateProductCode()
}

type Widget string

func (Widget) privateProductCode() {}

type Gizmo string

func (Gizmo) privateProductCode() {}

// OrderQuantity は数量を抽象化した型です。以下の2つが存在します。
//
//   - [UnitQuantity]
//   - [KilogramQuantity]
type OrderQuantity interface {
	privateOrderQuantity()
}

type UnitQuantity int

func (UnitQuantity) privateOrderQuantity() {}

type KilogramQuantity float64

func (KilogramQuantity) privateOrderQuantity() {}

func toOrderQuantity(productCode ProductCode, quantity float64) (OrderQuantity, error) {
	switch productCode.(type) {
	case Widget:
		return UnitQuantity(quantity), nil
	case Gizmo:
		return KilogramQuantity(quantity), nil
	default:
		return nil, errors.New("invalid product code")
	}
}

func toValidatedOrderLine(unvalidatedOrderLine *UnvalidatedOrderLine, checkProductCodeExists CheckProductCodeExists) (*ValidatedOrderLine, error) {
	orderLineID, err := ParseOrderLineID(unvalidatedOrderLine.OrderLineID)
	if err != nil {
		return nil, err
	}
	productCode, err := checkProductCodeExists(unvalidatedOrderLine.ProductCode)
	if err != nil {
		return nil, err
	}
	quantity, err := toOrderQuantity(productCode, unvalidatedOrderLine.Quantity)
	if err != nil {
		return nil, err
	}
	return &ValidatedOrderLine{
		OrderLineID: orderLineID,
		ProductCode: productCode,
		Quantity:    quantity,
	}, nil
}
