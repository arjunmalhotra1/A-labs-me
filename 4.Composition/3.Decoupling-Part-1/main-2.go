/*
	For our redability review, in this case we can start from the "higher Level API",
	Copy.
	Question. We have PullStorer interface, but how many implementations of PullSotrer are we
	going to have in this program?

	Answer. Rememebr one of the changes we made is that now the system type we can inject data inside of it
	If we can inject data inside of system then we don't need more than one system type through out the entire
	program. Now if we only need one system type, then we are not going to have multiple implementations of
	"PullStorer". If we are not going to have multiple implementations of "PullStorer" then that's just
	TypePollution. So we can get rid of the PullStorer interface and just get the concrete System pointer
	and we would still be as decoupled as before and with one less type.

	"func Copy(ps PullStorer, batch int) error" is changed to
	"func Copy(sys *System, batch int) error"

	It was a cool idea when we created PullStorer when we were creating a bigger interface but since all the
	data injections were happening inside the system, that's the only type we are going to have.

	There's another readability aspect here that's important, it's about precision of an API.
	Having an API that's is precise helps with testing, readability. It will help with
	maintainability and it will help with debugging.
	API precision is everything. It is usually during these readability reviews we can think about "Precision".
	Encapsulation is an area where we create this kind of precise semantic or precise behavior, where we can
	minimize any misuse or fraud.

	There is actually a little potential misuse of fraud on this API call of "Copy" function.
	It is only asking for "system", it is not clear what "system" has to be initialized to
	in order for this program to work and at an API level it is not really clear what system is all about.
	So if we change the "Copy" API to not ask for system but ask for the actual behaviors that we need.

	"func Copy(p Puller, s Storer, batch int) error"
	Question. Is this not a much better precise API? Isn't this much clearer?
	Isn't this saying I need a concrete piece of data that knows how to pull, I need a concrete piece of data
	knows how to store.
	The earlier idea of doing data injection was cool but it isn't really making the program any easier
	to understand.
	Because of this precision change, one of the things we can now do is, simplify even main.

*/
// Sample program demonstrating struct composition.
package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================

// Data is the structure of the data we are copying.
type Data struct {
	Line string
}

type Puller interface {
	// Puller is based on Pull behavior
	Pull(d *Data) error
}

type Storer interface {
	// Storer is based on Store behavior
	Store(d *Data) error
}

// type PullStorer interface {
// 	Puller
// 	Storer
// }

// =============================================================================

// Xenia is a system we need to pull data from.
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia.
func (*Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("error reading data from Xenia")

	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil
	}
}

// Pillar is a system we need to store data into.
type Pillar struct {
	Host    string
	Timeout time.Duration
}

// Store knows how to store data into Pillar.
func (*Pillar) Store(d *Data) error {
	fmt.Println("Out:", d.Line)
	return nil
}

// =============================================================================
// type System struct {
// 	Puller
// 	Storer
// }

// =============================================================================

// pull knows how to pull bulks of data from Xenia.
func pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data into Pillar.
func store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// Copy knows how to pull and store data from the System.
//func Copy(sys *System, batch int) error {
func Copy(p Puller, s Storer, batch int) error {
	data := make([]Data, batch)

	for {
		// i, err := pull(sys, data)
		i, err := pull(p, data)
		if i > 0 {
			// if _, err := store(sys, data[:i]); err != nil {
			if _, err := store(s, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {
	// sys := System{
	// 	Puller: &Xenia{
	// 		Host:    "localhost:8000",
	// 		Timeout: time.Second,
	// 	},
	// 	Storer: &Pillar{
	// 		Host:    "localhost:9000",
	// 		Timeout: time.Second,
	// 	},
	// }
	/*
		Now we just construct Xenia and construct Pillar, Now during
		initialization we are focussing only on the concrete data like we should be.
		And we can pass xenia nad pillar to Copy.
		We can pass "bob", "alice", we have made main a lot easier to code and initialize.
		This code is still completely decoupled.

		"Code readability" review stuff is all about validating the choices we made while prototyping,
		because we were trying to make things to work. Validating those variable names,
		validating function names, validating are the APIs precise, validating pollution, like we
		don't need "system" any more, like we don't need "PullStorer" anymore.

		It was great while we were drafting, code but at the end of the day once we have a working solution,
		we started realizing that it really wasn't necessary.

		Remember, we are not going to get tho this state on Day 1. We will need iterations.

		"A Good API is not easy to use but hard to misuse."
		This is where we can make APIS more precise and making them easy to understand,
		and hence minimize runtime issues by catching them at compile time.

		"We can always Embed what we can't decompose."
		We created "PullStorer" interface by embedding the "Puller" and "Storer" interfaces.
		This was the composition at the interface level.

		"We don't design with interfaces, we discover them."
		Create that concrete implementation first. If we can put the code in production even better.
		If not then with a good strong working prototype at least we have validated the ideas.
		Validated our thoughts, we even validated some potential performance issues, because
		we can debug and test it.

		Then we can ask, how can we decouple it?
		How do we identify the interfaces what we need and where they go.

		Finally the idea,
		"Duplication is far cheaper than wrong abstraction"
		We saw this earlier in the grouping exercise.(I don't remember which one)
	*/
	x := Xenia{
		Host:    "localhost:8000",
		Timeout: time.Second,
	}
	p := Pillar{
		Host:    "localhost:9000",
		Timeout: time.Second,
	}

	if err := Copy(&x, &p, 3); err != io.EOF {
		fmt.Println(err)
	}
}
