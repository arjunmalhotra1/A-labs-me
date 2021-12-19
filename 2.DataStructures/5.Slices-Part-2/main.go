// Sample program to show how to grow a slice using the built-in function append
// and how append grows the capacity of the underlying array.
package main

import "fmt"

/*
	This doesn't represent the zero value for all types but var does.
	u:= User{}

	So if we do
	u := []string{}
	This doesn't represent a zero value slice. It represents an empty slice.
	See pic 1.png. On left is a zero value slice and right is an empty slice. Length and capacity are 0
	but there is a pointer.
	Question arises then that why do we need the semantic for empty?
	Remember that slices represent collection. There are times when collection is nil (doesn't exist) but there
	are times when the collection might be empty.
	Say, we do a query and the query says give me all the users older than 100 years old. And the query comes
	back with no users. That query didn't fail. It's not that the collection is nil or doesn't exist.
	It's an empty collection. Hence we need the semantic to represent empty.
	Hence we always want to use var for zero value and non-var for everything else.

	Question. In empty slice, there is a pointer. Where is the pointer pointing to?
	Answer. There is another type in Golang. That type is called the "empty struct".
	That type is called the empty struct.
	Here, name of the type uses the empty literal notation to represent the fact that it is an empty struct.
	var es struct{} This is an empty struct and it doesn't have anything in it. This is a zero allocation type.
	It is embedded in runtime and it has a location. If we were to construct a hundred thousand empty structs
	they would all have the same address. Checkout 2.png.


	If we use empty literal construction to construct an empty struct we use:
	es := struct{}{}

*/

func main() {

	// Declare a nil slice of strings.
	// Whenever a reference type is set to a zero value they are considered to be "nil". Just like
	// a pointer set to it's zero value.
	// We use value semantics to move the reference types in our program but
	// while reading and writing out of them we use pointer semantics.
	// Pointer means that in their zero value state they are "nil".
	// Here we are starting with nil slice. We are starting with the idea that there is no backing array.
	// This code is appending 100,000 strings into a slice.
	// Remember slice is a dynamic array.
	var data []string

	//data := make([]string,0,1e5)

	// Capture the capacity of the slice.
	lastCap := cap(data)

	// Append ~100k strings to the slice.
	for record := 1; record <= 1e5; record++ {

		// Use the built-in function append to add to the slice.
		value := fmt.Sprintf("Rec: %d", record)
		data = append(data, value)
		// Notice append takes in it's own copy of the slice value.
		// Remember it needs to be designed in this way, because we move slices around our program by
		// creating copies of the slice. Everybody gets their own copy. We do not want to pollute the
		// heap with these values we are leveraging the stack.
		// append get's it's own cpy of the slice. See 3.png.
		// When append gets it's own copy of the slice value it really just want to check one thing.
		// It just wants to know if the length and the capacity the same.
		// If they are the same then it will create a new backing array with an element.
		// See 4.png
		// And then update the slice, with new length and new capacity.
		// Append then returns the new version of the slice value back. 5.png
		// On return the copy append created disappears. 6.png.
		// this is called value semantic mutation api.
		// When the length of the current backing array is less than 1000,
		// the size of the backing array will double.
		// When it's over a 1000 elements the sie will increase by 25%.
		// Data from the previous backing array is copied to the new backing array that was created.
		// See 7.png and 8.png.
		// next time when the garbage collector runs it will notice that there is nothing referencing the
		// earlier array and it will get swept away by the garbage collector.

		// What is memory leak?
		// Memory leak in any language like C is a situation like when the user decides where to construct
		// the value in the heap or the stack. If we decide on the heap then we use functions like
		// "malloc" or "new". Memory leak is when we forget to call the "delete" or "free" on that memory
		// location after the call to malloc or new.

		// What is a memory leak in Golang?
		// Memory leak in golang is a situation where there is value on the heap and there is a reference
		// to it and it really shouldn't be. We have a value on the heap there is a reference to it therefore
		// it's not being released but reality is that that reference shouldn't exist.
		// What we can do when we have a memory leak in go is that we ensure that we have a memory leak by using
		// "gc trace" a program, shown later. Let's say that "gc trace" shows a leak then we
		// look at 4 things when there is a leak in our program.
		// 1. First thing we look for are the go routines that are either leaking or blocking when they shouldn't.
		// This is the number 1 reason for memory leak in go.
		// We are creating extra Go routines for whatever reason and we are expecting those go routines to terminate
		// and they don't. May eb they block on a channel or some synchronization, they never die. So then not only
		// it's the go routine leaking but anything that Go routine is still referencing.= is leaking too.
		// 2. Maps a Lot of people use maps as a cache. Reality is that we have to delete keys if we are using
		// a map as a cache. We have several option to delete the keys:
		// 		2a. Size based cache. We set a max size/ceiling on how much memory the cache would use.
		// 		2b. Time based cache. May be we have leaks there because we said that we will keep the data in the
		// cache for 5 mins and 5 mins have passed but the data is still in the cache.
		// Lot of times we have memory leaks with maps when a map is based on an event.
		// Like a system where a socket connection comes in, we mantain a cache for that socket connection.
		// When a socket drops and we are clearing that key.
		// So we will look at the maps and see if we are missing a point of deleteing the keys.
		// 3. Append call. Append calls can be dangerous.
		//  data = append(data, value)
		// When the "data" going in the
		// append call is not replaced on the way out. That is bad. Because we could be holding the extra
		// references to the old backing array that we wouldn't want to.
		// 4. Finally we look at the API calls that probably have a close function that we forgot to call.
		// Say like http get, put post calls and forgetting to call the "close()"
		// function could result in memory leaks but also result in running out of file handles.

		// Moving on with the append, with A,B the capacity is still full so we double the size again.
		// We add "C" and now the length is 3 and capacity is 4. Note 9.png
		// We have extra length of capacity for growith. The garbage collector would come in and claim
		// the "A,B" Now ehn we add "D" we see the length is not the same as capacity.
		// So we simply add "D", 10.png.
		//************************************************************************//

		// Since we already know we are going to append 100,000 times. Hence it is really silly of
		// us to be starting with a nill slice in this case. "var data []string" on line 48.
		// So if we replace that with
		// data := make([]string,0,1e5)
		// We get no tracing output with this code. This is because we don't have to do any extra allocations.
		// AS we set upfront a slice with length "0" and capacity 1e5. And it points to a very large array of
		// those 1e5 empty elements. 12.png.
		// The reason we are setting length to 0 and capacity to 1e5 is because append works out of length.
		// If we just set the length and we are using length. Crazy things happen. The very first append will not
		// now happen at index 0 but at index 100,000.
		//data := make([]string,1e5)
		// Now when we have "data := make([]string,1e5)" since we know exactly what I need.
		// We can change our for loop to be "0" index and get rid of the entire append call.
		/*
			This is even more efficient.
			for record := 0; record < 1e5; record++ {
				// Use the built-in function append to add to the slice.
				value := fmt.Sprintf("Rec: %d", record)
				data[record] = value
		*/
		// Take away
		// Append has been designed around value semantics because we use our value semantics to move that
		// slice around our program and append is a great example of a value semantic mutation API.
		// A function cannot dictate how to it works with data the data tells the API how it's supposed to work.

		// When the capacity of the slice changes, display the changes.
		if lastCap != cap(data) {

			// Calculate the percent of change.
			capChg := float64(cap(data)-lastCap) / float64(lastCap) * 100

			// Save the new values for capacity.
			lastCap = cap(data)

			// Display the results.
			fmt.Printf("Addr[%p]\tIndex[%d]\t\tCap[%d - %2.f%%]\n",
				&data[0],
				record,
				cap(data),
				capChg)
		}
	}
}
