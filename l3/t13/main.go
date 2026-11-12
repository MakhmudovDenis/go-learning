// Пул подключений к БД с использованием `sync.Cond`
package main

import (
	"fmt"
	"sync"
	"time"
)

type Connection struct {
	ID int
}

type ConnectionPool struct {
	mu   sync.Mutex
	cond *sync.Cond

	free  []*Connection
	inUse map[*Connection]struct{}
}

func NewConnectionPool(size int) *ConnectionPool {
	p := &ConnectionPool{
		free:  make([]*Connection, 0, size),
		inUse: make(map[*Connection]struct{}, size),
	}
	p.cond = sync.NewCond(&p.mu)

	for i := range size {
		p.free = append(p.free, &Connection{ID: i})
	}

	return p
}

func (p *ConnectionPool) Get() *Connection {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.free) == 0 {
		p.cond.Wait()
	}

	last := len(p.free) - 1
	conn := p.free[last]
	p.free[last] = nil
	p.free = p.free[:last]

	p.inUse[conn] = struct{}{}
	return conn
}

func (p *ConnectionPool) Release(conn *Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.inUse[conn]; !ok {
		panic(fmt.Sprintf("pool: releasing connection %v that is not in use", conn))
	}

	delete(p.inUse, conn)
	p.free = append(p.free, conn)
	p.cond.Signal()
}

func main() {
	pool := NewConnectionPool(3) // Пул на 3 подключения

	for i := range 10 {
		go func(id int) {
			conn := pool.Get()
			defer pool.Release(conn)

			fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)
			time.Sleep(2 * time.Second) // Имитация работы
		}(i)
	}

	time.Sleep(10 * time.Second)
}
