/*
	"Polymorphism means that you write a certain program and it behaves differently
	on the data that it operates on"
		- Tom Kurtz(inventor of BASIC)

*/

// Sample program to show how polymorphic behavior with interfaces.
package main

import "fmt"

// reader is an interface that defines the act of reading data.
/*
	Because reader is a type
	we could do "var r reader", defining a variable of type reader.
	This cannot be done because "interface" types are not concrete types.
	They don't define any sought of data, they only define a method set of behavior.
	That means interfaces are value lists. There is nothing concrete about them.
	"r" is not real, we cannot read it nor can we write it. We cannot even manipulate it.
	We have to keep this idea that the interfaces are value less.

	Polymorphism means a piece of code changes it's behavior depending upon the concrete data
	it is operating on.
*/

/*
	May be we could have written the method like this,
	read (n int)([]byte, error)
	and may be the method would have been simpler for the developer to use.
	We could have said "Well pass me the number of bytes you want me to read and
	I will return back a slice of data"

	Question. why is "read (n int)([]byte, error)" a horrific API design in Go?
	API Design has to be about making things easy to understand first and not simple to do.
	But, also API design has to take into account the impact the API is going to have on the caller's program.
	Answer. "read (n int)([]byte, error)" this API has a HUGE cost compared to
	"read(b []byte) (int, error)".
	Everytime "read (n int)([]byte, error)" gets called, there has to be an allocation that occurs on the heap
	there's no way around it.
	First allocation will come from the construction of the slice that we will have to make.
	b:= make([]byte,n) because the length is going to be set using a variable "n".
	If the compiler doesn't know the size of something at compile time then that must be
	allocated/constructed on the heap.
	Here "b:= make([]byte,n)" the compiler doesn't know the size of the backing array it is based on the variable
	"n". So this allocation has to occur on the heap.

	But what if we do,
	read() (int, error){
		b:= make([]byte,4096)
	}
	We can rid of the allocation based on the idea that the compiler doesn't know the size, because now it does.
	But at some point we will have to return the slice back up the call stack.

	read() (int, error){
		b:= make([]byte,4096)
		return b
	}

	When we do this then the pointer to the backing array says this can't be on the frame anymore.
	There's no way around having the allocation on the heap with this sought of API design.

	So remember one of the things we have to be careful about in our API design is not just about the
	data semantics. We also have to understand the impact the API is going to have on
	the user/caller of the program and design APIS that are sympathetic with every thing we are doing.
	In this case the garbage collector.

*/

/*
	Reader is an interface type which makes it a valueless type.
	So even if we declare a variable
	"var r reader"
	it doesn't make it concrete from our programming model, "r" doesn't really exist.
*/
type reader interface {
	read(b []byte) (int, error)
}

// file defines a system file.

/*
	But here "file" is a struct a concrete type and is not an interface type.
	This is real data, we can read, write and move around our program.
	This represents a file on our file system.
*/
type file struct {
	name string
}

// read implements the reader interface for a file.
/*
	read method matches the exact name and signature as the read method defined inside the reader interface.
	The concrete type "file" implements the reader interface using value semantics.
*/
func (file) read(b []byte) (int, error) {
	s := "<rss><channel><title>Going Go Programming</title></channel></rss>"
	copy(b, s)
	return len(s), nil
}

// pipe defines a named pipe network connection.
type pipe struct {
	name string
}

// read implements the reader interface for a network connection.
/*
	Pipe represents stream. pipe is a concrete type.
	It's a piece of data we can read and manipulate around in our program.
	Pipe also implements the reader interface also using value semantics.

	We now have 2 concrete data declarations that implements the same read behavior defined by the
	reader interface.
*/
func (pipe) read(b []byte) (int, error) {
	s := `{name: "bill", title: "developer"}`
	copy(b, s)
	return len(s), nil
}

func main() {

	// Create two values one of type file and one of type pipe.
	f := file{"data.json"}
	p := pipe{"cfg_service"}

	// Call the retrieve function for each concrete type.
	// retrieve will receive it's own copy of the file value we will say the value semantics are at play.
	/*
		Where does the decoupling come in?
		Even though in our program "r" is value less it still has a compiler implementation.
		See 1.png.
		Interface value is a two word data structure both with pointers.
		Second pointer of the interface value that the compiler is defining is one for storage.
		Also remember when we have decoupling it means that copy of "f" has to allocate in the heap memory.

		The first word of the interface value, points to a special table called the "iTable"
		"interface table"
		"iTable" let's use define the type of data that's being stored, a file value is being stored.
		Then the itable also gives us a pointer to the implementation of "read (for the file)"
		3.png

		When we call r.read() then we do an "itable" lookup and identify that we have got a file value that is
		stored in "r" and then we call the read implementation against our copy.

		This is where now polymorphism comes in.
		Now, when we call "retrieve(p)" we now not passing a copy of "f" but passing in a copy of "p".
		We are still using our value semantics here. "itable" now says pipe.
		4.png.
		Now when we call "r.read()" we call the pipe's implementation of read.

		Polymorphism means that a piece of code changes it's behavior depending on the concrete data
		file values or pipe values that we pass through.

		Moving forward we don't really care about the iTable.
		So moving forward even though we know the first word of the interface is pointer,
		we will say that first word is going to always describe to us what type of data is being stored
		and gives us the ability to access th concrete implementation of that data and the second word
		stores the pointer which will give us the storage which makes the interface value "r" not longer
		valueless.
		But remember we are never passing around "r" in our program. We only pass the things we can pass
		the concrete data.

	*/
	retrieve(f)
	retrieve(p)
}

// retrieve can read any device and process the data.
/*
	Retrieve will accept a value of type reader.
	But as we said values of type reader do not exist.
	Interfaces are valueless. This function is not asking for a reader value.
	It really means - I will accept any piece of concrete data or any value or any pointer that
	implemnets the reader behavior. That exhibits the full method set behavior defined by reader interface.
	This polymorphic function is creating the level of precise decoupling because it is not asking for the
	data based on what it is, but it is asking for data based on what it can do.

	This is one case where we choose a method over a function because, we have a polymorphic function
	that says I don't want data based on what it is but on what it does.
	Polymorphism means that a piece of code changes it's behavior depending on the concrete data it is
	operating on.
*/
func retrieve(r reader) error {
	data := make([]byte, 100)

	/*
		Here we are calling the read behavior based on our value less type.
		If we pass a file value into this function, then we would be implementing the file's implementation of
		read. If we pass in pipe value through then we would be executing pipe's executing.
		This is our polymorphism this line "len, err := r.read(data)"
		changes it's behavior based on the concrete data that we pass through this function.

		Next, we go to the main fucntion.
	*/
	len, err := r.read(data)
	if err != nil {
		return err
	}

	fmt.Println(string(data[:len]))
	return nil
}
