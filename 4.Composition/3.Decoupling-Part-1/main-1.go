/*
	We would like to start with eth high level API. In this case it would be "copy".
	But sometimes it's better to start with the lower level API.
	Because sometimes when we start with the higher level API, it's good. But our higher level API is
	focussing on two behaviors and our lower level API is only focussing on one behavior.
	As JBD says,
	"Start with the smaller interfaces and compose the larger ones."
	So we wil focus on the lower level first and discover the smaller interfaces and then we will compose the
	larger ones.

	When we look at "Pull" we see that Pull is focussed on one behavior "pull" and it accepts "Xenia".
	Right now the function says that "I only work on what the data is."
	What we want to do is we want this function to focus on now, not what the concrete data is but
	what it does.
	We already know that what "Xenia" does is pull.

	So if we change the below code to say let's not base on "what is" but "what does"
	This func pull(x *Xenia, data []Data) (int, error) is changed to
	func pull(p Puller, data []Data) (int, error) now we could define an interface called "Puller".
	that would mean that we call the behavior through the puller interface.
	Now this code is completely decoupled and is precise.

	Now the code doesn't care what the system is, Xenia, Bob, Puller etc.
	As long as it (system) knows how to pull I can work with it.

	func pull(p Puller, data []Data) (int, error) {
		for i := range data {
			if err := p.Pull(&data[i]); err != nil {
				return i, err
			}
		}
		return len(data), nil
	}

	We can do the same with the "store". It's not about the "pillar" implementation but about any concrete data
	that knows how to store. And then execute that behavior in a decoupled state through the store interface "s".

	func store(s Storer, data []Data) (int, error) {
		for i := range data {
			if err := s.Store(&data[i]); err != nil {
				return i, err
			}
		}
		return len(data), nil
	}

	We have just discovered the two interfaces based on the fact that we already have our concrete working
	implementation.

	We are not done yet, because our high level API "Copy" uses the concrete type system.
	copy needs 2 behaviors, not just pull and not just store.
	It needs both pulling and storing.
	We do this a lot in Go, is we do this interface composition.
	We define a new interface composition called, "PullStorer" interface. It could be any concrete type that
	knows both how to pull and store. Which means we can pass in "ps" to both pull(ps, data) and
	store(ps, data[:i])

	func Copy(ps PullStorer, batch int) error {
		data := make([]Data, batch)
		for {
			i, err := pull(ps, data)
			if i > 0 {
				if _, err := store(ps, data[:i]); err != nil {
					return err
				}
			}

			if err != nil {
				return err
			}
		}
	}

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

/*
	This is what JBD is trying to tell to us, to compose these larger interfaces.
	We can leverage composition and say that PullStorer is the composition of both Puller and Storer.
	Now any concrete type that knows how to both pull and store is a pullStorer.
*/
type PullStorer interface {
	Puller
	Storer
}

/*
	We know already that a system is a pullStorer because of the embedding of Xenia and Pillar.
	Remember we are embedding for behavior, system knows how to pull. System knows how to store.
	Which means we can pass system pointer to copy.
*/

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

// System wraps Xenia and Pillar together into a single system.
// type System struct {
// 	Xenia
// 	Pillar
// }

type System struct {
	Puller
	Storer
}

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
func Copy(ps PullStorer, batch int) error {
	data := make([]Data, batch)

	for {
		/*
			pull(p Puller, data []Data) (int, error)
			Question. ps is of type "PullStore", but pull function takes in a puller.
			How is this still working? How is the compiler okay with this?
			When we learnt our interface mechanics, we learnt that
			Does PullStore exist? Does Puller exist?
			Does Store exist?
			The answer is no, they are interfaces, they are valueless.
			Next,
			Question. Is the pull function really asking for "Puller" value?
			Is store function really asking for "Storer" value?
			Answer. No, they can't since those values don't exist.
			These functions are saying, to pass in any piece of concrete data, any value of pointer
			that knows how to behave like a storer or behave like a puller.

			Question. Isn't it true that any concrete data sitting inside of the "ps" interface value "PullStorer"
			knows how to pull and store?
			Answer. It does.

			What's happening here is that we are not passing "ps" to the pull function.
			We are passing the system pointer that is embedded inside of "ps".
			Always keep in mind that it's only the concrete data that can be passed across the
			program boundaries.
			So now the lower level and the higher level APIs have been decoupled.
			They are not completely decoupled yet.
			In main we are constructing our system.
			System has Xenia nad Pillar. We have embedded these values inside of System.
			Xenia knows how to pull and Pillar knows how to store.
			Which means system knows how to pull and system knows how to store.
			See 6.png
			The Embedding is adding the layer of indirection and it is this layer of indirection that
			gives us the decoupling.

			Now when we call "copy", "copy" is saying I want a PullStorer.
			"func Copy(ps PullStorer, batch int)"
			But in main we are calling "Copy" like - "Copy(&sys, 3)"
			We are passing the address of system across the program boundary.
			When we pass the address of system across the program boundary
			See 7.png
			We get our ps variable our valueless "ps", which has our system address into "ps"
			We also know that "ps" also knows how to pull and store.
			We just added another layer of indirection here.
			When we "p" off of "ps" we are actually calling "p" off of "sys" which in turn is calling
			"p" off of Xenia.
			All of this indirection is giving us all of this decoupling.

			When we call pull and we pass "ps", we are not passing "ps" since "ps" isn't real.
			We are passing the concrete address of "sys" again.
			So now when we get to pull this also will have a pointer to system and it will also point to
			to "sys" as well. Except the "Puller" only knows how to "pull".
			See 8.png
			We will have the same exact situation with "Store".

			This is how we can use the interface variables on these function calls, because we are not really
			passing the interace values, we are passing the concrete data inside it because that is the only data
			that can be moved around our program.

			We do have a problem since we are not fully decoupled yet.
			Because if we go into main and if we want to use "Bob" instead of "Xenia".
			System is tied to Bob directly so, this poses a problem.
			We can't be breaking System everytime we want to use "Bob". But may be we want to run a
			couple of copies at the same time.
			We do have a PullStorer interface, copy function doesn't care what the concrete data is as
			long as it knows how to pull and store.

			One solution could be that we make multiple system types for all combinations, but this
			is an irrational solution. As the number of serves grow there are going to be too many combinations.

			type System struct {
				Xenia
				Pillar
			}

			type System2 struct {
				Bob
				Pillar
			}

			type System3 struct {
				Bob
				Alice
			}

			Question. How do we solve this problem without breaking the type?
			Answer. Composition, what's nice about composition is that the system type doesn't have to be
			just embedded around concrete values, system could also be the composition of interface types.

			type System struct {
				Puller
				Storer
			}

			This will really allow us to inject data into our concrete system. This will allow us to decouple
			behavior at the type level.
			Now in main we are not initializing the embedded "Xenia", we are initializing the
			embedded interface. For both "Puller" and "Storer".

			Now if we want to work with Bob data, we can simply do:
			"Puller: &Bob{" this is something we can do dynamically. Now our code really is
			truly truly decoupled from all the concrete. We now completely decoupled from the change.

			System is now at another level of indirection. System is now based on Puller and Storer. See 9.png

			During initialization we will embed a "Xenia" value inside Puller and "Pillar" value inside
			our Storer. See 10.png. And we have added an extra level of indirection.
			Now we have a fully decoupled program.

			Remember we are always drafting code. One of the most important thing we need to do after
			we get a piece of code working is to perform a "readability" review.
			To validate mental models, to validate the code we have. Clean it up. Variable names.
			The ordering of things, where things are relocated.

			See main-2.go


		*/
		i, err := pull(ps, data)
		if i > 0 {
			/*
				store(s Storer, data []Data) (int, error)
			*/
			if _, err := store(ps, data[:i]); err != nil {
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
	sys := System{
		/*
			Now if we ant to work with "Bob" data, all we can do is:
			"Puller: &Bob{"
		*/
		Puller: &Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Storer: &Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}

	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
