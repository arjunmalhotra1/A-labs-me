/*
Benchmarks can be misleading. It's very important that we validate our results and that we have some reasonable expectation. If a benchmark is not correct then we could be making some really bad engineering decisions.

Here we have 1 million numbers and 3 functions,BenchmarkSingle,BenchmarkUnlimited & BenchmarkNumCPU.
The 3 functions implement the mergeSort.
BenchmarkSingle - will use one go routine sequentially from the beginning to the end.
BenchmarkUnlimited - Every time we split the list in half we will throw a go routine at each half.
We will have as many go routines as there are lists to sort.
BenchmarkNumCPU - We will only have as many lists as our number of CPUs. So we wil not split the list in more than 8 parts.
We eventually only split the list in 8 parts and then do the sorting.
We run this by
"go test -run none -bench ."
See "1.png" for results.
"Single" (sequentially) took 94 ms.
"Unlimited" took significantly longer(probably in seconds.)
"CPU-8" took 114 ms.
But now Bill ran only the NumCPU.
"go test -run non -bench NumCPU"
Results now were different.
see 2.png
suddenly "NumCPU" is running faster than a single go routine.
What's happening here is when we run all the benchmarks together, we are violating a major rule of thumb, which is "A machine must be idle."
The test for "unlimited" used som many go routines, and so much chaos that all of that hadn't been cleared out yet by the time "NumCPU" benchmark ran.
Therefore the cleanup from "unlimited" were effecting the performance for the "CPU-8".
Once we let the machine go idle and let "NumCPU" run by itself everything started to be more accurate.
This doesn't mean that we have to always run every benchmark in isolation. But, if we know we have a benchmark like "unlimited", it's going to put a lot of stress on the machine. We absolutely don't have to run that in conjunction with other benchmarks.
This is a classic example where we have to have some expectation, we have to validate the result.
If the result don't seem right either then something is wrong, could be assumptions, or the machine isn't idle. Hence we should always validate the benchmarks.
*/

// Sample program to show you need to validate your benchmark results.
package main

import (
	"math"
	"runtime"
	"sync"
	"testing"
)

// n contains the data to sort.
var n []int

// Generate the numbers to sort.
func init() {
	for i := 0; i < 1000000; i++ {
		n = append(n, i)
	}
}

func BenchmarkSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		single(n)
	}
}

func BenchmarkUnlimited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unlimited(n)
	}
}

func BenchmarkNumCPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numCPU(n, 0)
	}
}

// single uses a single goroutine to perform the merge sort.
func single(n []int) []int {

	// Once we have a list of one we can begin to merge values.
	if len(n) <= 1 {
		return n
	}

	// Split the list in half.
	i := len(n) / 2

	// Sort the left side.
	l := single(n[:i])

	// Sort the right side.
	r := single(n[i:])

	// Place things in order and merge ordered lists.
	return merge(l, r)
}

// unlimited uses a goroutine for every split to perform the merge sort.
func unlimited(n []int) []int {

	// Once we have a list of one we can begin to merge values.
	if len(n) <= 1 {
		return n
	}

	// Split the list in half.
	i := len(n) / 2

	// Maintain the ordered left and right side lists.
	var l, r []int

	// For each split we will have 2 goroutines.
	var wg sync.WaitGroup
	wg.Add(2)

	// Sort the left side concurrently.
	go func() {
		l = unlimited(n[:i])
		wg.Done()
	}()

	// Sort the right side concurrenyly.
	go func() {
		r = unlimited(n[i:])
		wg.Done()
	}()

	// Wait for the spliting to end.
	wg.Wait()

	// Place things in order and merge ordered lists.
	return merge(l, r)
}

// numCPU uses the same number of goroutines that we have cores
// to perform the merge sort.
func numCPU(n []int, lvl int) []int {

	// Once we have a list of one we can begin to merge values.
	if len(n) <= 1 {
		return n
	}

	// Split the list in half.
	i := len(n) / 2

	// Maintain the ordered left and right side lists.
	var l, r []int

	// Cacluate how many levels deep we can create goroutines.
	// On an 8 core machine we can keep creating goroutines until level 4.
	// 		Lvl 0		1  Lists		1  Goroutine
	//		Lvl 1		2  Lists		2  Goroutines
	//		Lvl 2		4  Lists		4  Goroutines
	//		Lvl 3		8  Lists		8  Goroutines
	//		Lvl 4		16 Lists		16 Goroutines

	// On 8 core machine this will produce the value of 3.
	maxLevel := int(math.Log2(float64(runtime.NumCPU())))

	// We don't need more goroutines then we have logical processors.
	if lvl <= maxLevel {
		lvl++

		// For each split we will have 2 goroutines.
		var wg sync.WaitGroup
		wg.Add(2)

		// Sort the left side concurrently.
		go func() {
			l = numCPU(n[:i], lvl)
			wg.Done()
		}()

		// Sort the right side concurrenyly.
		go func() {
			r = numCPU(n[i:], lvl)
			wg.Done()
		}()

		// Wait for the spliting to end.
		wg.Wait()

		// Place things in order and merge ordered lists.
		return merge(l, r)
	}

	// Sort the left and right side on this goroutine.
	l = numCPU(n[:i], lvl)
	r = numCPU(n[i:], lvl)

	// Place things in order and merge ordered lists.
	return merge(l, r)
}

// merge performs the merging to the two lists in proper order.
func merge(l, r []int) []int {

	// Declare the sorted return list with the proper capacity.
	ret := make([]int, 0, len(l)+len(r))

	// Compare the number of items required.
	for {
		switch {
		case len(l) == 0:
			// We appended everything in the left list so now append
			// everything contained in the right and return.
			return append(ret, r...)

		case len(r) == 0:
			// We appended everything in the right list so now append
			// everything contained in the left and return.
			return append(ret, l...)

		case l[0] <= r[0]:
			// First value in the left list is smaller than the
			// first value in the right so append the left value.
			ret = append(ret, l[0])

			// Slice that first value away.
			l = l[1:]

		default:
			// First value in the right list is smaller than the
			// first value in the left so append the right value.
			ret = append(ret, r[0])

			// Slice that first value away.
			r = r[1:]
		}
	}
}
