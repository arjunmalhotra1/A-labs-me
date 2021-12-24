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
	// person user // Not Embedding.

	user // Embedding Type. We pull the person field out of the admin type and just leave the type "user".
	/*
		This is embedding, there are two ways we can look at embedding. We can now say,
		that any value of type admin also has a value of type user embedded inside of it.
		Value of type admin contains value of type user.

		A better way to look at this is to look at this idea of an inner type and an outer type.
		Outer type admin has embedded inside of it the inner type "user".
		When we have this inner and outer type relationship, we have what is called inner type promotion.

		Inner type promotion says that anything related to the inner type promotes to the outer type.
		All the fields and methods associated with the inner type now become part of the outer type.

		This does not make admin a user. We cannot pass around admin values like they are user values.

		Note: This is not Inheritance. There is no subtyping. Admins are Admins and Users are Users.
		All we are doing here is, allowing fields and methods associated with user type to promote up to the
		outer type.
	*/
	level string
}

func main() {

	// Create an admin user, using the struct literal form.
	/*
		Note the difference in construction. TO access the inner value we just access the
		name of the type.
	*/
	ad := admin{
		//Person: user{
		// "user" looks like a field and acts like a field but let's not look at it like a field.
		// It's the name of the value that's embedded we are using the type's name.
		user: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
	}

	// We can still access the inner type's method directly
	// ad.person.notify() //earlier
	ad.user.notify()

	// We can access fields methods. WE can access the inner type's behavior directly through the
	// outer type without having to do ad.user.notify()
	ad.notify()
}
