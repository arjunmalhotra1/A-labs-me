// Package users provides support for user management.
package users

/*
	user is unexported type but two fields Name and ID are exported.
	Manager is exported type with "Title" also exported and "user" embedded within Manager.

	Outside of users package we can create a value of type Manager as it is exported and we can
	even set the field title.

	Question. What is going to happen with "user".
	We can't access "user" directly that is a lower case "u".
	But because of inner type promotion we can access the "Name" and "ID" not during construction but after.
	See example5.go
*/

// ********** NOTE THIS IS AVERY VERY VERY VERY BAD CODE.****************
// Code review will stop because these ideas of partial construction always result in bugs.
// So he would ask to make the "user" type exported.

// However this is common having exported type fields within the unexported types.
// This is becasue of marshalling and unmarshalling. Our marshallers & decoders, they require fields
// to be exported in order to access them because that's the way reflection to work.
// We wouldn't want the entire type to be exported because it's still data that should be
// unexported or internal to the package that is using it.

// user represents information about a user.
type user struct {
	Name string
	ID   int
}

// Manager represents information about a manager.
type Manager struct {
	Title string

	user
}
