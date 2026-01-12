package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Capture output from main
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Sum of 3 and 5 is 8\n"
	if output != expected {
		t.Errorf("got %q, want %q", output, expected)
	}
}
