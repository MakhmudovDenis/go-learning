// Обработка паники в многоуровневых вызовах функций
package main

import "fmt"

func level1() {
	level2()
}

func level2() {
	level3()
}

func level3() {
	panic("паника в level3")
}
