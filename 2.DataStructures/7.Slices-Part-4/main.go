// Sample program to show how one needs to be careful when appending
// to a slice when you have a reference to an element.
package main

import "fmt"

type user struct {
	likes int
}

// Say we have crated a new application and this application tracks user likes.
// Say for now we have only 3 users in our user system.
// Say we use a slice to store those users.

func main() {

	// Declare a slice of 3 users.
	users := make([]user, 3)

	// Share the user at index 1.
	// Say, we share the user in index1 with the handler. 1.png.
	shareUser := &users[1]

	// Add a like for the user that was shared.
	shareUser.likes++

	// Display the number of likes for all users.
	for i := range users {
		fmt.Printf("User: %d Likes: %d\n", i, users[i].likes)
	}

	// Next say we get on Hacker news and we start seeing some new registrations.
	// Add a new user.
	// Since we are using a slice we decide we are going to append.
	// Since length is same as capacity this means we need new allocation.
	// Now we allocate from 3 to 6.
	// Note that even though the backing array has changed the pointer is still pointing to the original
	// backing array. 4.Png.
	// the new likes for Ed are going to the original array.
	// We have to remember that when we are dealing with slices (reference type) we are
	// moving it around using the value semantics. But anytime we are reading and writing
	// we are using the pointer semantics.
	// Suggestion -  If you see an append cal in code and if the code is not inside of
	// a function like "decode" or "unmarshall", then we should stop and check the use of append.
	// Anytime we have a situation where we have an established data structure we have been sharing it
	// and now we call append and the backing array can be replaced, it might be erroneous.
	// When using value semantic mutations we are safe. Mutations are happening in an isolated space.
	// When using Pointer semantic mutations crazy things can happen.
	users = append(users, user{})

	// Add another like for the user that was shared.
	shareUser.likes++

	// Display the number of likes for all users.
	fmt.Println("*************************")
	for i := range users {
		fmt.Printf("User: %d Likes: %d\n", i, users[i].likes)
	}

	// Notice the last like has not been recorded.
}
