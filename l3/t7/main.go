// декоратор
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type ServerMetric struct {
	Name  string  // Название метрики (например, "memory_usage")
	Value float64 // Значение в байтах
}

func transform[T any](input <-chan T, action func(T) T) <-chan T {
	output := make(chan T)
	go func() {
		defer close(output)
		for val := range input {
			output <- action(val)
		}
	}()

	return output
}

func bytesToMegabytes[T any](m ServerMetric) ServerMetric {
	m.Value = m.Value / 1024 / 1024
	return m
}

func main() {
	letters := []byte("abcdefghjklmnopqrstuvwxyz") // имена метрик

	metrics := make(chan ServerMetric)
	var wg sync.WaitGroup

	wg.Go(func() {
		defer close(metrics)
		for i := range 10 {
			metrics <- ServerMetric{
				Name:  string(letters[i]),
				Value: rand.ExpFloat64() * 999_777_555,
			}
		}
	})

	resChan := transform(metrics, bytesToMegabytes[ServerMetric])

	wg.Go(func() {
		for v := range resChan {
			fmt.Println(v)
		}
	})

	wg.Wait()
}
