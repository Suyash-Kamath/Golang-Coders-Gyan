
package main

import "fmt"
import "time"

func worker(url string)  {
    
    time.Sleep(time.Millisecond*50)
    fmt.Printf("Image processed: %s\n",url)
    
}


func main(){
    // fmt.Println("Welcome to Go Concurrency")
    
    startTime:=time.Now()
    
   go worker("image_1.png") // error : goroutine ke andhar code run karte hai toh value return nahi kar sakte iss function se 
    go worker("image_2.png")
    
    
    
    fmt.Printf("It took %s ms.\n",time.Since(startTime))
    
}
