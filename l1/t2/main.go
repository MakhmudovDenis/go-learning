package main

import (
	"fmt"
	"math/rand"
)

func main() {
}

type PaymentProcessor interface {
	ProcessPayment(float64) error
}

type (
	SberPP  struct{}
	AlfaPP  struct{}
	TbankPP struct{}
)

func (p SberPP) ProcessPayment(amount float64) error {
	return nil
}

func (p AlfaPP) ProcessPayment(amount float64) error {
	return nil
}

func (p TbankPP) ProcessPayment(amount float64) error {
	return nil
}
