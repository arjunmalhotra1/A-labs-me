/*
	In the previous case, just knowing that an error exists is enough context, but that is not always the situation.

	What if the function we saw previously "WebCall()" returned more than one error, now we
	need to know was it error A or error B.
	The context to know that with the error value stored inside an error interface is not enough.

	In this code here WebCall can return more than one error.
	When we have a function that can return more than one error, what we should do is use these error
	variables. A lot of times like here, we put the error variables in a "var" block just to keep them together.

	When we use an error variable, we use this "Err" naming convention.
	They will be exported variables and be defined at the top of the code that's using them (or the file/package using them)
	"Err" is the naming convention.

	ErrBadRequest = errors.New("Bad Request")

	ErrPageMoved = errors.New("Page Moved")

	Note we are using errors.New() actual errors package and the factory function to construct them.

	At this line
	"if err := webCall(true); err != nil {"

	webCall() return the error interface value with or without a concrete piece of data stored inside of it,
	the error string. Next we check if there is a concrete value stored inside of err with "err != nil"
	If there is then we are switching. Is the error variable is based on "ErrBadRequest" or "ErrPageMoved".
	Then a default case if it's neither of those 2.

	Errors variable has to be the next choice when just knowing if there's an error in there is not enough.

	We don't want to simply just create error variables just for the sake of it.
	We don't want to be doing anything just for the sake of doing anything. That's just type pollution.
	If error variable is not going to give us any context then we start thinking about the custom error types.

*/

// Sample program to show how to use error variables to help the
// caller determine the exact error being returned.
package main

import (
	"errors"
	"fmt"
)

var (
	// ErrBadRequest is returned when there are problems with the request.
	ErrBadRequest = errors.New("Bad Request")

	// ErrPageMoved is returned when a 301/302 is returned.
	ErrPageMoved = errors.New("Page Moved")
)

func main() {
	if err := webCall(true); err != nil {
		switch err {
		case ErrBadRequest:
			fmt.Println("Bad Request Occurred")
			return

		case ErrPageMoved:
			fmt.Println("The Page moved")
			return

		default:
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Life is good")
}

// webCall performs a web operation.
func webCall(b bool) error {
	if b {
		return ErrBadRequest
	}

	return ErrPageMoved
}
