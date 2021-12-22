// Sample program to show how maps behave when you read an
// absent key.
package main

import "fmt"

func main() {

	// Create a map to track scores for players in a game.
	scores := make(map[string]int)

	// We haven't stored any data in the map yet.
	// Read the element at key "anna". It is absent so we get
	// the zero-value for this map's value type.
	// Map will return successfully and is going to return the zero value for that value type.
	// Integer in this case. Integer set to it's zero value, 0.
	score := scores["anna"]

	fmt.Println("Score:", score)

	// If we need to check for the presence of a key we use
	// a 2 variable assignment. The 2nd variable is a bool.
	// If "ok" is true then the value existed and if "ok" is false then the value didn't exist.
	score, ok := scores["anna"]

	fmt.Println("Score:", score, "Present:", ok)

	// We can leverage the zero-value behavior to write
	// convenient code like this:
	scores["anna"]++

	// Without this behavior we would have to code in a
	// defensive way like this:
	if n, ok := scores["anna"]; ok {
		scores["anna"] = n + 1
	} else {
		scores["anna"] = 1
	}

	score, ok = scores["anna"]
	fmt.Println("Score:", score, "Present:", ok)
}

/*
	There are constraints on the key. The only type that can be used as a key is something that can be used
	as a conditional operation. If we can't use something int he if statement then it can't be used as a key.
*/
