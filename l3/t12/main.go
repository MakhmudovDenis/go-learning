// Моделирование работы ресторана с использованием `sync.Cond`
package main

import (
	"fmt"
	"sync"
)

type Restaurant struct {
	notFull  sync.Cond
	notEmpty sync.Cond
	mu       sync.Mutex

	cap            int
	occupiedTables int
}

func NewRestaurant(cap int) *Restaurant {
	q := &Restaurant{
		cap:            cap,
		occupiedTables: 0,
	}

	q.notEmpty = *sync.NewCond(&q.mu)
	q.notFull = *sync.NewCond(&q.mu)

	return q
}

func (r *Restaurant) OccupyTable() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for r.occupiedTables == r.cap {
		r.notFull.Wait()
	}

	r.occupiedTables++

	r.notEmpty.Signal()
	fmt.Println("occupied, now", r.occupiedTables)
}

func (r *Restaurant) ReleaseTable() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for r.occupiedTables == 0 {
		r.notEmpty.Wait()
	}

	r.occupiedTables--

	r.notFull.Signal()
	fmt.Println("released, now", r.occupiedTables)
}

func main() {
	visitorsCount := 100
	r := NewRestaurant(5)

	var wg sync.WaitGroup
	wg.Go(func() {
		for range visitorsCount {
			r.OccupyTable()
		}
	})

	wg.Go(func() {
		for range visitorsCount {
			r.ReleaseTable()
		}
	})

	wg.Wait()
}
