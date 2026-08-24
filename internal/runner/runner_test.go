package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	config, err := parseCommandLine(
		[]string{"--input", "cases.txt", "--output", "actual.jsonl", "1", "-2"},
		Options{InputFile: "input.txt", OutputFile: "results.jsonl"},
	)
	if err != nil {
		t.Fatalf("parseCommandLine returned an error: %v", err)
	}

	want := Config{
		InputFile:  "cases.txt",
		OutputFile: "actual.jsonl",
		Numbers:    []string{"1", "-2"},
	}
	if !reflect.DeepEqual(config, want) {
		t.Fatalf("parseCommandLine() = %#v, want %#v", config, want)
	}
}

func TestParseCommandLineRejectsMissingFlagValue(t *testing.T) {
	_, err := parseCommandLine(
		[]string{"-out"},
		Options{InputFile: "input.txt", OutputFile: "results.jsonl"},
	)
	if err == nil {
		t.Fatal("parseCommandLine should reject -out without a file path")
	}
}

func TestReadInputCases(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte("1 2 3\n\n10 -2 30\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := readInputCases(path)
	if err != nil {
		t.Fatalf("readInputCases returned an error: %v", err)
	}

	want := [][]int{{1, 2, 3}, {10, -2, 30}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("readInputCases() = %v, want %v", got, want)
	}
}

func TestReadInputCasesReportsLineWithInvalidNumber(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte("1 2\n3 nope\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := readInputCases(path)
	if err == nil || !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("readInputCases error = %v, want an error containing line 2", err)
	}
}

func TestRunReadsInputAndAppendsOutput(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	outputPath := filepath.Join(tempDir, "results.jsonl")
	if err := os.WriteFile(inputPath, []byte("1 2 3 1\n1 2 3 4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	existingResult := "{\"input\":[9],\"output\":false}\n"
	if err := os.WriteFile(outputPath, []byte(existingResult), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	err := Run(
		nil,
		Options{InputFile: inputPath, OutputFile: outputPath},
		func(numbers []int) bool {
			seen := make(map[int]bool)
			for _, number := range numbers {
				if seen[number] {
					return true
				}
				seen[number] = true
			}
			return false
		},
		&stdout,
	)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	wantFile := existingResult +
		"{\"input\":[1,2,3,1],\"output\":true}\n" +
		"{\"input\":[1,2,3,4],\"output\":false}\n"
	gotFile, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotFile) != wantFile {
		t.Fatalf("output file = %q, want %q", gotFile, wantFile)
	}

	if !strings.Contains(stdout.String(), "Output: true") ||
		!strings.Contains(stdout.String(), "Output: false") {
		t.Fatalf("stdout does not contain both results: %q", stdout.String())
	}
}
