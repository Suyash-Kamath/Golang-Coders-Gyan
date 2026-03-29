package main

import "fmt"

// Basic function
func add(a, b int) int {
    return a + b
}

// Multiple return values
func getLanguages() (string, string, string) {
    return "Golang", "JavaScript", "C"
}

// Mixed return types
func getInfo() (string, bool) {
    return "Golang", true
}

// Function accepting a function as parameter
func processIt(f func(int) int) {
    result := f(1)
    fmt.Println(result)
}

// Function returning a function
func makeMultiplier() func(int) int {
    return func(a int) int {
        return a * 2
    }
}

func main() {

    // Basic function call
    result := add(3, 5)
    fmt.Println(result)  // 8

    // Shorthand parameters — same result
    fmt.Println(add(4, 6))  // 10

    // Multiple return values
    lang1, lang2, lang3 := getLanguages()
    fmt.Println(lang1, lang2, lang3)  // Golang JavaScript C

    // Ignoring a return value
    l1, l2, _ := getLanguages()
    fmt.Println(l1, l2)  // Golang JavaScript

    // Print all at once
    fmt.Println(getLanguages())  // Golang JavaScript C

    // Passing function as argument
    fA := func(a int) int {
        return 2
    }
    processIt(fA)  // 2

    // Returning function from function
    multiplier := makeMultiplier()
    fmt.Println(multiplier(6))  // 12
}
