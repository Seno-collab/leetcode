package array

func GetConcatenation(nums []int) []int {
	result := make([]int, 0)
	result = append(result, nums...)
	result = append(result, nums...)
	return result
}
