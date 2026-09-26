// Кастомные ошибки в Go
package main

import (
	"errors"
	"fmt"
)

type MyError struct {
	Code int
	Msg  string
}

func (err MyError) Error() string {
	return err.Msg
}

func SimpleError() error {
	return errors.New("простая ошибка")
}

func FormattedError(age int) error {
	return fmt.Errorf("%w: возраст %d недопустим", errors.New("ошибка"), age)
}

func StructError() error {
	return MyError{Code: 404, Msg: "не найдено"}
}

func main() {
	var err error = SimpleError()

	fmt.Println(err)
	fmt.Println(FormattedError(12))
	fmt.Println(StructError())
}
