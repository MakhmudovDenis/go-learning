// Обработка паники при делении на ноль :wq
package main

import "fmt"

func SafeDivide(a, b int) (ret int) {
	defer func() {
		if r := recover(); r != nil {
			ret = 0
		}
	}()

	if b == 0 {
		panic("деление на ноль")
	}

	return a / b
}

func main() {
	fmt.Println(SafeDivide(10, 2))
	fmt.Println(SafeDivide(10, 0))
}
