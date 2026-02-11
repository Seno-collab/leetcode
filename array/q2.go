package array

func shuffle(nums []int, n int) []int {
	ans := make([]int, 0)
	for i := range n {
		ans = append(ans, []int{nums[i], nums[i+n]}...)
	}
	return ans
}
