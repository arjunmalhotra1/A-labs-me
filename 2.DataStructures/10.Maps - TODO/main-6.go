/*
	Double says that I want my own copy of the map.
	Had it been "*map[string]int" the code review would  have stopped.
	func double(scores *map[string]int, player string)

	Remember we move the reference type around the program using the value semantics.
	So double will get it's own copy of the map.

	Note - When it comes to maps, channels and functions.
	That map value, channel value and function value is just a pointer. So we are already just
	copying a pointer or making a copy of that address to tha map internal data structure.
*/

// Sample program to show how maps are reference types.
package main

import "fmt"

func main() {

	// Initialize a map with values.
	scores := map[string]int{
		"anna":  21,
		"jacob": 12,
	}

	// Pass the map to a function to perform some mutation.
	double(scores, "anna")

	// See the change is visible in our map.
	// So even thought we are using the value semantics to move this map around our program,
	// remember any sought of raeding or writing we are using the pointer semantics.
	// So when we get to the multi threaded code we need to be cautious that accessing a map is not safe through
	// multi threads.
	// It has to be synchronized. Map is not safe in multi threaded we need to take care of that.
	// We are using pointer semantics when we read and wite to a amap.
	fmt.Println("Score:", scores["anna"])
}

// double finds the score for a specific player and
// multiplies it by 2.
func double(scores map[string]int, player string) {
	scores[player] = scores[player] * 2
}
