/*
	This is a read lock with a write unlock.
	mu.RLock() -
	In this piece of code we don't need mu.RLock() because we need synchronization happening all the time.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	In golang there is no one single line of code that is atomic.
	value++
	is a read ,modify write operation in itself is 3 lines of assembly code.
	This is not atomic. At the hardware level any of the 3 lines of code, read, modify or write could have
	the preemptive context switch at any given time.

	If we don't tell the compiler there needs to be synchronization it can take special liberties underneath
	which can produce random data races.

	There is no one line of code that is synchronous in itself in Go. We will need those
	atomic instructions those mutexes to make sure there is synchronization when needed.
	This is Synchronization.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	Next thing is Orchestration, what is the mechanism when we need two or more go routines to talk
	to each other. We don't want to be using atomic instructions or mutexes for that.
	Golang gave us channels to do orchestration.

*/
// go build -race

// Sample program to show how to create race conditions in
// our programs. We don't want to do this.
package main

import (
	"fmt"
	"sync"
)

// counter is a variable incremented by all goroutines.
// var counter int
var counter int32

func main() {

	// Number of goroutines to use.
	const grs = 2

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(grs)

	var mu sync.Mutex

	// Create two goroutines.
	for g := 0; g < grs; g++ {
		go func() {
			for i := 0; i < 2; i++ {
				// mu.RLock() - This is a read lock and write unlock. We need synchronization here all the time.
				mu.Lock()
				{
					// Capture the value of Counter.
					value := counter

					// Increment our local value of Counter.
					value++

					fmt.Println("logging")

					// Store the value back into Counter.
					counter = value
				}
				mu.Unlock()

			}

			wg.Done()
		}()
	}

	// Wait for the goroutines to finish.
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

/*
==================
WARNING: DATA RACE
Read at 0x0000011a5118 by goroutine 7:
  main.main.func1()
      /Users/bill/code/go/src/github.com/ardanlabs/gotraining/topics/go/concurrency/data_race/example1/example1.go:33 +0x4e

Previous write at 0x0000011a5118 by goroutine 6:
  main.main.func1()
      /Users/bill/code/go/src/github.com/ardanlabs/gotraining/topics/go/concurrency/data_race/example1/example1.go:39 +0x6d

Goroutine 7 (running) created at:
  main.main()
      /Users/bill/code/go/src/github.com/ardanlabs/gotraining/topics/go/concurrency/data_race/example1/example1.go:43 +0xc3

Goroutine 6 (finished) created at:
  main.main()
      /Users/bill/code/go/src/github.com/ardanlabs/gotraining/topics/go/concurrency/data_race/example1/example1.go:43 +0xc3
==================
Final Counter: 4
Found 1 data race(s)
*/
