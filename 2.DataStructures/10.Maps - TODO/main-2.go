/*
	There are constraints on the key. The only type that can be used as a key is something that can be used
	as a conditional operation. If we can't use something int he if statement then it can't be used as a key.
*/

// Sample program to show how only types that can have
// equality defined on them can be a map key.
package main

import "fmt"

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

// users defines a set of users.
type users []user

/*
	Slice can't be used as a key because you can't compare a slice as it is.
	Hence users as a key is invalid. As we cannot use it in a conditional logic.

*/
func main() {

	// Declare and make a map that uses a slice as the key.
	u := make(map[users]int)

	// ./example3.go:22: invalid map key type users

	// Iterate over the map.
	for key, value := range u {
		fmt.Println(key, value)
	}
}
