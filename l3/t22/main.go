// Исправление синхронизации горутин
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type logic struct{}

var Logic logic

func (l *logic) UpdateDB(ctx context.Context, item *Item) error {
	return nil // Заглушка
}

func (l *logic) FetchItems(ctx context.Context) ([]*Item, error) {
	return []*Item{
		{Value: 5},
		{Value: 15},
		{Value: 7},
	}, nil // Заглушка
}

type Item struct {
	Value int
}

func processItem(wg *sync.WaitGroup, item *Item) {
	defer wg.Done()
	time.Sleep(time.Second)
	if item.Value > 10 {
		fmt.Printf("ERROR: item %d can't be more than 10\n", item.Value)
		return
	}

	err := Logic.UpdateDB(context.Background(), item)
	if err != nil {
		fmt.Println("ERROR: can't process item")
	}
}

func DoBusinessLogic(wg *sync.WaitGroup) error {
	items, err := Logic.FetchItems(context.Background())
	if err != nil {
		return err
	}
	wg.Add(len(items))

	for _, item := range items {
		go processItem(wg, item)
	}

	return nil
}

func main() {
	var wg sync.WaitGroup
	err := DoBusinessLogic(&wg)
	if err != nil {
		fmt.Println("Error:", err)
	}

	wg.Wait()
	fmt.Println("All items processed")
}
