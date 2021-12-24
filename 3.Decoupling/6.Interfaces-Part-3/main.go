// Sample program to show how the concrete value assigned to
// the interface is what is stored inside the interface.
package main

import "fmt"

// printer displays information.
type printer interface {
	print()
}

// cannon defines a cannon printer.
type cannon struct {
	name string
}

/*
	Cannon implements printer interface using the print function, using value semantics.
	Value semantics mean that we can store both values and addresses inside that interface.
	Pointer semantics mean that we can only share those values inside the interface.

*/
// print displays the printer's name.
func (c cannon) print() {
	fmt.Printf("Printer Name: %s\n", c.name)
}

// epson defines a epson printer.
type epson struct {
	name string
}

// print displays the printer's name.
func (e *epson) print() {
	fmt.Printf("Printer Name: %s\n", e.name)
}

func main() {

	// Create a cannon and epson printer.
	c := cannon{"PIXMA TR4520"}
	e := epson{"WorkForce Pro WF-3720"}

	// Add the printers to the collection using both
	// value and pointer semantics.

	// This is not a collection of concrete data based on what it is, but it's a collection of
	// interface types based on what it does. Decoupling is based on behavior.
	// See 1.png, in the picture replace "User" with "Printer"
	// Backing array is 2 interface values.
	// At 0th index the value semantics are at play and we are storing a value of type "cannon" and we are
	// making a copy of cannon.
	// Index 1 is leveraging pointer semantic.
	printers := []printer{

		// Store a copy of the cannon printer value.
		c,

		// Store a copy of the epson printer value's address.
		&e,
	}

	// Change the name field for both printers.
	c.name = "PROGRAF PRO-1000"
	e.name = "Home XP-4100"

	// If you look at the 1.png, we will see that the only index in our collection that will see the change
	// is the one using pointer semantics. Since value semantics is operating on it's own copy.

	// Iterate over the slice of printers and call
	// print against the copied interface value.
	// Any time we range over the collection we should ask, should we be using the value semantic
	// or pointer semantics form of the for range?
	// Here we are using the value semantic form of the for range, why?
	/*
		What is the entity?
		It is a collection.

		What is it a collection of?
		Printer

		What is a printer?
		Interface

		What is an interface?
		It's a reference type.

		Ah! And we use value semantics for the reference type.
	*/

	for _, p := range printers {
		p.print()
	}
	// At the first iteration, see 2.png, "e" in the picture is "p" in our code
	// "p" is it's own copy and "p" is also pointing ot he copy, and when we call print against "p"
	// We will just being seeing "Bill".

	/*
		On the next iteration, "p" becomes a copy of index 1 which will say "I am a pointer of type epson"
		And the pointer points to original "epson" and when we call print we will see the changed "p".
		See 3.png
	*/

	// When we store a value, the interface value has its own
	// copy of the value. Changes to the original value will
	// not be seen.

	// When we store a pointer, the interface value has its own
	// copy of the address. Changes to the original value will
	// be seen.
}
