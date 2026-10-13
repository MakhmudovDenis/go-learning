// Оптимизация HTTP-обработчика с sync.Pool
package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type RequestData struct {
	Headers map[string][]string
	Body    []byte
}

var reqPool = sync.Pool{
	New: func() any {
		return &RequestData{
			Headers: make(map[string][]string),
			Body:    make([]byte, 0, 4096),
		}
	},
}

func handleRequest(rw http.ResponseWriter, req *http.Request) {
	data := reqPool.Get().(*RequestData)
	defer func() {
		req.Body.Close()
		data.Reset()
		reqPool.Put(data)
	}()

	if req.ContentLength > int64(cap(data.Body)) {
		http.Error(rw, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	buf := data.Body[:cap(data.Body)]
	n, err := req.Body.Read(buf)
	if err != io.EOF {
		http.Error(rw, "request reading failed", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	data.Body = buf[:n]

	for key, value := range req.Header {
		data.Headers[key] = append(data.Headers[key], value...)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("success"))
}

func (rd *RequestData) Reset() {
	for k := range rd.Headers {
		delete(rd.Headers, k)
	}

	rd.Body = rd.Body[:0]
}

func main() {
	http.HandleFunc("/t18", handleRequest)
	fmt.Println("server started at :8080")
	http.ListenAndServe(":8080", nil)
}
