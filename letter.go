package main

type HTMLString string

type CreateOrderAcknowledgmentLetter func(order *PricedOrder) (HTMLString, error)

type OrderAcknowledgment struct {
	EmailAddress EmailAddress
	Letter       HTMLString
}

type SendOrderAcknowledgment func(ack *OrderAcknowledgment) error

type AcknowledgeOrderConfig struct {
	CreateOrderAcknowledgmentLetter CreateOrderAcknowledgmentLetter
	SendOrderAcknowledgment         SendOrderAcknowledgment
}

func (config *AcknowledgeOrderConfig) AcknowledgeOrder(order *PricedOrder) (bool, error) {
	letter, err := config.CreateOrderAcknowledgmentLetter(order)
	if err != nil {
		return false, err
	}
	ack := &OrderAcknowledgment{
		EmailAddress: order.CustomerInfo.EmailAddress,
		Letter:       letter,
	}
	if err := config.SendOrderAcknowledgment(ack); err != nil {
		return false, err
	}
	return true, nil
}
