/*
	Now we shall learn how the Go scheduler semantics work and how it sits on top of the OS scheduler.
	When our go program starts up, the runtime will identify how any hardware threads we have on the host machine.

	Say we have 8 hardware threads on our machine, that means any program that start ups on our machine we will get
	8 logical processors. One per hw thread on this machine.

	Every logical processor is given a real live operating system thread, that operating system thread
	will execute on the hardware. OS is the one responsible for making sure that the OS threads get on the hardware.

	We have one more thing which is the Go routine. Go routine is an application level thread.
  
	Go routine thread/Application level thread is almost identical to those behaviors that the OS level threads have.
	The only difference is that the Go routines doesn't have priorities.
	One of the reasons why GO routines don't have priorities, because we are not dealing with events in our application
	space. The OS has to deal with events, we don't.

	So everything we have talked about applies at the Go scheduler level even though we are a layer above the OS.
	Except for the priorities.

	So far we have seen single threaded Go program, like the Go playground.
	The Go scheduler is a cooperating scheduler, which means that we need events, that are occurring at the application
	level to perform scheduling.
	We have to find the safe points in the code that we can yield the Go routines on and off of "M" in 2.png
	Till 1.11 the safe points happen during function call transitions.
	We are waiting for the function call events to occur in Go.

	In 1.12/1.13 to add some preemptive techniques to the scheduler, to make it a little more random. Similar
	to the way the OS works.

	For now, when the video was recorded the Go scheduler, is a cooperative scheduler,
	that means the Go routines have to yield their time and instead of asking developers to call yield, the
	Go scheduler is the one that is calling yield, on behalf of the Go routines under the idea that there is a minimal
	time slice, etc.

	========================================================================

	There are 4 classes of events that can occur, with function call transitions that give the scheduler an opportunity to make a context switch.

	We use the keyword "go" anytime we are creating a GO routine.
	There are 2 more Data structures in 2.png. Every "P" has a "local run queue", LRQ. This is for the Go routines in the
	runnable state.

	When we create a Go routine, we end up with more Go routines in the runnable state.

	Go Routines in the runnable state end up in the Local Run Queue.

	There's another Data Structure called the "Global Run Queue". GRQ.
	This is another queue for Go routines which are in the runnable state but haven't be assigned to "P" yet.
	See 5.png
	Any time our garbage collector runs is going to be a lot of chaos and you will see a lot of scheduling
	happening.

	Anytime there is a system calls, logging, producing output, you are hitting the network. Almost anything you do
	outside of your app, is going to be a system call in Go.
	That will give the scheduler an opportunity to make a scheduling decision. They are opportunities for
	scheduler to yield Goroutines on and off the "M" in pic 2.png.

	Taking a Go routines on and off of "M" is what we call "context switching in Go".
	In OS the context switch was 1-2 micro second, a 1000-2000 nanoseconds,

	Here in Go scheduler the same context switch is going to be roughly, 200 nanoseconds.
	A significant difference at Go level that will help us as well.

	There are 2 types of system calls, Asynchronous and Synchronous.
	There are other types of blocking calls that are going to happen in our software. In multi threaded software
	the blocking calls could be synchronization and orchestration, Mutexes, Atomic instruction, we also have potentially
	C-GO issues, C-GO is the c compiler for Go so that we can call c libraries and there could be cases that
	we have c functions that could be blocking the "M".

	So we have our 4 classes of events that are occurring at the scheduler that give the scheduler an opportunity.
	We shall go a little deeper into 2 system call and blocking calls to know how those work.

	===========================================================================================================

	Go has something called the network poller, the network poller starts out with a single thread.
	e-poll or cave vent on macos or IOCP on windows, those kind of thread pooling technologies.

	That's what being used here with the network poller.
	Network poller job is to handle Asynchronous system calls. We call it a network poller is because the only
	Asynchronous calls we really work today are the Networking calls.
	All 3 operating systems have an excellent support for Asynchronous networking.
	Unfortunately the only operating systems except of windows doesn't have support for file I/O in an asynchronous way.

	So anytime we are doing the file I/Os we are doing the system synchronous calls.
	Networking will be using the network system calls and the network poller.

	================================================================================================================

	Say the Go Routine in 6.png is wanting to perform a "Read" operation. It's going to perform a "Read" on the network.

	So, as soon as that happens the scheduler is going to say "Since the Go routine is about to perform
	a read, we will take you off of "M" and put you on the network poller."
	We are going to issue that Asynchronous system call and you are going to wait here, till the Operating system can
	complete the read operation.

	We have taken "G" out of "M" and put it in the pool of "M"s (which are responsible for networking requests) which frees up original "M" to do more work.

	So now the scheduler can now choose another GO routine from the Local Run Queue for that "P",
	and we can start getting some more work done.
	Remember there's a 200 ns context switch and it is helping us here because that "M"
	gets to stay busy doing more work.

	Once the network system call is complete, we take that Go routine off of the Network poller and
	bring it back into the "Local Run Queue".
	This is a nice workflow because we are efficiently using the threads that we have.

	But what if this recall wasn't a network call but file I/O call?
	Now we have a little bit of more complications, we can't use the Network Poller anymore. See 10.png

	So what are we going to do? Since the last thing we would want is the Go routine to block the "M".
	If the go routines makes the read call to open up the file. It's going to take a minimum of
	few microseconds of time, we are not getting any work done during that blocking.
	We don't want that. So what scheduler is going to do is
	identify that the Go routine is about to do some file I/O read, and will
	take the "M" and the "g" and will detach it off the "P", now the go routines blocks
	the original "M" (M1) on the side, we move it off. See 11.png.

	But now still we are not getting any extra work done, we will bring a brand new "M"(M2), see 12.png
	If there's an "M" parked and waiting we will bring it in, but if not we will create a new "M".

	Now we will take a "G" off the LRQ, and attach it to this new "M" (M2) and get it running.
	We are still able to get work done.

	Question. What happens once the original Go routine's attached to "M1" system call is
	complete?
	Answer. We will take that "G" off of "M1" and attach it to the LRQ. See 13.png
	And now we will park this "M" (Mp) on the side if we have to do it again.

	============================================================================================

	When GO was first released (8 years ago) a bulk of third party libraries like
	Kafka etc weren't there so Go had to rely on C libraries early on in Go to be able to write programs.
	Since Go manages memory in it's run time, C isn't part of the runtime.
	So the problem is that as soon as we call into a C function, we are kind of out of our managed state.

	Question. So how does a Go scheduler knows, a Go routine making a C function call isn't going to be blocked for a long
	period of time?
	Answer. There's another part of the Go scheduler called "System monitor" and System monitor,
	is monitoring these go routines & the work they are doing.
	So let's say what's happening is that the Go routine (See 16.png) is making a "Read" call,
	But it's from some sought of "C" library. See 17.png

	The system monitor is going to see how long it's been since that Go routine has executed an instruction.
	Say for example 20 ns be the base time period of inactivity. So essentially is a system monitor, that the go routine,
	has been inactive for 20 ns the system monitor will assume that it's blocked and then it will go through the
	same cycle we saw before.
	We wil move the "M2" off with the "Go" routine that is making the C call and is blocked.
	We will being another "M" in and pull another go routine from the LRQ (Local Run Queue).

	So when system monitor, identifies that something is stalled for about 20 ns, it's the same activity
	we saw earlier. The goal is to keep the Go routines running that are in the runnable state.
	==============================================================================================================

	Another efficiency, is when we run Go routines in parallel.

	Say now we have multi threaded Go program. See 19.png
	We have 2 Ps, 2 OS threads (M) and each OS thread, on it's own Hardware thread,a nd the two go  routines
	G1 and G2 can run in parallel.

	Another interesting thing here when we have got these Ps is that the Go scheduler is also a work
	stealing scheduler.

	Currently we have 2 go routines running in parallel, we have 3 Go routines in the LRQ in one processor and
	1 Go Routines in one LRQ. And say we have 2 Go routines in the Global Run Queue (GRQ), they haven't be
	assigned to a "P" yet.
	See 21.png

	Say the "P" on the right finishes all of its's work. Say there's not more work for this "P" to do.
	See 22.png

	Go is a work stealing scheduler, we have on the right CPU capacity but is not being leveraged and we
	want to use it.

	So there's a little algorithm around work stealing. And part of the algorithm says
	"Go see if there's another `P` with a lot of work in it?, If there is then steal half of the work it has."
	In this case,  the left "P" has 4 Go routines, so quite possibly,
	the Right "P" will steal 2 of the Go routines, put one in the RQ and start executing the
	other one.
	The last thing we want is the right "M" to go into the waiting state, because it will be pulled off the hardware,
	and we will not be able to execute the GO routines as quickly as we otherwise could.

	Now say the left "P" suddenly finishes it's work, see 24.png. Now it has no work to do, may be the left "P" wants to steal
	some work out of the right "P" but there really isn't a lot.
	Another part of the algorithm says, another place you can look is the "Global Run Queue".

	In this case the left "P" will steal work from the Global Run queue, and will get that work going.

	================================================================================================================

	May be this a program running in C. 26.png
	Say we have 2 threads at OS level and they will be passing messages to each other. It's not important that
	how they'll be doing it but it's important that this is I/O bound work.

	Say T1 is first waiting for some context switch waiting to get in the running state.
	Then it does a little bit of work, and then sends the message across to T2. See 27.png

	When T1 sends the message across to the other thread, it context switches off and will be in the waiting state.
	Say T1 was on core 1 and now will be off the core. See 28.png

	On T2 when the message reaches T2, we will have an event, then there will be ctx switch on T2 and T2 will
	move form Waiting state to the running state, then there would be a context switch to move T2 to an executing state,

	Then do some work with "M" and sed that message back to T1. See 29.png.
	Then T2 will move to the waiting state again.

	Now there is an event on T1, T1 goes from waiting to the running state and then will go into the executing sate.
	And then again some work and then send message to T2.

	Every time any thread passes a message back and forth, we will have context switches. Say T2 was running on
	Core 2.
	And after than T1 ends up being on core 3. See 30.png.

	We are also going to be bouncing cores because we don't really know which core is idle at a given time.

	So you can see that with this I/O bound workload, we have context switches going from waiting to running,
	to executing. We do a little bit of work we passing messages we are moving around.
	We are bouncing off on cores.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	Now let' switch around ot he Go scheduler, the 30.png is going to remain the same
	but we are going to add a couple changes.
	Replacing Ts with Gs to represent the Go routines, G1 and G2.
	Also let's have single "P" with a single "M" running on a single hardware thread say Core 1(C1).
	See  31.png.

	Now when we do all of this on 31.png the only thing that further changes is "What core we are running on."
	Because the only core we'll be running on is core 1. see22.png

	The rest of the diagram stays the same. Because the Go routines have to transition in the same states as
	the OS threads we saw earlier.

	The only difference is that now instead of these context switches taking 1-2 micro seconds
	(1000 nano seconds - 2000 nano seconds) they are now only taking 200 nano. See 33.png
	Because the context switches are happening on the Go thread but not on the hardware.
	See 34.png

	So what GO has really done has turned the IO bound work into CPU bound work at the OS level.

	Context switches are happening at the application level on the thread.

	From the OS point of view, The "M" is ALWAYS running in an executing/running state.
	From the OS point of view we are doing CPU bound work even though it's really IO bound work.

	This is why we don't need more threads in our Go programs than we have Hardware threads
	because the scheduler is converting all of the work that we do CPU and IO bound work
	always into CPU bound work at the OS level.

	And we saw in the previous video that the most efficient work you can do, is CPU bound so also the simplest.
	Because we know that we only need one thread per hardware, and if we are doing concurrent work
	then we know know to maximize the efficiency.

	So Go has reduced a huge amount of the load off the operating system because, the OS basically just
	maintains a level of CPU bound work loads we can use a lot less threads which could cost us more.
	Since we are using less the cost is less.

	Since the context switches cost less we have gone from 1000-2000ns to 200 ns.
	Then even the transition between states happen a lot faster.

	+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	Even though we have the scheduler doing all the wonderful stuff we still have to be sympathetic with the
	scheduler, in other words we still have to know what our work loads is.
	More Go routines, then we have "M" for CPU bound workloads are not going to make us get our work done faster
	because we still have 200 ns that is 2400 instruction losses that we incur if we have context switch.

	If we have CPU bound workload like adding numbers that's still 2400 instructions loss. That's like 2400
	number we could have added. That we otherwise are not adding because of the context switches.

	when it comes to IO bound work, we will have real efficiencies because the context switch is less.
	Now we will have more go routines than we will have "M"s.

	What's really brilliant that the application thread, the Go routine is so light weight that
	for the first time in our programming model we could build a web server that could throw a Go routines 
	at every request.
	That means we could have 50,000 go routines in flight.

	So now instead of us to have a magic number of Go routines for our thread pool, 
	Go is taking off some cognitive load off our backs.
	Instead of saying let's find some magic number let's throw a Go routines at every 
	problem that comes into the server at every task, every request and let the scheduler handle it.
	
	// I didn't understand this part.
	This may not be the fastest way to solve the problems but is fast enough. Remember we need to balance
	the performance with the complexity. If it's fast enough, and you have reduced a huge amount of complexity, 
	that is a big win. If you are looking for pure performance then Go is not for you.
	You need to start looking at "C", "Assembly language", even "Rust".

	