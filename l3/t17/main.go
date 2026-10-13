// Оптимизация обработки строк с sync.Pool
package main

import (
	"fmt"
	"strings"
	"sync"
)

const ULDelta byte = 32 // осознанно поддерживаем только юникод

var bufferPool = sync.Pool{
	New: func() any {
		return make([]byte, 0, 64)
	},
}

func ProcessString(s string) string {
	buf := bufferPool.Get().([]byte)
	buf = buf[:0]

	for i := range len(s) {
		char := s[i]

		if char >= 'a' && char <= 'z' {
			char -= ULDelta
		}

		buf = append(buf, char)
	}

	res := string(buf)
	bufferPool.Put(buf)

	return res
}

var s string = strings.ToUpper("asdsa")

func main() {
	examples := []string{
		"hello, world!",
		"gopher",
		"lorem ipsum dolor sit amet",
	}
	fmt.Println('C')

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}
}
