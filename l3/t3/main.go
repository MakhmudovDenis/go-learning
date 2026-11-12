package main

import (
	"fmt"
	"sync"
)

func generator(n int) <-chan int {
	ch := make(chan int)

	go func() {
		for i := range n {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func squarer(naturals <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		for val := range naturals {
			output <- val * val
		}
		close(output)
	}()

	return output
}

func printer(wg *sync.WaitGroup, squares <-chan int) {
	defer wg.Done()
	for i := range squares {
		fmt.Println(i)
	}
}

func main() {
	var wg sync.WaitGroup
	n := 50
	naturals := generator(n)
	squares := squarer(naturals)
	wg.Add(1)
	go printer(&wg, squares)
	wg.Wait()
}
