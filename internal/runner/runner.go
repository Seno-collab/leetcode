package runner

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	InputFile  string
	OutputFile string
	Numbers    []string
}

type Options struct {
	InputFile  string
	OutputFile string
}

// DefaultOptions keeps generated files beside the main.go that calls it.
func DefaultOptions() Options {
	directory := "."
	if _, sourceFile, _, ok := runtime.Caller(1); ok {
		directory = filepath.Dir(sourceFile)
	}

	return Options{
		InputFile:  filepath.Join(directory, "input.txt"),
		OutputFile: filepath.Join(directory, "results.jsonl"),
	}
}

type result[T any] struct {
	Input  []int `json:"input"`
	Output T     `json:"output"`
}

// Run reads every test case, applies solve, prints the result, and writes one
// JSON object per line to the output file.
func Run[T any](args []string, options Options, solve func([]int) T, stdout io.Writer) error {
	config, err := parseCommandLine(args, options)
	if err != nil {
		return err
	}

	var cases [][]int
	if len(config.Numbers) > 0 {
		numbers, err := parseNumbers(config.Numbers)
		if err != nil {
			return err
		}
		cases = [][]int{numbers}
	} else {
		cases, err = readInputCases(config.InputFile)
		if err != nil {
			return err
		}
	}

	results := make([]result[T], 0, len(cases))
	for _, numbers := range cases {
		record := result[T]{
			Input:  numbers,
			Output: solve(numbers),
		}
		results = append(results, record)
		fmt.Fprintf(stdout, "Input:  %v\nOutput: %v\n", record.Input, record.Output)
	}

	return writeResults(config.OutputFile, results)
}

func parseCommandLine(args []string, options Options) (Config, error) {
	config := Config{
		InputFile:  options.InputFile,
		OutputFile: options.OutputFile,
	}

	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "-in", "--input":
			if index+1 >= len(args) {
				return Config{}, fmt.Errorf("%s requires a file path", args[index])
			}
			index++
			config.InputFile = args[index]
		case "-out", "--output":
			if index+1 >= len(args) {
				return Config{}, fmt.Errorf("%s requires a file path", args[index])
			}
			index++
			config.OutputFile = args[index]
		default:
			config.Numbers = append(config.Numbers, args[index])
		}
	}

	if config.InputFile == "" {
		return Config{}, errors.New("input file path is required")
	}
	if config.OutputFile == "" {
		return Config{}, errors.New("output file path is required")
	}

	return config, nil
}

func parseNumbers(values []string) ([]int, error) {
	if len(values) == 0 {
		return nil, errors.New("provide at least one number")
	}

	numbers := make([]int, 0, len(values))
	for _, value := range values {
		number, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%q is not a valid integer", value)
		}
		numbers = append(numbers, number)
	}

	return numbers, nil
}

func readInputCases(path string) ([][]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cases [][]int
	scanner := bufio.NewScanner(file)
	const maxTestCaseSize = 10 * 1024 * 1024
	scanner.Buffer(make([]byte, 64*1024), maxTestCaseSize)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}

		numbers, err := parseNumbers(fields)
		if err != nil {
			return nil, fmt.Errorf("%s line %d: %w", path, lineNumber, err)
		}
		cases = append(cases, numbers)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("%s does not contain any test cases", path)
	}

	return cases, nil
}

func writeResults[T any](path string, results []result[T]) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, record := range results {
		if err := encoder.Encode(record); err != nil {
			return err
		}
	}

	return nil
}
