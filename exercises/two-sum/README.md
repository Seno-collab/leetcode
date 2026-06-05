# LeetCode 01: Two Sum

## Goal

Given an integer slice `nums` and an integer `target`, return the indices of the two numbers such that they add up to `target`.

Each input has exactly one answer, and the same element cannot be used twice.

## Example

```go
nums := []int{2, 7, 11, 15}
target := 9
```

Expected result:

```go
[]int{0, 1}
```

Because:

```go
nums[0] + nums[1] == 2 + 7 == 9
```

## Function To Implement

```go
func TwoSum(nums []int, target int) []int
```

## Practice Path

1. Solve it with two loops first.
2. Explain why that is `O(n^2)`.
3. Then solve it with a map.
4. Explain why the map version is `O(n)`.

