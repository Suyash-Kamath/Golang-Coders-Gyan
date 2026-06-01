package main

import "fmt"

/*
========================================================================
  GO INTERFACES — the full journey from the transcript
  Use case: plugging payment gateways into a payment system
========================================================================

WHY INTERFACES EXIST (the story the tutor builds up):

  STAGE 0 — Naive, no abstraction:
    makePayment creates the gateway itself and calls it directly:

        func (p Payment) makePayment(amount float32) {
            razorpay := Razorpay{}   // hard-wired dependency
            razorpay.Pay(amount)
        }

    Works fine... until a requirement changes.

  STAGE 1 — The problem (Open/Closed violation):
    "Switch Razorpay -> Stripe." You must EDIT makePayment:

        // razorpay := Razorpay{}
        // razorpay.Pay(amount)
        stripe := Stripe{}
        stripe.Pay(amount)

    You modified already-tested code. Open/Closed Principle says code
    should be OPEN for extension, CLOSED for modification. Violated.

  STAGE 2 — Dependency injection, but with a CONCRETE field type:
        type Payment struct { gateway Razorpay }   // <- still concrete!
    Better, but the field is locked to Razorpay, so:
      - switching providers means editing the struct definition, AND
      - you CANNOT inject a fake gateway for unit testing.

  STAGE 3 — The cure (the code below): type the field to an INTERFACE.
    An interface is a CONTRACT of behaviour. Any type with the matching
    method(s) satisfies it AUTOMATICALLY — Go's implicit / structural
    typing (no `implements` keyword). Now you EXTEND (add a gateway)
    without MODIFYING makePayment, and you can inject a fake.
    This is the Dependency Inversion Principle (the "D" in SOLID).
========================================================================
*/

// ----------------------------------------------------------------------
// THE INTERFACE = the contract.
// "Anything that has a Pay(float32) method IS a PaymentGateway."
// Convention: single-method interfaces are often named -er
// (Reader, Writer, Stringer). PaymentGateway is descriptive and fine.
// ----------------------------------------------------------------------
type Paymenter interface {
	Pay(amount float32)
	// To grow the contract, add more signatures here, e.g.:
	//     Refund(amount float32, account string)
	// The moment you do, EVERY implementer must add it too or it stops
	// satisfying the interface (checked at compile time, at point of use).
	// Keep interfaces small -> Interface Segregation Principle.
}

// ----------------------------------------------------------------------
// CONCRETE IMPLEMENTATIONS (the low-level "details").
// NONE of them declare "implements PaymentGateway" — in Go it's IMPLICIT.
// They satisfy the interface purely by HAVING a Pay(float32) method.
// ----------------------------------------------------------------------

type Razorpay struct{}

func (r Razorpay) Pay(amount float32) {
	// Real code: call Razorpay's API here.
	fmt.Println("Making payment using Razorpay", amount)
}

type Stripe struct{}

func (s Stripe) Pay(amount float32) {
	fmt.Println("Making payment using Stripe", amount)
}

type Paypal struct{}

func (p Paypal) Pay(amount float32) {
	fmt.Println("Making payment using Paypal", amount)
}

// Fakepayment is a TEST DOUBLE. Because it also has Pay(float32) it
// satisfies PaymentGateway, so it can be injected in unit tests instead
// of a real provider — no network calls, no real charges.
type Fakepayment struct{}

func (f Fakepayment) Pay(amount float32) {
	fmt.Println("FAKE PAYMENT", amount)
}

// ----------------------------------------------------------------------
// PAYMENT SERVICE LAYER (the high-level "policy").
// It depends on the ABSTRACTION (PaymentGateway), NOT on any concrete
// provider. This one field is the entire point of the lesson.
// ----------------------------------------------------------------------
type Payment struct {
	gateway Paymenter // <- the "seam": accepts ANY gateway
}

func (p Payment) makePayment(amount float32) {
	// --- STAGE 0/1 (the WRONG way) would have been: ---
	//     razorpay := Razorpay{}   // dependency created INSIDE
	//     razorpay.Pay(amount)     // = hard-wired, untestable
	//
	//     stripe := Stripe{}
	//     stripe.Pay(amount)
	//
	//     paypal := Paypal{}
	//     paypal.Pay(amount)
	//
	// --- STAGE 3 (the RIGHT way): just use whatever was injected. ---
	// makePayment neither knows nor cares which gateway it received.
	p.gateway.Pay(amount)
}

// ----------------------------------------------------------------------
// CONSTRUCTOR (optional, but idiomatic): makes the dependency mandatory
// so a caller can't forget to set the gateway.
// ----------------------------------------------------------------------
func NewPayment(gw Paymenter) Payment {
	return Payment{gateway: gw}
}

func main() {
	// ==================================================================
	// DEPENDENCY INJECTION HAPPENS HERE — not inside makePayment.
	// We build the gateway OUTSIDE and hand it in. The decision of WHICH
	// provider to use lives in main (the high-level wiring point), so
	// makePayment stays untouched no matter how many gateways we add.
	// ==================================================================

	// Production: swap freely, makePayment never changes:
	//     p := Payment{gateway: Razorpay{}}
	//     p := Payment{gateway: Stripe{}}
	//     p := Payment{gateway: Paypal{}}

	// Here we inject the FAKE (exactly what you'd do in a test):
	fake := Fakepayment{}
	p := Payment{gateway: fake} // <-- THIS line is the injection
	// (or with the constructor:  p := NewPayment(fake) )

	p.makePayment(100) // -> FAKE PAYMENT 100
}
