// Sample program to show how to walk through a map by
// alphabetical key order.
package main

import (
	"fmt"
	"sort"
)

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

func main() {

	// Declare and initialize the map with values.
	users := map[string]user{
		"Roy":     {"Rob", "Roy"},
		"Ford":    {"Henry", "Ford"},
		"Mouse":   {"Mickey", "Mouse"},
		"Jackson": {"Michael", "Jackson"},
	}

	// Pull the keys from the map.
	var keys []string // We make this for the sort.
	for key := range users {
		keys = append(keys, key)
	}

	// Sort the keys alphabetically.
	// See 2.png
	// Sorts the slice, keys.
	sort.Strings(keys)

	// Walk through the keys and pull each value from the map.
	for _, key := range keys {
		fmt.Println(key, users[key])
	}
}
