// fan out means multiple goroutines are doing the same work concurrently and then we wait for all of them to finish before we move on to the next step.

package main

import "fmt"
import "time"
import "sync"


func worker(url string, wg *sync.WaitGroup)  {
    
    defer wg.Done()
    
    time.Sleep(time.Millisecond*50)
    fmt.Printf("Image processed: %s\n",url)
    
}


func main(){
    // fmt.Println("Welcome to Go Concurrency")
    var wg sync.WaitGroup
    startTime:=time.Now()
    
   wg.Add(5)
   go worker("image_1.png",&wg ) // error : goroutine ke andhar code run karte hai toh value return nahi kar sakte iss function se 
    go worker("image_2.png",&wg)
    go worker("image_3.png",&wg)
    go worker("image_4.png",&wg)
    go worker("image_5.png",&wg)
    
    
    wg.Wait()
    fmt.Printf("It took %s ms.\n",time.Since(startTime))
    
}
