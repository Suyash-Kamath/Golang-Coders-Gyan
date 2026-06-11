// channels : data ko ek goroutine se dusre goroutine me transfer karne ke liye use hota hai.
// channels ke through hum data ko synchronize kar sakte hai, matlab ek goroutine se dusre goroutine me data transfer kar sakte hai without using locks or other synchronization primitives.
// channels ke through hum multiple goroutines ke beech me communication kar sakte hai, matlab ek goroutine se dusre goroutine me data transfer kar sakte hai without using locks or other synchronization primitives.


// Channel is memory ke andhar shared space where go routines has access of it and does communication through memory 
package main

import "fmt"
import "time"
import "sync"


func worker(url string, wg *sync.WaitGroup, resultChan chan string)  {
    
    defer wg.Done()
    
    time.Sleep(time.Millisecond*50)
    fmt.Printf("Image processed: %s\n",url)
    
    resultChan <-url
    
}

// Channel is memory ke andhar shared space where go routines has access of it and does communication through memory 

func main(){
    // fmt.Println("Welcome to Go Concurrency")
    var wg sync.WaitGroup
    resultChan:=make(chan string,5)
    
    startTime:=time.Now()
    
   wg.Add(5)
   go worker("image_1.png",&wg , resultChan ) // error : goroutine ke andhar code run karte hai toh value return nahi kar sakte iss function se 
    go worker("image_2.png",&wg,resultChan)
    go worker("image_3.png",&wg,resultChan)
    go worker("image_4.png",&wg,resultChan)
    go worker("image_5.png",&wg,resultChan)
    
    
    wg.Wait()
    close(resultChan)
    
    // Channel ke upar read karte jaa rahe ho , close karo warna deadlock aa jayega
    for result:=range resultChan{
        fmt.Printf("Received: %s\n",result)
    }
    
    fmt.Printf("It took %s ms.\n",time.Since(startTime))
    
}
