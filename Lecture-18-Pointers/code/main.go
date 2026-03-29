package main

import "fmt"

// Pass by value — original NOT changed
func changeNumByValue(num int) {
    num = 5
    fmt.Println("in changeNum (value):", num)
}

// Pass by reference — original IS changed
func changeNumByRef(num *int) {
    *num = 5
    fmt.Println("in changeNum (ref):", *num)
}

func main() {

    // Pass by value demo
    num1 := 1
    changeNumByValue(num1)
    fmt.Println("after value change:", num1)  // still 1

    // See memory address
    num2 := 1
    fmt.Println("memory address:", &num2)  // 0xc000...

    // Pass by reference demo
    num3 := 1
    changeNumByRef(&num3)                          // pass address
    fmt.Println("after ref change:", num3)         // 5 — changed!
}

