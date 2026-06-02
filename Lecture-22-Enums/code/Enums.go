/*
What an enum actually is
An enum is a type that can only hold one value out of a fixed, named set you defined. That's the whole concept. A traffic light is red | yellow | green and nothing else. An order is received | confirmed | prepared | delivered and nothing else.
The key correction to your mental model: enums are not "incrementing things." The incrementing is just how Go happens to fill in the underlying values using iota. That's a storage detail, not the point.
The actual point of an enum is restriction + naming:

A string can be infinitely many values.
An OrderStatus enum can be exactly one of the few you listed.

The numbers (0, 1, 2…) underneath are just labels Go uses internally so each name maps to something. You could back them with strings instead and there'd be zero incrementing — and it would still be an enum. So "incrementing" ≠ "enum." Restriction to a named set = enum.



Almost — one small fix. It's a specific list, not a range.
A range means everything between two points (like 1–100 = all 100 numbers). An enum is just a handful of specific named values you picked, not a continuous span:
goconst (
    Received  OrderStatus = iota // a specific value
    Confirmed                    // a specific value
    Delivered                    // a specific value
)
That's 3 chosen values — Received, Confirmed, Delivered — not "a range from Received to Delivered."
The fact that they happen to get numbered 0, 1, 2 underneath is just labeling; you're not saying "any value from 0 to 2 is allowed," you're saying "only these exact named members are valid."
So: enum = a data type whose values are a fixed, discrete set of named options you defined. Think list of specific choices, not range.




An enum is a fixed set of specific named values that you picked, bundled together as a data type. Only those values are valid for that type — nothing else.


"a fixed set" — once defined, the list is locked: exactly Received, Confirmed, Delivered. No more, no fewer (unless you edit the code).
"specific named values" — each option has a name (Received, not just a bare 0). They're specific picks, not a range.
"that you picked" — you chose these three. The language didn't; you decided what's valid.
"bundled together as a data type" — they all share one type, OrderStatus. That type is what ties them into a single group.
"only those values are valid — nothing else" — a variable of type OrderStatus is meant to hold one of those three. changeOrderStatus(Confirmed) works; passing a random string or unrelated value doesn't fit the type.
So the statement holds: enum = your hand-picked, named, fixed set of values, unified under one data type.


*/
package main
import "fmt"


// type OrderStatus int

// const (
//     Received OrderStatus =iota
//     Confirmed
//     Delivered
//     Prepared
//     )

type OrderStatus string

const (
    Received OrderStatus = "received"
    Confirmed = "confirmed"
    Prepared = "prepared"
    Delivered = "delivered"
    
    )


func changeOrderStatus(status OrderStatus){
    fmt.Println("Changing Order Status to ", status)
}

func main() {
    changeOrderStatus(Received)
}
