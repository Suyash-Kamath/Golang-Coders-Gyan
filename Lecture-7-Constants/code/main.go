package main

import "fmt"

// Package level constant
const age = 30

// Grouped constants
const (
    port = 5000
    host = "localhost"
)

func main() {

    // Simple constant
    const name = "Golang"
    fmt.Println(name)   // Output: Golang

    // Integer constant
    fmt.Println(age)    // Output: 30

    // Grouped constants
    fmt.Println(port)   // Output: 5000
    fmt.Println(host)   // Output: localhost

    // This would cause an ERROR:
    // name = "JavaScript"  ❌ cannot assign to name
}