package main


import "fmt"
import "sync"

type post struct{
    views int 
    mu sync.Mutex
    // good practice as you are making the mutex lock for this resource so put this into this struct
    // race condition is solved 
    
}

func (p *post) inc(wg *sync.WaitGroup){
    defer func(){
        wg.Done()
        p.mu.Unlock()
    } ()
    // error aa jaye kuch bhi ho jaaye  , defer function() call hoti hai function execution ke baad last me to Unlock ho jaayega 
    
    // good practice is use lock where exact modification happens , taaki utni hi jagaha bottleneck ban jaaye 
    p.mu.Lock()
    p.views+=1
    
}

func main(){
    
    var wg sync.WaitGroup
    
    myPost:=post{
        views:0,
    }
    
    for i:=0;i<100;i++{
        
        wg.Add(1)
         go myPost.inc(&wg)
         
    }
   
   wg.Wait()
    
    fmt.Println(myPost.views)
    
    
}
