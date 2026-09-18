// Анализ цепочек ошибок в Go error.Is
package t6

import (
	"errors"
	_ "fmt"
)

var (
	ErrNotFound  = errors.New("ресурс не найден")
	TimeoutError = errors.New("таймаут операции")
)

func SimulateRequest() error {
	return nil
}

func ProcessError(err error) {
}
