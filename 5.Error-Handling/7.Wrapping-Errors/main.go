/*
	We just saw the idea that errors have to provide context, we talked about the error itself being the context,
	we talked about the error variables being the context, we talked about using custom error
	types where the type itself is the context even better when we can use the interface type and can make
	behavior the context.
	Then we saw the bug where we saw that, we have to make sure that all our functions return the
	error interface for hte error using the interface value and not the custom error types themselves.

	Here we will see the design pattern that will really help  minimize the problems in code.
	That's the idea of handling an error.

	Handling an error mean 3 things:
	1. Errors is going to stop there, the code that's handling it. That is the error doesn't get
		propagated any more. either we are going to recover or shut down.
	2. The error will be logged. The error has to be logged with the full context. At that point
	we will also have to make sure that the code we are deciding is doing all this, can log an error.
	When we get into package oriented design we will talk about policies who can and who cannot
	log. Who can shut down apps and who cannot shut down apps.
	3. Then we will make the decsion whether can recover or shut down.

	There are a lot of inconsistencies in code bases (codes-out there in the world)
	when it comes to handling an error.

	Dave Trainey wrote his own error package and we will see how to use it.

	Here, we have defined a custom error type called "AppError". It has a state field.
	type AppError struct {
		State int
	}
	We are using the pointer semantics to implement error interface
	"func (c *AppError) Error() string {"

	The main function is making function calls to firstCall which calls SecondCall and
	which calls thirdCall.

	We see that the thirdCall() is failing and is retuning the concrete error value
	AppError with a value of 99.
	"&AppError{99}"
	1.png
	We can imaging that in 1.png is the AppError value which has 99 in it.
	This is the root error value. this function is obviously is not going to handle the error
	because this function is producing the error nad is going to return it back up the call stack.

	What happens next?
	The thirdCall() function as being executed by secondCall() function.
	Then when we check "err!=nil", the answer is yes. When we called the thirdCall()
	it didn't bring back the concrete value, it brough it back through the "err" interface.
	So we know since there is a concrete value inside it hence it is "true".

	But now are we going to handle it?
	Remember, Handling means:
	1. It stops here.
	2. Then we are going to log it.
	3. Then make a decision on recovery or shut down.

	but this code - secondCall() - is saying I am not going to handle the error.
	We are going to decide we are or we are not, and if we decide we are not going to handle the error,
	then we have one choice that is to wrap all the context we need and pass it up the call stack.
	What we are trying to do is, the idea of logging the same error over an over again in the logs
	just to get that callstack information in another context, this just pollutes the log and is something we
	would want to avoid.

	The whole idea of the wrap function that we get from Dave Chainey's error package
	is that we are going to wrap this error, around another error.
	2.png
	Then we will have a callcontext(we get this for free, we'll know what line of code we are on
	(all the good stuff)), see 3.png.

	We will have our "UserContext" (see 4.png) this "secondCall->thirdCall()"
	"return errors.Wrap(err, "secondCall->thirdCall()")"
	Now when there is an error and we have this stuff in our logs we should be able to
	fix the problem between the mental model of our code base, between this logging and our context.

	We wrap the error in bunch of stuff and then we return that error back up the call stack.
	Here we just returned the error back up the callstack to the firstCall().

	Then we go through the same pattern.
	"firstCall->secondCall(%d)"
	Then in the first call we check if the second gave an error, it did.
	We did get a concrete value stored inside the error interface and now we do another
	Wrap.

	"return errors.Wrapf(err, "firstCall->secondCall(%d)", i)"
	In this case we are wrapping the parameter(i) the value that we passed into the secondCall().
	So one more wrap on the error. More context, both call and user context.
	See 5.png

	Now somebody has to eventually handle the error.
	We see that in main(). We get back the error value.
	Was there a concrete value stored inside the error interface, after the call to firstCall()?
	Yes, there is.

	What's nice is that Dave added a function called "Cause()", cause function.
	What cause function does is that it allows us to unwind all the wrapping.
	And get the error value back out.
	See 6.png
	So this cause function is giving us back the error value that we started with from the begining.
	This then let's us use all those mechanics, that we saw earlier.

	As we can see here we get the root cause error using the cause function and then we do the type assertion.
	switch v := errors.Cause(err).(type) {
	case *AppError:

	Then we can ask was the root cause an "AppError" then we can do all the mechanics we saw earlier.
	"type as context", "behavior as context".

	This whole wrapped error value, in 5.png does implement the ability to do some custom
	formatting using %v and %+v formatter on both "fmt" and "log" this will gives us the stack trace.
	Now, if we use "+v" we are going to get the full stack trace, from every wrapping position.
	If we use only "%v" we will get custom context.
	This is something we will be doing in development to validate all the error handling/logging is working.

	The output of the program, in 7.png.
	We got the full stack trace, and notice we have the stack trace from those different wrapping calls.
	8.png and 9.png(see we passed the value of "10")
	So we have the stack trace from different point of views. If we don't want al that detail,
	we can just pull out the user level context 10.png and still see everything.

	This is a brilliant pattern which is going to keep your code simple
	because now as a developer on the team we basically say
	"Hey is there an error? Yes there is an error"
	"Am I handling the error? No I am not handling the error"
	then we just wrap and bring that error back up and code will be very clean

	And if
	"Hey ! I am hadling the error, no problem now we have to log it."
	"Now we have to make the decision, I know the error stops with me. Am I going to be able to recover or
	would we have to shut down?"

	And in the code review we can also validate that this philosophy and this guidelines are being followed.
	and then during developing and during unit testing, any testing we do, we should validate that
	we have enough context in the logs and if not we can enhance those string in the wrap calls to get
	everythng that we need.








*/

// Sample program to show how wrapping errors work with pkg/errors.
package main

import (
	"fmt"

	"github.com/pkg/errors"
)

// AppError represents a custom error type.
type AppError struct {
	State int
}

// Error implements the error interface.
func (c *AppError) Error() string {
	return fmt.Sprintf("App Error, State: %d", c.State)
}

func main() {

	// Make the function call and validate the error.
	if err := firstCall(10); err != nil {

		// Use type as context to determine cause.
		switch v := errors.Cause(err).(type) {
		case *AppError:

			// We got our custom error type.
			fmt.Println("Custom App Error:", v.State)

		default:

			// We did not get any specific error type.
			fmt.Println("Default Error")
		}

		// Display the stack trace for the error.
		fmt.Println("\nStack Trace\n********************************")
		fmt.Printf("%+v\n", err)
		fmt.Println("\nNo Trace\n********************************")
		fmt.Printf("%v\n", err)
	}
}

// firstCall makes a call to a second function and wraps any error.
func firstCall(i int) error {
	if err := secondCall(i); err != nil {
		return errors.Wrapf(err, "firstCall->secondCall(%d)", i)
	}
	return nil
}

// secondCall makes a call to a third function and wraps any error.
func secondCall(i int) error {
	if err := thirdCall(); err != nil {
		return errors.Wrap(err, "secondCall->thirdCall()")
	}
	return nil
}

// thirdCall create an error value we will validate.
func thirdCall() error {
	return &AppError{99}
}
