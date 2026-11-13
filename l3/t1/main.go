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
	for range capacity {
		doubles = append(doubles, rand.Intn(10))
	}
	uniqueIDs := make(chan int, capacity)
	wg := sync.WaitGroup{}
	for i := range capacity {
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
	close(uniqueIDs)
	for val := range uniqueIDs {
		fmt.Println(val)
	}
	fmt.Println(uniqueIDs)
}
