package main

import "fmt"

func main() {

    // Basic array declaration
    var nums [4]int
    fmt.Println(len(nums))  // 4

    // Adding elements
    nums[0] = 1
    nums[1] = 2
    fmt.Println(nums[0])    // 1
    fmt.Println(nums)       // [1 2 0 0]

    // Zero values — boolean
    var vals [4]bool
    fmt.Println(vals)       // [false false false false]
    vals[2] = true
    fmt.Println(vals)       // [false false true false]

    // Zero values — string
    var names [3]string
    names[0] = "Golang"
    fmt.Println(names)      // [Golang  ]

    // Single line declaration
    nums2 := [3]int{2, 3, 4}
    fmt.Println(nums2)      // [2 3 4]

    // 2D array
    grid := [2][2]int{{3, 4}, {5, 6}}
    fmt.Println(grid)       // [[3 4] [5 6]]

	var matrix = [2][2] int {{7,8},{9,10}}
	fmt.Println(matrix)
}