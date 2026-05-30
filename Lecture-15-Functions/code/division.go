package main

import "fmt"

func divide(a,b int) (int,error){
    if b==0{
        return 0,fmt.Errorf("Cannot divide by 0")
        
    }
    
    return a/b , nil
}

func main(){
    
    result,err := divide(10,0)
    if err!=nil{
        fmt.Println("Error ",err)
    }else{
        fmt.Println(result)
    }
}
