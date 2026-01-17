package main

import "fmt"

// IsInRangeAndPositive checks if a value is in range AND positive
func IsInRangeAndPositive(value, min, max int) bool {
	return value >= min && value <= max && value > 0
}

// IsValidOrDefault checks if input is valid OR uses default
func IsValidOrDefault(input string, useDefault bool) bool {
	return input != "" || useDefault
}

// ComplexLogic demonstrates nested logical operators
func ComplexLogic(a, b, c bool) bool {
	return (a && b) || (b && c) || (a && c)
}

func main() {
	fmt.Println("IsInRangeAndPositive(5, 1, 10):", IsInRangeAndPositive(5, 1, 10))
	fmt.Println("IsInRangeAndPositive(-5, 1, 10):", IsInRangeAndPositive(-5, 1, 10))
	fmt.Println("IsValidOrDefault(\"\", true):", IsValidOrDefault("", true))
	fmt.Println("IsValidOrDefault(\"\", false):", IsValidOrDefault("", false))
	fmt.Println("ComplexLogic(true, true, false):", ComplexLogic(true, true, false))
}
