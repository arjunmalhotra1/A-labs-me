// Sample program to show how the for range has both value and pointer semantics.
package main

import "fmt"

func main() {

	// Using the pointer semantic form of the for range.
	friends := [5]string{"Annie", "Betty", "Charley", "Doug", "Edward"}
	//Literal construction of String array.
	fmt.Printf("Bfr[%s] : ", friends[1]) // "Betty"

	// Pointer semantics.
	for i := range friends {
		friends[1] = "Jack" // Change betty to Jack for i=0.
		// Change "Jack" to "Jack" itself for i=1.

		if i == 1 {
			fmt.Printf("Aft[%s]\n", friends[1]) // Prints "Jack"
		}
	}
	// fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n")

	// Using the value semantic form of the for range.
	friends = [5]string{"Annie", "Betty", "Charley", "Doug", "Edward"}
	fmt.Printf("Bfr[%s] : ", friends[1]) // Displays "Betty"

	for i, v := range friends {
		friends[1] = "Jack" // Change "Betty" to "Jack".

		if i == 1 {
			fmt.Printf("v[%s]\n", v) // This will still say "Betty".
			// Why?
			// Value Semantics mean that every piece of code is operating on it's own copy of the data.
			// We are not talking about "v" though "v" is an aspect of the fact that we are making a copy.
			// My personal comments *** Think of Range as it's own function call in the stack with it's own copy of friends.
			// It is this copy of friends in the range function that is returning i and v.
			// But the friends we access inside the for loop (line29) is the original friends of the main function
			// And not the copy that is in the range function call stack.
			// Refer this code - https://goplay.space/#MZu_Oa_IGHU
			// Also worthy to note - https://goplay.space/#_3h7b7QBP9I
			// Also note here, this is interesting case - https://goplay.space/#LJH8Yt-VUS3
			// *** My personal comments
			// It's about the fact that the for range loop on line 28 isn't iterating over "friends" array
			// It is actually iterating over it's own copy.
			// We absolutely change "Betty" to "Jack" but that's not the array we are iterating over.
			// We are iterating over a copy of friends at this point.
			// So any changes to friends is irrelevant in terms of this iteration.
			// So when we say that it's the value semantic form of the for range.
			// We are saying that the for range is operating on it's own copy of the data.
			// "v" is the copy of the copy.

		}
	}

	// Using the value semantic form of the for range but with pointer
	// semantic access. DON'T DO THIS. Don't use "&freinds".
	friends = [5]string{"Annie", "Betty", "Charley", "Doug", "Edward"}
	fmt.Printf("Bfr[%s] : ", friends[1])

	for i, v := range &friends {
		friends[1] = "Jack"

		if i == 1 {
			fmt.Printf("v[%s]\n", v)
		}
	}
}

// OUTPUT
// Bfr[Betty] : Aft[Jack]
// Bfr[Betty] : v[Betty]
// Bfr[Betty] : v[Jack]
// Note he is printing on line 32 v and not friends[1], "fmt.Printf("v[%s]\n", v)"
// But tried something here:
// https://goplay.space/#nmxTvdrEGRd

// Original code but printing friends[v] in the value semantics.
// https://goplay.space/#obTIEw7Bvas
