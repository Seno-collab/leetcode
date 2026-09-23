package main

import (
	"fmt"
	"os"
	"sort"

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

// Replace this function with the solution for the current exercise.
func solve(numbers []int) [][]int {
	sort.Ints(numbers)
	ans := make([][]int, 0)
	for left := 0; left < len(numbers)-2; left++ {

		//skip duplicate
		if left > 0 && numbers[left] == numbers[left-1] {
			continue
		}
		middle := left + 1
		right := len(numbers) - 1
		for middle < right {
			sum := numbers[right] + numbers[left] + numbers[middle]
			if sum == 0 {
				ans = append(ans, []int{numbers[left], numbers[middle], numbers[right]})
				for middle < right && numbers[middle] == numbers[middle+1] {
					middle++
				}
				for middle < right && numbers[right] == numbers[right-1] {
					right--
				}
				middle++
				right--
			} else if sum < 0 {
				middle++
			} else {
				right++
			}
		}
	}
	return ans
}
