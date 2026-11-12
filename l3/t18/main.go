// Оптимизация HTTP-обработчика с sync.Pool
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type RequestData struct {
	UserID int               `json:"user_id"`
	Tags   []string          `json:"tags"`
	Meta   map[string]string `json:"meta"`
}

const maxBodyBytes = 1 << 20 // 1 MiB

var reqPool = sync.Pool{
	New: func() any {
		return &RequestData{
			Tags: make([]string, 0, 16),
			Meta: make(map[string]string, 8),
		}
	},
}

func (rd *RequestData) Reset() {
	rd.UserID = 0
	clear(rd.Tags)
	rd.Tags = rd.Tags[:0]
	clear(rd.Meta)
}

func handleRequest(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	data := reqPool.Get().(*RequestData)
	defer func() {
		data.Reset()
		reqPool.Put(data)
	}()

	req.Body = http.MaxBytesReader(rw, req.Body, maxBodyBytes)

	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			http.Error(rw, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(rw, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		"user_id":   data.UserID,
		"tag_count": len(data.Tags),
		"meta_keys": len(data.Meta),
	}
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		log.Printf("write response: %v", err)
	}
}

func main() {
	http.HandleFunc("/t18", handleRequest)
	fmt.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
