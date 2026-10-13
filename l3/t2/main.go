package main

import (
	"time"
)

func worker() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
	}()
	return ch
}

func main() {
	timeStart := time.Now()
	//_, _ = <-worker(), <-worker() // вторая горутина начинает выполняться только когда из канала от первого вызова worker приходит значение => 6 секунд
	ch1, ch2 := worker(), worker()
	_, _ = <-ch1, <-ch2
	println(int(time.Since(timeStart).Seconds()))
}
