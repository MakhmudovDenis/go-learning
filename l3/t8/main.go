// pipeline
package main

import (
	"fmt"
	"sync"
)

func generate(items ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, it := range items {
			out <- it
		}
	}()
	return out
}

func parse(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			out <- "parsed - " + v
		}
	}()
	return out
}

func split(in <-chan string, n int) []<-chan string {
	outs := make([]chan string, n)
	result := make([]<-chan string, n)
	for i := range n {
		outs[i] = make(chan string)
		result[i] = outs[i]
	}

	go func() {
		defer func() {
			for _, ch := range outs {
				close(ch)
			}
		}()

		i := 0
		for v := range in {
			outs[i] <- v
			i = (i + 1) % n
		}
	}()

	return result
}

func send(ins []<-chan string) <-chan string {
	out := make(chan string)

	var wg sync.WaitGroup
	for _, in := range ins {
		wg.Go(func() {
			for v := range in {
				out <- "sent - " + v
			}
		})
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	items := make([]string, 10)
	for i := range items {
		items[i] = fmt.Sprintf("item %d", i)
	}

	results := send(split(parse(generate(items...)), 4))

	for r := range results {
		fmt.Println(r)
	}
}
