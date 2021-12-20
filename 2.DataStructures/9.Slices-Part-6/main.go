// Sample program to show how the for range has both value and pointer semantics.
package main

import "fmt"

func main() {

	// Using the value semantic form of the for range.
	// Friends slice will look like 1.png.
	// Since this is the value semantic for range it means that this
	// range loop will be operating against it's own copy of the friends slice
	// 2.png
	// When we reach - "friends = friends[:2]"
	// this says let's create a new slice. We end up in a situation like in 3.png.
	// We are no slicing over the original but we are slicing over the copy of "friends" that was created
	// by the "range".
	friends := []string{"Annie", "Betty", "Charley", "Doug", "Edward"}
	for _, v := range friends {
		friends = friends[:2]
		fmt.Printf("v[%s]\n", v)
	}

	// Using the pointer semantic form of the for range.
	// In the pointer semantics we were looping over the original array and in the middle of the loop
	// we changed the original slice length from 5 to 2.
	// 4.png
	// Now there is a problem because now if we keep iterating and we check the length, we end up
	// beyond the bounds of the slice.
	// See 5.png for the output.
	// The idea is that the "value semantic" protects us from the mutations, because in value semantics
	// that is done in isolation. But while using a pointer semantics we have to be careful.
	// In value semantics we are iterating over the copy of the slice values but in the pointer semantics
	// we are iterating over the original.
	friends = []string{"Annie", "Betty", "Charley", "Doug", "Edward"}
	for i := range friends {
		friends = friends[:2]
		fmt.Printf("v[%s]\n", friends[i])
	}
}
