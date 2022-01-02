/*
	One of things common is that "Everything needs to be mocked."
	Go is very data oriented and one of the best tests are the ones that are data oriented as well.
	Concrete data in and concrete data out, when we do those testing.

	Not that we shouldn't and can't mock. Mocking is the last resort, a lot of times we can just refactor the
	code a little bit, to make it more testable.

	Another beautiful part of Go's compiler (the ability of having convention over configuration -
	don't know what this meant) is the responsibility for mocking has moved away
	from the developer writing the API to the developer using the API.

	Imagine a client who comes and says "Bill I need a go version of Pub-Sub API"
	We have an internal messaging bus that we use, we don't use Rabbit MQ or anybody's else.
	We have implemented our own and developers are allowed to interact with this bus in any language they want.
	Beautiful thing about message bus is that it doesn't matter what language we are talking about.
	It's all message passing.
	See 1.png

	They have created a message bus and they say that everybody who is writing services, connect to this bus.
	We already have programs written in Java/Python, C/C++, C#.
	What they want to do is write add "Golang". So you need to write a package in Go, that implements all the
	binary protocols that relates to the message bus.
	First question, Bill will ask is "How many message buses are there?"
	They would say "Just this one."

	So even though Bill will have a concrete of this first, we know upfront that we don't have a decoupling problem here.
	Since there's only one implementation that we have to have. There's only one physical bus,
	one physical set of protocols. What's nice is when we are done with the concrete implementation
	we don't have to worry about any interfaces and any decoupling.

	So we go ahead and denfine our pubsub type, we also define our factory function that constructs & returns that
	pub sub pointer. Then we implement say, 10 methods as part of the API.

	When we are done with this API then we would want to write tests.
	One of the first things we can do is that we could write mocks for our tests.
	But one we don't need mocks. and two, we are not going ot have real binary protocol.

	If we can write tests that talk to a real message bus, life would be much better for me.
	We wouldn't be guessing. Knowing this is working because our mocks are
	also reimplementation of those internals.

	One of the nice things about docker is that we can wrap this mesage bus, in a docker container
	and run our tests agains the docker container directly.

	Say we wrap our message bus in a docker container and we write tests that hit that message bus directly.
	And we go to that developer who would be writing those Go services and say, this package is ready for you and
	now you may go start writing your code.

	A couple of days later the developer wants to write tests and they come to us saying "Hey! Bill,
	I have to write my tests and I have to mock my tests because I cannot hit a real message bus in my test
	environment."

	And Bill says to them,
	"Mocking the mssage bus is the right thing for you because, you should assume that I am doing my job and my API.
	And you shouldn't worry that, all of those message bus calls you are going to make, those pub sub calls
	will not work. I guarantee you that they work. And so mocking for you is okay."

	They say that back to us,
	"Bill I appreciate you said that, can you give me an interface that I can leverage to write the tests."
	Reason behind them saying this is because they are coming from some other languages, those interfaces
	have to be declared against those concrete types like pub sub, in order for them to be able to do this
	type of mocking.

	Well, if you have a mocking need you can do the mocking yourself, you can actually define those interaces yourself.

	See main-1.go



*/

// Package pubsub simulates a package that provides publication/subscription
// type services.
package pubsub

// PubSub provides access to a queue system.
type PubSub struct {
	host string

	// PRETEND THERE ARE MORE FIELDS.
}

// New creates a pubsub value for use.
func New(host string) *PubSub {
	ps := PubSub{
		host: host,
	}

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.

	return &ps
}

// Publish sends the data for the specified key.
func (ps *PubSub) Publish(key string, v interface{}) error {

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Subscribe sets up an request to receive messages for the specified key.
func (ps *PubSub) Subscribe(key string) error {

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}
