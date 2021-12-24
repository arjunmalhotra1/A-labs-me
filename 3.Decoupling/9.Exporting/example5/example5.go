// Sample program to show how to create values from exported types with
// embedded unexported types.
package main

import (
	"fmt"

	"github.com/ardanlabs/gotraining/topics/go/language/exporting/example5/users"
)

func main() {

	// Create a value of type Manager from the users package.
	// We are not accessing "Name" and "ID" during construction. Since "user" is unexported.
	u := users.Manager{
		Title: "Dev Manager",
	}

	// Set the exported fields from the unexported user inner type.
	// but once we are done with construction we do have access to it.
	u.Name = "Chole"
	u.ID = 10

	fmt.Printf("User: %#v\n", u)
}
