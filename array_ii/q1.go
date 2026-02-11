package arrayii

func FindErrorNums(nums []int) []int {
	dup, missing := -1, -1
	for i := 0; i < len(nums); i++ {
		x := nums[i]
		if x < 0 {
			x = -x
		}
		idx := x - 1
		if nums[idx] < 0 {
			dup = x
		} else {
			nums[idx] = -nums[idx]
		}
	}
	for i := 0; i < len(nums); i++ {
		if nums[i] > 0 {
			missing = i + 1
			break
		}
	}
	return []int{dup, missing}
}

//  1 2 2 4
//  1 1
