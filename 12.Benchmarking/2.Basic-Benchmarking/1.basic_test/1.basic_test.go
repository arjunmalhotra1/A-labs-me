/*
Idea of benchmark is to be able to profile the code, both for CPU and memory.
We will be focussing on how well the code is performing and what are those allocations that are being produced.
What is slowing down our program is either the allocations - is causing the garbage collector to run, harder, longer or more stressful, or we don't have the algorithmic efficiencies that we otherwise could.
When we CPU benchmark a program, we are tying the Go program into the operating system on the SIGPROF event, which means every 10 milliseconds this code will stop and program counter information will be gathered and the program will start again.
So the code will run a bit slower, when we do the CPU profiles.
Compiler will build the test code and will build into a test binary. Now since the compiler is
getting involved one of the things we have to worry about are optimizations that could affect our benchmarks. So we want to make sure that the code is, as close to, as identical as possible to the production code.
We also want to make sure that we are not putting things & benchmarks where the compiler could choose to throw code away.
So everything we do, we want to make sure that our testing is accurate.

When writing the unit tests we had functions that started with the word "Test" but for "Benchmarks" we have words that start with the word "Benchmark". It also has to be in a file that ends with "_test.go"
because the testing tool will run this file.
Instead of Testing "t" we will have Testing "b" pointer as the parameter.

Everything about this test will happen inside this loop

	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}

This loop will always be from 0 to "b.N". When we run this benchmark, the testing tool will be calling into this function. The very first time the testing tool calls this function "BenchmarkSprint" "b.N" will be 1.
After that it will increase "b.N" by some order of magnitude like say 10s, testing tool will increase it and the idea, is to increase or find the value of "b.N", that matches the bench-time for this benchmark. So we will have a bench-time that set. The default bench-time is 1 second. We have changed it here to 3 seconds.
"bench-time" means that this test isn't going to end until we find a value of "b.N" that let's that loop run for the full "bench-time". Even on a 3 second "bench-time". It doesn't mean that the test will run, from beginning till the end. It will take longer because it has to find that value of "b.N".
Here below we have 2 benchmark functions, we are trying to identify which function performs better.
"Sprint" or the formatted version "Sprintf".
Question is which benchmark will run faster?
and which one will allocate less.
"Sprint" or "Sprintf" for the string.
It seems that since the "Sprintf" or the formatting function, has to do more work around formatting, it will probably going to take longer.
"go test -run none" mean that if you find any
test function then don't run them.
There are times when we can have, packages that have combination of tests and benchmarks, and if we don't tell the test tooling to not to run the tests it will automatically run the tests.
"-bench" means to run the benchmarks.
"-benchmem" flag asks to show the memory allocation.
"go test -run none -bench . -benchtime 3s -benchmem"
My output of this command is in 1.png. His output is in 1-1.png.
We see that "Sprint" ran at 76 nanoseconds per op
with 1 allocation to the heap which was worth 5 bytes.
Notice that the "Sprintf", ran a little faster.
Nano seconds faster (16 nanoseconds). So virtually they are running at the same speed.
We originally thought that since "Sprintf" is a heavier function, it actually should run slower but it's not.
We can't assume about performance. If we are guessing about performance, we are going to be wrong that is so great about this tooling. We don't have to guess. So when we are writing code, we optimize for correctness, readability and integrity of things. Then when we have apiece of code that's working, we can come in, do the benchmarking and not guess. Know for sure how things are behaving. How things are performing.

Like sub tests, the benchmark has sub benchmarking.
See "basic_sub_test.go"
*/

// go test -run none -bench . -benchtime 3s -benchmem
// go test -run none -bench BenchmarkSprint/none -benchtime 3s -benchmem
// go test -run none -bench BenchmarkSprint/format -benchtime 3s -benchmem

// go test -run none -bench . -benchtime 3s -benchmem

// Basic benchmark test.
package basic

import (
	"fmt"
	"testing"
)

var gs string

// BenchmarkSprint tests the performance of using Sprint.
func BenchmarkSprint(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}

	gs = s
}

// BenchmarkSprint tests the performance of using Sprintf.
func BenchmarkSprintf(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}

	gs = s
}
