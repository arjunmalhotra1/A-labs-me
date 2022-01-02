/*
	We ideally want to have concrete implementation first and then
	discover what can be decoupled, apply the idea of what can be decoupled.
	That precision around decoupling.

	When we are polluting with interface, not only we are taking the allocation cost,
	but we are also making the code harder to read and harder to mantain. And we are not gainign any values.
	As a developer that' why we do the readability reviews. WE are trying to figure out, where are the smells and
	how are we going to minimize the pain and keep those mental models in place.

	Example of interface pollution. This is from an actual developer who was building a package to help people
	write TCP servers.

	Bill's review stopped as soon as he saw the three behaviors.

	type Server interface {
		Start() error
		Stop() error
		Wait() error
	}

	Problems with the above interface:
	1. First the name of the interface is "Server" it is a noun it is describing the interface. We talked about
	concrete data, what's real is what describes the person, places and things.
	Interfaces are supposed to be describing behavior. "Server" isn't even in the scope of describing a behavior.

	2. Other problem is that this interface is a 3 method interface. That's not necessarily wrong but
	these are discreet behaviors that could be composed later on to form the behaviors when needed.
	(I didn't understand this point)
	Bill - " I can't imagine that the every aspect of the code that we are about to look at,
	always needs all 3 behaviors." so there's smell there as well.

	3. Also since the developer said we are making this package to help people write servers.
	Bill's already thinking do you have 2 implementations of this support?
	Because remember we don't need an interface unless we need to decouple something.
	Decoupling means we will hav more than one implementation of something.
	Polymorphism means that a piece of code changes behavior depending upon concrete data it is operating on.
	But, if it's only operating on one piece of concrete data we don't need that decoupling.
	We can work with what is and not worry about "what does".

	One of the first question, bill asked the developer was "Do you have multiple implementations of the server
	or are you asking the developer (who is using the package) to write this behavior."
	In both the cases the answer was "No".
	Bill's friend - No I am not asking the developer to implement this. And No I only have one implementation.

	These were the signs that the developer was designing with the interface and not focussing on their concrete
	implementation first.
	---------------------------------------------------------------------------------------------------------------

	Even when we start to go further down the code we see that concrete type "server"
	which is the true implementation for what he was building, is starting off as an "unexported type".

	type server struct {
		host string
		// PRETEND THERE ARE MORE FIELDS.
	}

	But we also have a factory function that is returning not the concrete value but the interface value.

	func NewServer(host string) Server {

		// SMELL - Storing an unexported type pointer in the interface.
		return &server{host}
	}

	This is a smell, as the whole idea of the concrete function is to construct the concrete data initialize it for use
	and send it back up the call stack, whether its a value or pointer semantic.
	But this developer is sending back the interface, there are little to no value here.
	But there are some exceptions to these rules the network package with the con is a fantastic
	example that is kind of on the edge of these guidelines. But for the most part, these guidelines hold true.

	There's a big problem in main function because, as a user of this package do I care if
	srv is an interface or concrete type.
	'
		// Use the API.
		srv.Start()
		srv.Stop()
		srv.Wait()
	'
	the fact hat we aregetting an interface didn't make our life any better. It just added a level of indirection
	and allocation. We are taking costs here for zero gain.

	So after discussions they removed the interface. See main-1.go


*/

// This is an example that creates interface pollution
// by improperly using an interface when one is not needed.
package main

// Server defines a contract for tcp servers.
type Server interface {
	Start() error
	Stop() error
	Wait() error
}

// server is our Server implementation.
type server struct {
	host string

	// PRETEND THERE ARE MORE FIELDS.
}

// NewServer returns an interface value of type Server
// with a server implementation.
func NewServer(host string) Server {

	// SMELL - Storing an unexported type pointer in the interface.
	return &server{host}
}

// Start allows the server to begin to accept requests.
func (s *server) Start() error {

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Stop shuts the server down.
func (s *server) Stop() error {

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Wait prevents the server from accepting new connections.
func (s *server) Wait() error {

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

func main() {

	// Create a new Server.
	srv := NewServer("localhost")

	// Use the API.
	srv.Start()
	srv.Stop()
	srv.Wait()
}

// =============================================================================

// NOTES:

// Smells:
//  * The package declares an interface that matches the entire API of its own concrete type.
//  * The interface is exported but the concrete type is unexported.
//  * The factory function returns the interface value with the unexported concrete type value inside.
//  * The interface can be removed and nothing changes for the user of the API.
//  * The interface is not decoupling the API from change.
