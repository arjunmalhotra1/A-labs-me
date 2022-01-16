/*
	We declare our mutex like:
	 	var mu sync.Mutex

	Mutex is usable in it's zero value state just like the wait group.

	Remember if you add a mutex or a wait group to a struct. Then we cannot make a
	copy of the struct anymore. Different copies of a mutex would be technically different mutexes.

	We could create an artificial code block with a mutex.
	With the artificial code block the visualization becomes easier.

	Now only one go routine at a time gets to execute that block.
	We guarantee the synchronization.
	See 14.png for the output.
	We get the output as 4. We have now synchronized these 4 lines of code.

	Remember:
	Mutex creates back pressure with latency for the waiting go routines.
	The longer the go routines wait to get in the longer will be the latency and worse will be our throughput overtime.

	Hence anytime we use mutexes we have to make sure that we are doing the bare minimum.
	Like here we can say that to have line:
	"fmt.Println("logging")" inside of a mutex is unnecessary.
	This shouldn't be inside the mutex. This is adding tens of micro seconds of latency inside the mutex.
	Which is going to add up to the latency of the other Go routines which are waiting to enter the code block.

	We want to make sure we are doing the bare minimum.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	There's also a read write mutex in Golang.
	REad Write mutex allows us to have multiple readers at the same time.
	Remember reading is free. REading is not our problem.
	It's the mutation of memory that's our problem.
	Remember data race doesn't occur unless there's one go routine performing a write.

	Say, you might decide one day to use a map as a cache.
	While there are Go routines can read the cache at the same time that's not a problem.
	But if suddenly say a Go routine wants to add a value to the map then we got to stop all the reads.

	See main-3.go for an example.


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
