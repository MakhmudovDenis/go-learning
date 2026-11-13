// fan-out
package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
)

const dataDir string = "data"

type result struct {
	file  string
	words int
	err   error
}

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

func worker(files <-chan string, results chan<- result) {
	for f := range files {
		n, err := getWordCount(f)
		results <- result{f, n, err}
	}
}

func getFiles() (<-chan string, error) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении директории: %w", err)
	}

	output := make(chan string)

	go func() {
		defer close(output)

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			output <- fmt.Sprintf("%s/%s", dataDir, file.Name())
		}
	}()

	return output, nil
}

func main() {
	files, err := getFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	results := make(chan result)

	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Go(func() { worker(files, results) })
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	wordCount := 0
	for r := range results {
		if r.err != nil {
			fmt.Printf("%s: %v", r.file, r.err)
			continue
		}
		fmt.Printf("%s: %d\n", r.file, r.words)
		wordCount += r.words
	}
	fmt.Println("всего слов:", wordCount)
}
