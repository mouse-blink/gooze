package main

import "testing"

func TestCalculate(t *testing.T) {
	if Calculate(5) != 10 {
		t.Error("Calculate(5) should be 10")
	}
}
