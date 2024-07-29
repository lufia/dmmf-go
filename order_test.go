package main

import (
	"errors"
	"testing"
)

func TestValidateOrder_checkProductCode(t *testing.T) {
	order := &UnvalidatedOrder{
		OrderID: "aaa",
		CustomerInfo: &UnvalidatedCustomerInfo{
			FirstName:    "firstName",
			LastName:     "lastName",
			EmailAddress: "foo@example.com",
		},
		ShippingAddress: &UnvalidatedAddress{
			AddressLine1: "karasumaoike",
			City:         "kyoto",
			ZipCode:      "000-0000",
		},
		Lines: []*UnvalidatedOrderLine{
			{
				OrderLineID: "123",
				ProductCode: "g123",
				Quantity:    1.5,
			},
		},
	}
	config := &ValidateOrderConfig{
		CheckAddressExists: func(ua *UnvalidatedAddress) (*CheckedAddress, error) {
			return (*CheckedAddress)(ua), nil
		},
	}

	t.Run("success", func(t *testing.T) {
		c := *config
		c.CheckProductCodeExists = func(s string) (ProductCode, error) {
			return Gizmo(s), nil
		}
		if _, err := c.ValidateOrder(order); err != nil {
			t.Errorf("ValidateOrder: got %v", err)
		}
	})
	t.Run("error", func(t *testing.T) {
		c := *config
		c.CheckProductCodeExists = func(s string) (ProductCode, error) {
			return nil, errors.New("test: invalid product code")
		}
		if _, err := c.ValidateOrder(order); err == nil {
			t.Errorf("ValidateOrder: want an invalid error")
		}
	})
}
