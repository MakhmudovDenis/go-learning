// Реализация очереди с ограниченной емкостью на sync.Cond
package main

import (
	"fmt"
	"sync"
)

type BoundedQueue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	tasks    []any
	cap      int
	shutdown bool
}

func NewBoundedQueue(cap int) *BoundedQueue {
	q := &BoundedQueue{
		tasks: make([]any, 0, cap),
		cap:   cap,
	}

	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)

	return q
}

func (q *BoundedQueue) Put(task any) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.tasks) == q.cap && !q.shutdown {
		q.notFull.Wait()
	}

	if q.shutdown {
		return
	}

	q.tasks = append(q.tasks, task)
	q.notEmpty.Signal()
}

func (q *BoundedQueue) Get() any {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.tasks) == 0 && !q.shutdown {
		q.notEmpty.Wait()
	}

	if len(q.tasks) == 0 {
		return nil
	}

	task := q.tasks[0]
	q.tasks[0] = nil
	q.tasks = q.tasks[1:]

	q.notFull.Signal()
	return task
}

func (q *BoundedQueue) Shutdown() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.shutdown {
		return
	}

	q.shutdown = true

	q.notEmpty.Broadcast()
	q.notFull.Broadcast()
}

func main() {
	q := NewBoundedQueue(2)

	var wg sync.WaitGroup

	for i := range 2 {
		wg.Go(func() {
			defer wg.Done()

			for {
				task := q.Get()
				if task == nil {
					fmt.Println("consumer", i, "stopping")
					return
				}

				fmt.Println("consumer", i, "got", task)
			}
		})
	}

	for i := range 5 {
		q.Put(i)
	}

	q.Shutdown()

	wg.Wait()
}
