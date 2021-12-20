/*
	https://blog.golang.org/strings

	Go lang strings are always based on UTF-8. UTF-8 character set is a multi tier character set.
	It starts with bytes, then turns into code point and code points turn into characters.

	Go source code is always UTF-8.
	A string holds arbitrary bytes.
	A string literal, absent byte-level escapes, always holds valid UTF-8 sequences.
	Those sequences represent Unicode code points, called runes.
	No guarantee is made in Go that characters in strings are normalized.

	----------------------------------------------------------------------------

	Multiple runes can represent different characters:

	The lower case grave-accented letter à is a character, and it's also a code
	point (U+00E0), but it has other representations.

	We can use the "combining" grave accent code point, U+0300, and attach it to
	the lower case letter a, U+0061, to create the same character à.

	In general, a character may be represented by a number of different sequences
	of code points (runes), and therefore different sequences of UTF-8 bytes.
*/

// Sample program to show how strings have a UTF-8 encoded byte array.
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {

	// Declare a string with both chinese and english characters.
	// the chinese characters require 3 bytes of allocation for that character to be represented.
	// This string is 18 bytes in length.
	// We need 3 bytes to represent one code point of the chinese character.
	// Code points are int32 so max bytes needed to represent a code point is 4 bytes.
	// Hence first 3 bytes are first chinese character and next 3 bytes are next chinese character
	// Beyond that, since everything is in english that would be one byte allocation.
	// Pic 1.png.
	s := "世界 means world"

	// UTFMax is 4 -- up to 4 bytes per encoded rune.
	// UTFMax represents the max number of bytes we need for any UTF-8 based code point.
	// UTFMax is 4.
	var buf [utf8.UTFMax]byte

	// Iterate over the string.
	// We are using the value semantic form of the for range. Why?
	// Because String is a collection of bytes and bytes is a built in type. Hence we use the value sematic
	// form of the for range.
	// When we range over a string we are not going to range by byte, we are not going to range by character
	// we range by "codepoint".
	// Golang has a data type called "rune". Which is really just an alias for int32. int32 represents the
	// 32 bit or 4 byte code point value.
	// 2.Png
	// When we range over a string, we don't range over by byte or character but we range over by codepoint.
	for i, r := range s {

		// Capture the number of bytes for this rune.
		// Here we have gotten the code point value, and ask the UTF-8 package about what is the length
		// or number of bytes that we need to represent that code point value.
		// for i = 0
		// our rune length will be 3.
		rl := utf8.RuneLen(r)

		// Calculate the slice offset for the bytes associated
		// with this rune.
		// for i=0
		// si = 0 + 3 = 3
		si := i + rl

		// Copy of rune from the string to our buffer.
		// This is the built in function copy. Copy takes in a destination and then a source.
		// Only destination we can use for copy could be a slice since slice is the only data structure
		// we can write into. For the source we can choose a string or a slice because string is similar to slice,
		// it is a built-in type and has a backing array.
		// What's also interesting about string is that we can slice a string.
		// Here, s[i:si] we are slicing a string. When we slice a string it will create a new string value.
		// Strings are immutable so we can't be creating a slice values when we slice a string.
		// for i=0, the length is going to be 0 to si=3. 3.png.
		// Remember copy function only works with slices.
		// We don't ave a slice we have an array, "buf" array.
		// "Every array in Go is just a slice waiting to happen."
		// by using the slicing syntax for the array what will happen is,
		// is that we would be constructing a new slice value.
		// But we will be able to use the "buf" array as a backing array fro the slice that is created.
		// 4.png.
		// Basically we are telling copy that take the first 3 bytes of the string(first chinese character)
		// 3 bytes as our source and copy into our buf array by converting the buf array into a slice.
		// Every array is a slice waiting to happen.
		// Again for the display we are converting the buff array into a slice
		// "buf[:rl]" That is we create a new slice value which starts from the begining of the buff array
		// but rune length will be 3 and capacity of 4.
		// On the next iteration "i" will be 3. As we have already processed the byte 1 and byte 2.
		// As the i has to move to the next code point in the string.
		// As we know next is also a chinese character so r will be 3 and si will be 3 + 3 = 6 now.
		// Again a slice is made of the buf array.
		// Now "i" is 3, rune length is 3. So we go from index 3 till index 6.
		// We cody those 3 bytes into our buf array.
		// After this everything will be incrementing by 1.
		// Output - 7.png.
		copy(buf[:], s[i:si])

		// Display the details.
		fmt.Printf("%2d: %q; codepoint: %#6x; encoded bytes: %#v\n", i, r, r, buf[:rl])
	}
}
