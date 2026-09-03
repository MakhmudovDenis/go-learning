package main

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrInvalidAmount       = errors.New("некорректная сумма платежа")
	ErrProviderUnavailable = errors.New("провайдер недоступен")
)

func main() {
	fmt.Println(rand.Intn(2) == 0)
}

type PaymentProcessor interface {
	ProcessPayment(float64) error
}

type (
	SberPP  struct{ APIKey string }
	AlfaPP  struct{ APIKey string }
	TbankPP struct{ APIKey string }
)

func (p SberPP) ProcessPayment(amount float64) error {
	if isProvicerAvailable() {
		return ErrProviderUnavailable
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}

	return nil
}

func (p AlfaPP) ProcessPayment(amount float64) error {
	return nil
}

func (p TbankPP) ProcessPayment(amount float64) error {
	return nil
}

func isProvicerAvailable() bool {
	return rand.Intn(2) == 0
}
