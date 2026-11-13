// Моделирование работы ресторана с использованием `sync.Cond`
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Restaurant struct {
	mu   sync.Mutex
	cond *sync.Cond

	tables   int
	occupied int

	nextTicket int
	nowServing int
}

func NewRestaurant(tables int) *Restaurant {
	r := &Restaurant{tables: tables}
	r.cond = sync.NewCond(&r.mu)
	return r
}

func (r *Restaurant) OccupyTable() {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket := r.nextTicket
	r.nextTicket++

	for ticket != r.nowServing || r.occupied == r.tables {
		r.cond.Wait()
	}

	r.occupied++
	r.nowServing++

	r.cond.Broadcast()
}

func (r *Restaurant) ReleaseTable() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.occupied == 0 {
		panic("restaurant: ReleaseTable called with no occupied tables")
	}

	r.occupied--
	r.cond.Broadcast()
}

func main() {
	const visitors = 20
	r := NewRestaurant(5)

	var wg sync.WaitGroup
	for id := range visitors {
		wg.Go(func() {
			r.OccupyTable()
			defer r.ReleaseTable()

			fmt.Printf("visitor %2d sat down\n", id)
			time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
			fmt.Printf("visitor %2d left\n", id)
		})
	}

	wg.Wait()
}
