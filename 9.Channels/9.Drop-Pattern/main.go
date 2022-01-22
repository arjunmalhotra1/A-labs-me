/*
	Drop pattern will help us reduce/minimize things like back pressure in our software.
	Idea of software writing is not when life is good but when life is bad.
	Being able to detect problems and detect failures and be able to keep the software running.

	Idea of drop pattern:
	Imagine you are writing a DNS server, the chance of the DNS server getting attacked is all day long.
	Idea of DNS attack is to flood our DNS server with requests. So many requests that when we try to handle them
	our server implodes.
	What we need to do to defend ourselves is that we need to be realistic and realize that we can't handle
	every request and that we have to know what our capacity is.
	If we know our capacity then we will take in what we can fulfill the request we can and we drop everything else.

	Drop pattern is knowing everything about out capacity.
	What's the number of requests, tasks, things that we can do when we get that many pending things
	when we are at capacity then we drop.

	"const cap = 100" what we are saying is that our capacity is 100.
	Once we get to 100 requests then we are at capacity everything else is going to drop.
	We will use a buffered channel to know when we are at capacity.

	Next we create our go routine and our go routine is now waiting on a receive to receive work.
	Note: This could also be a pool of go routines handling our dns requests. Just to keep the code simple
	we have used 1 go routine here.

	go func() {
		for p := range ch {
			fmt.Println("child : recv'd signal :", p)
		}
	}()

	The go routine is waiting for work here,
	"for p := range ch {"

	Here "const work = 2000" we decide we are going to send 2000 DNS requests to our server.
	Next we decide to feed our go routine.
	Remember the idea is that if we hit capacity we need to detect it and we need to drop those requests.
	We do that using "select", we use the select case statement.
	the select case allows us to perform multiple channel operations, sends or receives at the same time on
	the same go routine.

	select {
		case ch <- "data":
			fmt.Println("parent : sent signal :", w)
		default:
			fmt.Println("parent : dropped data :", w)
		}
	}
	This case is saying perform the signal send.
	"case ch <- "data":"
	If we are not at capacity, then we can do that send.
	But if we are at capacity then the channel send operation would block because there is no room to complete
	the send.
	This is where we do the drop because the default case is where we win this game.

		default:
			fmt.Println("parent : dropped data :", w)

	Because the select case is saying is if there is a default and all the cases are going to block,
	don't take the blocking lane and jump to the default. This allows us to detect if we are at capacity.
	If there is no default then "select" is a blocking call.
	We are going to block on all the different cases.

	We know if we get tot he default case we are at capacity and we are being flodded with requests more than
	we can handle and we will have to start rejecting anything that is coming into our server.

	Eventually when we get through that 2000 pieces of work, we will close that channel. Which will terminate
	this loop:

	for p := range ch {
			fmt.Println("child : recv'd signal :", p)
		}

	And the go routines will end.









*/

// func drop() {
// 	const cap = 100
// 	ch := make(chan string, cap)

// 	go func() {
// 		for p := range ch {
// 			fmt.Println("child : recv'd signal :", p)
// 		}
// 	}()

// 	const work = 2000
// 	for w := 0; w < work; w++ {
// 		select {
// 		case ch <- "data":
// 			fmt.Println("parent : sent signal :", w)
// 		default:
// 			fmt.Println("parent : dropped data :", w)
// 		}
// 	}

// 	close(ch)
// 	fmt.Println("parent : sent shutdown signal")

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
