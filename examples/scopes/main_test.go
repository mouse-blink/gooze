package main

import "testing"

func TestCalculate(t *testing.T) {
	// Test addition branch
	result := Calculate(5, 3)
	if result != 8 {
		t.Errorf("Calculate(5, 3) = %d, want 8", result)
	}

	// Test subtraction branch
	result = Calculate(3, 5)
	if result != -2 {
		t.Errorf("Calculate(3, 5) = %d, want -2", result)
	}
}

func TestInit(t *testing.T) {
	// Verify init ran
	if counter != 10 {
		t.Errorf("counter = %d, want 10 (init should set this)", counter)
	}
}
