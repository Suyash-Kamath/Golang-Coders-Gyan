package main


import "time"
import "fmt"
// import "math/rand"


func helper(numChan chan int){
    
    result:=<-numChan
    
    fmt.Println(result)
}

func main(){
    
    numChan:=make(chan int)
    /*
    
    Your deadlock comes from channel blocking + no goroutine being able to proceed.

Let’s walk through your code step by step:

❌ Your code
func helper(numChan chan int){
    
    result := <-numChan
    fmt.Println(result)
}

func main(){

    numChan := make(chan int)

    numChan <- 5   // ❌ PROBLEM HERE

    go helper(numChan)

    time.Sleep(time.Second*2)
}
💥 What caused the deadlock?
🚨 1. Unbuffered channel blocks immediately
numChan := make(chan int)

This is an unbuffered channel, meaning:

send (<-) blocks until someone receives
receive blocks until someone sends
🚨 2. This line blocks forever
numChan <- 5

At this moment:

No goroutine is receiving yet
So main() is stuck here permanently

👉 It never reaches the next line:

go helper(numChan)
💣 Result
main goroutine is stuck on send
helper goroutine never starts
nobody receives
Go detects:

“all goroutines are asleep → deadlock”

🧠 Key Insight
You wrote:
send first → then start receiver

But Go requires:

receiver must exist BEFORE send (for unbuffered channels)
    
    
    */
    
    
    go helper(numChan)
    
    numChan<-5
    
    
    time.Sleep(time.Second*2)
    
    
    
    
    
    
    
    
}
