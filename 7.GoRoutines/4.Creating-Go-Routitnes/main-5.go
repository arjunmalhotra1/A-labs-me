/*
	This is like the first program we started with but with some changes.
	Big change we make is in "init()"
	"runtime.GOMAXPROCS(2)"

	What we are doing is that we are running these Go routines in parallel.
	2 go routines executing their instructions at the same time on 2 different hardware threads.

	So we move from a single threaded to a multi threaded go program.
	Same program but we got rid of the named functions and now just have function literals (anonymous functions).

	See 9.png for run results.
	We will see a mixing of the outputs.
	Again we will not be able to predict what these outputs will look like because not only now we
	have Go level randomness but we also have hardware level randomness.

	We have two threads executing the system calls, and between the OS and the hardware, we have mix output
	and we are seeing it running in parallel.

	The waitgroup is our first kind of view around orchestration.
	This helps us create the guarantee point in this case "wg.Wait()".
	The guarantee point says we are not going to move on from  "wg.Wait()" unless these 2 go routines
	complete their work and they weill report us that they completed the work by calling
	"Done()".

	We cannot look at the perceived behavior (our local runs) as a guarantee.
	This is where multi threaded software gets complicated specially when we are dealing with both
	concurrency and parallelism.


*/

// Sample program to show how to create goroutines and
// how the goroutine scheduler behaves with two contexts.
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {

	// Allocate two logical processors for the scheduler to use.
	runtime.GOMAXPROCS(2)
}

func main() {

	// wg is used to wait for the program to finish.
	// Add a count of two, one for each goroutine.
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Start Goroutines")

	// Declare an anonymous function and create a goroutine.
	go func() {

		// Display the alphabet three times.
		for count := 0; count < 3; count++ {
			for r := 'a'; r <= 'z'; r++ {
				fmt.Printf("%c ", r)
			}
		}

		// Tell main we are done.
		wg.Done()
	}()

	// Declare an anonymous function and create a goroutine.
	go func() {

		// Display the alphabet three times.
		for count := 0; count < 3; count++ {
			for r := 'A'; r <= 'Z'; r++ {
				fmt.Printf("%c ", r)
			}
		}

		// Tell main we are done.
		wg.Done()
	}()

	// Wait for the goroutines to finish.
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}
