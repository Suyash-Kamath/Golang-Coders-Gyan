package main

import "fmt"

func main() {

    // Range over slice — value only
    nums := []int{6, 7, 8}
    for _, num := range nums {
        fmt.Println(num)
    }

    // Range over slice — sum
    sum := 0
    for _, num := range nums {
        sum = sum + num
    }
    fmt.Println("Sum:", sum)  // 21

    // Range over slice — index and value
    for i, num := range nums {
        fmt.Println(i, num)
    }

    // Range over map
    m := map[string]string{
        "fName": "John",
        "lName": "Doe",
    }
    for k, v := range m {
        fmt.Println(k, v)
    }

    // Range over map — keys only
    for k := range m {
        fmt.Println(k)
    }

    // Range over string — Unicode values
    for i, c := range "Golang" {
        fmt.Println(i, c)
    }

    // Range over string — actual characters
    for i, c := range "Golang" {
        fmt.Println(i, string(c))
    }

	for i, c := range "G😀lang" {
    fmt.Println(i, string(c))
}

}