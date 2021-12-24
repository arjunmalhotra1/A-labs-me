// Sample program to show how unexported fields from an exported struct
// type can't be accessed directly.
package main

import (
	"fmt"

	"github.com/ardanlabs/gotraining/topics/go/language/exporting/example4/users"
)

func main() {

	// Create a value of type User from the users package.
	u := users.User{
		Name: "Chole",
		ID:   10,

		password: "xxxx", // compiler will throw an error as this field "password" is not exported.
	}

	// ./example4.go:21: unknown users.User field 'password' in struct literal

	fmt.Printf("User: %#v\n", u)
}
