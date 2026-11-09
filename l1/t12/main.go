package main

import "fmt"

// RemoveUnordered удаляет элемент по индексу без сохранения порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveUnordered[T any](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		return s
	}

	s[i] = s[len(s)-1]

	var zero T
	s[len(s)-1] = zero

	return s[:len(s)-1]
}

// RemoveOrdered удаляет элемент по индексу с сохранением порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveOrdered[T any](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		return s
	}

	copy(s[i:], s[i+1:])

	var zero T
	s[len(s)-1] = zero

	return s[:len(s)-1]
}

// RemoveAllByValue удаляет все вхождения указанного значения.
func RemoveAllByValue[T comparable](s []T, value T) []T {
	result := s[:0]

	for _, v := range s {
		if v != value {
			result = append(result, v)
		}
	}

	var zero T
	for i := len(result); i < len(s); i++ {
		s[i] = zero
	}

	return result
}

// RemoveDuplicates оставляет только уникальные элементы (сохраняет порядок).
func RemoveDuplicates[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	result := s[:0]

	for _, v := range s {
		if _, exists := seen[v]; exists {
			continue
		}

		seen[v] = struct{}{}
		result = append(result, v)
	}

	var zero T
	for i := len(result); i < len(s); i++ {
		s[i] = zero
	}

	return result
}

// RemoveIf удаляет элементы, удовлетворяющие условию predicate.
func RemoveIf[T any](s []T, predicate func(T) bool) []T {
	result := s[:0]

	for _, v := range s {
		if !predicate(v) {
			result = append(result, v)
		}
	}

	var zero T
	for i := len(result); i < len(s); i++ {
		s[i] = zero
	}

	return result
}

// RemoveOrderedWithNil удаляет элемент по индексу (для слайса указателей),
// обнуляя удаляемый элемент для предотвращения утечек памяти.
func RemoveOrderedWithNil[T any](s []*T, i int) []*T {
	if i < 0 || i >= len(s) {
		return s
	}

	copy(s[i:], s[i+1:])
	s[len(s)-1] = nil

	return s[:len(s)-1]
}

// ShrinkCapacity сокращает вместимость слайса, если она превышает
// удвоенную длину после удаления элементов.
func ShrinkCapacity[T any](s []T) []T {
	if cap(s) <= 2*len(s) {
		return s
	}

	result := make([]T, len(s))
	copy(result, s)

	return result
}

func main() {
	numbers := []int{10, 20, 30, 40, 50}

	fmt.Println("RemoveUnordered:", RemoveUnordered(numbers, 1))

	numbers = []int{10, 20, 30, 40, 50}
	fmt.Println("RemoveOrdered:", RemoveOrdered(numbers, 1))

	numbers = []int{1, 2, 2, 3, 2, 4}
	fmt.Println("RemoveAllByValue:", RemoveAllByValue(numbers, 2))

	numbers = []int{1, 2, 2, 3, 1, 4, 3}
	fmt.Println("RemoveDuplicates:", RemoveDuplicates(numbers))

	numbers = []int{1, 2, 3, 4, 5, 6}
	fmt.Println("RemoveIf:", RemoveIf(numbers, func(n int) bool {
		return n%2 == 0
	}))

	a, b, c := 10, 20, 30
	pointers := []*int{&a, &b, &c}
	fmt.Println("RemoveOrderedWithNil:", RemoveOrderedWithNil(pointers, 1))

	numbers = make([]int, 3, 20)
	numbers[0] = 10
	numbers[1] = 20
	numbers[2] = 30
	fmt.Println("Before shrink: len=", len(numbers), "; cap=", cap(numbers))
	numbers = ShrinkCapacity(numbers)
	fmt.Println("After shrink: len=", len(numbers), "; cap=", cap(numbers))
}
