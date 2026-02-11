package arrayii

import (
	"slices"
)

func FindDisappearedNumbers(nums []int) []int {
	ans := make([]int, 0)
	slices.Sort(nums)

	n := len(nums)
	for i := 0; i < n-1; i++ {
		if nums[i+1]-nums[i] == 1 || nums[i+1]-nums[i] == 0 {
			continue
		} else {
			idx := nums[i+1] - nums[i]
			for j := 1; j < idx; j++ {
				ans = append(ans, nums[i]+j)
			}
		}
	}
	return ans
}
