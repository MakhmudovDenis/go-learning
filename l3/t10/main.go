// Параллельная загрузка данных из нескольких источников
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Comment struct {
	ID     int
	UserID int
	Text   string
}

type User struct {
	ID   int
	Name string
}

type Session struct {
	ID string `json:"session_id"`
}

type Attachment struct {
	ID   int
	Name string
}

type DB struct{}

var (
	dbOnce sync.Once
	dbConn *DB
)

func getDB() *DB {
	dbOnce.Do(func() {
		fmt.Println("подключение к бд")
		time.Sleep(200 * time.Millisecond)
		dbConn = &DB{}
	})
	return dbConn
}

func (db *DB) FetchComments() ([]Comment, error) {
	time.Sleep(1 * time.Second)
	return []Comment{
		{ID: 1, UserID: 10, Text: "я первый"},
		{ID: 2, UserID: 20, Text: "круто"},
		{ID: 3, UserID: 10, Text: "+1"},
		{ID: 4, UserID: 30, Text: "ыыыыы"},
	}, nil
}

func (db *DB) FetchUser(id int) (User, error) {
	time.Sleep(1 * time.Second)
	return User{ID: id, Name: fmt.Sprintf("user-%d", id)}, nil
}

func (db *DB) FetchSession() (Session, error) {
	time.Sleep(1 * time.Second)

	raw := `{"session_id":""}`
	if rand.Intn(2) == 1 {
		raw = `{"session_id":"abc-123"}`
	}

	var s Session
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return Session{}, fmt.Errorf("декод сессии: %w", err)
	}
	return s, nil
}

func (db *DB) FetchAttachments(sessionID string) ([]Attachment, error) {
	time.Sleep(1 * time.Second)
	return []Attachment{{ID: 123, Name: "photo.png"}}, nil
}

func fetchUsers(comments []Comment) ([]User, error) {
	seen := make(map[int]bool)
	var ids []int
	for _, c := range comments {
		if !seen[c.UserID] {
			seen[c.UserID] = true
			ids = append(ids, c.UserID)
		}
	}

	users := make([]User, len(ids))
	errs := make([]error, len(ids))

	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Go(func() {
			users[i], errs[i] = getDB().FetchUser(id)
		})
	}
	wg.Wait()

	return users, errors.Join(errs...)
}

func main() {
	start := time.Now()

	var (
		comments    []Comment
		users       []User
		commentsErr error

		session     Session
		attachments []Attachment
		sessionErr  error
	)

	var wg sync.WaitGroup

	wg.Go(func() {
		var err error
		comments, err = getDB().FetchComments()
		if err != nil {
			commentsErr = fmt.Errorf("комменты: %w", err)
			return
		}
		users, commentsErr = fetchUsers(comments)
	})

	wg.Go(func() {
		var err error
		session, err = getDB().FetchSession()
		if err != nil {
			sessionErr = fmt.Errorf("сессия: %w", err)
			return
		}
		if session.ID == "" {
			return
		}
		attachments, sessionErr = getDB().FetchAttachments(session.ID)
	})

	wg.Wait()

	if err := errors.Join(commentsErr, sessionErr); err != nil {
		fmt.Println("ошибки:", err)
	}

	fmt.Printf("комментарии: %v\n", comments)
	fmt.Printf("юзеры:       %v\n", users)
	fmt.Printf("скессия:     %q\n", session.ID)
	fmt.Printf("вложения:    %v\n", attachments)

	fmt.Printf("время: %v\n", time.Since(start).Round(100*time.Millisecond))
}
