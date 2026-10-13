// Синхронизация горутин с использованием `sync.WaitGroup`
package main

import (
	"fmt"
	"sync"
)

func main() {
	cnt := 100
	var wg sync.WaitGroup

	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}

	wg.Wait()
}
