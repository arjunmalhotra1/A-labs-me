// Sample program to show how to understand method sets.
package main

import "fmt"

// notifier is an interface that defines notification
// type behavior.
type notifier interface {
	notify()
}

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements the notifier interface with a pointer receiver.
/*
	concrete type user implements notify using pointer semantics.
	"interface" are value less types and structs are concrete data types.

	interface area valueless type and structs are concrete data types.
*/
func (u *user) notify() {
	fmt.Printf("Sending User Email To %s<%s>\n",
		u.name,
		u.email)
}

func main() {

	// Create a value of type User and send a notification.
	u := user{"Bill", "bill@email.com"}

	// Values of type user do not implement the interface because pointer
	// receivers don't belong to the method set of a value.

	/*
			sendNotification(u)	- gives us an error, saying that the "user" value
			doesn't implement the notifier interface.
			"user" data doesn't exhibit the data that we need.

			This has to do with a set of rules, defined in the spec about method sets.
			"Method sets" define what method behavior that you are defining gets attached to
			either the values or pointers.
			So not all the methods that we define get attached to the data.

			*** Spec definition ***
				If the value we are working with is a value of some type "T". Like user on line 34.
				Then only those methods implemented with the value receivers are going to be attached to that value.

				If the data we are workign with is an address, we are sharing something, "*T" (pointer of T)
				Then the pointer semantic methods and the value semantic methods they get attached to the address
				data.
				See 1.png

			So when we are working with the address of user, all the methods we define get attached with the
			address. But only the methods that are attached to value semantics get attached with values.

		That's the error we are getting below. with "sendNotification(u)"
		We have implemented the notifier interface using the pointer semantics. See pic 2.png
		So the only method we have in play right now is the pointer semantic method.
		But we are passing a value of type user. "sendNotification(u)"

		Question. Why is the compiler doing this?
		Answer. The compiler loves you. It is trying to help us solve a massive integrity issue.
		There are 2 integrity issues here.
		1. Minor - "Not every value that we work with has an address"
		If this is true then that value, that doesn't have an address can be shared.
		So if it can't be shared it can't be used in methods that use pointer receivers.

		See main-1.go & main-2.go
		******************************************************************************************
		Integrity is about having something 100% of the time or not at all.
		So the compiler from minor perspective is saying "I can't attach the pointer receiver methods to values
		because I can't assume that a 100% of the time I can get the address of that value of type 'T'"
		********************************************************************************************

		2. Major - the major issue is in the 4.png marked with green.
		If we read the green part form right to left.
		We focus on behavior. Remember anytime we talk about decoupling we should be focussing on the behavior.
		It's the different between saying "I care about the data based on what it is or I care about the
		data based on what it does." And when we talk about decoupling we are talking about
		"data based on what it does"
		The green chart is telling us that if we choose pointer semantics,
		for the implementation of the concrete data.
		Remember we choose the semantic immediately after we define the data.
		So if the semantic we choose is pointer semantics then
		only thing we are allowed to do moving forward is share that data.

		Remember how he said earlier the rule that we are not allowed to go from pointer semantics to
		value semantics.

		That green chart is saying the same thing. If we choose pointer semantics
		the only thing we are allowed to do is share the data.

		Once we are in the pointer semantic mode we a   re not allowed to make
		the copies of the value pointer is pointing to.

		The complete chart also says that when value semantics is chosen, to make copies.

		There are times/small cases - decoding and unmarshalling - that even though we are in value semantic
		mode it might be necessary to share.

		So if we are in pointer semantic mode then the only thing we are allowed to do is share,
		we are not allowed to make copies of data when we are in the pointer semantic mode.
		So we cannot be including the pointer semantic methods for those values because that would violate
		this massive integrity issue in our code. this is what the compiler is trying to prevent us from.

	*/
	// Value semantics are at play here. WE are trying ot make a copy of "T".
	// sendNotification(u)

	// So if we change this to use pointer semantics then the code works.
	// Because now we are in line with the right level of integrity, we are in line with the
	// data semantics.
	sendNotification(&u)

	// ./example1.go:36: cannot use u (type user) as type notifier in argument to sendNotification:
	//   user does not implement notifier (notify method has pointer receiver)
}

// sendNotification accepts values that implement the notifier
// interface and sends notifications.

/*
	Note that notifier values do not exist. An interface type is valueless.
	But what this means is that "I will accept any concrete piece of data, any value or any pointer
	that implements the full method sets of behavior defined by notifier which is notify() method."

	"n.notify" will change it's behavior depending on what concrete data we pass through our function.

*/
func sendNotification(n notifier) {
	n.notify()
}
