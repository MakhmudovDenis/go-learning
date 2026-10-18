package main

import (
	"testing"
)

func TestCarEngine(t *testing.T) {
	c := &Car{Brand: "Test"}
	if c.GetEngineStatus() {
		t.Error("двигатель должен быть изначально выключен")
	}
	if err := c.StartEngine(); err != nil {
		t.Errorf("ошибка StartEngine: %v", err)
	}
	if !c.GetEngineStatus() {
		t.Error("двигатель должен быт выключен после StartEngine")
	}
}

func TestPolymorphism(t *testing.T) {
	var v Vehicle = NewTruck("Test", 10.0)
	if v.GetInfo() == "" {
		t.Error("GetInfo должна вернуть непустую строку")
	}
}

func TestElectricCarLowBattery(t *testing.T) {
	ev := NewElectricCar("Test", 3)
	if err := ev.StartEngine(); err != ErrLowBattery {
		t.Errorf("ОР: ErrLowBattery, ФР: %v", err)
	}
}

func TestHonk(t *testing.T) {
	const expectedCarHonk = "Beep beep!"
	const expectedTruckHonk = "Honk Honk!"

	car := &Car{Brand: "X"}
	truck := NewTruck("Y", 5.0)

	actualCarHonk := car.Honk()
	actualTruckHonk := truck.Honk()

	if actualCarHonk != expectedCarHonk {
		t.Errorf("ОР: %v, ФР: %v", expectedCarHonk, actualCarHonk)
	}

	if actualTruckHonk != expectedTruckHonk {
		t.Errorf("ОР: %v, ФР: %v", expectedTruckHonk, actualTruckHonk)
	}
}
