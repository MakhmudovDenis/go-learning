// Возврат ошибки без дополнительных пакетов
package main

import (
	"fmt"
)

type MyError struct {
	message string
}

func (e MyError) Error() string {
	return e.message
}

func handle() error {
	// return fmt.Errorf("шутка")

	return MyError{message: "произошла чудовищная ошибка"}
}

func main() {
	var err error = handle() // проверка имплементации интерфейса
	fmt.Println(err)
}
