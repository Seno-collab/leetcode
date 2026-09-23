package main

import (
	"fmt"
	"os"

	"leetcode/internal/runner"
)

func main() {
	err := runner.RunJSON(
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

type TwoIntegerSumIIInput struct {
	Numbers []int `json:"numbers"`
	Target  int   `json:"target"`
}

// Replace this function with the solution for the current exercise.
func solve(input TwoIntegerSumIIInput) []int {
	return TwoIntegerSumII(input.Numbers, input.Target)
}

func TwoIntegerSumII(numbers []int, target int) []int {
	left := 0
	right := len(numbers) - 1
	for left < right {
		for left < right && numbers[left]+numbers[right] > target {
			right--
		}
		for left < right && numbers[left]+numbers[right] < target {
			left++
		}
		if numbers[right]+numbers[left] == target {
			return []int{left + 1, right + 1}
		}
	}
	return []int{-1, -1}
}
