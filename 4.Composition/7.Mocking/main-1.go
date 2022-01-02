/*
	Bill reply to the dev "
	With Golang, if you have a mocking need you can do the mocking yourself,
	you can actually define those interfaces yourself.
	If you have a mocking need you can do the mocking yourself, you can actually define those interfaces yourself.

	If you go ahead and do what I am doing here:
	type publisher interface {
		Publish(key string, v interface{}) error
		Subscribe(key string) error
	}
	why don't you define your own publisher interface and lay out just the methods from the API that you are using.
	You can define your entire own set of interfaces, that my concrete type satisfies.
	Then throughout your entire add you use that publisher interface, or different versions of the publisher interface
	through composition. So our APIs are precise. Bill don't need to do that anymore as a dev you can do it yourself."

	Now the developer can go ahead and write their mock. They can implement a mock implementation of pub sub for their
	test. They can use Bill's concrete implementation for their actual apps.

	(I didn't understand this)
		pubs := []publisher{
		pubsub.New("localhost"),
		&mock{},
	}
	What's nice is they can do the construction of either/or depending on where they are in their code base.

	This is one of Go's legacy. The idea that the developer, writing the API doesn't have to think about the user's
	test. Doesn't have to think about the decoupling needs for someone else.

	So if bill is adding an interface into his (pub/sub) package it's because he needs that interface.
	Bill need to a start asking for data based on not what it is but based on what it does. Bill doesn't have to think about
	what the user of the API needs, they can satisfy their own requirements.

	This will keep the code minimized, and rally clean and concise thorughout our entire codebase.
	BEcause everybody can just take focus and responsibility for what they  need

	Interface can provode the mocking support, we want to use it as a last resort.
	If we can use the docker container, to put real tech in them to test let's do that.
	We don't have to worry about someone elses test and other mocking ends, they can wite their own interfaces.

	We are not going to design with interfaces we are going to discover them adn everybody can discover the
	interfaces they need for themselves.




*/
// Sample program to show how you can personally mock concrete types when
// you need to for your own packages or tests.
package main

import (
	"github.com/ardanlabs/gotraining/topics/go/design/composition/mocking/example1/pubsub"
)

// publisher is an interface to allow this package to mock the pubsub
// package support.
type publisher interface {
	Publish(key string, v interface{}) error
	Subscribe(key string) error
}

// mock is a concrete type to help support the mocking of the pubsub package.
type mock struct{}

// Publish implements the publisher interface for the mock.
func (m *mock) Publish(key string, v interface{}) error {

	// ADD YOUR MOCK FOR THE PUBLISH CALL.
	return nil
}

// Subscribe implements the publisher interface for the mock.
func (m *mock) Subscribe(key string) error {

	// ADD YOUR MOCK FOR THE SUBSCRIBE CALL.
	return nil
}

func main() {

	// Create a slice of publisher interface values. Assign
	// the address of a pubsub.PubSub value and the address of
	// a mock value.
	pubs := []publisher{
		pubsub.New("localhost"),
		&mock{},
	}

	// Range over the interface value to see how the publisher
	// interface provides the level of decoupling the user needs.
	// The pubsub package did not need to provide the interface type.
	for _, p := range pubs {
		p.Publish("key", "value")
		p.Subscribe("key")
	}
}
