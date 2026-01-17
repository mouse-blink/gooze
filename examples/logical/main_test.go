package main

import "testing"

func TestIsInRangeAndPositive(t *testing.T) {
	tests := []struct {
		value, min, max int
		expected        bool
	}{
		{5, 1, 10, true},
		{1, 1, 10, true},
		{10, 1, 10, true},
		{0, 1, 10, false},
		{-5, 1, 10, false},
		{15, 1, 10, false},
		{-1, -10, 10, false},
	}

	for _, tt := range tests {
		result := IsInRangeAndPositive(tt.value, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("IsInRangeAndPositive(%d, %d, %d) = %v, expected %v",
				tt.value, tt.min, tt.max, result, tt.expected)
		}
	}
}

func TestIsValidOrDefault(t *testing.T) {
	tests := []struct {
		input      string
		useDefault bool
		expected   bool
	}{
		{"hello", false, true},
		{"hello", true, true},
		{"", true, true},
		{"", false, false},
	}

	for _, tt := range tests {
		result := IsValidOrDefault(tt.input, tt.useDefault)
		if result != tt.expected {
			t.Errorf("IsValidOrDefault(%q, %v) = %v, expected %v",
				tt.input, tt.useDefault, result, tt.expected)
		}
	}
}

func TestComplexLogic(t *testing.T) {
	tests := []struct {
		a, b, c  bool
		expected bool
	}{
		{true, true, false, true},
		{true, false, true, true},
		{false, true, true, true},
		{true, true, true, true},
		{false, false, false, false},
		{true, false, false, false},
		{false, true, false, false},
		{false, false, true, false},
	}

	for _, tt := range tests {
		result := ComplexLogic(tt.a, tt.b, tt.c)
		if result != tt.expected {
			t.Errorf("ComplexLogic(%v, %v, %v) = %v, expected %v",
				tt.a, tt.b, tt.c, result, tt.expected)
		}
	}
}
