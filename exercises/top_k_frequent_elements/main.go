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

type topKInput struct {
	Numbers []int `json:"numbers"`
	K       int   `json:"k"`
}

func solve(input topKInput) []int {
	return topKFrequent(input.Numbers, input.K)
}

func topKFrequent(numbers []int, k int) []int {
	if k <= 0 {
		return []int{}
	}

	freq := make(map[int]int, len(numbers))
	for _, number := range numbers {
		freq[number]++
	}

	uniques := make([]int, 0, len(freq))
	for number := range freq {
		uniques = append(uniques, number)
	}
	sort.Slice(uniques, func(i, j int) bool {
		left := uniques[i]
		right := uniques[j]

		if freq[left] == freq[right] {
			return left < right
		}

		return freq[left] > freq[right]
	})

	if k > len(uniques) {
		k = len(uniques)
	}

	return uniques[:k]
}
