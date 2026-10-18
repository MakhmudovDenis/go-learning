package main

import (
	"errors"
	"fmt"
)

var (
	ErrEngineAlreadyRunning = errors.New("двигатель уже работает")
	ErrEngineOff            = errors.New("двигатель не запущен")
	ErrLowBattery           = errors.New("низкий заряд батареи")
)

type Vehicle interface {
	StartEngine() error
	StopEngine() error
	GetInfo() string
}

type Car struct {
	Brand    string
	engineOn bool
}

func (c *Car) StartEngine() error {
	if c.engineOn {
		return ErrEngineAlreadyRunning
	}
	c.engineOn = true
	return nil
}

func (c *Car) StopEngine() error {
	if !c.engineOn {
		return ErrEngineOff
	}
	c.engineOn = false
	return nil
}

func (c *Car) GetInfo() string {
	status := "stopped"
	if c.engineOn {
		status = "running"
	}
	return fmt.Sprintf("%s Car, engine %s", c.Brand, status)
}

func (c *Car) GetEngineStatus() bool {
	return c.engineOn
}

func (c *Car) Honk() string {
	return "Beep beep!"
}

type Truck struct {
	Car
	cargoCapacity float64 // encapsulated
}

func NewTruck(brand string, cargoCapacity float64) *Truck {
	return &Truck{
		Car:           Car{Brand: brand},
		cargoCapacity: cargoCapacity,
	}
}

func (t *Truck) GetCargoCapacity() float64 {
	return t.cargoCapacity
}

func (t *Truck) Honk() string {
	return "Honk Honk!"
}

func (t *Truck) GetInfo() string {
	return fmt.Sprintf("%s, cargo capacity %.1f tons", t.Car.GetInfo(), t.cargoCapacity)
}

type ElectricCar struct {
	Car
	batteryLevel int // encapsulated
}

func NewElectricCar(brand string, batteryLevel int) *ElectricCar {
	return &ElectricCar{
		Car:          Car{Brand: brand},
		batteryLevel: batteryLevel,
	}
}

func (e *ElectricCar) GetBatteryLevel() int {
	return e.batteryLevel
}

func (e *ElectricCar) StartEngine() error {
	if e.batteryLevel <= 5 {
		return ErrLowBattery
	}
	return e.Car.StartEngine()
}

func (e *ElectricCar) GetInfo() string {
	return fmt.Sprintf("%s, battery %d%%", e.Car.GetInfo(), e.batteryLevel)
}

func main() {
	vehicles := []Vehicle{
		&Car{Brand: "Toyota"},
		NewTruck("Volvo", 20.5),
		NewElectricCar("Tesla", 80),
	}

	for _, v := range vehicles {
		fmt.Println(v.GetInfo())
		_ = v.StartEngine()
		fmt.Println("Engine started:", v.(interface{ GetEngineStatus() bool }).GetEngineStatus())
	}

	// Unique behaviors
	car := &Car{Brand: "Honda"}
	truck := NewTruck("Ford", 15.0)
	ev := NewElectricCar("Nissan", 3)

	fmt.Println(car.Honk())
	fmt.Println(truck.Honk())
	fmt.Println(ev.StartEngine())
}
