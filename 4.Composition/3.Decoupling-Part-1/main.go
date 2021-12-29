/*
	Code doesn't have to be perfect the first time around, second time around.
	We will follow ROb pike's philosophy, "Discover the interfaces".

	Only way to discover the interfaces is to design code in the concrete first.
	Some concrete implementation first, even if it's just a prototype.

	We will focus on concrete implementation then we will discover those interfaces.
	Discover those coupling points.

	Decoupling is all about polymorphism, and polymorphism is all about code changin it's behavior depending on
	concrete data it is operating on.

	We will focus on how to write code in layers. We want to make sure that
	we have well defined layers in our API.
	Each layer testable and each layer focussed on what's important which is data.

	Most testable code is the code which takes concrete data in and concrete data out and,
	that data can be tested in different scenarios.

	So we always want to be focussing our testability around as much as around concrete data as possible.

	-------------------------------PROBLEM STATEMENT--------------------------------------------------------
	We have Xenia - AS400, with a proprietary database. We also have another system called Pillar
	with another Database.

	We want to transfer the data from the database in Xenia - AS400 to Pillar Database.
	See 1.png
	--------------------------------------------------------------------------------------------------------

	The first thing Bill will tell his client is "If you tell me this data has to be moved real time
	then I am already walking away."

	We have 10 mins to see if the data hasn't been moved yet and then move it.

	As a developer we must know how to break the bigger problem into smaller problems.

	1. How to connect to the proprietary database.
	2. Once we connect to the database, we have to identify how to get the data identify and out of the
	Xenia, proprietary database.
	If we can;t figure out how to solve these 2 then nothing else matters.
	3. Next we have to figure out how to connect with the Pillar database.
	4. Finally we need to get the data transformed from Xenia db to Pillar db and store it in Pillar.

	So before we write an entire program we might have to do some prototyping, to make sure that
	we can do 1 & 2, then we can think about the production level code.
	First if we can solve these 4 primitive problems then we start building the API that kind of is layered.

	One of the things we have got to do is think about what that "primitive" API is. What are the big core
	problems we are trying to solve and can we write an API that solves those problems.

	Then we build a lower level API on top of those primitives.
	See 2.png

	Because if these primitives are workign and are unit testing well then we proceed to Lower level API.
	Only after that we move to Higher level API.

	This is how our APIs are being layered and being tested.
	Developers don't focus on this layering.

	One problems we see when we don't so this layering is that once we start having APIs/functions that kind of
	cross the barriers(green box in 3.png) we usually have big problems figuring out how to test that code.

	A lot of developers come to Bill and say "I have written a function but I don't know how to test it
	I guess I need an interface to mock everything."
	Bill says " There is a place for mocking but it shouldn't be our number one choice."
	Our number one choice should be can I pass concrete data in and get concrete data out & test it that way.

	A lot of times when we have functions that are not testable, is because the code is crossing too many of the
	layers. A lot of times we need to re think what's in the functions and break it back out.

	Note the below code is not for implementation, we should notice how we are structuring the code, how
	we are structuring the layer, and how we are refactoring things to discover the interfaces.
*/
// ----------------------------------------------------------------------
// This is a concrete implementation of the code and if you were workign with Bill and he puts you on this
// project, then he would be telling us that he wants a concrete implementation of this first,
// and not to give him an interface at all in this code base. If you give an interface in this codebase then
// you and him are going to have a problem.
// Focus on the concrete problem first, and right now it's just xenia and pillar.
// When we think about Xenia, we think that Xenia is a system, it has a state that we have to mantain.
// We have an activity which is to pull out data, we have to know the last thing we pulled out.
// There's also state involved in pillar.
// So since there's something stateful about Xenia and Pillar.

// We can define a type and then create a type based API for both Xeni and Pilar.
// We don't care about the data so much here, so we will keep the data simple here.
//------------------------------------------------------------------------
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

// =============================================================================

/*
	This is our primitive layer. Xenia is going to need host information and time out information.
	We assumed fornow that we have solved our connectivity issue, and we will focus on our data problems.
	Since this is going to be our primitive API, we are trying to figure things out and this is our
	first draft of code. We decide that we will define a method of behavior called pull.
	Pulls job is just to find the next piece of data.
	So we take the data value that we want to pull, we share it down the call stack so we don't have any allocations
	there.
	Then pull will worry about just one thing that is find the next piece of data.
	Question. For every piece of data we want to call pull what if there are thousand pieces of data, then
	are we going o call pull a thousand times, isn't that inefficient?
	Answer. Remember thi sis our first draft of code, we are just trying to get this o work.
	Also we are not optimizing for performance right now, we are optimizing for correctness.
	The correct thing to do is to solve this primitive problem, and primitive problem is finding a piece of data.
	If we can find one piece of data then we can find more data.
	Readability is also about writing a piece of code where average team member can mantain it.
	Once we build all of this and it's not fast enough we will try to do some optimizations.
*/
// Xenia is a system we need to pull data from.
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia.
/*
	Here we have solved our second problem. We now know how to find a piece of data out of Xenia.
*/
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

/*
	To keep with the symmetry of the API we keep pillar as a type since Pillar is a system too.
	Keeping with the symmetry of the "Pull" method, we create a method called, "Store"
	Store takes the same data value we just pulled.

	We have 10 mins to make the transfer.
	We don't care about performance rightnow, we just care about making it work.
	By keeping the primitive layer API symmetric, it wil make our API
	readable and maintainable.

	We don't know yet if this is fast enough or not, we are not going to get hung up.
	So now ee have our primitive layer API working.

*/
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
/*
	Next we focus on lower level API. We will build those on top of the primitive function.
	WE have two systems, one knows how to pull and one knows how to store. It might be interesting to create
	a new concrete type called "system" that knows both how to pull and store.
	That will get it's behavior through the embedding of the concrete types.

	We embed "Xenia" and "Pillar" inside of "System".
	Not so "system" has the state but it has the behavior.

	We don't know yet if we are going to use this.
	Some times when we are drafting code, an idea comes to mind we don't know if
	it's an viable idea yet, we just don't want to lose it. So we might just throw the type
	out there. Doesn't mean we will stick with it.
*/

// System wraps Xenia and Pillar together into a single system.
type System struct {
	Xenia
	Pillar
}

// =============================================================================
/*
	We decide to use function based API. Remember, the rule is use functions unless it's not
	reasonable or practical to do so.
	We use method based API for Xenia and Pillar because they were state.
	But now we don't have to mantain anymore statea because that's state can be carried forward
	through Xenia and Pillar.
	That is why we can go back to our function based API.

	WE decide that the pull will take Xenia as a value and will also take a
	collection of data this will maybe help us set up a "batch" pull.
	This will help us setup maybe a batch pull so that we can set up a batch of pulls or a batch of stores.
*/

// pull knows how to pull bulks of data from Xenia.
func pull(x *Xenia, data []Data) (int, error) {
	/*
		We will range over the collection using pointer semantics.
		Then we call pull off of xenia. Trying to pull as much as batch as possible till we
		get that eof or an error.
		Then we return the number of actual values we store inside of that collection.
	*/
	for i := range data {
		if err := x.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data into Pillar.
/*
	Store takes in the pillar and also the batch of data that we just pulled.
	We range over the data again using pointer semantics. We call store to store that data.
*/

func store(p *Pillar, data []Data) (int, error) {
	for i := range data {
		if err := p.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

/*
	Until here we have our lower level API done. The lower level API is sitting on top of primitive API.
	because lower level is calling pull and store on top of Xenia and Pillar.

	We do have an opportunity for the high level API too.
	In high level API we can do pulling and storing in one big batch until there's no more data.
*/
//-------------------------------------------------------------------------------------------------

/*
	This is high level API. We define a function "Copy" but copy needs 2 behaviors it needs both
	systems(Xenia (pull) and Pillar(store) to copy). We composed system with both Xenia and Pillar.
	With one passing of the system pointer we get both pulling and storing from Xenia and Pillar.
	We also set up the batch size.
	Then we go into an endless loop passing in Xenia and Pillar thorugh Pull and Store.
	We do that until the entire set of data is moved for that run.
*/
// Copy knows how to pull and store data from the System.
func Copy(sys *System, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(&sys.Xenia, data)
		if i > 0 {
			if _, err := store(&sys.Pillar, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================
/*
	Finally we build our application. We construct our system value, we initialize
	Xenia nad Pillar. Then we call Copy by sharing the system.
	Every time we run this it happens way within 10 mins. We now have concrete version of our code.
	This is(say) tested and we can put this in production.

	But we are not "done" yet.
	"Done" means 2 things:
	1. We have test coverage of at least 70% if not more. 70% and above is good.
	2. We have implemented our concrete implementation adn then we ask is this decoupled.
	Do we know what changes are coming(changes with respect to the data) and is it decoupled.

	So far we do have concrete version of our code, but we are not done.
	Question. What changes are coming to this program?
	Answer. Very quickly when we put this in prodcution, the sales team that got us involved in this
	project at the begining is going to be super excited, and say
	we have been working on a project named "Bob which is also an AS400" and idea is to migrate data
	off of "Bob". 4.png
	Also there's another system named, "Alis" similar to "Pillar". And we'll are going to
	be able to move some data off of Xenia and Bob to Alis. See 5.png

	What's now going to change?
	the system's are now going to change. So what we'll now do is, ask ourselves,
	is this something we want to decouple ourselves from now or later.

	Let's do the decoupling now. Remember we only have, a concrete working version of this code.
	But now what we can do is, decouple it by discovery.
	See main-1.go

*/
func main() {
	sys := System{
		Xenia: Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Pillar: Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}

	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
