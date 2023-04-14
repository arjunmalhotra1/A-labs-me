package logger

import (
	"fmt"
	"io"
	"sync"
)

/*
	First thing we would need to solve this is a new type.
	We will call it the logger type.
	Logger type will be data that knows how to log and give us the ability to detect the problems.
	For bare minimum we know that the logger needs a channel.
	We can use string because the log data is string base.
	Since we will be launching a go routine that will write the logs,
	we would want the go routine to have the ability to be start out but also shut down cleanly.
	Hence we also give the logger a wait group.
	That wait group will give us the ability to start and shut down that go routine cleanly.

	type logger struct {
		ch chan string
		wg sync.WaitGroup
	}

	Next what we need is a factory function. We have got to be able to construct a logger
	in such a way that it's usable.

	In order to construct a logger we need to be able to construct this channel
	"ch chan string" we know this is going to be based on a capacity so we need a buffer size.
	We also know we need to be able to write to a device.
	WE would need to know what device we want to write to.
	We can use the IO writer interface againt to say
	"I don't care what the device is as long as we can write to it"
	and we need a capacity value, so we know what that capacity is, so that we have a good number that
	doesn't create these false positives but also isn't too large where we are not detecting the failures
	soon enough.
	"func New(w io.Writer, cap int) *Logger {"
	Now we can construct our logger, because we have everything we need already.
	Wait groups are usable in their 0 value state.
	But channels can be constructed based on the capacity value.

	l := Logger{
		ch: make(chan string, cap),
	}

	Finally we can return the logger using pointer semantics. This
	isn't something we would want to make copies of, we want to create one logger and share it.
	So we use the pointer semantics.
	"return &l"

	Next thing is that we need to create the go routine.
	We need a go routine that is going to be able to receive off that channel and write to the device.
	So we set up a literal function adn say this is where the device thing is going to be.

	go func() {

	}()

	If we are going to have the literal function then we need to set up our wait group
	and add 1 to our wait group.
	We are going to say we know there's one go routine in flight that we know about.
	When the go routine terminates then we would want to call "Done()" to decrement that value.
	Note we are using closures here to get to the waitgroup that is based on the logger.
	Now we have to be able to feed off the buffered channel.
	So what we can do is we can range off the channel. And every value we receive off of the channel
	we can write it to the device.
	What's nice is that fmt package has an Frintln function which takes as it's first parameter
	any concrete data that knows how to write and we write "v".

	go func() {
		for v := range l.ch {
			fmt.Fprintln(w, v)
		}

		l.wg.Done()
	}()

	Next we will make our go routine how to receive off the channel, everytime
	the go routine sees logging data it will write it tot he device that we specify.
	And will keep going and eventually we should be able to shut that go routine down.
	Once we shut the go routine down we decrement the wait group and we are good.

	Next we focus on the shut down.
	Now we need a method that knows how to shut down our logger.
	How do we shut down the go routine?
	What happens if we just close the channel?

	func (l *Logger) Shutdown() {
		close(l.ch)
	}

	This will cause the for loop
	"for v := range l.ch {" to terminate once it is flushed.
	We have to remember that if we close the channel we can receive from it but no longer send.

	Therefore the application developer has to make sure that if they call
	shutdown, then no more go routines will try to do anymore logging.
	But that is the developer responsibility to have a clean shutdown and it is a good way of finding
	other problems in your code if you can't shut down completely.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	After close, we still need to wait to know that all the go routines shutdown.

	func (l *Logger) Shutdown() {
		close(l.ch)
		l.wg.Wait()
	}

	closes the channel signals the channel without data tell the go routine to terminate
	nd once the go routine reports that it is done,
	"l.wg.Done()"
	we will be able to move forward beyond "l.wg.Wait()"
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	there's one more thing we have to implement  that is the call to Println
	"l.Println(fmt.Sprintf("%d: log data", id))" this is what he had in his code in the video.
	This is the API that all the go routines will be logging against.

	Hence what we need to do is create another method, "Println"
	we will keep it simple and just say that it takes just a string to log.

	func (l *Logger) Println(v string){

	}

	It is Println() that was called by the 50,000 go routines. This is where we have to detect
	whether we have a problem or not. This is where we detect whether we have a capacity issue.
	We use select statement

	func (l *Logger) Println(v string){
		select {
		}
	}

	Let's attempt to send the log "v" over the channel.
	And if we can then great and if we can't then we will go into our default which will be our
	drop situation.
	For here we will just log the drop string so that we can see that we have moved
	on to the default case nad we are not blocking.

	func (l *Logger) Println(v string){
		select {
		case l.ch <- v:
		default:
			fmt.Println("Drop")
		}
	}











*/
type Logger struct {
	ch chan string
	wg sync.WaitGroup
}

func New(w io.Writer, cap int) *Logger {
	l := Logger{
		ch: make(chan string, cap),
	}

	l.wg.Add(1)
	go func() {
		l.wg.Done()
		for v := range l.ch {
			fmt.Fprintln(w, v)
		}

		l.wg.Done()
	}()

	return &l
}

func (l *Logger) Shutdown() {
	close(l.ch)
}

func (l *Logger) Println(v string) {
	select {
	case l.ch <- v:
	default:
		fmt.Println("Drop")
	}
}
