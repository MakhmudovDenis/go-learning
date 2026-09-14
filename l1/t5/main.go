package main

import "fmt"

var people map[string]int

func init() {
	people = make(map[string]int)

	AddPerson("Иван1", 10)
	AddPerson("Иван2", 20)
	AddPerson("Иван3", 30)
}

func AddPerson(name string, age int) {
	people[name] = age
}

func GetAge(name string) int {
	age, ok := people[name]
	if !ok {
		return 0
	}
	return age
}

func DeletePerson(name string) {
	delete(people, name)
}

func PrintAll() {
	for name, age := range people {
		fmt.Printf("%s: %d\n", name, age)
	}
}

func main() {
	fmt.Println("Инит:")
	PrintAll()
	fmt.Println()

	AddPerson("Иван4", 40)
	fmt.Println("После добавления четвертого Ивана")
	PrintAll()
	fmt.Println()

	fmt.Printf("Возраст первого Ивана: %d\n", GetAge("Иван1"))
	fmt.Printf("Возраст Андрея: %d\n", GetAge("Андрей"))
	fmt.Println()

	DeletePerson("Иван2")
	DeletePerson("Иван3")
	fmt.Println("После удаления Иванов 2 и 3:")
	PrintAll()
}
