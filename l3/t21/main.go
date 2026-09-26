package main

import (
	"fmt"
	"net/http"
	"time"
)

func fetchUrl(url string) error {
	_, err := http.Get(url)
	return err
}

func main() {
	urls := []string{
		"https://www.lamoda.ru",
		"https://www.yandex.ru",
		"https://www.mail.ru",
		"https://www.google.ru",
	}
	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s....\n", url)
			err := fetchUrl(url)
			if err != nil {
				fmt.Printf("Error feaching %s: %v\n", url, err)
				return
			}
			fmt.Printf("Fetched %s\n", url)
		}(url)
	}
	fmt.Println("All request launched!")
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Program finished")
}
