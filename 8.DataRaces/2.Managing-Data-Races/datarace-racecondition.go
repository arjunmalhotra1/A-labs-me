// From the book
package main

import (
	"fmt"
	"sync"
)

func dataRace() {
	counter := 0

	go func() {
		counter++
	}()

	go func() {
		counter++
	}()

	fmt.Println(counter)

}

// A race condition occurs when the behavior depends on the sequence or the timing of the
// the events that can't be controlled. Here timing of events is the go routines's execution order.
func raceCondition() {
	counter := 0
	var mu sync.Mutex

	go func() {
		mu.Lock()
		counter = 10
		mu.Unlock()

	}()

	go func() {
		mu.Lock()
		counter = 9
		mu.Unlock()
	}()

	fmt.Println(counter)
}
