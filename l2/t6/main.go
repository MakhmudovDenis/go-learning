// Анализ цепочек ошибок в Go error.Is
package main

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrNotFound  = errors.New("ресурс не найден")
	TimeoutError = errors.New("таймаут операции")
)

func SimulateRequest() error {
	randInt := rand.Intn(10)

	switch {
	case randInt < 5:
		return fmt.Errorf("запрос не выполнен: %w", TimeoutError)
	case randInt < 8:
		return fmt.Errorf("ошибка: %w", ErrNotFound)
	default:
		return errors.New("неизвестная ошибка")
	}
}

func ProcessError(err error) {
	if errors.Is(err, TimeoutError) {
		fmt.Println("Требуется повторная попытка")
		return
	}

	if errors.Is(err, ErrNotFound) {
		fmt.Println("Ресурс не найден")
		return
	}

	fmt.Println("Неизвестная ошибка")
}

func main() {
	for range 10 {
		ProcessError(SimulateRequest())
	}
}
