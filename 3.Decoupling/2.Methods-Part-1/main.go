/*
	1. For Built-in types we will use value semantics to move that data around in our program.
	Numeric, String or Bool. Only time we'll be okay for the pointer of type int string or bool
	will be when we have them as a field in a struct where we have to represent the idea of null.
	For the built in data types all the moving around and also the reading/writing is done by
	value semantics as well.

	2. For reference types we should be using value semantics to move the data around in our program.
	Everyone gets their own copy of the interface, maps, slice, function. They all get their own copy of it.
	*Note* when we are reading and writing with the reference type we are already in the pointer semantic mode.

	3. For the struct types we have to decide. At the time we define the types we have to decide,
	am I using the pointer semantics or value semantics and mantain strict consistency with this.

*/

import (
	"sync/atomic"
	"syscall"
)

// Sample code to show how it is important to use value or pointer semantics
// in a consistent way. Choose the semantic that is reasonable and practical
// for the given type and be consistent. One exception is an unmarshal
// operation since that always requires the address of a value.

// *****************************************************************************

// These is a named type from the net package called IP and IPMask with a base
// type that is a slice of bytes. Since we use value semantics for reference
// types, the implementation is using value semantics for both.

// These 2 types are defined in the net package in the standard library.
// They are custom type but are based on an existing data structure which in this case is slice.
// These are now reference types.
// This means we are using the value semantics to move them around our program. Whereas we use 
// the pointer semantics to read and write.
type IP []byte
type IPMask []byte

// Mask is using a value receiver and returning a value of type IP. This
// method is using value semantics for type IP.

// Note here "Mask" is a value sematic, receiver is value based and so is the return type.
// We should never be confused about how to design the API once we know
// what the input is and what the output is. The data is itself defining how 
// it comes in and out of these functions.
// func "(ip IP)" Mask(mask IPMask) IP 
// Mask is attached with a value semantic "(ip IP)".
// If we were using the pointer semantics "(ip *IP)" for Mask, we would be using the wrong semantics.
// Because slices don't use pointer semantics they use value semantics. 
// That is why we see a value receiver and a value return. Hence it is a value semantic API.
func (ip IP) Mask(mask IPMask) IP {
	if len(mask) == IPv6len && len(ip) == IPv4len && allFF(mask[:12]) {
		mask = mask[12:]
	}
	if len(mask) == IPv4len && len(ip) == IPv6len && bytesEqual(ip[:12], v4InV6Prefix) {
		ip = ip[12:]
	}
	n := len(ip)
	if n != len(mask) {
		return nil
	}
	out := make(IP, n)
	for i := 0; i < n; i++ {
		out[i] = ip[i] & mask[i]
	}
	return out
}

// ipEmptyString accepts a value of type IP and returns a value of type string.
// The function is using value semantics for type IP.

// This function also takes in a "value" of IP and returns a "value" of type string.
// Why? Because these are the data that use value semantics.
func ipEmptyString(ip IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}

// *****************************************************************************
/* 
	There are exceptions, one of the exceptions for both slices and map, 
	there may be a time when we need to share the address of a slice or a map, is when we are sharing
	a slice or a map down a call stack and into the functions that are either called "decode" or "unmarshalling".
	"Decoding" and "Unmarshalling" require pointer semantics.

*/
// *****************************************************************************

// Should time use value or pointer semantics? If you need to modify a time
// value should you mutate the value or create a new one?

/*
	Sometimes its obvious that something should be using a value semantic or a pointer semantic.
	"We are not making decision based on performance. 
	We are not making decision based on idea that something could or could not be copied.
	We are doing it based on correctness. What is the reasonable thing to do for this type of data.
	What is the expectation. What are people going to assume when they look at the data like this"

	Here is a question we can ask ourselves as developers.
	Say, we have a value of type time, and say it represents 5opm and if we add 5 minutes to it.
	Then is it the same time value just 5 minutes later. Or does "5:05pm" represents a different piece of data.
	this is what we can ask ourselves, that if a mutation occurs to the data, 
	if it the same piece of data modified, or is it a brand new piece of data. 
	When it comes to time, we can see that "5:05pm" and "5:00pm" are two distinct pieces of data.
	With that in mind value semantics tend to make more sense.
	
	We can think the same with the string. If we have a string and if we modify or mutate the string
	we get a new string.

	Thinking about a User, say the current users name is "Jane" and we decide to change "Jane's" name to
	"Sally".
	Does that change who "Jane" is. It doesn't change who "Jane" is. There is just a modification in the name.
	When we have data that doesn't really change with modificaion then usually pointer semantics is the
	correct way to go.
	When not a 100% sure what semantic to use then we go for pointer semantics.
	Because we cannot assume that all data can be copied.
	Example if we have a mutex, waitgroup as fields then that data can no longer be copied.
	When not sure safest way is pointer semantics to share within our code.
	But we have to be extra careful when reading an writing.

*/

/*
	Best way to know what semantic is at play is to see the "Factory" functions.
	Type, then factory functions and then methods would be in this order when reviewing code.
	Don't scatter the methods thorughout the code base across multiple files.
	Keep these things together. Also we don't need a separate file for every type.

	Now is a factory function.
	Looking at the return type tells us lot about the data semantic that the developer chose.
	Now() returns a value of type time.
	This tells us that value semantics are at play. Any function that accepts a time value should
	be done "by value". Return to time should also be done "by value"
	A field in the struct should be done "by value" as well.
	An exception would be when we need to represent null. 
	When we need to unmarshall or decode would be an exception.

	Value semantics for time make sense because if we are moving time values around by themselves then
	we can guarantee that we wouldn't be polluting our heap.

	Next, look at the Add method.
*/

type Time struct {
	sec  int64
	nsec int32
	loc  *Location
}

// Factory functions dictate the semantics that will be used. The Now function
// returns a value of type Time. This means we should be using value
// semantics and copy Time values.

func Now() Time {
	sec, nsec := now()
	return Time{sec + unixToInternal, nsec, Local}
}

// Add is using a value receiver and returning a value of type Time. This
// method is using value semantics for Time.

// Add is a mutation operation. Add uses value semantic mutation API. Why?
// Because time has to use value semantics. API must respect the DATA not other other way around. 
// We have to know what the data is and what the data is defining in terms of semantic and then we can design the
// API to provide that functionality it needs to provide.
// Hence we have a value receiver and the return is also a new time "Value".
func (t Time) Add(d Duration) Time {
	t.sec += int64(d / 1e9)
	nsec := int32(t.nsec) + int32(d%1e9)
	if nsec >= 1e9 {
		t.sec++
		nsec -= 1e9
	} else if nsec < 0 {
		t.sec--
		nsec += 1e9
	}
	t.nsec = nsec
	return t
}

// div accepts a value of type Time and returns values of built-in types.
// The function is using value semantics for type Time.
// This is a function.
func div(t Time, d Duration) (qmod2 int, r Duration) {
	// Code here
}


/*
	There are times when we would need to go from value semantics to pointer semantics.
	There can be exceptions. Normally when we go from value to pointer that is change semantics, 
	it's going to be done in a small scope and normally around the ide of decoding nad unmarshaling
	These below methods part of the standard library.
	Look the API has switched and defined these methods with a pointer semantics.

*/
// The only use pointer semantics for the `Time` api are these
// unmarshal related functions.

func (t *Time) UnmarshalBinary(data []byte) error {
func (t *Time) GobDecode(data []byte) error {
func (t *Time) UnmarshalJSON(data []byte) error {
func (t *Time) UnmarshalText(data []byte) error {

// *****************************************************************************
/*
	It is okay to go from value semantic to pointer semantics in a very small scope. 
	But it is never okay to go from the pointer semantics to value semantics. 
	We are never allowed to make a copy of the value our pointer is pointing to. 

	Value to Pointer allowed in small situations.
	Pointer to Value never allowed. As it violates the integrity.
*/
// *****************************************************************************

// Factory functions dictate the semantics that will be used. The Open function
// returns a pointer of type File. This means we should be using pointer
// semantics and share File values.

/* 
	Open takes a value of type string. String uses a value semantic but what is this function returning?
	A pointer of type file. Here is a factory function using pointer semantic. 
	Here OS package is saying that regardless of the operation we share the file value.
	Look at chdir. change Dir is using pointer semantics on the receiver why?
	Because we must share the file values. Even if the method doesn't mutate anything, it must be shared.
*/

func Open(name string) (file *File, err error) {
	return OpenFile(name, O_RDONLY, 0)
}

// Chdir is using a pointer receiver. This method is using pointer
// semantics for File.

func (f *File) Chdir() error {
	if f == nil {
		return ErrInvalid
	}
	if e := syscall.Fchdir(f.fd); e != nil {
		return &PathError{"chdir", f.name, e}
	}
	return nil
}

// epipecheck accepts a pointer of type File.
// The function is using pointer semantics for type File.
// Err or is an interface it is a reference type hence we are sharing the value.
// Value semantic.

func epipecheck(file *File, e error) {
	if e == syscall.EPIPE {
		if atomic.AddInt32(&file.nepipe, 1) >= 10 {
			sigpipe()
		}
	} else {
		atomic.StoreInt32(&file.nepipe, 0)
	}
}