package arrayii

import "slices"

func SmallerNumbersThanCurrent(nums []int) []int {
	ans := make([]int, 0)
	sortNums := make([]int, len(nums))
	copy(sortNums, nums)
	slices.Sort(sortNums)
	for index, value := range nums {
		j := slices.Index(sortNums, value)
		if j != -1 {
			ans[index] = j
		}
	}
	return ans
}
