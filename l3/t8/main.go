// pipeline
package main

import (
	"fmt"
	"sync"
)

type data[T any] struct {
	value  T
	parsed bool
	sent   bool
}

func parse[T any](input <-chan T) <-chan data[T] {
	output := make(chan data[T])

	go func() {
		defer close(output)

		for val := range input {
			output <- data[T]{
				value:  val,
				parsed: true,
			}
		}
	}()

	return output
}

func fanOut[T any](input <-chan data[T], n int) []<-chan data[T] {
	output := make([]chan data[T], n)
	for i := range output {
		output[i] = make(chan data[T])
	}

	go func() {
		idx := 0
		for v := range input {
			output[idx] <- v
			idx = (1 + idx) % n
		}

		for _, ch := range output {
			close(ch)
		}
	}()

	res := make([]<-chan data[T], n)
	for i, ch := range output {
		res[i] = ch
	}

	return res
}

func send[T any](input <-chan data[T]) <-chan data[T] {
	output := make(chan data[T])

	go func() {
		defer close(output)

		for val := range input {
			val.sent = true
			output <- val
		}
	}()

	return output
}

func main() {
	chanCount := 4
	input := make(chan int)
	var wg sync.WaitGroup

	wg.Go(func() {
		defer func() {
			close(input)
			wg.Done()
		}()
		for i := range 30 {
			input <- i
		}
	})

	wg.Add(1 + chanCount)

	go func() {
		defer wg.Done()
		for i, ch := range fanOut(parse(input), chanCount) {
			sent := send(ch)
			go func() {
				defer wg.Done()
				for v := range sent {
					fmt.Printf("получено %v из канала %d\n", v, i)
				}
			}()
		}
	}()

	wg.Wait()
}
