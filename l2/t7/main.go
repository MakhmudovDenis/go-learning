// Универсальный стек (LIFO) на Go с дженериками
package main

import (
	"fmt"
)

type Stack[T any] struct {
	elements []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

func (s *Stack[T]) Dump() {
	fmt.Println(s.elements, len(s.elements), "items")
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack[T]) Pop() (retVal T, ok bool) {
	if s.IsEmpty() {
		return
	}
	stackSize := len(s.elements)

	retVal = s.elements[stackSize-1]
	ok = true
	s.elements = s.elements[:stackSize-1]

	return
}

func (s *Stack[T]) Peek() (retVal T, ok bool) {
	if s.IsEmpty() {
		return
	}

	return s.elements[len(s.elements)-1], true
}

func main() {
	s := NewStack[int]()
	fmt.Println(s.IsEmpty())
	s.Push(12)
	fmt.Println(s.IsEmpty())
	s.Push(4)
	s.Dump()
	fmt.Println(s.Pop())
	s.Dump()
	fmt.Println(s.Pop())
	s.Dump()
	fmt.Println(s.Pop())
}
