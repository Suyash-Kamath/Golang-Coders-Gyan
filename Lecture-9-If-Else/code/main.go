package main

import "fmt"

func main() {

    // Basic if/else
    age := 16
    if age >= 18 {
        fmt.Println("Person is an adult")
    } else {
        fmt.Println("Person is not an adult")
    }

    // if / else if / else
    if age >= 18 {
        fmt.Println("Adult")
    } else if age >= 12 {
        fmt.Println("Teenager")
    } else {
        fmt.Println("Kid")
    }

    // Logical operators
    role := "admin"
    hasPermission := false

    if role == "admin" || hasPermission {
        fmt.Println("Access granted (OR)")
    }

    if role == "admin" && hasPermission {
        fmt.Println("Access granted (AND)")
    }

    // Variable declared inside if
    if age := 15; age >= 18 {
        fmt.Println("Adult", age)
    } else if age >= 12 {
        fmt.Println("Teenager", age)
    } else {
        fmt.Println("Kid", age)
    }
}