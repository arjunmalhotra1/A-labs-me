// All garbage collectors have the same job.
// Their job is to walk through the heap and identify the values that are no longer needed or in use and then
// to sweep that memory free so that, that memory can be re-used again.
// go has a concurrent garbage collector, so we should be able to
// get work done as well while the garbage collector is running.
// Pacer is the big important thing of the garbage collector. Pacer's job is to figure out, when to start a collection
// How long the collection is going to take. Make sure that it starts the collection at the very last moment, yet
// It can finish in time before we run out of heap space. Since this is a concurrent collector.
// "Pressure on heap" is defined as how quickly we fill the heap up in between these garbage colletion
// we have to perform.
// Say after the garbage collection 2 MEGs of memory is in-use.
// Golang's compiler has a configuration option, called
// GC (percent) which by default is equal to 100%.
// We can get to this value by using an environmental variable "GOGC".
// GC=100% means that if we end up with a 2MEG heap then we are going to size the entire heap to 4MEG.
// That is we create a 2MEG "GAP".
// Say if after the next collection we end up with a 3MEG in use space. Then we will end up with a 3MEG gap,
// which will give us a 6MEG heap.
// That's what going to happen after the end of every garbage collection.
// How much in use space do we have and then let's create a "GAP" of equal size that will represent the 
// total size of the heap.
// "Stress on heap" means how quickly does it take to fill in the remaining gap of memory that was 
// added in the last grabage collection. 
// Quicker the gap is filled quicker the garbage collection has to run again.
// Number 2 reason on why the program is not fast enough otherwise it could be is
// Internal latency 
// Grabage collection is a big reason for internal latency. 
// Grabage collector algorithms try to minimize "Stop the world" time to less than 100 micro seconds.
// There's a few points of "Stop the world" that we will be experienceing during the garbage collection.
// 1. Is to bring all the current go routines to a stop the world point on of the reasons for that is that we got to run thw write barrier on.
/* Say we have 4 logical processors.
We will have to stop all the Go routines to turn on the write barrier.
Write barrier is the littlepiece of code that has to run while the garbage collector is running.
Garbage collector is discussed in detail in the next section - 
"Garbage Collection Addendum" of the course. Which is not in this repository.
