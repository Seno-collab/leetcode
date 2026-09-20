package main

import (
	"reflect"
	"testing"
)

func TestTopKFrequent(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int
		k       int
		want    []int
	}{
		{name: "two most frequent", numbers: []int{1, 1, 1, 2, 2, 3}, k: 2, want: []int{1, 2}},
		{name: "one value", numbers: []int{1}, k: 1, want: []int{1}},
		{name: "stable tie", numbers: []int{4, 1, -1, 2, -1, 2, 3}, k: 2, want: []int{-1, 2}},
		{name: "zero k", numbers: []int{1, 1}, k: 0, want: []int{}},
		{name: "k exceeds unique values", numbers: []int{2, 2}, k: 3, want: []int{2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := topKFrequent(test.numbers, test.k)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("topKFrequent(%v, %d) = %v, want %v", test.numbers, test.k, got, test.want)
			}
		})
	}
}
