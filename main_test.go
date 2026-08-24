package main

import (
	"testing"
)

func TestSolve(t *testing.T) {
	if got, want := solve([]int{10, -2, 30}), 38; got != want {
		t.Fatalf("solve() = %d, want %d", got, want)
	}
}
