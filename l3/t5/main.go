// fan-out
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const dataDir string = "data"

func getWordCount(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("ошибка при чтении файла: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("ошибка при подсчете: %w", err)
	}

	return count, nil
}

func getFiles() (chan string, int, error) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при чтении директории: %w", err)
	}

	output := make(chan string)

	go func() {
		for _, file := range files {
			output <- fmt.Sprintf("%v/%v", dataDir, file.Name())
		}
		close(output)
	}()

	return output, len(files), nil
}

func fanOut(inputCh chan string, n int) []chan string {
	output := make([]chan string, n)
	for i := range n {
		output[i] = make(chan string)
	}

	idx := 0
	go func() {
		for v := range inputCh {
			output[idx] <- v
			idx = (idx + 1) % n
		}

		for _, ch := range output {
			close(ch)
		}
	}()

	return output
}

func worker(wg *sync.WaitGroup, ch chan string, res *sync.Map) {
	defer wg.Done()
	filename := <-ch
	wordCount, err := getWordCount(filename)
	if err != nil {
		log.Fatalf("ошибка при подсчёте слов в %v:%v", filename, err)
		return
	}

	res.Store(filename, wordCount)
}

func printSyncMap(m *sync.Map) {
	out := make(map[string]int)
	m.Range(func(key, val any) bool {
		_, filename, found := strings.Cut(key.(string), "/")
		if !found {
			return false
		}
		out[filename] = val.(int)
		return true
	})
	fmt.Println(out)
}

func main() {
	var wg sync.WaitGroup

	mainChan, chanCount, err := getFiles()
	if err != nil {
		log.Fatal(err)
		return
	}
	splitChans := fanOut(mainChan, chanCount)

	var res sync.Map
	wg.Add(chanCount + 1)

	go func() {
		defer wg.Done()
		for _, ch := range splitChans {
			go worker(&wg, ch, &res)
		}
	}()

	wg.Wait()
	printSyncMap(&res)
}
