package main

import (
	"fmt"
	"sort"
	"strings"
)

func WordFrequency(text string) map[string]int {
	freqMap := make(map[string]int)
	words := strings.FieldsSeq(text)

	for word := range words {
		freqMap[word]++
	}

	return freqMap
}

func PrintWordFrequency(freqMap map[string]int) {
	type wordFreq struct {
		word  string
		count int
	}

	var wordFrequencies []wordFreq
	for word, count := range freqMap {
		wordFrequencies = append(wordFrequencies, wordFreq{word, count})
	}

	sort.Slice(wordFrequencies, func(i, j int) bool {
		return wordFrequencies[i].count > wordFrequencies[j].count
	})

	for _, wordFrequency := range wordFrequencies {
		fmt.Printf("%s: %d\n", wordFrequency.word, wordFrequency.count)
	}
}

func main() {
	text1 := "golang is great and golang is fast"
	text2 := "he hee he hee heee he he he ha"

	fmt.Printf("text1: %v\n", text1)
	freqMap1 := WordFrequency(text1)
	PrintWordFrequency(freqMap1)
	fmt.Printf("text2: %v\n", text2)
	freqMap2 := WordFrequency(text2)
	PrintWordFrequency(freqMap2)
}

