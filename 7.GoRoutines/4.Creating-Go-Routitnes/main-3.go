/*
	Question. What happens when we forget to call "wg.Done" after "lowercase()"?
	Answer. Now we not decrementing the waitgroup.
	As in the output see 4.png, we now have a deadlock situation.

	What's happening here is that both Go routines got to run but only one decremented that wait group
	from 2 to 1.
	The decrement going from 1 to 0 doesn't happen.

	And the run time is able to identify that on "wg.Wait()" we will never be able to move on.
	So now we have some deadlock detection. Essentially every single Go routine,
	has to be stuck in a waiting state with no opportunity for it to move on.
	then the Go scheduler can identify that we have this deadlock.

	So if you have go routines in timer loops, which Bill doesn't recommend and tends to avoid
	timer loops in specially long runing applications. With timer loops you are walking away from the
	deadlock detection which we wouldn't want to lose it.

	Say we only do Add() for only one Go routine. It's really not predictable what the program is going to do
	because it's got concurrency or out of order execution.
	So if we do wg.Add(1).
	See 5.png for output on bill's machine.
	Only uppercase() got to run decremented the waitgroup to 0. The scheduler instead of giving the other Go routine
	that's had no time, a chance to run. The scheduler went back to the main Go routine that was runing and allowed
	it to finish the program.
	Nothing is predictable.

	It is important that essentially we think that every Go routine that's in a
	runnable state is running at the same time.

	If you don't do these things you will make really bad engineering decisions
	because you are not going to place those guarantees in the code, like we don't have guarantee
	for go routines when we do, wg.Add(1) instead of wg.Add(2).

	---------------------------------------------------------------------------------------------------------

	"A lot of Go developers say why are you using a wait group you should be using a channel."

	Bill - "We are always going to go back to this idea of complexity."
	Wait groups are a lot simpler than channels and since we don't need any data back from these Go routines,
	we just need to track them, waitgroup is the right orchestration here.

	Don't walk away from orchestration or synchronization because you think they are too basic or too simple.

	One think we should try to do is, keep "Add()" and "Done" in the same line of sight.
	Like we have here. It is going to do reduce a large number of bugs. We can't always do this but it's nice if we
	can.
	6.png

	For example-2 - we will discuss the idea of time slice. See main-4.go



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
