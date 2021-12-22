/*
	When we know the data semantics, we know the behavior. When we know the behavior
	we know the cost. When we know the cost we are engineering.
*/

// Sample program to show how to declare function variables.
package main

import "fmt"

// data is a struct to bind methods to.
type data struct {
	name string
	age  int
}

/*
	"displayName()" is using value semantics and "setAge" is using pointer semantic or pointer receiver.
	Note: We do not write getters and setters in Golang.
	Code review stops whenever we see methods like "setAge" or "getAge".
	These are not the APIs that provide something. We don't want APIs that just allow the user the ability to
	just manipulate state. We don't want o do these things in Golang.

	Technically all we have to a program is a state and the behavior thorugh functions.
	We have our data and we have our functions.
*/
// displayName provides a pretty print view of the name.
func (d data) displayName() {
	fmt.Println("My Name Is", d.name)
}

// setAge sets the age and displays the value.
func (d *data) setAge(age int) {
	d.age = age
	fmt.Println(d.name, "Is Age", d.age)
}

func main() {

	// Declare a variable of type data.
	d := data{
		name: "Bill",
	}

	fmt.Println("Proper Calls to Methods:")

	// How we actually call methods in Go.
	d.displayName()
	d.setAge(45)

	fmt.Println("\nWhat the Compiler is Doing:")

	// This is what Go **compiler** is doing underneath.
	// The method call "d.displayName()" is changed to actually a function call "data.displayName(d)"
	// And the "d" is the parameter.
	// We are not going to call any methods like this.
	// This is just to show that emthod are just syntacttic sugar behind a function calls.
	// In stack ctrace that is what we see.
	data.displayN ame(d)
	(*data).setAge(&d, 45)

	// =========================================================================
	/*
		Whole idea is to understand cost and we can't do that unless we know the behavior.
		Data semantics are the big driver to understand the behavior.
	*/
	// =========================================================================
	fmt.Println("\nCall Value Receiver Methods with Variable:")

	// Declare a function variable for the method bound to the d variable.
	// The function variable will get its own copy of d because the method
	// is using a value receiver.
	// We are assigning the method value to f1, we are not making a method call.
	/*
		f1 is function variable it is  from the family of reference types.
		f1 is just a pointer like channels and maps. 1.png
		What is the pointer pointing to, it is pointing to a 2 word data structure.
		First word is a pointer to the code for the function, in this case "displayName".
		See 2.png. so that first pointer is pointing to the implementation if displyName.
		Because "d" has 2 methods to it, DisplayName and SetAge.

		At the end of the day we would want to call the displayName like "f1()"
		That's the whole point of decoupling. Decoupling comes from indirection.
		f1 is indirectly going to execute behavior of "d".
		Question, now what about the data?
		We cannot call displayName() without any data

		1. displayName is a mehtod. 
		2. displayname is using a value receiver which means that it is using value semantics.
		3. The rules for value semantics are that every piece of code operates on it's own copy the data.
		What does all this mean for f1? Since display name is using value semantics 
		that is everyone gets their own copy, this extends to f1 as well.
		The second word in f1 will point to "d" it's own copy. See 3.png

		Now when we call f1() the displayName() method but it will be called against our copy of 
		d, which means if we change "Bill" to "Joan" and we call f1() again. We will not be seeing the change.
		Because anytime we call displayName through f1() we will be operating on the own copy of f1.
		See 4.png

		this indirection is the3 cost but the other cost is the allocation of "copy of d".
		This value now has to allocate because we made a copy at the runtime.
		Remember one of the rules for values being on the stack is that if the compiler does not know the size
		of a value at compile time. It cannot be on the stack. Because frames are fixed at compile time.
		The size of the frame are designed at compile time. So if the compiler doesn't know the 
		size of something at the compile time it cannot be on the stack.
		Compiler does not know the value of the copy of "d" for "f1" at compile time because 
		we are doing this copy at run time. that means it has to be on the heap.

		Note:
		When value semantics are at play and decoupling is being applied. That copy is going to be made and 
		that is going to cost us indirect and allocation of memory as well.


	*/
	f1 := d.displayName

	// Call the method via the variable.
	f1()

	// Change the value of d.
	d.name = "Joan"

	// Call the method via the variable. We don't see the change.
	f1()

	// =========================================================================

	fmt.Println("\nCall Pointer Receiver Method with Variable:")

	// Declare a function variable for the method bound to the d variable.
	// The function variable will get the address of d because the method
	// is using a pointer receiver.

	/*
		Here we are going to assign "setAge" to the variable f2.
		5.png.
		f2 has a pointer to the code for set age.
		What happens at the second work?
		setAge is using pointer semantics. Pointer semantics mean shared access. Hence we are not going to
		make a copy. The second word will point to the original value.
		Now when we call "f2()" we will see "Joan".
		Next we change that to "Sammy".  
		We will see "Sammy" when we call "f2(45)" because we are sharing the original data.
		Question. Since all of this code is in the one frame of the "main" function call.
		6.Png. Does that mean that now since we are sharing "d" and all the code is isolated int he frmae,
		does that mean that "d" has to allocate?
		Answer. No, because it is not being shared in such a way that it has to come off the frame.
		But, there are flaws in escape analysis and pne opf the flaws is that when there is 
		double indirection to data. In this case it is. See pic 7.png
		Compiler says that I cannot track this anymore through static code analysis and this now 
		has to end up on the heap.
		The cost of decoupling - is always allocation and indirection regardless of the data semantic we
		are working with. So if we need the blazing fast go code we cannot be decoupling it
		since even one allocation on the heap is going to have the garbage collector to come and do some work.
		We already saw the latencies behind that.
		Reality is that we don't need the most blazing fast software ever.
		We need things just to be fast enough. We want that allocation on heap to happen if it makes our lives 
		better specially  around multi threaded software development.

		If we are using interfaces for the sake of using interfaces and now we got that level of decoupling 
		and now we are allocating and there is no real value behind it, that is a non productive allocation.
		This is decoupling thorugh function based reference type.
	*/
	f2 := d.setAge

	// Call the method via the variable.
	f2(45)

	// Change the value of d.
	d.name = "Sammy"

	// Call the method via the variable. We see the change.
	f2(45)
}
