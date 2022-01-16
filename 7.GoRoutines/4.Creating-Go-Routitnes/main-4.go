/*
	The idea of this program is to show that Go will time slice the Go routines that have a little bit more work to do.
	This will again be a single threaded Go program. We will have our waitgroup to help with our orchestration.

	In this example the Go routines will be printing hashes.
	The Go routines will be printing 50,000 hashes.

	This line will provide us with an opportunity of context switch.
	"fmt.Printf("%s: %05d: %x\n", prefix, i, sum)"
	When we make the call after a certain time slice has been leveraged.

	The idea is to see how many times we will see a context switch between "A" go routine & the "B" go routine.

	See 7.png for the output and we can see that we have 8 context switches.
	On the second run see 8.png, 10 context switches.

	Everytime we run the program we see a different number of context switches. And ths' the idea that we can't
	// really predict what the output is going to be, because the scheduler truly looks and feel
	preemptive.
	Even though it is a cooperative scheduler and we are caling yield() we are still based on events that are
	occurring may be application level events like the system call to fmt.Println() but they are still
	based on events so it is very very unpredictable.

	See main-5.go for one more program.



*/
// This means we will run example2 and we will pipe that output and cut that and then grep on "A" & "B"
// and see the unique transitions. See 7.png for the output.
// $ ./example2 | cut -c1 | grep '[AB]' | uniq

// Sample program to show how the goroutine scheduler
// will time slice goroutines on a single thread.
package main

import (
	"crypto/sha1"
	"fmt"
	"runtime"
	"strconv"
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

	fmt.Println("Create Goroutines")

	// Create the first goroutine and manage its lifecycle here.
	go func() {
		printHashes("A")
		wg.Done()
	}()

	// Create the second goroutine and manage its lifecycle here.
	go func() {
		printHashes("B")
		wg.Done()
	}()

	// Wait for the goroutines to finish.
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("Terminating Program")
}

// printHashes calculates the sha1 hash for a range of
// numbers and prints each in hex encoding.
func printHashes(prefix string) {

	// print each has from 1 to 10. Change this to 50000 and
	// see how the scheduler behaves.
	for i := 1; i <= 50000; i++ {

		// Convert i to a string.
		num := strconv.Itoa(i)

		// Calculate hash for string num.
		sum := sha1.Sum([]byte(num))

		// Print prefix: 5-digit-number: hex encoded hash
		fmt.Printf("%s: %05d: %x\n", prefix, i, sum)
	}

	fmt.Println("Completed", prefix)
}
