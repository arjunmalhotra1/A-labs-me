// Sample program to show how to takes slices of slices to create different
// views of and make changes to the underlying array.
package main

import "fmt"

func main() {

	// Create a slice with a length of 5 elements and a capacity of 8.
	slice1 := make([]string, 5, 8)
	slice1[0] = "Apple"
	slice1[1] = "Orange"
	slice1[2] = "Banana"
	slice1[3] = "Grape"
	slice1[4] = "Plum"

	inspectSlice(slice1)

	// Question. Why is a slice considered to be a slice?
	// We will use the slice1 to crate slice2. We will reuse the backing array.

	// Take a slice of slice1. We want just indexes 2 and 3.
	// Parameters are [starting_index : (starting_index + length)]
	// Refer 2.png
	// The new slice2 will start at the index 2 which is "a". "b" value which is 4, means not including 4.
	// Another way to look at index b is - instead of looking b as a separate number, we can look at "b" as
	// a+length. If we know our length ahead of time we will not make any mistakes on slicing.
	// Say we only want length of 2. We want only "Banana" and "Grape". We start at index 2 and we need
	// length of 2. We have our "b" as 2+2=4.
	// Our capacity is the total elements from that point in position which is 2, 6.
	// WE wnd up with a slice of length 2 and capacity of 6. See 3.png.
	// [a:a+length]

	// Problem with this is that we are reading and writing using the pointer semantics.
	// Anytime we use pointer semantics we have to be aware of the side affect. Which is when the data gets
	// mutated or changed and we didn't expect it.
	// See slice2 example below.
	slice2 := slice1[2:4]
	inspectSlice(slice2)

	fmt.Println("*************************")

	// Change the value of the index 0 of slice2.
	slice2[0] = "CHANGED"
	// Problem is that we have two slices sharing the backing array memory location in two different ways.
	// slice2 index 0 is the same as the slice1 index 2.
	// So now both slices see the change.
	// So if slice1 didn't expect this data to change then we have created a potential bug.
	// Let's say someone says can't we do this:
	// slice2 = append(slice2,"CHANGED")
	// Length and capacity for slice 2 is not the same. Length is 2 and capacity is 6.
	// Now we append "CHANGED" to the index 2 of slice2 which is index 4 of slice1. See pic 4.png.
	// "Plum" gets changed to "CHANGED". Same side affect. Now on index 4 of slice1.
	//
	// In another scenario say if the new slice2 was of length 2 and capacity 2. Pic 5.png.
	// This would mean that now if we call append. Since no extra capacity so now append will create a copy.
	// we do this by
	// slice2:=slice1[2:4:4]
	// See picture 6.PNG for the output.
	// Slice2 now has length and capacity as 2. Earlier it had length 2 and capacity 6.
	// Now when we call append we now have to make a copy. The length now becomes 3 and capacity doubles to 4.
	// We now also have different addresses. What happens is that slice2 is not sharing the backing array anymore.
	// Slice2 has now it's own copy of everything. See 7.png.
	// Hence we should always be worried when reading and writing with pointer semantics.
	/*
		For an array, pointer to array, or slice a (but not a string), the primary expression
			a[low : high : max]
			constructs a slice of the same type, and with the same length and elements as the
			simple slice expression a[low : high].
			Additionally, it controls the resulting slice's capacity by setting it to max - low.
			Only the first index may be omitted; it defaults to 0. After slicing the array a

			a := [5]int{1, 2, 3, 4, 5}
			t := a[1:3:5]
			the slice t has type []int, length 2, capacity 4, and elements

			t[0] == 2
			t[1] == 3
	*/

	// Display the change across all existing slices.
	inspectSlice(slice1)
	inspectSlice(slice2)

	fmt.Println("*************************")

	// Make a new slice big enough to hold elements of slice 1 and copy the
	// values over using the builtin copy function.
	slice3 := make([]string, len(slice1))
	copy(slice3, slice1)
	inspectSlice(slice3)
}

// inspectSlice exposes the slice header for review.
func inspectSlice(slice []string) {
	fmt.Printf("Length[%d] Capacity[%d]\n", len(slice), cap(slice))
	for i, s := range slice {
		fmt.Printf("[%d] %p %s\n",
			i,
			&slice[i],
			s)
	}
}
