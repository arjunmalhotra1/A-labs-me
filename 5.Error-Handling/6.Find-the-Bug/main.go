/*
	This code has amistake that people writing code since a long time have made.

	"if _, err = fail(); err != nil {"
	We call fail() and ask if there is any concrete value sotred inside of err.
	At runtime the output is
	"log.Fatal("Why did this fail?")"
	which means there is a concrete value stored inside of "err".

	What's gong on we return "nil", we return on failure. Yet from the application point of view,
	there is a concrete value in there.

	What's happening?
	When we call "var err error" we are constructing an error interface value set to it's zero value. 1.png

	There is no doubt that "nil" is being returned as a value.
	But the interesting part is,
	"func fail() ([]byte, *customError)"

	The return type for fail is "*customError". This is the bug.
	We are using the concrete value directly in the signature of the function.
	Which means we are getting back "nil" but we are getting back "nil" as a pointer of
	customError type. "*customError", see 3.png.
	Therefore there really is a concrete value stored inside. Because we got the zero value for pointer.

	The bug is that anytime we are writing code, and we are going to be working with errors,
	we are not using the custom error types.

	We do not use - "func fail() ([]byte, *customError) {"
	we use the error interface.
	"func fail() ([]byte, error) {"

	It means now when we returned "nil" for the error from the "fail()".
	Both the value and the pointer will actually be "nil".
	Because we are returning the zero value for the error interface.

	Remember "nil" takes on the type it needs, in this context "nil" now is a zero valued
	interface value.
	"nil" now represents a nil pointer.
	2 different values of nil, 2 different types which leads to this bug.




*/

// Sample program to show see if the class can find the bug.
package main

import "log"

// customError is just an empty struct.
type customError struct{}

// Error implements the error interface.
func (c *customError) Error() string {
	return "Find the bug."
}

// fail returns nil values for both return types.
//func fail() ([]byte, *customError) {
func fail() ([]byte, error) {
	return nil, nil
}

func main() {
	var err error
	if _, err = fail(); err != nil {
		log.Fatal("Why did this fail?")
	}

	log.Println("No Error")
}
