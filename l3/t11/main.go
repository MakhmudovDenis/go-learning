// Реализация очереди с ограниченной емкостью на sync.Cond
package main

import (
	"errors"
	"fmt"
	"sync"
)

var ErrQueueClosed = errors.New("queue is shut down")

type BoundedQueue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond

	// теперь ринг буфер
	tasks    []any
	head     int
	count    int
	shutdown bool
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{
		tasks: make([]any, capacity),
	}

	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)

	return q
}

func (q *BoundedQueue) Put(task any) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.count == len(q.tasks) && !q.shutdown {
		q.notFull.Wait()
	}

	if q.shutdown {
		return ErrQueueClosed
	}

	tail := (q.head + q.count) % len(q.tasks)
	q.tasks[tail] = task
	q.count++

	q.notEmpty.Signal()
	return nil
}

func (q *BoundedQueue) Get() (task any, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.count == 0 && !q.shutdown {
		q.notEmpty.Wait()
	}

	if q.count == 0 {
		return nil, false
	}

	task = q.tasks[q.head]
	q.tasks[q.head] = nil
	q.head = (q.head + 1) % len(q.tasks)
	q.count--

	q.notFull.Signal()
	return task, true
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
			for {
				task, ok := q.Get()
				if !ok {
					fmt.Println("consumer", i, "stopping")
					return
				}

				fmt.Println("consumer", i, "got", task)
			}
		})
	}

	for i := range 5 {
		if err := q.Put(i); err != nil {
			fmt.Println("put failed:", err)
		}
	}

	q.Shutdown()

	if err := q.Put(100); err != nil {
		fmt.Println("put after shutdown:", err)
	}

	wg.Wait()
}
