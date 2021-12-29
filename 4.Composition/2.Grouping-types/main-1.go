/*
	Now we get rid of the Animal type altogether.
	"Duplication is far cheaper than the wrong abstraction" - Sandi Metz
	As per Bill, Animal type is the wrong abstraction. It is not addign any value.
	So following what "Sandi Metz" says we'll copy "IsMammal" and "Name" fields inside the Dog
	type.
	Dog already implements the speaker interface.
	same thing with cat, we add two fields to cat too.
	But what if we want to add new field?
	Answer. Then we add another field.
	What if we forget to add the field in the "cat" type?
	Answer. Well then that wasn't important to be there and now the types are
	much more precise around what they need to be.
	Cat also implements the speaker interface.
*/

// This is an example of using composition and interfaces. This is
// something we want to do in Go. We will group common types by
// their behavior and not by their state. This pattern does
// provide a good design principle in a Go program.
package main

import "fmt"

// Speaker provide a common behavior for all concrete types
// to follow if they want to be a part of this group. This
// is a contract for these concrete types to follow.
type Speaker interface {
	Speak()
}

// Dog contains everything a Dog needs.
type Dog struct {
	Name       string
	IsMammal   bool
	PackFactor int
}

// Speak knows how to speak like a dog.
// This makes a Dog now part of a group of concrete
// types that know how to speak.
func (d *Dog) Speak() {
	fmt.Printf(
		"Woof! My name is %s, it is %t I am a mammal with a pack factor of %d.\n",
		d.Name,
		d.IsMammal,
		d.PackFactor,
	)
}

// Cat contains everything a Cat needs.
type Cat struct {
	Name        string
	IsMammal    bool
	ClimbFactor int
}

// Speak knows how to speak like a cat.
// This makes a Cat now part of a group of concrete
// types that know how to speak.
func (c *Cat) Speak() {
	fmt.Printf(
		"Meow! My name is %s, it is %t I am a mammal with a climb factor of %d.\n",
		c.Name,
		c.IsMammal,
		c.ClimbFactor,
	)
}

func main() {

	// Create a list of Animals that know how to speak.
	/*
		Now when we make a collection of dogs and cats, we don't do it on the basis of what we think they are,
		like animals. But we do it on the fact that they can speak.
		this really lends itself a lot of flexibility in the code.

		The developers defining these concrete types don't have to think about the future use,
		they just define the features they need.

		Now we group dogs and cats not by what they are but on basis of what they do.

	*/
	speakers := []Speaker{

		// Create a Dog by initializing its Animal parts
		// and then its specific Dog attributes.
		&Dog{
			Name:       "Fido",
			IsMammal:   true,
			PackFactor: 5,
		},

		// Create a Cat by initializing its Animal parts
		// and then its specific Cat attributes.
		&Cat{
			Name:        "Milo",
			IsMammal:    true,
			ClimbFactor: 4,
		},
	}

	// Have the Animals speak.
	for _, spkr := range speakers {
		spkr.Speak()
	}
}

// =============================================================================
/*
	Declare types that represent something new or unique.
	Around this idea we see a lot of developers do,

	type Handle int
	func foo(h Handle)

	We are doing this wrong.
	Because if Bill asks us what a "handle" is, then we would say that
	handle represents an integer value.

	If it's just an integer then why do we need a type?
	Answer. It gives us compiler protection.

	First of all it's not the handle that represents anything new it is just an integer.
	And second, we are not going to get any compiler protection, because
	we can always call foo like this, "foo(10)".

	Because 10 is a constant of kind int and is compatible with handle the compiler will do the
	implicit conversion.

	The correct way to do this is,
	"type Handle int
	func foo(handle int)
	foo(10)"

	So anytime when Bill asks us what that type is and if we answer back with a name of the base type.
	We really don't a type that is new or unique.
	Remember we don't get any compiler protection because we can convert those built in types based
	on a kind.
*/

// =============================================================================

/*
	Here is a type which is unique and new.
	We would find it in the time package called "Duration"

	type Duration int64

	If we asked Bill, what is Duration, then Bill would say that represents a nano second of time.
	Note we are not saying the "base" type to tell us what "Duration" represents.
	It truly does represent something new or unique.

	Much better than handle.

	Another way we can valdate whether or not it is a good type, new or unique, is if it has a method set.
	Most of the time when we see types declared off of a base, we want to see a method set or
	behavior that's reasonable and practical.

	Animal wasn't fitting in any of that so we got rid of it.
	We want to validate that if we do define a type we need values of that type throughout the scope of
	program somewhere. If not that is pollution.

*/

// =============================================================================
// NOTES:

/* Here are some guidelines around declaring types:
// 	* Declare types that represent something new or unique.
		** The animal type wasn't describing anything new or unique.
// 	* Validate that a value of any type is created or used on its own.
// 	* Embed types to reuse existing behaviors you need to satisfy.
		** We always need embedding types for behavior and not for reusable state.
		We would need embedding to share that behavior.
// 	* Question types that are an alias or abstraction for an existing type.
// 	* Question types whose sole purpose is to share common state.

Remember there are exception to every rule.

Composition is about different types of data together in collectiona dn groups.
*/
