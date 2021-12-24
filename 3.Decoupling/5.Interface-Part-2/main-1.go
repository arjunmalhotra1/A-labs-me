// Sample program to show how you can't always get the address of a value.
package main

import "fmt"

// duration is a named type with a base type of int.
type duration int

// notify implements the notifier interface.
func (d *duration) notify() {
	fmt.Println("Sending Notification in", *d)
}

func main() {
	/*
		Here we take the literal value of "42", convert it to value of type duration and call the
		notify method. Compiler throws an error. See 3.png
		As it cannot take the address of 42 even though it has been converted to a value of type duration.
		Why?
		In constants we talked that the literal value of "42" is a constant of kind int.
		Even though we have converted it to a constant of type duration it is still a constant which means
		it only exists in compile time.
		If it only exists in the compile time it is not on the stack nor on the heap. If it's not in memory
		it's not addressable.
		So the side affect how we implement constants created a situation where we can't
		necessarily take the address of every value that we are working with.
		Integrity is about having something 100% of the time or not at all.
		So the compiler from minor perspective is saying "I can't attach the pointer receiver methods to values
		because I can't assume that a 100% of the time I can get the address of that value of type 'T'"
	*/
	duration(42).notify()

	// ./example3.go:18: cannot call pointer method on duration(42)
	// ./example3.go:18: cannot take the address of duration(42)
}
