// Sample program to show how constants do have a parallel type system.
package main

import "fmt"

const (
	// Max integer value on 64 bit architecture.
	maxInt = 9223372036854775807

	// Much larger value than int64.
	// Had "bigger" been of type int this program wouldn't have compiled.
	// As we can see on line 19. But this line 16 does compile.
	// Because we have 256 bit precision at the constants level.
	// hence if we change that number to be of type int64
	// compiler complains that it overflows. on line 19.
	bigger = 9223372036854775808543522345
	// This shows constants have a parallel type system.
	// They have 256 bits of precision and they are not variables at all.

	// Will NOT compile
	// biggerInt int64 = 9223372036854775808543522345
)

func main() {
	fmt.Println("Will Compile")
}
