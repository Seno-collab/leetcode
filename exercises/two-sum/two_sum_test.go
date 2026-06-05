package twosum

import "testing"

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{
			name:   "example",
			nums:   []int{2, 7, 11, 15},
			target: 9,
			want:   []int{0, 1},
		},
		{
			name:   "uses different indices with same value",
			nums:   []int{3, 3},
			target: 6,
			want:   []int{0, 1},
		},
		{
			name:   "with negative number",
			nums:   []int{-3, 4, 3, 90},
			target: 0,
			want:   []int{0, 2},
		},
		{
			name:   "answer at the end",
			nums:   []int{1, 5, 8, 12},
			target: 20,
			want:   []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum(tt.nums, tt.target)

			if len(got) != 2 {
				t.Fatalf("expected 2 indices, got %v", got)
			}

			if got[0] != tt.want[0] || got[1] != tt.want[1] {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

