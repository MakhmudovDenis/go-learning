// Tee
package main

import (
	"fmt"
	"sync"
	"time"
)

// Реплика БД (имитация)
func dbReplica(name string, in <-chan int) {
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

func Tee(mainCh chan int, output []chan int) {
	go func() {
		for v := range mainCh {
			for _, ch := range output {
				ch <- v
			}
		}
		for _, ch := range output {
			close(ch)
		}
	}()
}

func main() {
	input := make(chan int) // Канал для входящих данных
	replicas := []chan int{ // Реплики БД (каналы)
		make(chan int),
		make(chan int),
		make(chan int),
	}

	var wg sync.WaitGroup
	wg.Add(len(replicas) + 1)

	go func() {
		defer wg.Done()
		for i := range 10 {
			input <- i
		}
		close(input)
	}()

	go Tee(input, replicas)

	for i, ch := range replicas {
		go func() {
			defer wg.Done()
			dbReplica(fmt.Sprintf("db_%d", i), ch)
		}()
	}
	wg.Wait()
}
