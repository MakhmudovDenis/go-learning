// Паники и их обработка в Go
package main

import "fmt"

func CausePanic() {
	panic("Что-то пошло не так")
}

func HandlePanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("паника перехвачена")
		}
	}()

	CausePanic()
}

func main() {
	// CausePanic()
	HandlePanic()
}
