package set

// 幂集。编写一种方法，返回某集合的所有子集。集合中不包含重复的元素。
// 说明：解集不能包含重复的子集。
// https://leetcode-cn.com/problems/power-set-lcci/
func Subsets(nums []int) [][]int {
	res := make([][]int, 0)
	res = append(res, []int{})
	for k1 := 0; k1 < len(nums); k1++ {
		a1 := []int{nums[k1]}
		res = append(res, a1)
		loop(nums, &res, a1, k1)
	}
	return res
}

func loop(nums []int, res *[][]int, a1 []int, k int) {
	for k1 := k + 1; k1 < len(nums); k1++ {
		a2 := make([]int, len(a1))
		copy(a2, a1)
		a2 = append(a2, nums[k1])
		(*res) = append((*res), a2)
		loop(nums, res, a2, k1)
	}
}
