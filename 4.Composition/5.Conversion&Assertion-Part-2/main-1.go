/*
	We have a car, we are implementing method "String()" using value semantics of "car".
	Question. With "func (car) String() string" how come there is no variable defined with the
	receiver?
	Answer. That's because in this case we are not using the receiver variable for anything
	the compiler says we don't have to declare that variable.
	Value of type car also now implements behavior named string.
	There is an interface in both "log" and the "fmt" package named "Stringer", here
	we are implementing the "Stringer" interface which now means that we can pass a
	value or the address of a car value to "fmt" or "log" and override the logging behavior.

	We have done the same with Cloud as well.
	so now we have two pieces of concrete data, each exhibiting the string behavior.
	Polymorphism means a piece of code changes it's behavior depending upon the concrete data
	it is operating on.

	The method "String()" is giving us the data the behavior really focussed
	around the idea of "Polymorphism". And "Polymorphism" is giving us those ideas around that thin
	layer of precise decoupling.

	Here,
	mvs := []fmt.Stringer{
		car{},
		cloud{},
	}
	We are generating a collection of concrete data not based in what it is, but based on what it does.
	In this case we are storing two pieces of concrete data that know, how to call String() .
	They know how to display themselves in a custom way.

	Then we run a loop picking a random number, and depending on the index that  number can represent randomly,
	we do the assertion.
	"mvs[rn].(cloud)"
	We are saying "If the value we are going to type assert against, is a cloud", so if index 0 or 1
	if there is a cloud, then we print "Got Lucky".

	We can notice that the type assertions happen in runtime.
	Even though at compile time we are validating that we are still working with the right type,
	this is still a run time check and not a compile time check.
	See main-2.go for another example.

*/
// Sample program demonstrating that type assertions are a runtime and
// not compile time construct.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// car represents something you drive.
type car struct{}

// String implements the fmt.Stringer interface.
func (car) String() string {
	return "Vroom!"
}

// cloud represents somewhere you store information.
type cloud struct{}

// String implements the fmt.Stringer interface.
func (cloud) String() string {
	return "Big Data!"
}

func main() {

	// Seed the number random generator.
	rand.Seed(time.Now().UnixNano())

	// Create a slice of the Stringer interface values.
	mvs := []fmt.Stringer{
		car{},
		cloud{},
	}

	// Let's run this experiment ten times.
	for i := 0; i < 10; i++ {

		// Choose a random number from 0 to 1.
		rn := rand.Intn(2)

		// Perform a type assertion that we have a concrete type
		// of cloud in the interface value we randomly chose.
		if v, is := mvs[rn].(cloud); is {
			fmt.Println("Got Lucky:", v)
			continue
		}

		fmt.Println("Got Unlucky")
	}
}
