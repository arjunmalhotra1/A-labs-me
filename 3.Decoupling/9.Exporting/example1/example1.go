// Sample program to show how to access an exported identifier.
package main

import (
	"fmt"

	// Remember import paths are physical locations on the disk.
	"github.com/ardanlabs/gotraining/topics/go/language/exporting/example1/counters"
	//"github.com/arjunmalhotra1/A-labs-me/3.Decoupling/9.Exporting/example1/counters/counters.go"
)

func main() {

	// Create a variable of the exported type and initialize the value to 10.
	/*
		Because there is a capital letter associated with "AlertCounter". We can take the value of kind in 10.
		Convert it into "counters.AlertCounter"
		Also notice that the name of the package also becomes the name space for accessing inside that API.
	*/
	counter := counters.AlertCounter(10)

	fmt.Printf("Counter: %d\n", counter)
}
