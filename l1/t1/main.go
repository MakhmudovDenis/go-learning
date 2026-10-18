package main

import (
	"fmt"
	"math"
)

func main() {
	c := Circle{r: 4}
	r := Rectangle{a: 2, b: 12}

	dump(c)
	dump(r)
}

func dump(s Shape) {
	fmt.Println(s)
	fmt.Println(s.Area())
	fmt.Println(s.Perimeter())
}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	r float64
}

func (c Circle) Area() float64 {
	return c.r * c.r * math.Pi
}

func (c Circle) Perimeter() float64 {
	return c.r * 2 * math.Pi
}

type Rectangle struct {
	a, b float64
}

func (r Rectangle) Area() float64 {
	return r.a * r.b
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.a + r.b)
}
