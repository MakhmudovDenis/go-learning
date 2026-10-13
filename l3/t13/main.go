// Пул подключений к БД с использованием `sync.Cond`
package main

import (
	"fmt"
	"sync"
	"time"
)

type Connection struct {
	ID     int
	isBusy bool
}

type ConnectionPool struct {
	mu           sync.Mutex
	hasFreeConns sync.Cond
	hasBusyConns sync.Cond

	conns     []*Connection
	busyCount int
}

func NewConnectionPool(cap int) *ConnectionPool {
	p := &ConnectionPool{
		conns: make([]*Connection, cap),
	}

	p.hasFreeConns = *sync.NewCond(&p.mu)
	p.hasBusyConns = *sync.NewCond(&p.mu)

	for i := range cap {
		p.conns[i] = &Connection{i, false}
	}

	return p
}

func (c *ConnectionPool) Get() *Connection {
	c.mu.Lock()
	defer c.mu.Unlock()

	for c.busyCount == len(c.conns) {
		c.hasFreeConns.Wait()
	}

	c.hasBusyConns.Signal()

	for _, conn := range c.conns {
		if !conn.isBusy {
			conn.isBusy = true
			c.busyCount++
			return conn
		}
	}

	return nil
}

func (c *ConnectionPool) Release(conn *Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for c.busyCount == 0 {
		c.hasBusyConns.Wait()
	}

	c.hasFreeConns.Signal()

	conn.isBusy = false
	c.busyCount--
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
