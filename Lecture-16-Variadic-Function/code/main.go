package main

import "fmt"

// Variadic function — accepts unlimited ints
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total = total + num
    }
    return total
}

// Variadic function — accepts any type
func printAll(values ...interface{}) {
    for _, v := range values {
        fmt.Println(v)
    }
}

func main() {

    // Passing individual values
    result := sum(3, 4, 5, 6)
    fmt.Println(result)  // 18

    // Passing different counts
    fmt.Println(sum(1, 2))      // 3
    fmt.Println(sum(10))        // 10
    fmt.Println(sum())          // 0

    // Passing a slice — use spread operator
    nums := []int{3, 4, 5, 6}
    result2 := sum(nums...)
    fmt.Println(result2)  // 18

    // Any type variadic
    printAll(1, "hello", true, 3.14)
}
