package main

type Stacker interface {
	Push(v int)
	Pop() int
}

type stack struct {
	items []int
}

func (s *stack) Push(v int) {
	s.items = append(s.items, v)
}

func (s *stack) Pop() int {
	if len(s.items) == 0 {
		panic("pop from empty stack")
	}

	last := len(s.items) - 1
	v := s.items[last]
	s.items = s.items[:last]

	return v
}

func New() *stack {
	return &stack{}
}
