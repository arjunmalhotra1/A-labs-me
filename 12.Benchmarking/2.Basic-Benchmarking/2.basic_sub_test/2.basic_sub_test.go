/*
Here, instead of "t.Run" we have "b.Run".
"-bench ." means run all the benchmarks.
Do not run benchmarks in parallel.
Machine must be idle, when we run benchmarks.
We need a quiet machine when running benchmarks.
If it's not quiet we will not get accurate results.
"go test -run none -bench . -benchmem"
check "1.png"
*/

// go test -run none -bench . -benchtime 3s -benchmem
// go test -run none -bench BenchmarkSprint/none -benchtime 3s -benchmem
// go test -run none -bench BenchmarkSprint/format -benchtime 3s -benchmem

// Basic sub-benchmark test.
package basic_sub_test

import (
	"fmt"
	"testing"
)

var gs string

// BenchmarkSprint tests all the Sprint related benchmarks as
// sub benchmarks.
func BenchmarkSprint(b *testing.B) {
	b.Run("none", benchSprint)
	b.Run("format", benchSprintf)
}

// benchSprint tests the performance of using Sprint.
func benchSprint(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}

	gs = s
}

// benchSprintf tests the performance of using Sprintf.
func benchSprintf(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}

	gs = s
}
