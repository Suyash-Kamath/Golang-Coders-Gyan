package main

import "fmt"

func runes() {

    // Basic rune
    var r rune = 'G'
    fmt.Println(r)          // 71 — Unicode code point
    fmt.Println(string(r))  // G — actual character

    // Rune with emoji
    var emoji rune = '😀'
    fmt.Println(emoji)          // 128512
    fmt.Println(string(emoji))  // 😀

    // String byte length vs rune length
    s := "Go😀"
    fmt.Println(len(s))           // 6 — bytes
    fmt.Println(len([]rune(s)))   // 3 — actual characters

    // Convert string to rune slice
    runes := []rune("Golang")
    fmt.Println(runes)  // [71 111 108 97 110 103]

    // Range over string — getting runes
    for i, c := range "Golang" {
        fmt.Println(i, string(c))
    }

    // Rune vs byte
    var b byte = 'A'
    var ru rune = 'A'
    fmt.Println(b)   // 65
    fmt.Println(ru)  // 65 — same for ASCII
}

