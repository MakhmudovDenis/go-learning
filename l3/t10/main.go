// Параллельная загрузка данных из нескольких источников
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Attachment struct {
	data int
}

type SessionData struct {
	SessionID int `json:"session_id"`
}

func FetchComments(wg *sync.WaitGroup) <-chan string {
	output := make(chan string)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(output)
		}()
		time.Sleep(1 * time.Second)
		output <- "comments"
	}()

	return output
}

func FetchSession(wg *sync.WaitGroup) <-chan string {
	output := make(chan string)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(output)
		}()
		time.Sleep(1 * time.Second)
		output <- fmt.Sprintf("{\"session_id\":%d}", rand.Int()%2)
	}()

	return output
}

func FetchUserData(wg *sync.WaitGroup, comments string) <-chan string {
	output := make(chan string)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(output)
		}()
		time.Sleep(1 * time.Second)
		output <- fmt.Sprintf("user data based on %s", comments)
	}()

	return output
}

func FetchAttachments(wg *sync.WaitGroup) <-chan Attachment {
	output := make(chan Attachment)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(output)
		}()
		time.Sleep(3 * time.Second)
		output <- Attachment{123}
	}()

	return output
}

func main() {
	var wg sync.WaitGroup
	commentsChan := FetchComments(&wg)
	sessionDataChan := FetchSession(&wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		commentsData := <-commentsChan
		fmt.Printf("comments data: %s\n", commentsData)

		userDataChan := FetchUserData(&wg, commentsData)
		wg.Add(1)
		go func() {
			defer wg.Done()
			userData := <-userDataChan
			fmt.Println(userData)
		}()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sessionDataString := <-sessionDataChan
		var sessionData SessionData

		err := json.Unmarshal([]byte(sessionDataString), &sessionData)
		if err != nil {
			fmt.Printf("json decoding error: %v", err)
		}

		fmt.Printf("session data: %v\n", sessionData)

		if sessionData.SessionID == 0 {
			return
		}

		attachmentsChan := FetchAttachments(&wg)
		wg.Add(1)

		go func() {
			defer wg.Done()
			attachment := <-attachmentsChan
			fmt.Printf("attachments: %v\n", attachment)
		}()
	}()

	wg.Wait()
}
