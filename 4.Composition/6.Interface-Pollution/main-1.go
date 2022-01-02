/*
	1. After discussing they removed the interface.

	2. Then made the concrete type as an exported type.
	type Server struct {
		host string

		// PRETEND THERE ARE MORE FIELDS.
	}

	3. The function "Start()" now returns the address of the concrete data.

	The main looks identical to what we saw earlier.

	Now when are working with the concrete data, it gives the developer the ability to do  more and
	to get rid of allocations if that's something that ends up being important to the developer.
*/

// This is an example that removes the interface pollution by
// removing the interface and using the concrete type directly.
package main

// Server is our Server implementation.
type Server struct {
	host string

	// PRETEND THERE ARE MORE FIELDS.
}

// NewServer returns an interface value of type Server
// with a server implementation.
func NewServer(host string) *Server {

	// SMELL - Storing an unexported type pointer in the interface.
	return &Server{host}
}

// Start allows the server to begin to accept requests.
func (s *Server) Start() error {

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Stop shuts the server down.
func (s *Server) Stop() error {

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Wait prevents the server from accepting new connections.
func (s *Server) Wait() error {

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
// 1. Either we need some implementation detail from the caller from the user.
// 2. We are going to have multiple implementations of something, or
// 3. There are very obvious and clear piece in our code, that need to be decoupled (idea of polymorphism).
// NOTES:
// These are the guidelines that we should use when we are going to use an interface.
// Here are some guidelines around interface pollution:
// * Use an interface:
//      * We need an interface when users of our API need to provide an implementation detail that we can't do ourselves.
// 			We cannot know all the implementations that's the core idea behind polymorphism.
//      * If our APIs have multiple implementations that need to be maintained.
//			If we are going to have multiple implementations of something,then we need that indirection.
//			The indirection the decoupling gives us ability to have multiple implementation.
//      * When parts of the API that can change have been identified and require decoupling.

// Everytime we see an interface and if it's not obvious and clear why is it there, question it.
//	1. If the developer is saying that the interface only exists so that we can write test, that's really bad.
// We are writing code for users, we are writing code to build apps not to build tests.
// So if we are having struggles writing tests then we are not going to make changes tot he API to make it testable.
// We want to make sure that the code is usable first. We may make changes to make the code testable but for
// Bill it doesn't mean that we go ahead and start add interfaces and new types to make testing possible.
// It's more about how do we break down these functions that are not testable, into functions that can take
// concrete data in and give concrete data out. Reorganize somethings in the code and not make the API more generalized.
//
// We are lookign to make sure that the interface is there, then it's obvious that it's providing decoupling support
// if not obvious then probably polluting the code. If it's not clear that an interface is making the code better,
//  it probably isn't.
//
// * Question an interface:
//      * When its only purpose is for writing testable API’s (write usable API’s first).
//      * When it’s not providing support for the API to decouple from change.
//      * When it's not clear how the interface makes the code better.
