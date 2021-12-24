/*
	the idea of exporting and unexporting also happens at the type level. WE can see that this type is
	exported "User" and the 2 fields are also exported. but "password" is not exported.
	That means only the code inside the users package can access this code directly.
*/

// Package users provides support for user management.
package users

// User represents information about a user.
type User struct {
	Name string
	ID   int

	password string
}
