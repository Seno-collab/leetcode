package main

import (
	"fmt"
	"os"

	"leetcode/internal/runner"
)

func main() {
	err := runner.RunInts(
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
func solve(heights []int) int {
	left, right := 0, len(heights)-1
	ans := 0
	for left < right {
		max := MaxArea(heights[left], heights[right]) * (right - left)
		if max > ans {
			ans = max
		}
		left++
	}
	return ans
}

func MaxArea(a, b int) int {
	if a > b {
		return a
	}
	return b
}
