// package main

// import (
// 	"fmt"
// 	"time"
// )

// // order struct

// type order struct{

// 	id 			string
// 	amount 		float32
// 	status 		string
// 	createdAt 	time.Time // precision is nanosecond 

// }

// // receiver type , attaches the function with struct , o is the first letter of strucut name , ek convention hai 
// // struct ko pass kya hai lekin pointer pass nahi kiya , it means it acts like the call by value kind of thing
// func (o *order) changeStatus(status string){
// o.status = status // yahape *o.status karne ki jarurat nahi hai kyuki strucut automatically dereferencing ka kaam karega 
// }

// func (o *order) getAmount()float32{
// 	return o.amount
// }

// func main(){
// 	myOrder:= order{
// 		id: "1",
// 		amount:50.00,
// 		status: "Received",


// 	}

// 	myOrder.changeStatus("Confirmed")
// 	fmt.Println(myOrder.getAmount())

// 	// myOrder.createdAt = time.Now()

// 	// fmt.Println(myOrder.amount)

// 	// myOrder2:=order{
// 	// 	id:"2",
// 	// 	amount:100,
// 	// 	status:"Delivered",
// 	// 	createdAt: time.Now(),
// 	// }
// 	// myOrder.status="Paid"
// 	// fmt.Println("Order struct",myOrder)
// 	// fmt.Println("Order struct",myOrder2)
// }

package main

import (
    "fmt"
    "time"
)

// Struct definition
type Order struct {
    id        string
    amount    float32
    status    string
    createdAt time.Time
}

// Constructor function
func newOrder(id string, amount float32, status string) *Order {
    myOrder := Order{
        id:        id,
        amount:    amount,
        status:    status,
        createdAt: time.Now(),
    }
    return &myOrder
}

// Method — modifying (pointer receiver)
func (o *Order) changeStatus(status string) {
    o.status = status
}

// Method — reading (normal receiver)
func (o Order) getAmount() float32 {
    return o.amount
}

func main() {

    // Create instance using constructor
    myOrder := newOrder("1", 50.00, "received")
    fmt.Println(myOrder)

    // Access fields
    fmt.Println(myOrder.status)  // received

    // Call method — change status
    myOrder.changeStatus("paid")
    fmt.Println(myOrder.status)  // paid

    // Call method — get amount
    fmt.Println(myOrder.getAmount())  // 50

    // Multiple independent instances
    myOrder2 := newOrder("2", 100.00, "delivered")
    fmt.Println(myOrder2.status)  // delivered — independent

    // Anonymous struct
    language := struct {
        name   string
        isGood bool
    }{
        "Golang",
        true,
    }
    fmt.Println(language)  // {Golang true}
}
