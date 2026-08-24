package main

import (
	"fmt"
	"os"

	"leetcode/internal/runner"
)

func main() {
	err := runner.Run(
		os.Args[1:],
		runner.DefaultOptions(),
		solve,
		os.Stdout,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// Replace this function with the solution for the current exercise.
func solve(numbers []int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}
