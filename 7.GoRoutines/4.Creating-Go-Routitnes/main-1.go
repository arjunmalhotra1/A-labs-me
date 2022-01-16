/*
	Once we remove the call to wait()
	 // wg.Wait()
	 Now there is no guarantee if those go routines are going to run,
	 when they are going to run, how they are going to run,
	 if they are going to run. We just remove the guarantee point.

	Then when we run it, see 2.png
	We can see that there are no go routines in output, because the go routines
	got terminated even if they even had a chance to run.

	The main go routnine started,completed and terminated the program.

	Again there are no guarantees that this is how the output is going to be anytime we are going to run this program.
	That's because there is call to "fmt.Println("Waiting To Finish")"
	That is a system call. Technically there is an opportunity there for the scheduler to
	make a scheduling decision. Most likely we are not taking it because we are so little into our time slice that
	it let's the go routine continue to run.

	But remember we cannot predict what the scheduler is going do when all things are equal, so we can't
	really predict what's really going to happen here.

	There's another runtime call called, "runtime.Gosched()" see main-2.go



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
