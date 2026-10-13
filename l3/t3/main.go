package main

import (
	"fmt"
	"sync"
)

func generator(wg *sync.WaitGroup, cap int) chan int {
	ch := make(chan int)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := range cap {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func squarer(wg *sync.WaitGroup, naturals chan int, squares chan int) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for val := range naturals {
			squares <- val * val
		}
		close(squares)
	}()
}

func printer(wg *sync.WaitGroup, squares chan int) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := range squares {
			fmt.Println(i)
		}
	}()
}

func main() {
	var wg sync.WaitGroup
	cap := 50
	naturals := generator(&wg, cap)
	squares := make(chan int)
	squarer(&wg, naturals, squares)
	printer(&wg, squares)
	wg.Wait()
}
