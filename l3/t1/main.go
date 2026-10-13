package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// 1000 горутин конкурентно проверяют уникальность по одной мапе => использовать sync.Map
func main() {
	var alreadyStored sync.Map

	capacity := 1000
	doubles := make([]int, 0, capacity)
	for i := 0; i < capacity; i++ {
		doubles = append(doubles, rand.Intn(10))
	}
	uniqueIDs := make(chan int, capacity)
	wg := sync.WaitGroup{}
	for i := 0; i < capacity; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := doubles[i]

			if _, loaded := alreadyStored.LoadOrStore(val, struct{}{}); !loaded {
				uniqueIDs <- val
			}
		}()
	}
	wg.Wait()
	for val := range uniqueIDs {
		fmt.Println(val)
	}
	fmt.Println(uniqueIDs)
}
