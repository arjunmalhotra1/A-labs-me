/*
	We will see here a production level code.
	Bill was writing a code and there was a bug this video is about that bug and how they solved it.
	This is about the idea of signalling nad being able to detect failures quickly.

	Scenario - They were building software that was gong to be handling, over 50,000 users.
	Pounding servers with tons of requests.
	this software did a lot of logging. That was a requirement from the client see 1.png.

	After 3 days in the middle of the testing, after a minor change, they see all the logging getting stopped.
	2.png.
	They had 50,000 go routines trying to write a log and the entire server deadlocked.

	After research they identified what had happened. The host program that was
	feeding, standard out for the logs, that was writing data to disk, the disk filled up and
	then every write of log was blocked.
	All they had to do was clear the disk space and logging would start again.
	But this is a production environment. So we cannot have that.

	This is an interesting situation. They have standard library logger we let the 50,000
	go routines write to standard out but now we have a potential problem, the standard out could block
	and could cause all the go routines to deadlock and the web server is no longer useful.
	We can't have this, logging is important but in this particular app it's not more important than
	getting the work done.

	So now we have a problem. We have to detect when we cannot log anymore and no longer allow
	that logging to block. But more importantly not only we have to detect when we are having
	a logging issue so that we can keep going.
	Logging is important so we also have to detect when we can log again so that we can start that
	up and running.

	Below is the code we are using to simulate our issue. First we have defined a device.

	type device struct {
		problem bool
	}
	And here we have implemented the I/O writer interface
	"func (d *device) Write(p []byte) (n int, err error) {"
	So that we can simulate out own diskfull problems while writing to our logs.

	Everytime we execute "log.Print()" `Write` is going to get executed for the device.
	If the problem field is true,
	`for d.problem {

		// Simulate disk problems.
		time.Sleep(time.Second)
	}`

	We go into an endless loop every second simulating the disk blocking.
	If problem is false then we will just write that back out to standard out.
	This is how we are simulating our disk problems here.

	In the main go routine we have 10 go routines, we also create our device set to it's 0 value.
	Next we share the device variable with the factory function for the standard library logger
	because our device implements the I/O writer interface we can pass device into the
	logger.
	"l := logger.New(&d, grs)""

	Next, we launch the 10 go routines and every 10 milliseconds the 10 go routines write to the log.
	for i := 0; i < grs; i++ {
		go func(id int) {
			for {
				// He had here,
				// l.Println(fmt.Sprintf("%d: log data", id))
				l.Write(fmt.Sprintf("%d: log data", id))
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	He said everytime this line executes,
	"l.Write(fmt.Sprintf("%d: log data", id))"
	It is essentially executing the write method implemented by our device.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	To simulate the problem,
	we hook into the operating system
	"sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)"
	And then we are in this lop waiting for the Bill to hit `ctrl+C`

	for {
		<-sigChan

		// I appreciate we have a data race here with the Write
		// method. Let's keep things simple to show the mechanics.
		d.problem = !d.problem
	}

	When we hit `ctrl+C` we flip the problem from "false to true" and "true to false".
	Yes we have a data race here because we have the potential to read and write at the same time.
	But this is not about the data race, so we will just keep it simple.
	Everytime we hit `ctrl+C` we can flip the problem a little bit and that will cause the writes
	to either endlessly loop or eventually write.

	`for d.problem {

		// Simulate disk problems.
		time.Sleep(time.Second)
	}`

	That is why when the program was running and in the write, "log data".
	When we hit `ctrl+C` now we are in the
	endless loop (the write stops as it's on "time.Sleep(time.Second)") when we hit `ctrl+C` again,
	we are out of the endless loop.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	But what we need is, if we hit ctrl+C and if we are in the endless loop

	`for d.problem {

		// Simulate disk problems.
		time.Sleep(time.Second)
	}`

	we cannot have this code blocking.
	// He had here,
	// l.Println(fmt.Sprintf("%d: log data", id))
	l.Write(fmt.Sprintf("%d: log data", id))

	This is where we had those 50,000 go routines getting blocked.
	We can't block anymore on this call.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	So we want to minimize this complexity, we are going to be able to identify the failure,
	also identify when it's fixed and keep this program running.

	We can't use the standard library logger anymore we can't have 50,000 go routines
	writing to the log directly.
	so how are we going to be able to detect a problem, be able to essentially by pass that blocking call
	and get logging back on.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	For Bill this is a classic classic drop pattern problem. A capacity drop pattern problem.

	So what we ned to do is to write our own log package now.
	Idea is that instead of letting the 50,000 go routines, write to the device.
	What we will do is only have one go routine write to the device.
	If we keep it down to one go routine then we will have an easier time detecting
	when there is a problem on this write as opposed to detecting a problem on the write when we have
	50,000 go routines. See 3.png.

	That means we still have 50,000 go routines, that want access to that device. See 4.png.
	How are we going to do that?
	This is where we can use our signalling and create a buffered channel.
	To be able to signal to the go routine data we want writtin to the log.
	See 5.png.
	Idea of the buffered channel is again for our capacity.
	We need to find a capacity, where under normal load, may be we end up filling 80-90% of
	the buffer under normal load. See 6.png
	We only have the 10-20% extra buffer for room. Because the idea is that
	if everythong is working the go routine is receiving of the channel it is keeping the 10-2% gap
	fine. But if the device "D" blocks all of a sudden then the go routine writing ot the device
	is not going to receive off of the buffer any longer, and when that receive starts to fail,
	we will fill that 10-20% really quickly the 50,000 go routines will not be able to
	signal anymore, and we can drop those logs. See 7.png.

	The best part is if when we fix the disk write issue,
	the receive of logs by the go routine writign to the device starts happening again and be able to
	get back the 80-90% capacity.
	So we wil sue the drop pattern to solve this problem.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	// Number of goroutines that will be writing logs.
	const grs = 10

	To have our capacity number of our channel, we are not sure but we start with the 
	number of go routines.
	So for every go routine that we will be logging we will give one capacity value in that 
	bufferend channel. 

	var d device
	l := logger.New(&d, grs)

	So now when 50,000 go routines are writing to the disk we see 8.png
	Now hen he does ctrl+c to simulate the disk space is full we get nothing in the screen see 9.png.
	But we haven't hung the program we have not completely stopped the program.
	We are not deadlocked, we are running even though  we cannot log.

	Now if we go and fix the problem press 'ctrl+C' in our case.
	Then we see the logs again on the screen 8.png

	The go routine was able to pull enough data out of our buffer to get usour 
	gap (10-20%) again and now we are logging again.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	Using the channels semantics and the drop channel we are able to detect the problem quickly
	recover detect that the problem has been fixed, and keep logging.
	The channel gave us the orchestration primitives, is what is really allowing us to do this.
	 





*/

// This sample program demonstrates how the logger package works.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	//"github.com/ardanlabs/gotraining/topics/go/concurrency/patterns/logger"
	/*
		We replace the standard library logger
	*/
	log "C:\Users\Arjun\Documents\Udemy\Ardan-Labs-Ultimate-Go\10.Concurrency-pattern\2.Failure-detection\logger\logger.go"
)

// device allows us to mock a device we write logs to.
type device struct {
	problem bool
}

// Write implements the io.Writer interface.
func (d *device) Write(p []byte) (n int, err error) {
	for d.problem {

		// Simulate disk problems.
		time.Sleep(time.Second)
	}

	fmt.Print(string(p))
	return len(p), nil
}

func main() {

	// Number of goroutines that will be writing logs.
	const grs = 10

	// Create a logger value with a buffer of capacity
	// for each goroutine that will be logging.
	var d device
	l := logger.New(&d, grs)

	// Generate goroutines, each writing to disk.
	for i := 0; i < grs; i++ {
		go func(id int) {
			for {
				// He had here,
				// l.Println(fmt.Sprintf("%d: log data", id))
				l.Write(fmt.Sprintf("%d: log data", id))
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// We want to control the simulated disk blocking. Capture
	// interrupt signals to toggle device issues. Use <ctrl> z
	// to kill the program.

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		<-sigChan

		// I appreciate we have a data race here with the Write
		// method. Let's keep things simple to show the mechanics.
		d.problem = !d.problem
	}
}
