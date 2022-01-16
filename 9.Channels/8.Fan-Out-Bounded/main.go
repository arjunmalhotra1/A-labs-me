/*
This is a situation where we have a discreet amount of work, say 2000 pieces of work.
But we don't want to throw 2000 go routines at the problem we really want to
control the number of go routines that are going to work.
So here we will bind/bound the number of go routines that actually get any work done.
We are not going to create 2000 adn then limit only 8 like we saw with fanout-semaphore pattern.
We are really going to end up kind of in this pooling idea but it's not an unlimited pool we know exactly
how much work needs to get done.
We are just going to limit the number of go routines to do it.

Here we are defining a 2000 pieces of work to do.
"work := []string{"paper", "paper", "paper", "paper", 2000: "paper"}"

**NOTE**
we are setting an index position for the last piece of work, '2000: "paper"'
that will create a slice of length of 2000.

Again he used
"g:= runtime.NumCPU()" as an initial magic number for this pool.
So in our case we will have 8 go routines, that are going to do the full work of 2000 pieces of work that
we have.

Then we use wait group to orchestrate when all this work gets done.
"var wg sync.WaitGroup"
we add 8 go routines to the wait group in our case.
"wg.Add(g)"

Then we make a channel
"ch := make(chan string, g)"
This is a buffered channel. The idea is that we will be feeding these 8 go routines, the 2000 pieces of work.
We can really limit our memory. We don't have to create a buffered channel of 2000.
What if ther's a million pieces of work that we need to get done. We wouldn't want to create a buffered channel
of a million.
So what's nice here is that as long as we continue to keep the data flow in the channel it doesn't have to be a
large channel to feed off of.
We can see that with some benchmarking, using a buffer 8 or 2,000 we will get the same results.
"ch := make(chan string, g)" So we have a buffer that we will feed work into for our bounded fan out here.

Then we launch the 8 go routines.

	for c := 0; c < g; c++ {
		go func(child int) {
			defer wg.Done()
			for wrk := range ch {
				fmt.Printf("child %d : recv'd signal : %s\n", child, wrk)
			}
			fmt.Printf("child %d : recv'd shutdown signal\n", child)
		}(c)
	}

	These go routines are blocked on line "for wrk := range ch {"
	waiting for work to do.
	Q. Which go routine will do what work?
	A. It doesn't matter because we are dealing with concurrency here.

	"defer wg.Done()" means that execute this statement but not now, defer it for when the function or in this case the
	go routine terminates.
	We will decrement the wait group once there is no more work. That will happen once
	the for range terminates and the for range will terminate when we close it.

	We have  a pool of work and we fan out 8 go routines, they will be waiting for work.
	Eventally we will tell them there is no work.
	They will all terminate, we will decrement that wait group and we will be able to move on.

	Here we are raging over the collection of work.
	And we are sending the work into the channel.

	for _, wrk := range work {
 		ch <- wrk
 	}

	This will get the go routines busy doing work.
	for wrk := range ch {
		fmt.Printf("child %d : recv'd signal : %s\n", child, wrk)
	}

	Once the 8 go routines are busy we are going to be blocked here at "ch <- wrk"
	but eventually the go routine finishes the work, comes back into the range
	and we will keep the data flowing.
	Once the last piece of work is submitted, into the buffered channel we close the channel immediately.
	"close(ch)"
	We are signalling without data here.

	The go routines continue to flush the buffer out eventually there is no more data the close signal starts to fire.
	The go routines start to terminate, we start to decrement those wait groups.
	Once we finish the close we are now at the guarantee point
	"wg.Wait()" we are blocked, waiting for all the go routines to decrement the waitgroup count back to 0.
	And then we move on.

	This fan out bounded pattern let's us stage tremendous amount of work that we may not feel comfortable
	throwing a go routine at every task, and this also allows us to minimize the memory resources that we need.

	We still would need to be able to figure out a magic number for this kind of
	fan out pool that we have created.
	Normally Num.CPU() is an amazing starting point and ususally it does give us the results that we need.
	Fan out patterns is one of Bill's favorite.






*/

// func boundedWorkPooling() {
// 	work := []string{"paper", "paper", "paper", "paper", 2000: "paper"}

// g:= runtime.NumCPU()
// 	g := runtime.GOMAXPROCS(0)
// 	var wg sync.WaitGroup
// 	wg.Add(g)

// 	ch := make(chan string, g)

// 	for c := 0; c < g; c++ {
// 		go func(child int) {
// 			defer wg.Done()
// 			for wrk := range ch {
// 				fmt.Printf("child %d : recv'd signal : %s\n", child, wrk)
// 			}
// 			fmt.Printf("child %d : recv'd shutdown signal\n", child)
// 		}(c)
// 	}

// 	for _, wrk := range work {
// 		ch <- wrk
// 	}
// 	close(ch)
// 	wg.Wait()

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
