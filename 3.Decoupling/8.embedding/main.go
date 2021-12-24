// Sample program to show how what we are doing is NOT embedding
// a type but just using a type as a field.
package main

import "fmt"

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements a method notifies users
// of different events.
func (u *user) notify() {
	fmt.Printf("Sending user email To %s<%s>\n",
		u.name,
		u.email)
}

// admin represents an admin user with privileges.
type admin struct {
	person user // NOT Embedding. This is just a field based on a different concrete type.
	level  string
}

func main() {

	// Create an admin user, using the struct literal form.
	ad := admin{
		person: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
	}

	// We can access fields methods.
	ad.person.notify()
}

// No embedding in this code. We make a few changes to this code in "main-1.go".
