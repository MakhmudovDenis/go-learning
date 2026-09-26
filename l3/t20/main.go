// Синхронизация горутин с использованием `sync.WaitGroup`
package main

import "fmt"

func main() {
	cnt := 100
	for i := 0; i < cnt; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}
