// Ленивая инициализация подключения к базе данных с использованием `sync.Once`
package main

import (
	"fmt"
	"sync"
	"time"
)

type Connection struct {
	id int
}

type Database struct {
	conn Connection
	sync.Once
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) GetConnection() *Connection {
	db.Do(func() {
		time.Sleep(2 * time.Second)
		db.conn = Connection{0}
	})
	return &db.conn
}

func main() {
	db := NewDatabase()

	var wg sync.WaitGroup

	for i := range 5 {
		wg.Go(func() {
			conn := db.GetConnection()
			fmt.Printf("goroutine %d: conn %d\n", i, conn.id)
		})
	}

	wg.Wait()
}
