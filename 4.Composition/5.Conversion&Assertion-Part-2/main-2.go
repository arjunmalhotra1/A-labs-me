/*
	Here we create a user type, and again we are implementing the Stringer interface.
	We are using the pointer semantics, which means we can only share the user values with the interface.
	Now in main we construct a user value and call print.
	The idea that the Stringer interface exists in the fmt package and a nice thing about the print function
	is it can check and do the type assertion and ask, does this value that we are asked to be print, does it
	implement the Stringer interface?
	We know that the Stringer interface is implemented for
	"fmt.Println(&u)" because we used the pointer semantics "func (u *user) String() string {"
	but does is it implemented for the value call,
	"fmt.Println(u)"?
	Answer when we run this program, we see that the print function is still using the default formatter for
	the user value "fmt.Println(u)". But for the pointer "fmt.Println(&u)" the custom formatter is used.
	See 1.png
	But if we change the String() function to be value semantics, see main-3.go
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
func (u *user) String() string {
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
