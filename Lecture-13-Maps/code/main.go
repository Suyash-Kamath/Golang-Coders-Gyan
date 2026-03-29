package main

import (
    "fmt"
    "maps"
)

func main() {

    // Way 1 — make()
    m := make(map[string]string)
    m["name"] = "Golang"
    m["area"] = "Backend"
    fmt.Println(m["name"])   // Golang
    fmt.Println(m["area"])   // Backend
    fmt.Println(m["phone"])  // "" — key doesn't exist, zero value

    // Length
    fmt.Println(len(m))  // 2

    // Delete element
    delete(m, "area")
    fmt.Println(m)  // map[name:Golang]

    // Clear all
    clear(m)
    fmt.Println(m)  // map[]

    // Way 2 — map literal
    m2 := map[string]int{
        "price":  40,
        "phones": 3,
    }
    fmt.Println(m2)  // map[phones:3 price:40]

    // Check if key exists
    v, ok := m2["price"]
    if ok {
        fmt.Println("Found:", v)   // Found: 40
    } else {
        fmt.Println("Not found")
    }

    v2, ok2 := m2["phone"]
    if !ok2 {
        fmt.Println("Not found, zero value:", v2)  // Not found, zero value: 0
    }

    // Compare maps
    map1 := map[string]int{"price": 40}
    map2 := map[string]int{"price": 40}
    map3 := map[string]int{"price": 80}

    fmt.Println(maps.Equal(map1, map2))  // true
    fmt.Println(maps.Equal(map1, map3))  // false
}
