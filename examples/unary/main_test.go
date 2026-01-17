package main

import "testing"

func TestNegate(t *testing.T) {
	if Negate(5) != -5 {
		t.Error("Expected -5")
	}
	if Negate(-5) != 5 {
		t.Error("Expected 5")
	}
}

func TestIsNotTrue(t *testing.T) {
	if IsNotTrue(true) != false {
		t.Error("Expected false")
	}
	if IsNotTrue(false) != true {
		t.Error("Expected true")
	}
}

func TestBitwiseNot(t *testing.T) {
	if BitwiseNot(0) != -1 {
		t.Error("Expected -1")
	}
}

func TestPositive(t *testing.T) {
	if Positive(5) != 5 {
		t.Error("Expected 5")
	}
	if Positive(-5) != -5 {
		t.Error("Expected -5")
	}
}
