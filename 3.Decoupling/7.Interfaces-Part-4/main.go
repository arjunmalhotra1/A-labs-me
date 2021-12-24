// Sample program to show the syntax and mechanics of type
// switches and the empty interface.

// empty interface and type assertions.
/*
	Print function uses reflection to determine, what type of data is being passed in.
	"string", "float", "boolean" etc.

	What if we want to implement our own version of Println, say "myPrintln"
	WE do this using an empty interface. See
	myPrintln() function.

*/
package main

import "fmt"

func main() {

	// fmt.Println can be called with values of any type.
	fmt.Println("Hello, world")
	fmt.Println(12345)
	fmt.Println(3.14159)
	fmt.Println(true)

	// How can we do the same?
	myPrintln("Hello, world")
	myPrintln(12345)
	myPrintln(3.14159)
	myPrintln(true)

	// - An interface is satisfied by any piece of data when the data exhibits
	// the full method set of behavior defined by the interface.
	// - The empty interface defines no method set of behavior and therefore
	// requires no method by the data being stored.

	// - The empty interface says nothing about the data stored inside
	// the interface.
	// - Checks would need to be performed at runtime to know anything about
	// the data stored in the empty interface.
	// - Decouple around well defined behavior and only use the empty
	// interface as an exception when it is reasonable and practical to do so.
}

/*
Problem with empty interface is that it describes nothing.
It's not suggesting that the data has to be behaving in any sought of way.
It's a very generic form to bring the data in.
And whenever we are using empty interface we will have to leverage reflection package or something like
type assertion to figure out how to act on the concrete data.

Here we are passing the concrete data through the polymorphic function but not based on any data.

Now we have to inspect the data at run time, this is where type assertion comes in.

One way to do type assertion is to be explicit and say,
does the concrete data inside of "a" is it a string?
if "a"  is a strign then "v" will be copy of the data that we have stored in "a".
	v:=a.(string)
	But if "a" is not a string, this program will panic.

So times when we want to test something and are not 100% sure what is in there we can use the second form
with a boolean flag

v, ok = a.(string)
if "ok" is true there was a strign adn we get the value in "v". If "ok" was false, there wasn't a string in "a"
and we don't get a string back.
"v" will be of empty value and "ok" will be false.
Instead of doing conditional logic on the "ok" variable. We use the special mechanism on the switch.

Switch does the type assertion in a generic way. Let's type assert on the keyword "type".
And inside of the switch case we can ask,
if we perform this type assertion then is the data an int, string, float64 if any of these
then let "v" be a copy of data of these type and let's execute this code.
And iin the default case we are saying that it's unknown.
Here all this is done using the type assertions. We can also leverage the reflection package to do this.

Empty interface is a way of saying "I will accept all concrete data not by what it is, and also not
by what it does."
Now we know nothing about the data and we have to leverage the type assertions to get the copies of the data
back of that interface whether it's a value or pointer.
We can also levrage reflections as well.
*/

func myPrintln(a interface{}) {
	switch v := a.(type) {
	case string:
		fmt.Printf("Is string  : type(%T) : value(%s)\n", v, v)
	case int:
		fmt.Printf("Is int     : type(%T) : value(%d)\n", v, v)
	case float64:
		fmt.Printf("Is float64 : type(%T) : value(%f)\n", v, v)
	default:
		fmt.Printf("Is unknown : type(%T) : value(%v)\n", v, v)
	}
}
