/*
	Now that we know about the custom error type we can start talking about, "Type as Context".
	Suddenly the error variables are no longer working for us.
	There's not enough context in them, now we need our own custom error type.

	Here we take some context from the standard library because there are some packages,
	that use customary type when the default error type is not giving us enough context.

	Here we see "UnmarshalTypeError" this coming form the Json package.
	The unmarshal type error is a concrete type that we use when we have marshalling situation. Where
	we set the field as a string and over JSon it may be coming over as an int and we have type mismatch.
	WE get an error of type "UnmarshalTypeError".

	----------------------------------------------------------------------------------------------------------
	This is the implementation of the error interface, we are using the pointer semantics, and remember
	pointer semantics are the default semantics for the custom error types.
	func (e *UnmarshalTypeError) Error() string {
	return "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()

	Because "UnmarshalTypeError" implements the error interface, "UnmarshalTypeError" can now be used as
	an error value.
	Note on line 18 that custom error type implementing method uses all of the fields in the type.
	If they are not using all of the fields to log something then Bill would question the field's necessity
	to provide context.

	This is a second concrete type from same package json.
	type InvalidUnmarshalError struct {
		Type reflect.Type
	}

	The implementation of the Error interface using "InvalidUnmarshalError" has a lot going on, the message can change
	depending on a few things.
	We will get an error of this type if we forget ot pass the address of something into Unmarshall.
	Remember unmarshalling requires pointer semantics.

	We have two different error types here,
	"UnmarshalTypeError" and "InvalidUnmarshalError" their variables start with "Err" & these types themselves end
	with "Error", this is an idiomatic way in golang to define the custom error type.

	Notice the function "Unmarshal", Bill has here taken a little bit of code out of the Json package with the
	Unmarshal function.
	Note that the Unmarshal function is asking for a slice of byte, that's how it is asking for that string.
	As string is a slice of bytes. Then it's saying I will accept any piece of data, any value any pointer
	regardless of the behavior, that's what an empty interface is.
	Empty interface is a dangerous type because it tells nothing about the concrete data that's being passed in.
	We have to be very careful when we use empty interface because it tels us nothing.
	It's not about what it does because we don't know what it does. It's very generic
	and if we are using empty interface to move dat around our program be very sure that, you are doing it
	because at compile time it's not about what it is and what it does that's going to be at runtime.
	This is where reflection comes in, we want to be to unmarshal into any piece of data and then use the reflection
	package ot be able to do that work at runtime.
	But notice here,

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	We are asking is the concrete value stored inside of "v". We are asking if that concrete data
	is a pointer or not. If it is nil or not.
	If it'snot a pointer type or if it's not nil
	we will return a "&InvalidUnmarshalError{reflect.TypeOf(v)}"
	an address of a concrete value of "InvalidUnmarshalError".

	If not that then we will return an error of the type,
	"UnmarshalTypeError".

	What's interesting in terms of error handling now is that we are interested in now is the type of value,
	that's stored inside the error interface that's being returned.
	In this case when "func Unmarshal(data []byte, v interface{}) error" returns an error variable the
	concrete value will either be an address of type "InvalidUnmarshalError" or an address of "UnmarshalTypeError".

	That's what we have to do our conditional logic on.
	Go has a special way of doing this type as context type of conditional logic.
	It's with the switch in the main function.
	"switch e := err.(type) {"
	We call Unmarshal and we get back that error interface value, we ask
	if there is a concrete value stored inside of "err" (if err!= nil)
	If there is, then we do this generic type assertion.
	"switch e := err.(type) {"
	Then we ask,

	case *UnmarshalTypeError:
			fmt.Printf("UnmarshalTypeError: Value[%s] Type[%v]\n", e.Value, e.Type)
	case *InvalidUnmarshalError:
			fmt.Printf("InvalidUnmarshalError: Type[%v]\n", e.Type)

	if we type assert, to be an address of "UnmarshalTypeError" then we execute logic.
	Or if "e" ends up being an address of "InvalidUnmarshalError", if it is then "e" will be that
	and we will then execute that logic.

	The whole idea here is that our conditional logic is not based on value itself but is based on what type of
	value, was stored inside that "err" interface. Type as context.

	There are lot of practical uses of "Type as context" not just for error handling but other things.

	There's potential danger of using "type as context" as well.
	The whole idea of error handling here was to, maintain the decoupled state using the error interface.

	The more decoupled we can maintain ourselves in terms of error handling the more,
	change we can apply to ur code without it being cascading (I didn't understand this).

	And here when we are talking about "type as context" what's happening is that we are moving from a decoupled state
	back to that coupled state, we are moving from the interface back to that concrete data.
	Which means now if the concrete data changes all of this error handling code also changes and breaks.

	Bill - Even though he loves "type as context" but he tends to want to avoid it when it comes to error handling. He
	would rather work off of the interface than to work off with the concrete values when it comes to error handling.









}
*/

// http://golang.org/src/pkg/encoding/json/decode.go
// Sample program to show how to implement a custom error type
// based on the json package in the standard library.
package main

import (
	"fmt"
	"reflect"
)

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // description of JSON value
	Type  reflect.Type // type of Go value it could not be assigned to
}

// Error implements the error interface.
func (e *UnmarshalTypeError) Error() string {
	return "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

// Error implements the error interface.
func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

// user is a type for use in the Unmarshal call.
type user struct {
	Name int
}

func main() {
	var u user
	err := Unmarshal([]byte(`{"name":"bill"}`), u) // Run with a value and pointer.
	if err != nil {
		switch e := err.(type) {
		case *UnmarshalTypeError:
			fmt.Printf("UnmarshalTypeError: Value[%s] Type[%v]\n", e.Value, e.Type)
		case *InvalidUnmarshalError:
			fmt.Printf("InvalidUnmarshalError: Type[%v]\n", e.Type)
		default:
			fmt.Println(err)
		}
		return
	}

	fmt.Println("Name:", u.Name)
}

// Unmarshal simulates an unmarshal call that always fails.
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return &UnmarshalTypeError{"string", reflect.TypeOf(v)}
}
