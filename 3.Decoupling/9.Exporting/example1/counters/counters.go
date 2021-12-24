// Package counters provides alert counter support.
/*
	Here we name the package after the folder name. It is a good convention to follow
*/
package counters

// AlertCounter is an exported named type that
// contains an integer counter for alerts.

// Note first letter is capital.
// If the thing is named with a capital letter, then it's going to be exported from the package
// Any code outside of the package will have access to it.
// A lower case letter is unexported and it is only for internal use.
type AlertCounter int
