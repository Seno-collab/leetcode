package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
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
		Arguments:  []string{"1", "-2"},
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

	got, err := readInputCases(path, parseInts)
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

	_, err := readInputCases(path, parseInts)
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
	baseTime := time.Date(2026, time.August, 27, 0, 0, 0, 0, time.UTC)
	timestamps := []time.Time{
		baseTime,
		baseTime.Add(2 * time.Millisecond),
		baseTime.Add(10 * time.Millisecond),
		baseTime.Add(13 * time.Millisecond),
	}
	clockIndex := 0
	testClock := func() time.Time {
		currentTime := timestamps[clockIndex]
		clockIndex++
		return currentTime
	}

	err := run(
		nil,
		Options{InputFile: inputPath, OutputFile: outputPath},
		parseInts,
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
		testClock,
	)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	wantFile := existingResult +
		"{\"input\":[1,2,3,1],\"output\":true,\"duration_ns\":2000000}\n" +
		"{\"input\":[1,2,3,4],\"output\":false,\"duration_ns\":3000000}\n"
	gotFile, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotFile) != wantFile {
		t.Fatalf("output file = %q, want %q", gotFile, wantFile)
	}

	wantStdout := "=== Case 1 ===\n" +
		"Input      : [1,2,3,1]\n" +
		"Output     : true\n" +
		"Solve time : 2ms\n\n" +
		"=== Case 2 ===\n" +
		"Input      : [1,2,3,4]\n" +
		"Output     : false\n" +
		"Solve time : 3ms\n\n" +
		"=== Summary ===\n" +
		"Cases        : 2\n" +
		"Total time   : 5ms\n" +
		"Average time : 2.5ms\n" +
		"Results file : " + outputPath + "\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout = %q, want %q", stdout.String(), wantStdout)
	}
}

func TestRunJSONSupportsStructuredInputAndSliceOutput(t *testing.T) {
	type structuredInput struct {
		Numbers []int `json:"numbers"`
		K       int   `json:"k"`
	}

	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	outputPath := filepath.Join(tempDir, "results.jsonl")
	rawInput := "{\"numbers\":[1,1,1,2,2,3],\"k\":2}\n"
	if err := os.WriteFile(inputPath, []byte(rawInput), 0o644); err != nil {
		t.Fatal(err)
	}

	baseTime := time.Date(2026, time.August, 27, 0, 0, 0, 0, time.UTC)
	timestamps := []time.Time{baseTime, baseTime.Add(time.Microsecond)}
	clockIndex := 0
	var stdout bytes.Buffer

	err := run(
		nil,
		Options{InputFile: inputPath, OutputFile: outputPath},
		ParseJSON[structuredInput],
		func(input structuredInput) []int { return input.Numbers[:input.K] },
		&stdout,
		func() time.Time {
			currentTime := timestamps[clockIndex]
			clockIndex++
			return currentTime
		},
	)
	if err != nil {
		t.Fatalf("run returned an error: %v", err)
	}

	wantFile := "{\"input\":{\"numbers\":[1,1,1,2,2,3],\"k\":2}," +
		"\"output\":[1,1],\"duration_ns\":1000}\n"
	gotFile, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotFile) != wantFile {
		t.Fatalf("output file = %q, want %q", gotFile, wantFile)
	}

	if !strings.Contains(stdout.String(), "Input      : {\"numbers\":[1,1,1,2,2,3],\"k\":2}") {
		t.Fatalf("stdout does not contain structured JSON input: %q", stdout.String())
	}
}

func TestParseJSONRejectsUnknownFields(t *testing.T) {
	type structuredInput struct {
		Value int `json:"value"`
	}

	_, err := ParseJSON[structuredInput](`{"value":1,"unexpected":2}`)
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("ParseJSON error = %v, want unknown field error", err)
	}
}

func TestExecuteCasesClampsNegativeDuration(t *testing.T) {
	baseTime := time.Date(2026, time.August, 27, 0, 0, 0, 0, time.UTC)
	timestamps := []time.Time{baseTime, baseTime.Add(-time.Millisecond)}
	clockIndex := 0

	results, totalDuration := executeCases(
		[][]int{{1}},
		func([]int) int { return 1 },
		func() time.Time {
			currentTime := timestamps[clockIndex]
			clockIndex++
			return currentTime
		},
	)

	if results[0].DurationNS != 0 {
		t.Fatalf("duration = %d ns, want 0", results[0].DurationNS)
	}
	if totalDuration != 0 {
		t.Fatalf("total duration = %s, want 0", totalDuration)
	}
}
