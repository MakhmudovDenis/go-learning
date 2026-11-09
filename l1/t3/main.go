package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrUnsupported = errors.New("обновление недоступно")

type Device interface {
	UpdateOs(string) error
	GetInfo() string
}

type Smartphone struct {
	OSVersion string
	Model     string
}

func (s *Smartphone) UpdateOS(version string) error {
	currentOSMajor, _, _ := strings.Cut(s.OSVersion, ".")

	if majorInt, _ := strconv.Atoi(currentOSMajor); majorInt >= 12 {
		return ErrUnsupported
	}

	s.OSVersion = version

	return nil
}

func (s *Smartphone) GetInfo() string {
	return fmt.Sprintf("Модель: %v, ОС: %v", s.OSVersion, s.Model)
}

type Laptop struct {
	OSVersion string
	Model     string
}


func (l *Laptop) UpdateOS(version string) error {
	if !strings.HasPrefix(version, "Windows") {
		return ErrUnsupported
	}

	l.OSVersion = version

	return nil
}

func (s *Laptop) GetInfo() string {
	return fmt.Sprintf("Модель: %v, ОС: %v", s.OSVersion, s.Model)
}

type Smartwatch struct {
	OSVersion string
	Model     string
}

func (s *Smartwatch) UpdateOS(version string) error {
	if len(version) < 5 {
		return ErrUnsupported
	}

	s.OSVersion = version

	return nil
}


func (s *Smartwatch()) GetInfo() string {
	return fmt.Sprintf("Модель: %v, ОС: %v", s.OSVersion, s.Model)
}

func main() {
	sm := Smartphone{"1.2.3", "doooood"}
	fmt.Println((&sm).GetInfo())
}
