// Sample program to show how embedded types work with interfaces.
package main

import "fmt"

// notifier is an interface that defined notification
// type behavior.
type notifier interface {
	notify()
}

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements a method notifies users
// of different events.
// "user" already has the notifier method based on the pointer semantics.
// We can sya that the "user" type implements the notifier interface using the pointer semantics.
func (u *user) notify() {
	fmt.Printf("Sending user email To %s<%s>\n",
		u.name,
		u.email)
}

// admin represents an admin user with privileges.
type admin struct {
	user
	level string
}

func main() {

	// Create an admin user.
	ad := admin{
		user: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
	}

	/*
		Thanks to inner type promotion, admin also implements the notifier interface
		because the method notify promotes to the outer type.

		That is why we do not need subtyping in Go. It's not about making a "User" an "Admin".
		It's not about what this data is, it's about what the data does, always.
	*/
	// Send the admin user a notification.
	// The embedded inner type's implementation of the
	// interface is "promoted" to the outer type.
	sendNotification(&ad)

	// User implemented the notifier interface using the pointer semantics that means we can share
	// admin value as well thanks to the inner type promotion. And inner type implementation is what
	// gets executed when we call "n.notify()"
}

// sendNotification accepts values that implement the notifier
// interface and sends notifications.
// This is a "Polymorphic function"
func sendNotification(n notifier) {
	n.notify()
}
