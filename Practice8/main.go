package main

import "fmt"

func main() {
	fmt.Println("Practice 8 - Go Testing")
	result, err := Divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("10 / 2 = %d\n", result)
	fmt.Printf("10 - 3 = %d\n", Subtract(10, 3))
	fmt.Printf("2 + 3 = %d\n", Add(2, 3))
}
