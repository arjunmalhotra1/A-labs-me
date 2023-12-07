/*
	First we comment out these 4 lines of code:
		// value := counter
		// value++
		// fmt.Println("logging")
		// counter = value

	Since our shared state is "counter" variable.
	And use
	atomic.AddInt32()

	But if we do,
	atomic.AddInt32(&counter,1) - We get compiler error because count right now, is not a precision
	based integer. It is just based on int (var counter int) this is saying use the most efficient integer
	for the architecture that we are running on.

	So in this case we have to be very explicit about our precision 32 bit integers regardless of the architecture.
	"var counter int32" Then we get no error.

	atomic.AddInt32(&counter,1)
	The sharing is occuring at the address. So if we call any atomic instruction at that address regardless of
	if the operation is Add, Load or Store. Then the most goroutines will fall in line and synchronize.

	Now when we build it with
	"go build -race"
	And then run it, we will see that there is no race and now we get "4" as our "Final Counter"
	Just one line of code fixes all of it because this is a counter.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	But say for some reason, we want all of these instructions to be synchronized.
	We can no longer use the Atomic instruction anymore.
	Now we have to use a "mutex".

	// value := counter
	// value++
	// fmt.Println("logging")
	// counter = value

	Mutex gives us the ability to take multiple lines of code and to make the execution of those lines atomic.

	A Mutex allows us to create a room around our code. All go routines will be at the door and would want to be led
	in.
	See 11.png
	But it is the job of the scheduler to make sure that only one go routines gets in the room
	executing that code in any given time. Only when that go routine sis out of the code another go routine
	gets a chance to come in.
	See 12.png

	For Go routines to ask or get access to the room they will call a function called "lock"
	and when they leave they call unlock.

	If we have a Go routine that calls lock and doesn't call unlock then we will end up with a deadlock.
	Since the other Go routines will never be allowed in.

	Remember mutexes are not a queue. This is not "First Go routine here will be the first go routine that gets
	the access to the code."

	Mutexes use a fair scheduling algorithm.
	Say we have 1st go routine that shows up at the door and then there's a 5th go routine.
	The 5th go routine could be accessing the code before 1st. See 13.png.

	So don't think of mutexes as a queue but it is fair in that the first go routines shouldn't also be the
	last one that gets in.

	Remember the amount of time it takes to execute through the code window are the latencies for the go routines
	waiting to get in.

	So now we will create a room where only one Go routine will come and access the room.
	The longer the go routines wait will be the internal back pressure in the app.
	So we will have to call Unlock.

	See main-2.go to see the mutex code.
*/
// go build -race

// Sample program to show how to create race conditions in
// our programs. We don't want to do this.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
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

	// Create two goroutines.
	for g := 0; g < grs; g++ {
		go func() {
			for i := 0; i < 2; i++ {
				atomic.AddInt32(&counter, 1)

				// // Capture the value of Counter.
				// value := counter

				// // Increment our local value of Counter.
				// value++

				// fmt.Println("logging")

				// // Store the value back into Counter.
				// counter = value
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
