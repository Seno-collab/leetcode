package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"leetcode/internal/runner"
)

func main() {
	solution := &Solution{}

	err := runner.RunJSON(
		os.Args[1:],
		runner.DefaultOptions(),
		func(strs []string) encodeDecodeResult {
			encoded := solution.Encode(strs)

			return encodeDecodeResult{
				Encoded: encoded,
				Decoded: solution.Decode(encoded),
			}
		},
		os.Stdout,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type encodeDecodeResult struct {
	Encoded string   `json:"encoded"`
	Decoded []string `json:"decoded"`
}

type Solution struct{}

// Encode prefixes every string with its byte length and a separator. The
// length makes separator characters inside the original string unambiguous.
func (s *Solution) Encode(strs []string) string {
	var encoded strings.Builder

	for _, str := range strs {
		encoded.WriteString(strconv.Itoa(len(str)))
		encoded.WriteByte('#')
		encoded.WriteString(str)
	}

	return encoded.String()
}

func (s *Solution) Decode(encoded string) []string {
	strs := make([]string, 0)

	for index := 0; index < len(encoded); {
		separatorOffset := strings.IndexByte(encoded[index:], '#')
		if separatorOffset < 0 {
			return []string{}
		}

		separatorIndex := index + separatorOffset
		length, err := strconv.Atoi(encoded[index:separatorIndex])
		if err != nil || length < 0 {
			return []string{}
		}

		start := separatorIndex + 1
		end := start + length
		if end < start || end > len(encoded) {
			return []string{}
		}

		strs = append(strs, encoded[start:end])
		index = end
	}

	return strs
}
