package main


import "fmt"

// func printSlice[T int | string | bool](items []T){
// func printSlice[T comparable](items []T){

func printSlice[T comparable, V string](items []T, name V){
    for _,item:=range items{
        fmt.Println(item,name)
    }
}

// func printStringSlice(items []string){
//     for _,item:=range items{
//         fmt.Println(item)
//     }
// }

// LIFO 
// type Stack[T any] struct{
//     elements []T
    
// }

func main(){
    // nums:=[]int{1,2,3}
    // names:=[]string{"Golang","TypeScript"}
    vals:=[]bool{true,false,true}
    printSlice(vals,"Suyash")
    
    // stack:= Stack[int]{
    //     elements:[]int{1,2,3},
    // }
    
    // fmt.Println(stack)
    
    
    
    
}
