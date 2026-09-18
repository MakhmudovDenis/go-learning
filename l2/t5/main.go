// Кастомные ошибки в Go
package main

type MyError struct{}

func (err MyError) Error() string {
	return ""
}

func SimpleError() error {
	return nil
}

func FormattedError(age int) error {
	return nil
}
