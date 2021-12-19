// go test -run none -bench . -benchtime 3s
// "none" says not to run any test functions, but we wanat to run all the benchmark functions.
// Instead of "none" we could have used any test functions if we would have liked.
// Or any regular expression.
// It's good to run the benchmark functions a few times. AS the machine might not be idle.

// Row traversal is almost twice as fast as the column traversal.

// Tests to show how Data Oriented Design matters.
package caching

import "testing"

var fa int

// Capture the time it takes to perform a link list traversal.
// Becnhmarking is going to be relative not just to the Hardware but to also what's running in
// the machine. So when we run bechmarks the machine has to be idle.
func BenchmarkLinkListTraverse(b *testing.B) {
	var a int

	// This loop is executed b.N number of times.
	// Question is b.N?
	// There is not a set value for b.N. We have a set bench time.
	// Default bench time is 1 second. We will increase that to 3 seconds.
	// What bench time measn is that find a value of b.N such that we run this loop for full bench time.
	// In other words we will ask the bench mark funciton to find b.N that will run for 3 seconds.
	// b.N will start at 1. We will increase it at an increment of 10s and eventually finds a number
	// where the loop runs for the full 3 seconds.
	// Once that happes we will get some bench marking information.
	for i := 0; i < b.N; i++ {
		a = LinkedListTraverse()
	}

	fa = a //  This is jsut to make sure that there is no chance that the compiler ends up removing /
	// optimizing the code.
}

// Capture the time it takes to perform a column traversal.
func BenchmarkColumnTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = ColumnTraverse()
	}

	fa = a
}

// Capture the time it takes to perform a row traversal.
func BenchmarkRowTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = RowTraverse()
	}

	fa = a
}

/*
BenchmarkLinkListTraverse-8          132          33530739 ns/op
BenchmarkColumnTraverse-8             22         162789441 ns/op
BenchmarkRowTraverse-8               235          15519772 ns/op
*/
// Question how is it possible that the Row traversal is so much faster than the Column traversal.
// This is more to do witht he hardware. Mechanical sympathy issue - We recognize that the hardware
// is a platform and if we are not sympathetic with the hardware then it really doesn't matter
// what we do with our code.
// Here's a video of the different latency cost that we will have in our hardware.
// https://www.youtube.com/watch?v=WDIkqP4JbkE&t=1129s
// In pic 1.png
// 4 cores
// 2 harware threads per core
// Each core has it's own L1 and L2 cache.
// Hardware thread can only access data if it's in the L1 and the L2.
// In pic 2.png
// When we say we have a 3GHz clock we mean that we get 3 clock cycles evry nano seconds.
// clock is pushing electrons within that machine.
// Nothing can happen without that heart beat. That clock beat drives what's happening in the computer.
// With 3GHz clock we are saying that 3 heart beats every nano seconds.
// We should be able to execute on average 4 instructions per clock cycle.
// It's a good reasonable number to use. That means 12 instructions per nano cycles.
// That the number we will use to understand latency.
// For every micro second of time in latency we are losing 1200 instructions.
// For every 1 milisecond in latency we are losing appx 12 million instrucitons.
// Hence we say the network latencies costing minimum 1 ms in latency and
// OS/system call costing 1 millisecond of latency.
// Latency of 1 second costs us a billion of instructions lost.
// Now we can relate to the animation, in pic 2 we can see
// When we are getting data out of L1 cache we are losing 4 clock cycles
// to wait for the data to transfer out. That means we have stalled for 16 instructions.
// L2 we now have 12 clock cycles of latency. We stall for 48 instructions.
// At l3 we get even more latency. 40 cycles of latency and 160 instructions lost.
// In main memory we have 100 clock cycles of latency and 400 instructions lost.
// Industry defined latency in pic 3.png.
// Everytime we need data and if it is not in one of the caches. We take 100 ns or 1200 instrucitons hit.
// To get it off the main memory.
// We could have used these 1200 instructions hit to getting some work done.
// Only way to hide / fix the 1220 instructions cost is to get the data closer to the hardware threads.
// That's where caching was born to get the data clsoer to the hardware threads.
// Hardware also has special programs srunning like the "Pre-fetcher".
// Pre fetcher's job is to look at the data access patterns and try to figure out what data is needed
// before it's requested.It is also the part of the caching system.

/* In pic 1.png
Say the thread 0, T0 in core 2 wants the letter "M".
Hardware could take 1200 instruction hit to pull the letter in the cache but what would be bad is if the
next instruction neede the letter "a" in the word "Main",
we go back to the memory and take 1200 instrucitons hit again. It could be really bad if
every time we needed the next byte we take 1200 instruction hit loss.
So what happens is dife=ferent, the granularity is not down to a byte.
What cahcing system do is, use the granularity of 64 bytes. So we have a 64 byte cache line.
So when we wanted the letter "M" we will pull the letter "M" and the other 63 bytes. So "ain" of "Main"
are already coming along with it. With cahce lines moving on l1,l2,L3 AND L4 we would be taking 1200
instructions miss. We can have pre-fetchers help create preditable patterns in memory.
Pre fetcher has to take a look at the data access and
see if there's something predictable about the next one.
So how do we make predicatble access pattern to memory?
Best way is to allocate contigous block of memory.
If we can allocate a contigous block of memory and walk down that memory in a predctable stride.
Pre-fetcher can pick up on that. What data structure allows us to alocate contigous block of memory?
"Arrays". Arrays are the most important data structure as it relates to the hardware today.

Another cache not in the diagrams is "TLB"
Transalation lookaside buffer cache.
TLB is a little tabel that is able to convert virtual memory addresses to physical addresses in RAM
by caching operating system pages and offsets of where these virtual addresses are. Rememeber our software
works with virtual memory and not physical memory. So all the addresses hardware is dealing with are the
virtual memory addresses and there has to be a translation that is what TLB does.
So besides L1,L2,L3 and L4 caches being up-to-date we want the TLB to be up to date as well.
Worst case scenario is that TLB is not up-to-date because we are not creating a predictable access pattern.
When that happens the Harware has to ask the OS
"Hey! I don't know anyting about this virtual address can you please share with me where it is."
Now we to start scanning the paging tables at the operating system level.
Cloud machines/virtual machines have their own paging tables as well.
In a worst case scenario we will haev to scan the paging table on the VM. That has to scan the OS paging
table as well. Plus we have the cache line misses. We haev to talk about the major latency cost.
It can get very bad. Predictable acces patterns help us in both cache lines and the TLB.


Comin back to why the row was faster than column.

BenchmarkLinkListTraverse-8          132          33530739 ns/op
BenchmarkColumnTraverse-8             22         162789441 ns/op
BenchmarkRowTraverse-8               235          15519772 ns/op

When we are traversing the matrix by row. We are walking cache line by connected cacheline 
in that predictable stride. pic 4.png
That is so fast not because of the algorithm but becasue of the data access.
the pre fetcher is picking up our traversal and is caching these cache lines before we need them.

Column traversal is bad because, the elemnt in row 1 col 1 and the element in row 2 col 1 are in a 
totally seperate pages (since our matrix is quite big.) See pic 5.png
so every traversal in column is spanning pages so the TLB cannot be upto date which means we cannot have
any predictable access pattern to memory at all.
We are here with column traversal looking at random memory access. We are getting no mechanical 
sympathies form the hardware.
Some linux distributions use 2 MB page size to help keep more data at page level to minimize 
with these TLB misses.

Question why is linked list still faster then column but fairly close?
We can guarantee a couple of things:
1. Each of these linked list node probably fall on a different cache line.
When we load one cahce line we have only 1 out of the 64 bytes that are useful, 
and there is no predictable stride.
2. What's probably okay is that a lot of these nodes/bytes are on the same operating system page.
They are not getting slammed on the TLB misses. That is what probably making the linked list 
traversal faster than the column traversal.

What we are learning here is that sometime we have or we do not have mechanical sympathy with the hardware.
That is going to make bigger different than the algorithmic efficiency.
Performance today on the hardware that we use is not about pushing the clock, we do not wnat a 5GHz clock.
But they get real hot really fast. Even the 3 GHz might so,etime work at 4GHz but they will not 
run that speed very long because of the heat. So it can't be about making the clock run faster.
What it has to be about is "Getting the data into the processor more efficiently".
How efficient we can get the data into the processor. 
With row traversal we are as eficient as we can be
we have the predictable access pattern more of those data is on those cache lines,
one cache line that we bring in has a lot of the data for the next calls.  

Why does go only has arrays and slices and map?
Go is a language where the machine is our model. We do not have a virtual machine like in 
Java or C# that can come in and re organize the data to make it mechanically sympathetic.
So Go said instead of giving the users the data structures that are not 
mechanically sympathetic we will give the users the data structures that are mechanically sympathetic.
Array, slice(dyamic array use it to get the dynamic aspect which we do not get our of an array.)
But we get mechanical sympathy of efficiency in getting the data into the processor.
Map also tries to keep the data as contigious as possible.


Go could have used smaller stacks and then tried to link them 
(segmented stacks, where first stack stayed in place and new allocations linked together) 
but as we learnt from the hardware perspective that contigous memory is alwyas going to be better
so we will take the cost on our stack to allocate a new block to keep the data contigous 
to give us the best chance of the predictabel access pattern.

In the end it is all about the contigous data, the predictable access pattern,
it's about being efficient with getting data into the processor that's why Go ahs only gives us these 
3 data structures - arrays, slices and maps.
	

