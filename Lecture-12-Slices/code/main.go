package main

import (
    "fmt"
    "slices"
)

func main() {
// uninitialized slice is nil
	// var nums []int
	// fmt.Println(nums == nil)
	// fmt.Println(len(nums))

	// var nums = make([]int, 0, 5)
	// capacity -> maximum numbers of elements can fit
	// fmt.Println(cap(nums))
	// fmt.Println(nums == nil)

	// nums := []int{}

	// nums = append(nums, 1)
	// nums = append(nums, 2)

	// nums[0] = 3
	// nums[1] = 5
	// fmt.Println(nums)
	// fmt.Println(cap(nums))
	// fmt.Println(len(nums))

	// var nums = make([]int, 0, 5)
	// nums = append(nums, 2)
	// var nums2 = make([]int, len(nums))

	// // copy function
	// copy(nums2, nums)
	// fmt.Println(nums, nums2)

	// slice operator

	// var nums = []int{1, 2, 3, 4, 5}
	// fmt.Println(nums[0:1])
	// fmt.Println(nums[:2])
	// fmt.Println(nums[1:])

	// slices
	// var nums1 = []int{1, 2, 3}
	// var nums2 = []int{1, 2, 4}

	// fmt.Println(slices.Equal(nums1, nums2))

	// var nums = [][]int{{1, 2, 3}, {4, 5, 6}}
	// fmt.Println(nums)

	
    // Way 1 — Nil slice
    var nums []int
    fmt.Println(nums == nil)    // true
    fmt.Println(len(nums))      // 0

    // Way 2 — make()
    nums2 := make([]int, 0, 5)
    fmt.Println(cap(nums2))     // 5

    // Way 3 — short declaration
    nums3 := []int{1, 2, 3}
    fmt.Println(nums3)          // [1 2 3]

    // append()
    nums4 := []int{}
    nums4 = append(nums4, 1)
    nums4 = append(nums4, 2, 3)
    fmt.Println(nums4)          // [1 2 3]

    // copy()
    src := []int{10, 20}
    dst := make([]int, len(src))
    copy(dst, src)
    fmt.Println(dst)            // [10 20]

    // slice operator
    s := []int{1, 2, 3, 4, 5}
    fmt.Println(s[1:3])         // [2 3]
    fmt.Println(s[:2])          // [1 2]
    fmt.Println(s[3:])          // [4 5]

    // slices package
    a := []int{1, 2, 3}
    b := []int{1, 2, 3}
    fmt.Println(slices.Equal(a, b))  // true

    // 2D slice
    grid := [][]int{{1, 2}, {3, 4}}
    fmt.Println(grid)           // [[1 2] [3 4]]
}
