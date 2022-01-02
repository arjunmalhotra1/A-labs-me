/*
	Now we implement Stringer method using value semantics. We know that then
	both values and pointers satisfy the interface.
	The output is that we get custom formatter 7.png for both value semantics call and the pointer semantic call.

	But the reality is that we don't make these type of changes, the code
	"func (u user) String() string {" or this "func (u *user) String() string {"
	was dictated long before we even got the line where we call the method.

	This is just to show that we can write APIs that once we are in a decoupled state, that could really ask
	"Does this data I am operating on, does it implement something (assert) or some other behaviors"
	Or data we are passing in, implement some other behavior.

	In the io package of Golang, the copy function, here is the code for the copy function
	- https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/io/io.go;l=381

	We can see the the exported function "Copy" function calls the un exported function, "copyBuffer"
	Ans in the "copyBuffer" function
	Notice these 2 type assertions:

	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	It's very similar to what the fmt Print function is doing asking if the data implements, something else.
	BEcause copyBuffer has a default implementation for copying 32,000 bytes at a time.

	But notice that the "copyBuffer" above on line 21 asks if the concrete data passed being passed for src, if
	src implements the "WriterTo" interface -"src.(WriterTo)" - forget about my default implementation
	and use the custom implementation.

	And also on line 25, notice it says if the concrete data passed in for "dst" if it implements the "ReaderFrom"
	dst.(ReaderFrom) then also bypass the default and use the custom.
	Don't stress too much if you don't understand this example thoroughly.

	This gives the developer the ability to say, "I have got a working default implementation but there are times
	depending on the architecture and platform we might know, that we want to do something
	different or a specific implementation/customize"

	This idea of type assertions specially as part of API it's a really nice way of giving ability to overwrite the
	default behavior.


*/
// Sample program to show how method sets can affect behavior.
package main

import "fmt"

// user defines a user in the system.
type user struct {
	name  string
	email string
}

// String implements the fmt.Stringer interface.
func (u user) String() string {
	return fmt.Sprintf("My name is %q and my email is %q", u.name, u.email)
}

func main() {

	// Create a value of type user.
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	// Display the values.
	fmt.Println(u)
	fmt.Println(&u)
}
