// Sample program to show how the program can't access an
// unexported identifier from another package.
package main

import (
	"fmt"
)

func main() {

	// Create a variable of the unexported type and initialize the value to 10.
	// compiler gives us an error for this, since "alertCounter" is starting with a lower case letter.
	counter := counters.alertCounter(10)

	// ./example2.go:17: cannot refer to unexported name counters.alertCounter
	// ./example2.go:17: undefined: counters.alertCounter

	fmt.Printf("Counter: %d\n", counter)
}
