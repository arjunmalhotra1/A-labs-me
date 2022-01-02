/*
	We looked in the previous example, implicit type conversion.
	Imagine if we it wasn't the assignments (pointer assignments) that we cared about, what we really wanted to do was,
	get the copy of the data stored inside the interface.
	This is when we start talking about "Assertions".

	"b := m.(bike)"
	This is an assertion. We have an interface value "m""." and then in paranthesis we are asking
	one question, we are saying, "Hey! can you give me a copy
	(remember we are moving these copies across program boundaries) of bike value stored inside
	of 'm'?"
	If there is a bike stored inside of "m" (which is the case here), a new variable "b" will be a copy
	of that bike value that was stored in "m". 1.png

	But, what happens if there isn't a bike value inside of "m"?
	What happens if what we have is a "Car". 2.png
	Then this code will panic and say "You have got an integrity issue here.
	Because you think there is a bike in there but there isn't."
	But what if we didn't know if it was going to be a bike or a car.
	Then, if we are not sure then we have a second form of assertion.
	"b,ok := m.(bike)"

	Now if there is Bike, which we know there is see 3.png (he changed Car from 2.png
	back to Bike in 3.png).
	Now this code will be fine and "ok" would be true, we will get our copy of bike in "b".
	But if there was a Car, in "m" we would end up with "false" for "ok" and "b" would be
	set with zero value for whatever the type of "b" is.

	Now after assertion we can, once we get the copy of data out in "b". We could
	go back and now do the assignment.

	"b := m.(bike)
	ml = b"

	Sometimes people call this boxing and unboxing of a value inside of an interface.
	But since we know at compile time that "b" has all 3 behaviors and "bike" has all 3
	behaviors then that assignment would be fine at compile time.

	Let's look at more examples in
	main-1.go
	We can see that type assertion is runtime and not compile time operation.
	But still there is some integrity that compiler is guaranteeing that if we are
	working with the right data types, everything is going to be clean.
*/

// Sample program demonstrating when implicit interface conversions
// are provided by the compiler.
package main

import "fmt"

// Mover provides support for moving things.
type Mover interface {
	Move()
}

// Locker provides support for locking and unlocking things.
type Locker interface {
	Lock()
	Unlock()
}

// MoveLocker provides support for moving and locking things.
type MoveLocker interface {
	Mover
	Locker
}

// bike represents a concrete type for the example.
type bike struct{}

// Move can change the position of a bike.
func (bike) Move() {
	fmt.Println("Moving the bike")
}

// Lock prevents a bike from moving.
func (bike) Lock() {
	fmt.Println("Locking the bike")
}

// Unlock allows a bike to be moved.
func (bike) Unlock() {
	fmt.Println("Unlocking the bike")
}

func main() {

	// Declare variables of the MoveLocker and Mover interfaces set to their
	// zero value.
	var ml MoveLocker
	var m Mover

	// Create a value of type bike and assign the value to the MoveLocker
	// interface value.
	ml = bike{}

	// An interface value of type MoveLocker can be implicitly converted into
	// a value of type Mover. They both declare a method named move.
	m = ml

	// prog.go:68: cannot use m (type Mover) as type MoveLocker in assignment:
	//	   Mover does not implement MoveLocker (missing Lock method)
	ml = m

	// Interface type Mover does not declare methods named lock and unlock.
	// Therefore, the compiler can't perform an implicit conversion to assign
	// a value of interface type Mover to an interface value of type MoveLocker.
	// It is irrelevant that the concrete type value of type bike that is stored
	// inside of the Mover interface value implements the MoveLocker interface.

	// We can perform a type assertion at runtime to support the assignment.

	// Perform a type assertion against the Mover interface value to access
	// a COPY of the concrete type value of type bike that was stored inside
	// of it. Then assign the COPY of the concrete type to the MoveLocker
	// interface.
	b := m.(bike)
	ml = b

	// It's important to note that the type assertion syntax provides a way
	// to state what type of value is stored inside the interface. This is
	// more powerful from a language and readability standpoint, than using
	// a casting syntax, like in other languages.
}
