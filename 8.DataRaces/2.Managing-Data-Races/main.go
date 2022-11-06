/*
	We just saw how to create Go routines and little bit about orchestration.
	But the real problem with multi threaded software is that we have to manage nad control
	all those different paths of execution.

	Data races occur when we have 2 go routines, where one is doing a read and the other is doing a write.
	We need to have at least one write. They are trying to access the same memory location at the same time.

	This becomes a problem because there is no predictability on what's going to happen.
	So we cannot allow two go routines to access the shared memory location if at least one of those is write.

	This means now we have to synchronize.
	Synchronize means the go routines have to get in line and take a turn.

	In our below program this
		"var counter int"
	is a shared state.

	We will have 2 Go routines trying to access this shared state.

	In the loop we create the Go routines see 1.png.
	The 2 go routines will be each taking a turn and incrementing the shared state.
	In fact we are going to do it twice. Since we are starting at 0, once both the go routines finish
	their operation the value should end up as "4".

	When we build this program and run it, everytime we get a value of "4".

	Now say after few months we have a new developer who now starts logging.
	like,
	fmt.Println("logging")
	to start tracing the code.

	Once the developer is comfortable in seeing what the program does, he/she then forgets to pull out the logging.

	When we build and save this program with the logging we get to see 2.png as output.
	With the logging the output is "2". Everytime we run this program we see the value = "2".

	All we did was adding logging between those operationa and now, we have a different result.
	We have a classic data race problem.
	There were no guarantees in the code before and we were geting luck and now we are no longer getting lucky.

	Now we read the counter, we did the modification but before we do the write, now what happens is that print
	statement is allowing the scheduler to make a context switch. See 4.png.

	Now when we make the context switch, the second go routine comes in does it's read and does it's modification and
	context switch again. See 5.png.

	Now on G1(Go routine 1) we do our write, we read counter and then modify to 2 and now again context switch.
	And this time G2 has no idea tha the value it has is a dirty value.

	This is the problem with the data race. This is also a problem with our value semantics.
	This is when efficiency of pointers kicks in. BEcause we have no concept of that the value is dirty
	when we come back to G2. G2 writes 1 again adn then reads 1 and increments 1 to 2. The again Context switch.
	See 6.png
	Then on G1 it writes 2. We can see that the context switches are causing the values on our multi threaded
	Go program to now be dirty.
	We don't have any ability to know that they are dirty.

	So we need to do is have each Go routine take a turn so these 3 operations are atomic.
	They were atomic before we put the print statement in there.

	But now that hey have added the print statement we are seeing that each go routine is not running
	atomically. Technically before we added the print statements they weren't atomic and we were just getting lucky.

	One nice thing about Go is that we have a data race detector.

	See 7.png

	We do "go build -race" what this does is it builds a new binary but with the race detection.
	We still have to run the binary to see if there's a race.

	1. If the data race detector finds that you have a data race, don't argue with it. It;s there.
	2. If the data race detector doesn't find a data race, it doesn't mean you don't have one. It just means
	you haven't found it yet, may be you need to run the program a little bit more.

	Bill said that he normally uses the race detector with "test"
	"go test -race"
	He doesn't tend to build the binary with the race detector. Bill said he's heard that a binary with
	the race detector could run 25% slower. SO he uses "-race" with the "go test"

	He recommends that if we use the "go test -race" and if we are running these tests on a MAC, then
	we should probably set the CPU value to 3 times the number of hardware threads just to create a little more chaos.
	It will help find data races.

	"go test -race -cpu 24"

	Now when he runs see 8.png

	Since here he had already built it with the data race detector, "go build -race".
	He just ran it "./example1".
	The data race message pops up.
	The message tell us the line numbers where it found the data races.
	There are data races on line 41 and 33.

	If the race detector finds a race and if it's not obvious to you, we have probably lost some control
	over the control and the code might have got a little more complex, don't try to patch that and
	you may have to rethink about the architecture.

	The problem here is that this code is not atomic the race detector found it as well.
	We realized that we have the ability to have shared access to the same memory location
	where one of the Go routines is doing a write.

	If we take the print statement out, we are still going to have a data race.

	Question. How do we fix this? How do we make sure that these 4 lines of code are
	running atomically

	value := counter
	value++
	fmt.Println("logging")
	counter = value

	We have 2 options when it comes to synchronization:
	1. We have ability to use atomic instructions
	2. We have the ability to use Mutexes.

	Atomic instructions are the fastest way to do synchronization.
	That's because they are at the hardware level. Our hardware is going to help us with synchronization.
	Problem at the hardware level is that we can only synchronize a word or half a word of memory at a time.
	So that is going to limit our ability to use them.

	Luckily right now here we are dealing with an integer that's a word or half a word of memory depending upon
	our architecture, we can use atomic instructions in this case to create that synchronization.

	Atomic instructions are great for counters which is exactly what we have in our example here.
	See main-1.go








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
var counter int

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

				// Capture the value of Counter.
				value := counter

				// Increment our local value of Counter.
				value++

				fmt.Println("logging")

				// Store the value back into Counter.
				counter = value
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
