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

// order struct

type customer struct {
	name  string
	phone string
}

// composition
type order struct {
	id        string
	amount    float32
	status    string
	createdAt time.Time // nanosecond precision
	customer
}

// func newOrder(id string, amount float32, status string) *order {
// 	// initial setup goes here...
// 	myOrder := order{
// 		id:     id,
// 		amount: amount,
// 		status: status,
// 	}

// 	return &myOrder
// }

// // receiver type
// func (o *order) changeStatus(status string) {
// 	o.status = status
// }

// func (o order) getAmount() float32 {
// 	return o.amount
// }

func main() {
	// newCustomer := customer{
	// 	name:  "john",
	// 	phone: "1234567890",
	// }
	newOrder := order{
		id:     "1",
		amount: 30,
		status: "received",
		customer: customer{
			name:  "john",
			phone: "1234567890",
		},
	}

	newOrder.customer.name = "robin"
	fmt.Println(newOrder)

	// language := struct {
	// 	name   string
	// 	isGood bool
	// }{"golang", true}

	// fmt.Println(language)

	// myOrder := newOrder("1", 30.50, "received")
	// fmt.Println(myOrder.amount)
	// if you don't set any field, default value is zero value
	// int => 0, float => 0, string "", bool => false
	// myOrder := order{
	// 	id:     "1",
	// 	amount: 50.00,
	// 	status: "received",
	// }
	// myOrder.changeStatus("confirmed")
	// fmt.Println(myOrder)
	// myOrder.createdAt = time.Now()
	// fmt.Println(myOrder.status)

	// myOrder2 := order{
	// 	id:        "2",
	// 	amount:    100,
	// 	status:    "delivered",
	// 	createdAt: time.Now(),
	// }

	// myOrder.status = "paid"

	// fmt.Println("Order struct", myOrder2)
	// fmt.Println("Order struct", myOrder)
}