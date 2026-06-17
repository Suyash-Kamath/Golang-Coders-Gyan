package main

import "fmt"
import "time"

func main() {

    ch := make(chan int)

    select {
    case value := <-ch:
        fmt.Println(value)
        
    // Time.After() does NOT "become" a channel. It actually RETURNS a channel.
    
    case <- time.After(3*time.Second):
        fmt.Println("timeout")

    // default:
    //     fmt.Println("No data")
    }
}
