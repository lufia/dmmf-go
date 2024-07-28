package main

import (
	"context"

	"github.com/lufia/go-validator"
)

type UnvalidatedAddress struct {
	AddressLine1 string
	AddressLine2 string
	AddressLine3 string
	AddressLine4 string
	City         string
	ZipCode      string
}

type CheckedAddress UnvalidatedAddress

type CheckAddressExists func(address *UnvalidatedAddress) (*CheckedAddress, error)

type Address struct {
	AddressLine1 string
	AddressLine2 string
	AddressLine3 string
	AddressLine4 string
	City         string
	ZipCode      ZipCode
}

type ZipCode string

func ParseZipCode(s string) (ZipCode, error) {
	return ZipCode(s), nil
}

var validateZipCode = validator.New(func(s string) bool {
	_, err := ParseZipCode(s)
	return err == nil
})

var validateAddress = validator.Struct(func(s validator.StructRule, r *CheckedAddress) {
	validator.AddField(s, &r.AddressLine1, "address1", validateString50)
	validator.AddField(s, &r.AddressLine2, "address2", validateOptionString50)
	validator.AddField(s, &r.AddressLine3, "address3", validateOptionString50)
	validator.AddField(s, &r.AddressLine4, "address4", validateOptionString50)
	validator.AddField(s, &r.City, "city", validateString50)
	validator.AddField(s, &r.ZipCode, "zipCode", validateZipCode)
})

func toAddress(address *UnvalidatedAddress, checkAddressExists CheckAddressExists) (*Address, error) {
	checkedAddress, err := checkAddressExists(address)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := validateAddress.Validate(ctx, checkedAddress); err != nil {
		return nil, err
	}
	return &Address{
		AddressLine1: checkedAddress.AddressLine1,
		AddressLine2: checkedAddress.AddressLine2,
		AddressLine3: checkedAddress.AddressLine3,
		AddressLine4: checkedAddress.AddressLine4,
		City:         checkedAddress.City,
		ZipCode:      Must(ParseZipCode(checkedAddress.ZipCode)),
	}, nil
}
