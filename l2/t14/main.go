// Обработка паники в многоуровневых вызовах функций
package main

import "fmt"

func level1() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("паника обработана на уровне level1: %v", r)
		}
	}()
	level2()
}

func level2() {
	defer fmt.Println("завершаем level2")
	level3()
}

func level3() {
	panic("паника в level3")
}

func main() {
	level1()
}
