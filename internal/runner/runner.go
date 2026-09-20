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
	"time"
)

type Config struct {
	InputFile  string
	OutputFile string
	Arguments  []string
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

type Parser[I any] func(string) (I, error)

type result[I, O any] struct {
	Input      I     `json:"input"`
	Output     O     `json:"output"`
	DurationNS int64 `json:"duration_ns"`
}

type clock func() time.Time

// Run supports any input and output types. Each non-empty input line is parsed
// as one test case before solve is measured.
func Run[I, O any](
	args []string,
	options Options,
	parse Parser[I],
	solve func(I) O,
	stdout io.Writer,
) error {
	return run(args, options, parse, solve, stdout, time.Now)
}

// RunInts is the convenient runner for exercises whose input is only []int.
func RunInts[O any](args []string, options Options, solve func([]int) O, stdout io.Writer) error {
	return Run(args, options, parseInts, solve, stdout)
}

// RunJSON parses one JSON value per line. Structs, strings, matrices, slices,
// and nested combinations can all be used as exercise input.
func RunJSON[I, O any](args []string, options Options, solve func(I) O, stdout io.Writer) error {
	return Run(args, options, ParseJSON[I], solve, stdout)
}

func run[I, O any](
	args []string,
	options Options,
	parse Parser[I],
	solve func(I) O,
	stdout io.Writer,
	now clock,
) error {
	config, err := parseCommandLine(args, options)
	if err != nil {
		return err
	}

	cases, err := loadCases(config, parse)
	if err != nil {
		return err
	}

	results, totalDuration := executeCases(cases, solve, now)
	if err := writeResults(config.OutputFile, results); err != nil {
		return err
	}

	return printResults(stdout, results, totalDuration, config.OutputFile)
}

func loadCases[I any](config Config, parse Parser[I]) ([]I, error) {
	if len(config.Arguments) == 0 {
		return readInputCases(config.InputFile, parse)
	}

	input, err := parse(strings.Join(config.Arguments, " "))
	if err != nil {
		return nil, err
	}

	return []I{input}, nil
}

func executeCases[I, O any](cases []I, solve func(I) O, now clock) ([]result[I, O], time.Duration) {
	results := make([]result[I, O], 0, len(cases))
	var totalDuration time.Duration

	for _, input := range cases {
		startedAt := now()
		output := solve(input)
		duration := max(now().Sub(startedAt), 0)

		record := result[I, O]{
			Input:      input,
			Output:     output,
			DurationNS: duration.Nanoseconds(),
		}
		results = append(results, record)
		totalDuration += duration
	}

	return results, totalDuration
}

func printResults[I, O any](
	stdout io.Writer,
	results []result[I, O],
	totalDuration time.Duration,
	outputFile string,
) error {
	for index, record := range results {
		if _, err := fmt.Fprintf(
			stdout,
			"=== Case %d ===\nInput      : %s\nOutput     : %s\nSolve time : %s\n\n",
			index+1,
			formatValue(record.Input),
			formatValue(record.Output),
			time.Duration(record.DurationNS),
		); err != nil {
			return err
		}
	}

	averageDuration := time.Duration(0)
	if len(results) > 0 {
		averageDuration = totalDuration / time.Duration(len(results))
	}

	_, err := fmt.Fprintf(
		stdout,
		"=== Summary ===\nCases        : %d\nTotal time   : %s\nAverage time : %s\nResults file : %s\n",
		len(results),
		totalDuration,
		averageDuration,
		outputFile,
	)
	return err
}

func formatValue(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}

	return string(encoded)
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
			config.Arguments = append(config.Arguments, args[index])
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

func parseInts(raw string) ([]int, error) {
	values := strings.Fields(raw)
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

func ParseJSON[I any](raw string) (I, error) {
	var input I
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		return input, fmt.Errorf("invalid JSON input: %w", err)
	}

	var trailingValue any
	if err := decoder.Decode(&trailingValue); !errors.Is(err, io.EOF) {
		if err == nil {
			return input, errors.New("invalid JSON input: multiple values on one line")
		}
		return input, fmt.Errorf("invalid JSON input: %w", err)
	}

	return input, nil
}

func readInputCases[I any](path string, parse Parser[I]) ([]I, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cases []I
	scanner := bufio.NewScanner(file)
	const maxTestCaseSize = 10 * 1024 * 1024
	scanner.Buffer(make([]byte, 64*1024), maxTestCaseSize)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}

		input, err := parse(raw)
		if err != nil {
			return nil, fmt.Errorf("%s line %d: %w", path, lineNumber, err)
		}
		cases = append(cases, input)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("%s does not contain any test cases", path)
	}

	return cases, nil
}

func writeResults[I, O any](path string, results []result[I, O]) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	for _, record := range results {
		if err := encoder.Encode(record); err != nil {
			_ = file.Close()
			return err
		}
	}

	return file.Close()
}
