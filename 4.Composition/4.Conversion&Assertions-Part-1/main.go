/*
	We have 3 interfaces, Mover, Locker and MoveLocker.
	Mover interface has the "Move()" behavior.
	Loker interface has 2 behaviors, Lock() and Unlock().
	Then MoveLocker is an interface which is the composition of the other 2 interfaces.
	Which means that any piece of data, any value or any pointer that implements both
	Move(), Lock() and Unlock() can beMoveLocker.

	"bike" is a concrete data, we are using an empty struct. That implements the all 3 behaviors.
	"bike" knows how to move, how to lock and how to unlock.

	Now in main we have "MoveLocker" and "Mover" set to their zero value.
	See 1.png
	Then we construct a bike value and assign it to the MoveLocker.
	WE can do that because we know that bike implements the methods move, lock and unlock.
	2.png
	Then we do
	"m = ml"
	We can at compile time tell the compiler that I would live to assign,
	"ml" to "m".
	Remember we are not assigning "ml" (we cannot move interface values around our program).
	Only thing we can move around our program is concrete data.
	So what we are really saying is, can we assign or copy a bike value
	(which is currently stored inside of ml) into "m"?

	Question. Can we assign the bike value into "m", one that's currently stored in ml?
	Answer. Compiler will go off and look at the type information.
	We know that to be a "MoveLocker" we already need to be able to move.
	"MoveLocker" interface is the composition, of the Move(), Unlock() and Lock() behaviors.
	So compiler sees, "m = ml" ans says that's not a problem. Since "MoveLocker" has Move() it is
	of type "Mover". "Embedding - main-1.go"
	Remember when embedding,
	"Inner type promotion says that anything related to the inner type promotes to the outer type.
	All the fields and methods associated with the inner type now become part of the outer type."

	So now "m"'s pointer is also pointing to same "Bike" that ml is storing.
	"m" can store the same bike value that "ml" is. We use the pointers for efficiency for sharing.
	See 3.png

	When we reach to this line "ml = m".
	Question. Can we take the concrete value inside of "m" and assign it to "ml"?
	Can we go the other way?
	Answer. Compiler says we "Can't" go the other way. Because when it comes to mover,
	all we know is that "mover" only has the ability to Move().
	And if we want to store something inside "MoveLocker" you got to have the ability to "move()",
	"lock()" and "unlock()". You got to have all these 3 behaviors.
	In this particular case, all we know about is that any concrete piece of data stored inside of "m"
	only has the ability to "Move()".

	Here we are getting compile level integrity and the composition is going to make it
	easy for the compiler but also for us to understand, when moving the concrete data around or
	sharing it what we can and cannot do with the level of integrity.

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
