package main

import (
	"context"
	"errors"

	"github.com/lufia/go-validator"
)

// OrderID は注文を一意に特定する文字列です。ゼロ値は不正な値です。
type OrderID string

func ParseOrderID(s string) (OrderID, error) {
	switch {
	case s == "":
		return "", errors.New("order ID must not be empty")
	case len(s) > 50:
		return "", errors.New("order ID must not be more than 50 chars")
	default:
		return OrderID(s), nil
	}
}

type UnvalidatedOrder struct {
	OrderID         string
	CustomerInfo    *UnvalidatedCustomerInfo
	ShippingAddress *UnvalidatedAddress
	Lines           []*UnvalidatedOrderLine
}

type UnvalidatedCustomerInfo struct {
	FirstName    string
	LastName     string
	EmailAddress string
}

type CustomerInfo struct {
	Name         PersonalName
	EmailAddress EmailAddress
}

type PersonalName struct {
	FirstName string
	LastName  string
}

type EmailAddress string

func ParseEmailAddress(s string) (EmailAddress, error) {
	return EmailAddress(s), nil
}

type ValidatedOrder struct {
	OrderID         OrderID
	CustomerInfo    *CustomerInfo
	ShippingAddress *Address
	Lines           []*ValidatedOrderLine
}

type ValidateOrderConfig struct {
	CheckProductCodeExists CheckProductCodeExists
	CheckAddressExists     CheckAddressExists
}

var (
	validateEmailAddress = validator.New(func(s string) bool {
		_, err := ParseEmailAddress(s)
		return err == nil
	})
	validateCustomerInfo = validator.Struct(func(s validator.StructRule, r *UnvalidatedCustomerInfo) {
		validator.AddField(s, &r.FirstName, "firstName", validateString50)
		validator.AddField(s, &r.LastName, "lastName", validateString50)
		validator.AddField(s, &r.EmailAddress, "emailAddress", validateEmailAddress)
	})
)

func (config *ValidateOrderConfig) ValidateOrder(order *UnvalidatedOrder) (*ValidatedOrder, error) {
	orderID, err := ParseOrderID(order.OrderID)
	if err != nil {
		return nil, err
	}
	customerInfo, err := toCustomerInfo(order.CustomerInfo)
	if err != nil {
		return nil, err
	}
	shippingAddress, err := toAddress(order.ShippingAddress, config.CheckAddressExists)
	if err != nil {
		return nil, err
	}
	lines, err := toValidatedOrderLines(order.Lines, config.CheckProductCodeExists)
	if err != nil {
		return nil, err
	}
	return &ValidatedOrder{
		OrderID:         orderID,
		CustomerInfo:    customerInfo,
		ShippingAddress: shippingAddress,
		Lines:           lines,
	}, nil
}

func toCustomerInfo(info *UnvalidatedCustomerInfo) (*CustomerInfo, error) {
	ctx := context.Background()
	if err := validateCustomerInfo.Validate(ctx, info); err != nil {
		return nil, err
	}
	return &CustomerInfo{
		Name: PersonalName{
			FirstName: info.FirstName,
			LastName:  info.LastName,
		},
		EmailAddress: Must(ParseEmailAddress(info.EmailAddress)),
	}, nil
}

func toValidatedOrderLines(unvalidatedOrderLines []*UnvalidatedOrderLine, checkProductCodeExists CheckProductCodeExists) ([]*ValidatedOrderLine, error) {
	lines := make([]*ValidatedOrderLine, len(unvalidatedOrderLines))
	for i, u := range unvalidatedOrderLines {
		v, err := toValidatedOrderLine(u, checkProductCodeExists)
		if err != nil {
			return nil, err
		}
		lines[i] = v
	}
	return lines, nil
}
