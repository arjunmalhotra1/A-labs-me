// A drawback of having 2K stack is that no 2 go routines can share a value with another go routine that's on
// it's stack. So any pointers associated with a stack are going to be internal pointers to that Go routine only.
// If 2 go routines need to hsare a value, that value has to be on the heap.
// Value sematics are very powerful they allow the isolation model, they allow mutation to happen in a safe way.
// We are trying to replicate the picture attached in this folder.
// We share the string down the call stack 10 times.

// Sample program to show how stacks grow/change.
package main

// Number of elements to grow each stack frame.
// Run with 1 and then with 1024
const size = 1

// main is the entry point for the application.
func main() {
	s := "HELLO"
	stackCopy(&s, 0, [size]int{})
}

// stackCopy recursively runs increasing the size
// of the stack. Array is just to increase the size of the array for the stack to cross 2K space.
// At 1024 we get a new stack and the addresses are changed and pointers are fixed as part of the function call.
//go:noinline
func stackCopy(s *string, c int, a [size]int) {
	println(c, s, *s)

	c++
	if c == 10 {
		return
	}

	stackCopy(s, c, a)
}
