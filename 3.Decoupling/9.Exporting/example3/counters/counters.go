/*
	Note we have made this "alertCounter" unexported out of this package but we have added
	a factory function, "New".
	Notice new is constructing and returning an unexported type.
	This is not a good practice.

	This is not really about data access it's about name access.
	Notice in examply3.go we can access the factory function New() adn pass it the value of kind 10.
	The compiler is actually able to create a variable of the unexported type.

	If Bill saw this code then the code review would stop as we are not geting any level of encapsulation here.

*/
// Package counters provides alert counter support.
package counters

// alertCounter is an unexported named type that
// contains an integer counter for alerts.
type alertCounter int

// New creates and returns values of the unexported type alertCounter.
func New(value int) alertCounter {
	return alertCounter(value)
}
