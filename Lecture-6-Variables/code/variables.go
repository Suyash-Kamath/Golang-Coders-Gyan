package main

import "fmt"

func main() {

    // --- String ---
    var name1 string = "Golang"   // Method 1
    var name2 = "Golang"          // Method 2
    name3 := "Golang"             // Method 3

    // --- Boolean ---
    var isAdult1 bool = true
    var isAdult2 = true
    isAdult3 := true

    // --- Integer ---
    var age1 int = 30
    var age2 = 30
    age3 := 30

    // --- Float ---
    var price1 float32 = 50.5
    var price2 = 50.5
    price3 := 50.5

    fmt.Println(name1, name2, name3)
    fmt.Println(isAdult1, isAdult2, isAdult3)
    fmt.Println(age1, age2, age3)
    fmt.Println(price1, price2, price3)
}
