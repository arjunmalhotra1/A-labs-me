/*
	NOTE - At places he said milli second but I think he meant micro second.
	Concurrency - Means out of order execution.
	The idea that we have a sequence of executions from beginning to the end,
	but we are going to execute them out of order.

	Parallelism  - Means we are executing 2 or more instructions at the same time. It's about doing
	a lot of work at once.

	Let's start with a very simple hardware where all we have is a single processor.
	And that processor has a single hardware thread. See 1.png

	Let's represent that hardware thread as the ability to execute instructions. Every hardware thread we have
	in our machine gives us the ability to execute instructions.
	So if we have only one hardware thread then we can only execute one instruction at a time.
	If we have 2 hardware threads then we can execute 2 instructions at a given time that means execute instructions in parallel.

	Say, we have to code an OS. Hence we decide to code a OS level scheduler.
	For now imagine our hardware is a single hardware threaded. Now we have to in our OS
	level scheduler use this hardware thread. For that we will need an OS level thread.
	Think about threads as "paths of execution". So our programs are going to be a series of instructions,
	those instructions are laid out to execute in sequence. It is the job of our OS level thread to execute
	these instructions one at a time.
	OS level threads leverage the hardware level thread to get the execution in sequence.
	Hardware thread does the actual execution.

	OS level threads can be in 3 states:
	1. Running - When a thread is in a running state that means we have placed it on the Hardware thread and
		it is executing it's instructions it's responsible for.
	2. Runnable - When a thread is in a runnable state it just means that the OS thread wants time on the HW thread.
		It just has to wait for it's turn.
	3. Waiting state -  When a thread is in the waiting state we really consider it off the radar. It's gone and it is waiting for something to happen. It could be waiting for something on the network to come back in.

	From our OS scheduler's point of view we only care about the threads that are in "Running" state or "Runnable" state.
	Thread's that are waiting are not in our concern.

	Imagine the following scenario, we have a couple of program running, that is a few threads in play.
	The idea is that we have to create this illusion even with one hw thread,
	that all the threads we have (running or runnable) are actually running at the same time.

	So how can we create this illusion with our scheduler. One idea is that we will define what we call is the
	"scheduler's period".
	Scheduler period is defined as our ability to execute any thread, that is in the runnable state
	at the beginning of this period to get that thread execute within this time and if we can do that
	then at least we can create this illusion at least within the scope of the scheduler period that everything is
	running at the same time.

	Say Scheduler period (SP) is 100 ms
	What we want is to make sure that every thread that we have in a runnable state executes within this
	100 ms time frame to create the illusion that every thing is running at the same time.

	Say when we start off the scheduler we notice that we have only one thread in the runnable state.
	so if only have one thread in the runnable state, remember we have only one hardware thread which means
	only one thread can execute at any given time.
	thread A gets the full 100 ms to execute it's instructions.

	1 nano-second = 12 instructions.(Remember this)

	Now when the scheduling period is done we go back and look at how many threads we have in our runnable state.
	Say now we have two threads in the runnable state. This means we'll have both threads to execute within
	the same 100 ms of scheduler's period.
	Now thread A gets 50ms and thread B gets 50 ms see 3.png
	We did create the illusion that two threads ran within the scope of the scheduler's period.

	But now thread A's time is cut into a half. thread A is not executing as many instructions as in the previous run.

	Say now scheduler period is over and we now have 5 threads in the runnable state.
	this means we have to cut the time into 5 parts.
	thread A gets 20 ms, Thread B get's 20 ms and other thread C thread D and thread E get the remaining 60 ms of time.
	See 4.png
	Every time more threads show up the time foreach thread gets smaller.
	If we have 100 thread to run then each thread would get 1 ms of runtime.

	At a some point we get diminishing return on the time a given thread gets to do any work.
	It's mainly because of the context switch cost.
	There is some time cost associated with to pulling the thread out of the hardware. And now Thread B will go
	into the hardware.
	This is going to take some time. This is context switch. Taking one thread off an putting other thread on.

	On average at an operating system level a context switch is going to be 1 micro second - 2 micro second.
	That is 1000ns-2000ns. This mean 12,000-24,000 instructions lost. Instructions we could otherwise
	have been executing but we can't because of context switch.

	say if we have 1000 threads. Now we are talking about micro seconds of time, a thread gets to execute on top of
	1-2 micro second of context switch.

	What eventually happens is that doing more context switch work than actual work.
	So there is a point where the time slice for each thread, where we do more context switch work than actual work.
	And we are not really getting any throughput from our machine or our apps.

	So what we will have to do in our OS scheduler is define another parameter. We define the MTS (Minimal time slice).
	Minimal time slice(MTS) is the bare minimum time slice regardless of our scheduler period time that we will
	allow any thread to at least get some run time. Say we say 10ms is the point of diminishing return.
	That is we never go below 10 ms of a time slice.

	Because if we go below 10 ms of time slice our context switch will cause us problems.

	But now 10ms MTS and 100 threads means 1000ms = 1 second of schedule period.
	1 second is too long to create our illusion.

	Say we have 10 seconds of scheduled time that is 10 seconds to create the illusion.
	So now any given thread might have to wait for 10 seconds before it even gets to run again.
	====================================================================================================

	So it will always be about "less is more". The less threads we have to schedule, during any given
	scheduler period, that is each thread will get more time and that thread will be able to back into play sooner.
	===============================================================================================================

	There are 2 types of work loads that we have to make our early on if we have to
	make smart engineering decisions with our scheduler.

	1. CPU bound - CPU bound work load is a kind of workload where a thread naturally does not, move into
	a waiting state. An example of CPU bound work load will be an algorithm counting integers.

	Say we have 1 millions integers to be added. There's nothing in the process of adding that will cause the threads
	to wait or pause for anything. This is CPU bound work load.

	This means that every thread will get it's full time slice (from beginning till the end). So far we have seen CPU bound
	work load, where if thread A gets on the hardware thread it got it's full time slice and did not switch
	until the time slice was over.
	==================================================================================================================
	Opposite of CPU bound workloads are IO bound or blocking workload.

	IO bound work loads are the workloads where the thread never use it's actual time slice.
	At some time within the scope of it's time slice, the thread will make the call to the operating system over
	the network. It will do something that wil cause it to move from the running state into waiting state.

	Even in those situations even if we have 1000 threads the chances of the scheduling period getting to 10 seconds is
	pretty slim because most likely the thread in it's initial stages maybe within it's first millisecond or so,
	will make the call to the operating system to read something over the network, and them immediately will
	get context switched out.

	Context switches are good when it comes to the IO bound workloads. Since the thread will go in the waiting state
	and will keep the hardware thread idle, we can now leverage the 12,000 to 24,000 instructions to do the context
	switch and get another thread up and running.

	When it comes to CPU bound workloads the context switches are hurting us, because we take
	12,000 to 24,000 instruction hits for every context switch.
	And these instructions hits could have been the work we otherwise we could have been doing.

	Hence understanding the work load is so important.

	"BECAUSE IF IT'S CPU BOUND WORKLOAD THEN WE DON'T REALLY WANT EVER MORE THREADS THAN WE HAVE
	HARDWARE THREADS BECAUSE CONTEXT SWITCH WILL HURT US."
	But in I/O bound workload the context switches are our friends as we are not allowing the HW threads to stay idle.

	Our OS scheduler has to be a preemptive scheduler. That is it has to handle events.
	Think about it, the hardware is dealing with events all the time.
	Just to keep the clock up to date there is an event that is occurring on the machine.
	And every time an event occurs on the machine a thread has to respond to it.

	Part of scheduler also has to create a mapping of hardware events to OS threads and setting certain operating  system
	threads to have a high priority. So even if we are doing CPU bound work load, every time the clock has to run
	we will have to interrupt the workload get the thread to run the clock quickly, do it's thing and pull it off.

	Bill doesn't want the preemptive scheduler that's event based to get in our way of the mental model. But a preemptive
	scheduler does give to us that, "scheduling is unpredictable".
	Once all things are equal we never really know what the scheduler is going to do.
	We have to keep that in our heads when we design code because, no matter what we are doing when we have
	multi threaded software there has to be a guarantee. It is our job as a developer to have the
	synchronization adn orchestration. That is putting points of guarantee in the code so that we can guarantee,
	the algorithms (even if out of order) they do run correctly.

	For now let's not think about preemptive model and only focus on CPU bound and IO bound and what do they mean.

	=================================================================================================================
	CPU BOUND

	Say we have to write an algorithm that takes in a million numbers, and our task is to
	add all these number and give the value.

	Anytime we need to do anything like this specially in a multi threaded software, the first
	thing we need to do is to write a sequential version of algorithm (Add function here)

	Sequential version says start from beginning and keep adding nad then we are done.

	The next question after writing the sequential version is always going to be
	Question. "Is this kind of work that can be done concurrently?"
	Concurrency means out of order execution
	Does it matter in what order we add these integers in?
	Will that affect the final result?
	Answer. Here we are just adding and it doesn't matter.

	So we have a problem in front of us that can be done using concurrency.

	One idea could be to split the list into half. Then have 2 OS threads to add numbers from each halves.
	Then we can take the result of those 2 threads and get the result. See 6.png

	In fact we can break the array up into multiple lists, and multiple threads, adding small sections
	of numbers up and then o the add.
	the only reason to add this level of complexity is that we can get the work done faster, we can get some
	better throughput.

	Parallelism means executing instructions at the same time 2 or more.
	So far we had our hardware that had a single hardware thread.
	Which means we could only execute 1 OS thread at a given time.

	Say we have 4 sub lists and 4 threads.
	Question. Will running this algorithm sequentially be faster or slower than running it in parallel in 4 threads
	when we have a single hardware thread.
	Answer. Because this is CPU bound and we can only execute single thread at a time.
	Running the sequential version is going to be faster.
	Because every time we have to swap the thread off the hardware, we are taking 12,000 instruction losses.
	That's only to say we can get his work done within that time slice if we didn't then it's even more time.
	We could have used these instruction lost to add the integers that we are now not.

	But what if we change our hardware to have 2 hardware thread processes.
	Now we can run 2 threads in parallel. Now, we can break the list into say 2 halves and give each thread it's own
	hardware thread. They can run in parallel.
	We can get the work done at the same time, then we come up and we get our value.

	Now if we have 4 HW threads we could use, 4 OS threads and get work done even faster.
	When it is a CPU bound workload, and if we want to get more throughput, and the algorithm is concurrent.
	then we need parallelism.

	Without parallelism there is no efficiency in concurrency or out of order execution because context switches hurt us.

	======================================================================================================================

	Let's go back to single HW thread, let's say we don't have a list of integers anymore but a list of URLs, that
	we want to go fetch and process.

	This means that this is no longer CPU bound and is now IO bound, work. Which means at any given point when a thread
	makes a network call it will move to "Waiting" state.

	But now since the task is IO bound which mean now we can break this list up into multiple threads and even if we have
	one HW thread, and get work done faster because now the context switches are helping us.

	T1 starts and gets 1ms and then it makes the network call sand it goes into the waiting state.
	Now what happens? We take our 12k instruction hit to pull the T1 off and put T2 on the HW thread.
	As if not the HW thread would have been idle. Now T2 runs.
	Now all the threads are being able to get the work done while the processor was going to be idle
	because all the threads said "I don't need time anymore".

	But here's the interesting part about IO bound workloads, it's really complicated to figure out
	what's the most efficient number of threads, in other words there's still a point of diminishing return
	if we use too less or we use too more.

	Even if we have too less we will have a problem. We have to find what we call the magic number, the idea is that
	is there a magic number where, there's only one thread in the runnable state that's what we are looking for.

	We only have one runnable thread in the runnable stage to go.

	================================================================================

	We know for any CPU bound workload the number of threads we use, equals the number of HW threads we use.
	Their parallelism is important.

	But in I/O bound work load we don't need parallelism to see improvement with concurrency, because the OS threads
	can take turns on the HW threads and the context switches are helping us.
	Helping us to get more work done over time again we don't know how many.

	So when it comes to OS level stuff like this we always use "thread pools",
	imagine you are building a web service, and you decide to create an operating system thread for every
	request that came into the server and we have 50,000 requests coming in.
	This would cause a lot of problems.
	So we would want to have a pool of threads, that wwe would post work into and we have to find that magic number.
	We have to find the number where, if we use too many threads now the context switches are hurting us.
	That means we have too many threads in the runnable state, if we have too low of a number they are also
	now going to perform very well, because now we will have  more idle time on the processor as we will not have
	a thread ready to go.
	=================================================================================

	Concurrency means out of order execution.
	Parallelism means executing 2 or more instructions at the same time.
	We have to know what state of thread is in and we have to know if our work load is IO bound or CPU bound.
	CPU bounds work loads mean that the threads don't naturally move into waiting state.
	Which means we need parallelism to get any throughput with concurrency. I/O bound work loads mean
	that threads do naturally move into waiting state that means we do not need parallelism for the concurrency to give
	us throughput. The context switches help us get more work done even on a single threaded hardware.

	===============================================================================
	CPU bound workloads are the most efficient workloads we can do on the machine.

	We will now talk about Go scheduler that sits on top of the OS scheduler that is sitting on top of the hardware.



















*/