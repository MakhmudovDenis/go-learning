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

func sender(wg *sync.WaitGroup, ch chan int, cap, inc int) {
	defer wg.Done()
	for i := 0; i < cap; i += inc {
		ch <- i
	}
	close(ch)
}

func main() {
	var wg sync.WaitGroup
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	wg.Add(3)
	go sender(&wg, a, 10, 2)
	go sender(&wg, b, 20, 3)
	go sender(&wg, c, 30, 4)

	d := mergeChannels(a, b, c)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for val := range d {
			fmt.Println(val)
		}
	}()

	wg.Wait()
}
