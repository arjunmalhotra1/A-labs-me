/*
	So far we saw some core patterns. Now we will see some complex patterns.
	We have concurrency so we don't care what order those 2000 go routines finish in.
	The guarantee point is just waiting for it to finish.

	But there are situations when we don't want all the 2000 threads to be competing for the hardware threads
	at the same time. It would be nice to take some load off the scheduler by moving bulk of those
	go routines into their waiting state and bringing them alive little bit at a time.
	This is where semaphore count comes in.

	A really good practical use for this kind of pattern could be a situation where
	let's say we were going to do 2000 database transactions and if any of those database transaction
	failed we would want to roll back the entire transaction.
	We launch 2000 go routines to do individual inserts and if any of those inserts fail we roll them back.
	It would be horrible if the very first go routine that executes fails we would still have to do the other
	1,999 inserts before we can roll it back. (Didn't really understand this)

	Fan out semaphore is a practical way of batching a little bit of work at a time so we can check things.

	We create a buffer of 2000
	children := 2000
	"ch := make(chan string, children)"
	One slot for every go routine so the send side can complete without the receive.
	We will walk away from the guarantee here, and reduce the latency between the signalling send and receive.
	--------------------------------------------------------
	But we would want to batch the actual number of go routines that can execute at any given time.
	In his example in the video he used NumCPU as his initial magic number that is one go routine per hardware thread.
	That's how many we are going to let execute at any given time reduce some load on our scheduler.

	Then we create our semaphore channel. It is a buffered channel but it is based on a number of go routines
	that we would like to see in a running state that's why here we match our NumCPu to the channel.
	----------------------------------------------------------

	g:= runtime.NumCPU()
	// 	sem := make(chan bool, g)
	Next, we have a loop for these 2000 go routines and we create those 2000 go routines.
	Now we have 2000 go routines in a runnable state.
	But as the go routines move from runnable to running the scheduler chooses that.
	Each go routine has to signal a value inside the semaphore channel.
	"sem <- true"

	We know that only 8 go routines can do that any given time because then the channel will be full.
	Once a buffered channel is full, signalling on the the send side will block and no further go routines (beyond 8) would be able to
	go past the "sem <- true" statement.
	Just like if the buffer is empty the signalling on the send side will block.

	So this is going to guarantee that out of 2000 only 8 go routines can be doing the work inside this below block
	in any given time.

	sem <- true
	{
		t := time.Duration(rand.Intn(200)) * time.Millisecond
		time.Sleep(t)
		ch <- "data"
		fmt.Println("child : sent signal :", child)
	}

	Now we have 8 go routines in parallel in our machine, doing the work
	"t := time.Duration(rand.Intn(200)) * time.Millisecond" we are just simulating the work/unknown latency.
	Eventually these go routines finish their work, they signal their result
	on the other channel
	"ch <- 'data'"
	Remember only competition for that signal is is on the send side if multiple go routines finish at the same time.

	Then once they are done with their work they
	"<-sem" pull that value out. Once they pull a value out of the semaphore channel it will
	open a chance for other go routine to complete their signal.
	So our bulk of go routines block on "sem<-true" and
	as go routines come in and we move values out of the semaphore "<-sem"
	it opens up a slot for another random go routine.

	Finally we wait for all the children to get done.
	In the loop we are receiving values off the channel, we are decrementing children count until it gets back to 0.
	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent : recv'd signal :", children)
	}

	We saw the fanout pattern before now we add the semaphore channel so that we can control the number of
	go routines that are actually executing, in this case 8 out of the 2000 at any given time.



*/

// func fanOutSem() {
// 	children := 2000
// 	ch := make(chan string, children)

// g:= runtime.NumCPU()
// 	g := runtime.GOMAXPROCS(0)
// 	sem := make(chan bool, g)

// 	for c := 0; c < children; c++ {
// 		go func(child int) {
// 			sem <- true
// 			{
// 				t := time.Duration(rand.Intn(200)) * time.Millisecond
// 				time.Sleep(t)
// 				ch <- "data"
// 				fmt.Println("child : sent signal :", child)
// 			}
// 			<-sem
// 		}(c)
// 	}

// 	for children > 0 {
// 		d := <-ch
// 		children--
// 		fmt.Println(d)
// 		fmt.Println("parent : recv'd signal :", children)
// 	}

// 	time.Sleep(time.Second)
// 	fmt.Println("-------------------------------------------------")
// }

// This sample program demonstrates the basic channel mechanics
// for goroutine signaling.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	waitForResult()
	// fanOut()

	// waitForTask()
	// pooling()

	// Advanced patterns
	// 		fanOutSem()
	// 		boundedWorkPooling()
	// 		drop()

	// Cancellation Pattern
	// 		cancellation()

	// Retry Pattern
	// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 		defer cancel()
	// 		retryTimeout(ctx, time.Second, func(ctx context.Context) error { return errors.New("always fail") })

	// Channel Cancellation
	// 		stop := make(chan struct{})
	// 		channelCancellation(stop)
}

// waitForResult: In this pattern, the parent goroutine waits for the child
// goroutine to finish some work to signal the result.
func waitForResult() {
	ch := make(chan string)

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		ch <- "data"
		fmt.Println("child : sent signal")
	}()

	d := <-ch
	fmt.Println("parent : recv'd signal :", d)

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// fanOut: In this pattern, the parent goroutine creates 2000 child goroutines
// and waits for them to signal their results.
func fanOut() {
	children := 2000
	ch := make(chan string, children)

	for c := 0; c < children; c++ {
		go func(child int) {
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			ch <- "data"
			fmt.Println("child : sent signal :", child)
		}(c)
	}

	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent : recv'd signal :", children)
	}

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// waitForTask: In this pattern, the parent goroutine sends a signal to a
// child goroutine waiting to be told what to do.
func waitForTask() {
	ch := make(chan string)

	go func() {
		d := <-ch
		fmt.Println("child : recv'd signal :", d)
	}()

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	ch <- "data"
	fmt.Println("parent : sent signal")

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// pooling: In this pattern, the parent goroutine signals 100 pieces of work
// to a pool of child goroutines waiting for work to perform.
func pooling() {
	ch := make(chan string)

	// g:= runtime.NumCPU()
	g := runtime.GOMAXPROCS(0)
	for c := 0; c < g; c++ {
		go func(child int) {
			for d := range ch {
				fmt.Printf("child %d : recv'd signal : %s\n", child, d)
			}
			fmt.Printf("child %d : recv'd shutdown signal\n", child)
		}(c)
	}

	const work = 100
	for w := 0; w < work; w++ {
		ch <- "data"
		fmt.Println("parent : sent signal :", w)
	}

	close(ch)
	fmt.Println("parent : sent shutdown signal")

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// fanOutSem: In this pattern, a semaphore is added to the fan out pattern
// to restrict the number of child goroutines that can be schedule to run.
func fanOutSem() {
	children := 2000
	ch := make(chan string, children)

	// g:= runtime.NumCPU()
	g := runtime.GOMAXPROCS(0)
	sem := make(chan bool, g)

	for c := 0; c < children; c++ {
		go func(child int) {
			sem <- true
			{
				t := time.Duration(rand.Intn(200)) * time.Millisecond
				time.Sleep(t)
				ch <- "data"
				fmt.Println("child : sent signal :", child)
			}
			<-sem
		}(c)
	}

	for children > 0 {
		d := <-ch
		children--
		fmt.Println(d)
		fmt.Println("parent : recv'd signal :", children)
	}

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// boundedWorkPooling: In this pattern, a pool of child goroutines is created
// to service a fixed amount of work. The parent goroutine iterates over all
// work, signalling that into the pool. Once all the work has been signaled,
// then the channel is closed, the channel is flushed, and the child
// goroutines terminate.
func boundedWorkPooling() {
	work := []string{"paper", "paper", "paper", "paper", 2000: "paper"}

	g := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for c := 0; c < g; c++ {
		go func(child int) {
			defer wg.Done()
			for wrk := range ch {
				fmt.Printf("child %d : recv'd signal : %s\n", child, wrk)
			}
			fmt.Printf("child %d : recv'd shutdown signal\n", child)
		}(c)
	}

	for _, wrk := range work {
		ch <- wrk
	}
	close(ch)
	wg.Wait()

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// drop: In this pattern, the parent goroutine signals 2000 pieces of work to
// a single child goroutine that can't handle all the work. If the parent
// performs a send and the child is not ready, that work is discarded and dropped.
func drop() {
	const cap = 100
	ch := make(chan string, cap)

	go func() {
		for p := range ch {
			fmt.Println("child : recv'd signal :", p)
		}
	}()

	const work = 2000
	for w := 0; w < work; w++ {
		select {
		case ch <- "data":
			fmt.Println("parent : sent signal :", w)
		default:
			fmt.Println("parent : dropped data :", w)
		}
	}

	close(ch)
	fmt.Println("parent : sent shutdown signal")

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// cancellation: In this pattern, the parent goroutine creates a child
// goroutine to perform some work. The parent goroutine is only willing to
// wait 150 milliseconds for that work to be completed. After 150 milliseconds
// the parent goroutine walks away.
func cancellation() {
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	ch := make(chan string, 1)

	go func() {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		ch <- "data"
	}()

	select {
	case d := <-ch:
		fmt.Println("work complete", d)

	case <-ctx.Done():
		fmt.Println("work cancelled")
	}

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------")
}

// retryTimeout: You need to validate if something can be done with no error
// but it may take time before this is true. You set a retry interval to create
// a delay before you retry the call and you use the context to set a timeout.
func retryTimeout(ctx context.Context, retryInterval time.Duration, check func(ctx context.Context) error) {

	for {
		fmt.Println("perform user check call")
		if err := check(ctx); err == nil {
			fmt.Println("work finished successfully")
			return
		}

		fmt.Println("check if timeout has expired")
		if ctx.Err() != nil {
			fmt.Println("time expired 1 :", ctx.Err())
			return
		}

		fmt.Printf("wait %s before trying again\n", retryInterval)
		t := time.NewTimer(retryInterval)

		select {
		case <-ctx.Done():
			fmt.Println("timed expired 2 :", ctx.Err())
			t.Stop()
			return
		case <-t.C:
			fmt.Println("retry again")
		}
	}
}

// channelCancellation shows how you can take an existing channel being
// used for cancellation and convert that into using a context where
// a context is needed.
func channelCancellation(stop <-chan struct{}) {

	// Create a cancel context for handling the stop signal.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// If a signal is received on the stop channel, cancel the
	// context. This will propagate the cancel into the p.Run
	// function below.
	go func() {
		select {
		case <-stop:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Imagine a function that is performing an I/O operation that is
	// cancellable.
	func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.ardanlabs.com/blog/index.xml", nil)
		if err != nil {
			return err
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		return nil
	}(ctx)
}
