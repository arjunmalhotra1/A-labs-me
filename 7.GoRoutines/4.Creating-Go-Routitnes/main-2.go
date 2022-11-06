/*
	Bill says we should not be using this "runtime.GoSched()" function call. Other than in
	tests because what this does is it actually makes the call to Yield().
	Remember scheduler is making these calls during those safe points, around the function transitions.

	But Go does give us the option to yield() ourselves but there is no guarantee that things will
	yield() there's no guarantee on order.

	It's good for testing if we want to create chaos but not good to deal with orchestration and guarantees
	because there are none.

	When we run this program, see 3.png as output. We will see that it looks like that the program worked again. This is
	a total lie. That's why we have problems in multi threaded software.
	We see this output nad we think it works, may be it works on our machine all the time, but when we
	move it to the production environment everything will be behaving differently because
	there are no guarantees. Don't get tricked on perceived behaviors.

	Next see main-3.go. Without "runtime.GoSched()".
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
	// wg.Wait()
	runtime.GoSched()

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
