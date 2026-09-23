package main

import (
	"fmt"
	"os"
	"strings"

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
func solve(s string) bool {
	i, j := 0, len(s)-1
	s = strings.ToLower(s)
	fmt.Println(s)
	for i < j {
		if !checkCharacters(s[i]) { // nếu không phải thuộc dạng chữ số thì tăng đơn vị
			i++
			continue
		}
		if !checkCharacters(s[j]) { // nếu không phải thuộc dạng chữ số thì giảm xuống 1 đơn vị
			j--
			continue
		}
		if s[i] != s[j] {
			return false
		} else {
			i++
			j--
		}
	}
	return true
}

func checkCharacters(v uint8) bool {
	fmt.Println(v)
	if v >= 65 && v <= 90 {
		return true
	}
	if v >= 97 && v <= 122 {
		return true
	}
	if v >= 48 && v <= 57 {
		return true
	}
	return false
}
