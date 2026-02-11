package array

import "fmt"

func FindMaxConsecutiveOnes(nums []int) int {
	count := 0
	max := 0
	initId := 0
	for i := range nums {
		switch nums[i] {
		case 0:
			fmt.Println(initId, i)
			if count >= max {
				max = count
			}
		case 1:
			count++
		}
	}

	if count > max {
		return count
	}
	return max
}
