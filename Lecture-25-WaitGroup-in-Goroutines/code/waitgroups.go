package main

import "fmt"
import "sync"

func task(id int , w *sync.WaitGroup ){
    defer w.Done()
    fmt.Println("Doing task ",id)
}

func main(){
    var wg sync.WaitGroup
    
    
    for i:=0 ; i<=10;i++{
        wg.Add(1)
        go task(i,&wg)
    }
    
    wg.Wait()
    fmt.Println("Hi Suyash")
}
