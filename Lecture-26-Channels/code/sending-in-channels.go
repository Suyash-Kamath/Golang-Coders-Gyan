package main

import "fmt"

func processNum(numChan chan int){
    
    numChan<-11
    
}

func main(){


messageChannel:=make(chan int)

go processNum(messageChannel)

result:= <-messageChannel

fmt.Println(result)


}
