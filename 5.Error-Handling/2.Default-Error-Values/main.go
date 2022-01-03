/*
	Error handling in about showing the user of the API enough respect that you will give them
	informed piece of information about the state of their app.
	It's about the user having enough context, to understand whether their app is still in the level of
	integrity or not, so that they can make decision about to either correct it or shut down.

	 Systems cannot be developed assuming tat human beings will be able to write millions of lines of code
	 without making mistakes, and debugging alone is not an efficient way to develop reliable systems.
		-Al Aho (inventor of AWS)

*/

// Sample program to show how the default error type is implemented.
package main

import "fmt"

/*
	Error-Interface:
		// http://golang.org/pkg/builtin/#error
		type error interface {
			Error() string
		}
		error interface is defined inside of language, that's why it looks like it is unexported yet we can access it
		because it is built-in

		This is where the code will be located in the standard library -// http://golang.org/pkg/builtin/#error
		But there is no builtin package it is just for documentation.

		The error interface it has one active behavior, and is the core interface that we will be using to
		handle errors throughout Golang.
*/

// http://golang.org/pkg/builtin/#error
type error interface {
	Error() string
}

/*
	This is a concrete type named errorString.
	type errorString struct {
		s string
	}
	errorString comes from the errors package and is the most commonly used error value in Golang.

	Rob Pike - "Errors in Go are just values so they can be anything we would want them to be."
	He is right, they are just values and that error values are decoupled form the interface.

	In a lifetime of code the errors will change quite a bit and by working with errors in a decoupled
	state, it gives us the ability to change and improve error handling without crating cascading affects
	throughout our code base.
	"errorString" is most commonly used error type in go.
	Notice - "errorString" is unexported and also has a field that is unexported too.
*/

// http://golang.org/src/pkg/errors/errors.go
type errorString struct {
	s string
}

/*
	This (below) is the implementation of the error interface. it just returns the copy of string out.
	Implementation of the Error() interface is for two purposes.
	1. To allow that value to be used as an error value throughout our program.
	2. But the implementation of the error interface is also just for logging.

	So we are not expecting somebody to parse the string coming out of this error method to figure out what is
	going on, if they have to parse the string you have failed as a developer here in terms of the API, designing
	the error handling.
	Note the implementation is using pointer semantics.
*/

// http://golang.org/src/pkg/errors/errors.go
func (e *errorString) Error() string {
	return e.s
}

/*
	If we peak a little deeper in the error package we will see a factory function "New()"
	Note that "New()" returns an error interface value.
	Bill said earlier that the factory functions should be returning the concrete data,
	but when we are working with error handling we are always going to be working with the
	error interface.
	We want a level of decoupling all the time when it comes to handling errors so that we can make
	those changes, without the cascading affects.(I didn't  understand this)

	So here's an exception around errors when a factory function is returning interface not the concrete type.
	Note - it is doing the construction -"return &errorString{text}" using the pointer semantics.

	Below code is typical of how we do error handling in Golang.
	if err := webCall(); err != nil {
		fmt.Println(err)
		return
	}

	"if err := webCall(); err != nil "
	webCall() will return an error back, remember that the concrete value is decoupled from the
	error interface.

	Notice here how we are using the factory function to return the error.

	func webCall() error {
		return New("Bad Request")
	}

	Here, call to webCall() is always returning an error ("Bad Request").
	When we make the call to webCall() we get back an error interface value, see 1.png.
	It will say that we have a pointer to the errorString. "*eS"
	And it will point to the "errorString" value which is just a string which in this case
	is going to say, "Bad Request".
	See 1.png this is what is coming out of webCall.
	We get that errorString in "err" when we call "if err := webCall()".
	then we ask the question "err!=nil", what does nil mean?
	What's interesting about "nil" in Go that nil takes on the type it needs.
	Remember two classes of things can be "nil",

	1. Pointer set to their zero value.
	2. Reference types set to their zero value.

	An remember an interface is a reference type. So, an interface set to it's zero value will be nil.
	"err!=nil" is really asking is, "If there is a concrete value stored inside the error interface value."

	Because if err was "nil", or set to it's zero value. It would look like 2.png
	But here we see 1.png for our code. So if anytime we are saying "err!=nil" what we are really asking is,
	"Is there a concrete value, stored inside the error interface."
	If there is a concrete value stored then there is an error.

	We should always be using "err" for error.

	One nice thing baout this line of "if" code is that
	"if err := webCall(); err != nil" the error variable can be maintained within the score of the if statement.
	Which means if we had another call right wafterwards, then we can continue to follow this pattern.

	if err := webCall(); err != nil {
		fmt.Println(err)
		return
	}

	if err := webCall(); err != nil {
		fmt.Println(err)
		return
	}

	Since every if has it's won scope we can re use the "err" variable name.
	One more thing about this error handling is that we can mantain a happy path, by just keeping this line of
	sight. See 3.png

	We can keep the if statements tied to the error handlings like this. If there is an error we handle it and then
	leave, lot of exit points. This helps in code readability
	WE reduced the else clauses, because else clauses are harder to read and we start to put some positive
	paths, in these indents.


	Here bug point is that we are asking if there is a concrete error value stored inside the error interface
	if there eis then we have an error.
	Remember error handling is all about context, in this case our context is simply the fact that
	there is an error inside the error interface that is enough for us to make a decideion.
	Because this function only returns one type of error that is the error string "Bad Request"

*/

// http://golang.org/src/pkg/errors/errors.go
// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

func main() {
	if err := webCall(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Life is good")
}

// webCall performs a web operation.
func webCall() error {
	return New("Bad Request")
}
