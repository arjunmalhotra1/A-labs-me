// Sample program to show how to declare, initialize and iterate
// over a map. Shows how iterating over a map is random.
package main

import "fmt"

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

func main() {

	// Declare and initialize the map with values.
	// We are constructing here not by using make but literal construction.
	users := map[string]user{
		"Roy":     {"Rob", "Roy"},
		"Ford":    {"Henry", "Ford"},
		"Mouse":   {"Mickey", "Mouse"},
		"Jackson": {"Michael", "Jackson"},
	}

	// Iterate over the map printing each key and value.
	// Ranging over a map is random.
	for key, value := range users {
		fmt.Println(key, value)
	}

	fmt.Println()

	// Iterate over the map printing just the keys.
	// Notice the results are different.
	for key := range users {
		fmt.Println(key)
	}

	// The order of the ranging over a map is random.
	// See pic 1.png.
}
