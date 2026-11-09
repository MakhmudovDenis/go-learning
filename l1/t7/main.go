package main

import (
	"fmt"
)

func FilterByValue(m map[int]string, allowedValues []string) map[int]string {
	allowedSet := make(map[string]struct{}, len(allowedValues))
	for _, v := range allowedValues {
		allowedSet[v] = struct{}{}
	}

	result := make(map[int]string)
	for k, v := range m {
		if _, ok := allowedSet[v]; ok {
			result[k] = v
		}
	}

	return result
}

func InvertMap(m map[string]int) (map[int]string, error) {
	inverted := make(map[int]string)

	for k, v := range m {
		if _, exists := inverted[v]; exists {
			return nil, fmt.Errorf(fmt.Sprintf("duplicate value %d found for keys %q and %q", v, inverted[v], k))
		}
		inverted[v] = k
	}

	return inverted, nil
}

func main() {
	m := map[int]string{
		1: "odin",
		2: "dva",
		3: "odin",
		4: "chetyre",
	}
	allowed := []string{"odin", "chetyre"}
	filtered := FilterByValue(m, allowed)
	fmt.Println("Filtered:", filtered)

	orig := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	inverted, err := InvertMap(orig)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Inverted:", inverted)
}
