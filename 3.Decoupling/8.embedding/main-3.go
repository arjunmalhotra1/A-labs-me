/*
	We have the concrete type "user" implementing the notifier interface using pointer semantics.
	But we also have the admin implement the notifier interface.
	so now we have two implementations of the notifier interface for admin, one from the inner type
	and one from the outer type.
	But now that we have th outer type implement the same interface it overrides the promoted implementation
	from "user".
	Admin's implementation will always takes precedence.
*/

// Sample program to show what happens when the outer and inner
// type implement the same interface.
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

// notify implements a method notifies admins
// of different events.
func (a *admin) notify() {
	fmt.Printf("Sending admin Email To %s<%s>\n",
		a.name,
		a.email)
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

	// Send the admin user a notification.
	// The embedded inner type's implementation of the
	// interface is NOT "promoted" to the outer type.
	/*
		Now it will be the Admin's implementation that will be getting executed not the
		inner value user's implementation.
	*/
	sendNotification(&ad)

	// We can access the inner type's method directly.
	// We didn't lose access to the "user's" behavior.
	ad.user.notify()

	// The inner type's method is NOT promoted.
	ad.notify()
}

// sendNotification accepts values that implement the notifier
// interface and sends notifications.
func sendNotification(n notifier) {
	n.notify()
}

/*
	As long as there are no ambiguity the compiler will not care about these multiple implementations.
	Promotion doesn't happen automatically, promotion only happens when we try to access something
	via that promotion.

	Question. What if admin had another type say "customer" embedded, say "customer".
	type admin struct {
		user
		customer
		level string
	}

	And if customer also implemented the "notify" behavior, would we have a problem?

	Answer. No. We would only have a problem if we try to call notify in a way the compiler didn't know
	which one to use. Do I use user or customer.
	In this case since "admin" has overwritten anyway, we wouldn't really have any ambiguity or any issues.

	But in a code review if Bill saw this idea of "multiple Embeddings"- user and customer the code review would stop as
	this is not a good design pattern that we would be using in Go.
	Inner types can have inner types that can have inner types.
	Then also code reviews would stop as these are not the patterns we would want to use in Go.

	We are looking for readability, easy to understand and not easy to do.
	Such patterns don't help understand readability.

	It's again about "What a piece of data can do and not about what a piece of data does."
	We see it here in Embeding. We wil be leveraging embedding when we use composition.

*/
