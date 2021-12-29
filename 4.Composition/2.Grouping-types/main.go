/*
 "A good API is not just easy to use but also hard to misuse" - Jaana Dogan/JBD

"You can always embed, but you cannot decompose big interfaces once they are our there.
Keep interfaces small." - JBD/Jaana Dogan.

"Don't design with interfaces, discover them." - Rob Pike.

"Duplication is far cheaper than the wrong abstraction" - Sandi Metz
*/

// This is an example of using type hierarchies with a OOP pattern.
// This is not something we want to do in Go. Go does not have the
// concept of sub-typing. All types are their own and the concepts of
// base and derived types do not exist in Go. This pattern does not
// provide a good design principle in a Go program.
package main

import "fmt"

// Animal contains all the base fields for animals.
type Animal struct {
	Name     string
	IsMammal bool
}

// Speak provides generic behavior for all animals and
// how they speak.
func (a *Animal) Speak() {
	fmt.Printf(
		"UGH! My name is %s, it is %t I am a mammal\n",
		a.Name,
		a.IsMammal,
	)
}

// Dog contains everything an Animal is but specific
// attributes that only a Dog has.
type Dog struct {
	Animal
	PackFactor int
}

// Speak knows how to speak like a dog.
func (d *Dog) Speak() {
	fmt.Printf(
		"Woof! My name is %s, it is %t I am a mammal with a pack factor of %d.\n",
		d.Name,
		d.IsMammal,
		d.PackFactor,
	)
}

// Cat contains everything an Animal is but specific
// attributes that only a Cat has.
type Cat struct {
	Animal
	ClimbFactor int
}

// Speak knows how to speak like a cat.
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
	animals := []Animal{

		// Create a Dog by initializing its Animal parts
		// and then its specific Dog attributes.
		/*
			Compiler thrown an error - "That we cannot add Dogs and Cats in this collection of Animals."
			What is wrong here?

			What's happening is that the developer doing this, is used to the idea of sub typing.
			The idea that because a dog and a cat has an Animal embedded inside of it.
			It makes them an animal. Therefore we can group them by the fact that they are an
			animal.
			****
			Remember we talked about during our "Embedding" that "Embedding" does not create
			a sub typing relationship.
			That is Dogs are still Dogs. Cats are still Cats. Animals are not Dogs and Cats at all.

			When we are dealing with decoupling.
			When we are dealing with multiple types of data, we have to move from "what it is" down back to
			"what it does".

			What does a Cat and a Dog do, that's how we really want to group things together.
			They speak.

			Such a design we see in this code raises flags because the developer is focussing on
			"what it is" instead of "what it does".

			There are other smells in this code:

			1. Animal type is being used a as a reusable state. Most of the time this is a big smell.
			He didn't say the reason behind this.
			2. Other big smell we have is that the implementation of speak() is pretty useless.
			3. the other idea is that we want to have the types that truly represent something new and unique.
			Here, nt sure if Animal represent something new and unique because there's not part of the code that
			would need an Animal value. That's a pretty good indication as well that this a type
			that is kind of polluting our type system.

			See the end of this file for what do we tell this developer about.
			****
		*/
		Dog{
			Animal: Animal{
				Name:     "Fido",
				IsMammal: true,
			},
			PackFactor: 5,
		},

		// Create a Cat by initializing its Animal parts
		// and then its specific Cat attributes.
		Cat{
			Animal: Animal{
				Name:     "Milo",
				IsMammal: true,
			},
			ClimbFactor: 4,
		},
	}

	// Have the Animals speak.
	for _, animal := range animals {
		animal.Speak()
	}
}

// =============================================================================

// NOTES:

/* Smells:
// 	* The Animal type is providing an abstraction layer of reusable state.
		** Because when we embed we do not want to embed for state. We want to embed for behavior.
		As the decoupling comes from the behavior.
// 	* The program never needs to create or solely use a value of type Animal.
// 	* The implementation of the Speak method for the Animal type is a generalization.
		** The speak method is not going to be used but we will have to write a test for the
		test coverage which is a tola waste.
// 	* The Speak method for the Animal type is never going to be called.
*/

/*
	What do we tell this developer? How do they fix this code?
	How do we fix the mindset for, what's the better way to do this in Go?

	Answer. It's always going to be that focus on the behavior.
	So what we will tell the developer is, to stop thinking about what "Animals are."
	"What Dogs and Cats are and start thinking about what Dogs and Cats do."
	"Once we identify what cats and dogs does then we can go and define those interfaces."

	One of the things Bill looks for in the code review is that interfaces always describe the
	behavior things and not what things are.
	Our concrete types should describe those nouns, those persons, places things.
	Interfaces should describe the behavior those nouns those person, places and things can do.

	If interfaces are describing things that's also going to be a bad smell.
	"Don't design with interfaces, discover them." - Rob Pike.

	Because if we are designing around a behavior and if we don't have concrete implementation of something
	then we are guessing.

	We already have the concrete implementation of Dog and Cat.
	We don't have to guess about their behavior.
	It's there, it's called speak. Now notice "main-1.go" where we add the speaker interface with active
	behavior speak(). We now have that behavior well defined.

	Now through the interface, the valueless type we can start doing decoupling.


*/
