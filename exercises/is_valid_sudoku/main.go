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

// Replace this function with the solution for the current exercise.
func solve(board [][]byte) bool {
	rows := [9]map[byte]bool{}
	cols := [9]map[byte]bool{}
	boxed := [9]map[byte]bool{}
	for i := range 9 {
		rows[i] = make(map[byte]bool, 9)
		cols[i] = make(map[byte]bool, 9)
		boxed[i] = make(map[byte]bool, 9)
	}
	for r := range 9 {
		for c := range 9 {
			value := board[r][c]
			if value == '.' {
				continue
			}
			if cols[c][value] {
				return false
			}
			if rows[r][value] {
				return false
			}
			boxIndex := (r/3)*3 + c/3
			if boxed[boxIndex][value] {
				return false
			}
			cols[c][value] = true
			rows[r][value] = true
			boxed[boxIndex][value] = true
		}
	}
	return true
}
