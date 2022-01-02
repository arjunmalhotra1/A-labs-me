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
