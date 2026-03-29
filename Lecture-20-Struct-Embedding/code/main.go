package main

import "fmt"

// Inner struct
type Customer struct {
    name  string
    phone string
}

// Method on Customer
func (c *Customer) greet() {
    fmt.Println("Hello, I am", c.name)
}

// Outer struct with embedded Customer
type Order struct {
    id       string
    amount   float32
    status   string
    customer Customer
}

func main() {

    // Way 1 — separate customer then assign
    newCustomer := Customer{
        name:  "John",
        phone: "1234567890",
    }

    newOrder := Order{
        id:       "1",
        amount:   30,
        status:   "received",
        customer: newCustomer,
    }

    fmt.Println(newOrder)
    // {1 30 received {John 1234567890}}

    // Access embedded struct
    fmt.Println(newOrder.customer)
    // {John 1234567890}

    fmt.Println(newOrder.customer.name)
    // John

    // Change embedded struct field
    newOrder.customer.name = "Robin"
    fmt.Println(newOrder.customer.name)
    // Robin

    // Way 2 — inline customer
    newOrder2 := Order{
        id:     "2",
        amount: 50,
        status: "delivered",
        customer: Customer{
            name:  "Alice",
            phone: "9876543210",
        },
    }
    fmt.Println(newOrder2)

    // Call Customer method through Order
    newOrder2.customer.greet()
    // Hello, I am Alice
}
