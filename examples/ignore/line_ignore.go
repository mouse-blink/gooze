package main

func lineIgnored(a int) int {
	_ = a + 1 //gooze:ignore arithmetic
	_ = a + 1
	return a
}
