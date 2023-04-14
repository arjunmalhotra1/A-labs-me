/*
	What happens when we are in this context sinuation and you don't want to be,
	We have got to use custom error types, but we want to maintain the level of decoupling.

	This is where we will have to focus from type as context to behavior as context.

	Our context will move from type to behavior. From "what is" to "what does".

	"line, err := c.reader.ReadString('\n')"
	reader is an interface and will have a concrete value inside of it.
	In this case we will be reading a string may be over a network.

	Here we are blocking until one of the 2 things happen.
	1. either we that string, until the "\n(new line" is encountered.
	Or
	2. We get an error.

	Notice here "switch e := err.(type) {" we are still doing type as context.
	What we are saying is let's perform the type assertion on "err".
	And if the concrete value, stored inside of "err" is a pointer to the "OpError", then
	we execute that logic. If it's an address of "AddrError" then we execute that case.
	Since the dns package has a lots of different error types.
	"OpError" is the most common error type that we will use in Go.
	then he showed the implementation of the "OpError" type from the standard library.
	https://pkg.go.dev/net#OpError If you click that heading it will take you to,
	https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/net/net.go;l=435
	Note OpError type is an exported type, all the fields are exported.

	This is where the error interface is implemented - https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/net/net.go;l=464
	Note the pointer symantec implementation.

	In our below code what we are doing in this switch case.
	In every case we care about calling a method, called "temporary".
	Remember Errors are just values they can be anything we need them to be.
	That is not just a state it's also a behavior.
	It is nice to have behavior typed to the errors, like here
	"e.Temporary()". When we are told that whenever there is a networking error, call the
	"Temporary()" method.
	if "e.Temporary()" then we keep going, if not temporary, then we have an issue. May be socket connection died,
	may be listener died.

	// Not important,
	When we look at the "Temporary()" implementation for OpError,
	https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/net/net.go;l=515

	Mostly in every single case we are trying to just validate if the error was temporary or not.
	It's that behavior that we care about, not necessarily the fact that the error is of one
	particular concrete type or another.

	In the same package, Go did implement an interface called Temporary().
	https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/net/net.go;l=511

	type temporary interface {
		Temporary() bool
	}

	Notice that the "temporary interface is an Un-exported type."

	Question. Do we as users of the net package have to implement the temporary interface?
	Answer. No, since temporary interface is only inside the net package there are multiple implementations
	of temporary. So here' the interface that has multiple implementations but it's internal, it's not external.

	But it would have been nice if this would have been external because we could have used it to
	do our own type assertions.

	Thinking about it, when we talk about behavior's context what we are saying is we don't care if the error
	is OP, Add, or DNS. What we care about is the error value has a temporary method that's what we want to work with.

	So now, we can get rid of the code in "TypeAsContext()"
	We can define our own interface,

	type temporary interface {
		Temporary() bool
	}

	We can just call it anything.
	type bill interface {
		Temporary() bool
	}

	We are defining our own bill interface and isn't it true, that OpError, AddressError, DNSConfigError,
	all these error types, in the net package that have temporary method now implement the bill interface?
	Answer. It is absolutely true.
	Which means we can bring all the errors - OpError, AddressError, DNSConfigError,
	and bring it under one single case statement.

	switch e := err.(type)
	case bill:
			if !e.Temporary() {
				log.Println("Temporary: Client leaving chat")
				return
			}

	Now what we can say is the concrete value stored inside of error if
	it also implements the bill interface, then we are going to call "Temporary on it".

	Now we have that behavior as context. This also means that we don't have that breaking change anymore.
	If we add new concrete types, or make changes to those concrete types, this code is stable
	because it's not based on what teh concrete error value is, it's based on what it can do.
	Now we can handle bunch of cases in one place because of behavior.

	If we can focus in a decoupled state adn if it's adding value, it's going to mantain our code base and
	keep it stable longer.




*/

// Package example4 provides code to show how to implement behavior as context.
package example4

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// client represents a single connection in the room.
type client struct {
	name   string
	reader *bufio.Reader
}

// TypeAsContext shows how to check multiple types of possible custom error
// types that can be returned from the net package.
func (c *client) TypeAsContext() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			switch e := err.(type) {
			case *net.OpError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			case *net.AddrError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			case *net.DNSConfigError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			default:
				if err == io.EOF {
					log.Println("EOF: Client leaving chat")
					return
				}

				log.Println("read-routine", err)
			}
		}

		fmt.Println(line)
	}
}

// temporary is declared to test for the existence of the method coming
// from the net package.
type temporary interface {
	Temporary() bool
}

// BehaviorAsContext shows how to check for the behavior of an interface
// that can be returned from the net package.
func (c *client) BehaviorAsContext() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			switch e := err.(type) {
			case temporary:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			default:
				if err == io.EOF {
					log.Println("EOF: Client leaving chat")
					return
				}

				log.Println("read-routine", err)
			}
		}

		fmt.Println(line)
	}
}
