// fan-in
package main

import (
	"fmt"
	"sync"
)

func mergeChannels(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	retChan := make(chan int)
	for _, channel := range channels {
		go func() {
			defer wg.Done()
			for val := range channel {
				retChan <- val
			}
		}()
	}

	go func() {
		wg.Wait()
		close(retChan)
	}()

	return retChan
}

func sender(ch chan int, n, inc int) {
	defer close(ch)
	for i := 0; i < n; i += inc {
		ch <- i
	}
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go sender(a, 10, 2)
	go sender(b, 20, 3)
	go sender(c, 30, 4)

	for v := range mergeChannels(a, b, c) {
		fmt.Println(v)
	}
}
