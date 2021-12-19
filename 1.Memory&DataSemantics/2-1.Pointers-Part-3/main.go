// Sample program to teach the mechanics of escape analysis.
package main

// user represents a user in the system.
type user struct {
	name  string
	email string
}

func main() {
	// In go we do not have constructors but we have 'factory fucntions'.
	// Factory functions can construct a value, initialize it for property use, and then return it back to the caller.
	// Here we have 2 versions of the create user factory fucntions.
	// One of the first things we look for in a factory fucntion is the return type.
	// Return is going to tell us a lot about the data semantics.
	u1 := createUserV1()
	u2 := createUserV2()

	println("u1", &u1, "u2", u2)
}

// createUserV1 creates a user value and returns
// a copy back to the caller(own copy of the caller).
// This is a value semantics fucntion.
func createUserV1() user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V1", &u)

	return u
}

// createUserV2 creates a user value and shares
// the value with the caller. The caller will get the "shared" access to the
// value of the user this fucntion is constructing.
// This is pointer semantic.
func createUserV2() *user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V2", &u)

	return &u
}

// Go compiler can perfom, static code analysis.
// One type of static code analysiss that the compiler performs is called "escape analysis".
// "escape analysis" is trying to read the code during the compile time and trying to determine,
// where a value should be constructed in memory. Whether it's ont he stack or on the heap.
// When a value is constructed on the hep we call it an "escape". All values are constructed just once.
// "Escape analysis" is not about construction it's about how a value is shared.
// Construction on lines 41-43 doesn't tell us anything. We still don't know where that value needs to be constructed.
// Stack or the heap. On line 46 we are sharing hte value down to the print function.
// On line 48 we wil be "sharing up". Another fucntion call from the main might wipe out this V2's frame.
// With escape analysis if the compiler recognizes code where value is being shared up the call stack, it immediately
// decides not to construct the value on the stack and it's going to have to be on the heap.
// Line 48 is going to cause the escape construction on the heap.
// Now variable "u" is created outside the frame. Only way to access it is through a pointer.
// "u" represents a value not on the stack but on the heap. Unlike other laguages where it could only have been on the
// stack.
// On line 48 when we shrae it with main we are copying the heap address.
// Anything on the heap is the job of the Garbage collector to manage.
// anything on the stack is self cleaning. Garbage collector will add to the internal latency.
// Compiler looks at the code and sees how we are sharing (not constructing) and then makes a decision on to
// have the values constructed on the heap or the stack.
// Now as soon as we see code like line 48, we can say that's sharing up the call stack we know that
// it's going to be a cost on our software because that's going to create a latency cost in our software with the
// garbage collector.
// but if we see just say "return u" this doesn't tell us anything. We cannot say that this is only valuse
// sematics as the construction could have been done like this:

// 	u := &user{
//		name:  "Bill",
//		email: "bill@ardanlabs.com",
//	}

// We just lost readablility if we do something like above. We no longer know the cose without reading more code.

// GUIDELINES
// 1. When we do construction to a variable then we will be using value construction. Value construction
// that way we don't use the readability or the cost on line 48.
// But a pointer semantic on the return is good as below:
// 	u := return &user{
//		name:  "Bill",
//		email: "bill@ardanlabs.com",
//	}

// Or construction call inside a function call is good too.
// 	u := foo(&user{
//		name:  "Bill",
//		email: "bill@ardanlabs.com",
//	}) This is clear and we do not lose readability but can be formatted better.

// Point 1 will make our code a lot more readable.
// "Allocation" is reserved for construction of value in the heap.

/*

// We can use the gc flag, go compiler flag. using -m=2 the compiler will not build but produce the
// escape analysis report.
// See escape analysis and inlining decisions.
// Notice "./example4.go:39:2: u escapes to heap:" Which means that "u" is constructed on the heap
// becasue of the line " from &u (address-of) at ./example4.go:46:9"
// This could be used later on for profiling. As the profiler can show us, "what is allocating".
// For why, we need these compiler reports.
// Compiler tells us the decisions so that we do not have to guess about anything in golang.

$ go build -gcflags -m=2
# github.com/ardanlabs/gotraining/topics/go/language/pointers/example4
./example4.go:24:6: cannot inline createUserV1: marked go:noinline
./example4.go:38:6: cannot inline createUserV2: marked go:noinline
./example4.go:14:6: cannot inline main: function too complex: cost 132 exceeds budget 80
./example4.go:39:2: u escapes to heap:
./example4.go:39:2:   flow: ~r0 = &u:
./example4.go:39:2:     from &u (address-of) at ./example4.go:46:9
./example4.go:39:2:     from return &u (return) at ./example4.go:46:2
./example4.go:39:2: moved to heap: u

// See the intermediate representation phase before
// generating the actual arch-specific assembly.

$ go build -gcflags -S
CALL	"".createUserV1(SB)
	0x0026 00038 MOVQ	(SP), AX
	0x002a 00042 MOVQ	8(SP), CX
	0x002f 00047 MOVQ	16(SP), DX
	0x0034 00052 MOVQ	24(SP), BX
	0x0039 00057 MOVQ	AX, "".u1+40(SP)
	0x003e 00062 MOVQ	CX, "".u1+48(SP)
	0x0043 00067 MOVQ	DX, "".u1+56(SP)
	0x0048 00072 MOVQ	BX, "".u1+64(SP)
	0x004d 00077 PCDATA	$1,

// See bounds checking decisions.

go build -gcflags="-d=ssa/check_bce/debug=1"

// See the actual machine representation by using
// the disasembler.

$ go tool objdump -s main.main example4
TEXT main.main(SB) github.com/ardanlabs/gotraining/topics/go/language/pointers/example4/example4.go
  example4.go:15	0x105e281		e8ba000000		CALL main.createUserV1(SB)
  example4.go:15	0x105e286		488b0424		MOVQ 0(SP), AX
  example4.go:15	0x105e28a		488b4c2408		MOVQ 0x8(SP), CX
  example4.go:15	0x105e28f		488b542410		MOVQ 0x10(SP), DX
  example4.go:15	0x105e294		488b5c2418		MOVQ 0x18(SP), BX
  example4.go:15	0x105e299		4889442428		MOVQ AX, 0x28(SP)
  example4.go:15	0x105e29e		48894c2430		MOVQ CX, 0x30(SP)
  example4.go:15	0x105e2a3		4889542438		MOVQ DX, 0x38(SP)
  example4.go:15	0x105e2a8		48895c2440		MOVQ BX, 0x40(SP)

// See a list of the symbols in an artifact with
// annotations and size.

$ go tool nm example4
 105e340 T main.createUserV1
 105e420 T main.createUserV2
 105e260 T main.main
 10cb230 B os.executablePath
*/
