type example struct {
	flag    bool
	counter int16
	pi      float32
}

func main() {
	var e1 example
	// One would think that e1 is of size 7 bytes - 1 + 2 + 4
	// But we would be wrong, we do alignemnts that the compiler has to take care of.
	// Hence it is not a 7 byte value.
	// It is an 8 byte value, because of alignements we get extra one byte.

}

// For a struct like this,
// we will have 3 bytes of padding.
type example struct {
	flag    bool
	[3]byte padding
	counter int32 // Since int32 needs to fall of a 4-byte alignment.
	// int64 needs to fall on a 8-byte alignment. So for int64 it would be 7 bytes of padding.
	pi      float32
}

// Since this below is 1 byte + 7 bytes + 8 bytes + 4 bytes = 20 bytes.
// This entire struct needs to be aligned on an 8 byte boundary we will need to have another 4 bytes in the end to have this aligned on
// a 8 byte boundary.
type example struct {
	flag    bool
	[7]byte padding
	counter int64
	pi      float32
}

// So that becomes this below:
type example struct {
	flag    bool
	[7]byte padding
	counter int64
	pi      float32
	[4]bytes padding
}  

// We can re arrange the fields to be better optimized but we DO NOT do that.
// We order the fields in a way for better readability. We do not do this unless meemory is actually a problem.
type example struct {
	counter int64
	pi float32
	flag bool
}


