

package main
import "fmt"
import "math/rand"


func helper(numberChannel chan int){
    
    for result:=range numberChannel {
        fmt.Println("Processing number ,",result)
    }
}


func main() {
    
    
    numberChannel:=make(chan int)
    go helper(numberChannel)
    for{
        numberChannel<-rand.Intn(100)
    }
    
    
    
    
}
