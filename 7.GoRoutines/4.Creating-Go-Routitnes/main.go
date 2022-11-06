/*
	Synchronization means that the Go routines are waiting in line to take a turn.
	Orchestration is when there's going to be an interaction between 2 or more Go routines.

	Say you got order coffee and you want to talk to the server behind the register,
	but you can't because you have people waiting in front of you, you have to wait for your turn
	that's synchronization problem. There's some shared state.
	In this case the person behind the register can only talk to one person at a time. So you are stuck in line.

	Then when we get to the front of the line we have Orchestration problem now you have to talk,
	you got to pass money, that is an Orchestration problem.

	The worst thing we can do is when we have an orchestration problem but we leverage synchronization
	primitives.
	Worst case, if you have a synchronization problem, and we use orchestration primitives like,
	Channels.

	runtime.GOMAXPROCS(1)
	We are downgrading the number of threads that are available to us to execute Go routines on.
	This is going to be a single threaded Go program. One "P", one "M" to leverage one
	hardware thread on my machine. So we are going down tot he single threaded go program.

	----------------------------------------------------------------

	var wg sync.WaitGroup
	wg.Add(2)

	We are constructing a wait group setting it to it's zero value. Wait groups are usable in their zero value
	state.
	Wait Group is a synchronous counting semaphore. It let us maintain a count of number of Go routines at any given time.

	It has APIs of 'Add()', 'Done()' and 'Wait()'.

	WE should be calling Add() function just once. We shouldn't be saying
	Add(1) multiple times, that is a bad code smell.

	We should be knowing how many go routines we'll be creating before we create them.
	There's a basic rule of thumb,
	"You can't create a Go Routine unless you know when and how it's going to terminate."

	wg.Add(2) - means we will have 2 go routines and we will want to track them.
	wg.Done() decrement the number - 2
	Then
	wg.Wait() - Blocks and waits for the all Go routines to finish. This is the guarantee point.

	In multi threaded code we should always be talking about "guarantee point".
	Because if there is no guarantee somewhere then we end up with chaos.

	---------------------------------------------------------------

	// Create a goroutine from the lowercase function.
	go func() {
		lowercase()
		wg.Done()
	}()

	Here we are constructing a literal function - a function with no name.
	We are also making the function call "()". Then we are using the "go" keyword.
	This go routine is calling a named function called "lowercase()" and then calls
	wg.Done() to decrement the wait group.

	So we are declaring, calling and launching this function as a separate go routine.

	go func() {
		uppercase()
		wg.Done()
	}()


	We have 3 go routines associated with this program.
	1. Main Go routine, created by runtime executing this whole code,  and
	2. 2 Go routines that we construct ourselves.

	wg.Wait() - We hold the Go routines to wait, what wait is going todo is block the main
	go routine from causing the program to terminate.

	Basically when main program returns your program terminates. Any go routines that were executing
	instructions could terminate. Hence we wouldn't want our program to just terminate.
	We want to make sure that all the Go routines that we created are have reported that they are done.
	So we can have the clean startup and clean shutdown.

	So here we have 2 go routines on a single threaded Go program.
	See 1.png for the output.
	We do see concurrency here - out of order execution. The Go routine that was created second and which was
	calling "Uppercase" actually ran first.
	It didn't require a lot of time to complete it's work so it got to finish it's full time slice
	before the other Go routine finished as well.

	We have a guarantee point with "wg.Wait()" because we are not going to allow this program to terminate
	till both Go routines finish their work.

	Remember there is no predictability with this program. Every time it runs you get a different output.
	We juSt launched two go routines the order didn't matter.

	IF ORDER MATTERS THEN YOU SHOULDN'T BE USING CONCURRENCY.

	But,
	GUARANTEES MATTER. THAT IS WHY WE HAVE "wait()".
	WE need to guarantee that both functions complete the order doesn't matter and that's why we have the wait().

	If we remove - "wg.Wait()" see main-1.go



}

*/

// Sample program to show how to create goroutines and
// how the scheduler behaves.
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {

	// Allocate one logical processor for the scheduler to use.
	/*
		We are downgrading the number of threads that are available to us to execute Go routines on.
		This is going to be a single threaded Go program. One "P", one "M" to leverage one
		hardware thread on my machine. So we are going down tot he single threaded go program.

	*/
	runtime.GOMAXPROCS(1)
}

func main() {

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Start Goroutines")

	// Create a goroutine from the lowercase function.
	go func() {
		lowercase()
		wg.Done()
	}()

	// Create a goroutine from the uppercase function.
	go func() {
		uppercase()
		wg.Done()
	}()

	// Wait for the goroutines to finish.
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}

// lowercase displays the set of lowercase letters three times.
func lowercase() {

	// Display the alphabet three times
	for count := 0; count < 3; count++ {
		for r := 'a'; r <= 'z'; r++ {
			fmt.Printf("%c ", r)
		}
	}
}

// uppercase displays the set of uppercase letters three times.
func uppercase() {

	// Display the alphabet three times
	for count := 0; count < 3; count++ {
		for r := 'A'; r <= 'Z'; r++ {
			fmt.Printf("%c ", r)
		}
	}
}
