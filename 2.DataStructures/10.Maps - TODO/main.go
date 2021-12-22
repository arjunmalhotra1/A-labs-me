// Sample program to show how to initialize a map, write to
// it, then read and delete from it.
package main

import "fmt"

// user represents someone using the program.
type user struct {
	name    string
	surname string
}

func main() {

	// Declare and make a map that stores values
	// of type user with a key of type string.
	// If we have a map that is set to it's zero value we cannot use it.
	// We will have problems reading and writing to it.
	// So a map has to be made.
	// You might see people use empty literal construction for the map, again
	// user := map[string]user{}
	// Again we might want to avoid empty literal construction when we are settign to zero value.
	// Hence we are using the make call.
	// Note "user" here is value semantic. Map always stores it's own copy of whatever the data is.
	// Also it's the copies that are returned.
	users := make(map[string]user)

	// Add key/value pairs to the map.
	users["Roy"] = user{"Rob", "Roy"}
	users["Ford"] = user{"Henry", "Ford"}
	users["Mouse"] = user{"Mickey", "Mouse"}
	users["Jackson"] = user{"Michael", "Jackson"}

	// Read the value at a specific key.
	// We get a copy and not the original value that was stored in line 31.
	mouse := users["Mouse"]

	fmt.Printf("%+v\n", mouse)

	// Replace the value at the Mouse key.
	users["Mouse"] = user{"Jerry", "Mouse"}

	// Read the Mouse key again.
	fmt.Printf("%+v\n", users["Mouse"])

	// Delete the value at a specific key.
	delete(users, "Roy")

	// Check the length of the map. There are only 3 elements.
	fmt.Println(len(users))

	// It is safe to delete an absent key.
	delete(users, "Roy")

	fmt.Println("Goodbye.")
}
