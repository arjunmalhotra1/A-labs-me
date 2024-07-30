/*
	We have seen already how to use wait groups to achieve orchestration.
	My thoughts:
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	In the previous example when we used wait groups, the go routine was kind of like telling the main go routine that "Hey! I have completed my work and we can move ahead now."
	The difference between wait groups and channels is that when using channels the go routines can exchange information or data(signalling with data as mentioned below) with each other as well.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	But one of Go's biggest feature is the channel.

	It is channels that allow us to do orchestration in a much simpler way.

	There are basic patterns around channels that will help us with orchestration.

	A lot of developers make mistake by thinking of channels as a data structure/queue.
	Do not think of channels as a Data structure or as a queue.

	Channel provide one basic thing that is "signalling". A channel allows one go routine to signal to another
	go routine about some event. Sometimes we can signal with data sometime we can't.

	First thing we should think about when it comes to signalling is on the send side.
	He always like focussing on the send side.
	Question. Does the go routine performing the send, (remember this is signalling so we say send and receive
	not read and write), Is the go routine that is signalling performing the send, does that go routine need a guarantee
	that the signal being send has been received?
	Answer.
	Story line -
		Say Bill works with jack everyday. And in the morning Bill gives work to Jack.
		Jack usually has it done before lunch and Bill can go an pick it up.
		Normally, Jack is not always on the same time and it hasn't mattered before because Bill could just
		walk in and put over Jack's desk knowing that Jack is going to get it done.
		But a couple of days ago, when Bill left Jack some work, Jack couldn't get it done.
		In the end the things escalated and Bill got into trouble with his manager for
		not being able to get the work done. Jack told the manager that Bill never gave him the work and that is why
		Jack couldn't finish the work.

		Moving forward now Bill needs a guarantee that when Bill gives Jack work, Jack has received it.
		When Bill moved into the office next day, at 9am , Jack is not in.
		Bill is stuck at Jack's desk waiting for Jack to come into the office,
		because Bill needs a guarantee that Jack receives the work.

		When Jack comes in only when Jack pulls the work out of Bill's hands and says now I have it.
		That is where we get this guarantee of delivery. Jack pulls it out of Bills hand.
		Receive happens before the send.

		Then in the Afternoon Bill comes to get the done work and Jack is not at his desk.
		Bill is stuck at Jack's desk again Bill waits for Jack to come back from lunch.
		Then Jack perform his send, and Bill performs his receive and receive happens before the send.
		That's how Jack gets the guarantee.

		Guarantees are critically important they will help us with consistencies.

	As we saw guarantee is not free.
	If we want the guarantee, that a signal being sent has been received and that we have to wait for the receive
	to happen before the send. Then we have to live with the unknown latency cost.
	Remember latency is slowing down your program's throughput.
	To get the guarantees we have to consider that there may be unknown and high latencies.

	While you are waiting on the send side to be able to perform the send because the receiver has to come along.
	Or even on the other side if receiver comes first but we are waiting for the sender.
	Sender and receive has to come along.

		Okay so since Bill and Jack are losing a lot of time on this guarantee.
		So now when Bill comes to the office and when Jack is not at his desk, Bill doesn't care and
		can perform his signal/send place the piece of work on Jack's desk and walk away.
		The send happens before the receive here. Now Bill doesn't have to incur the latencies.
		But now Bill has a cost of risk (of guarantee) as Bill doesn't know when Jack is coming in,
		when he is going to take the work.

		We have to judge the situations and patterns. Some times we can take the latencies sometimes we can
		walk away from them and still make sure that we have high levels of integrity.

	There has to be a guarantee somewhere if it's not going to be on the channel signalling then
	it's got to be somewhere in the code so that we know some concurrent operation has been completed.

	So coming back to the question.
	Question. Do we need a guarantee that the signal being sent has been received?

	Other question we need to talk about is:
	Question. Do we want to signal with data or without data?
	Answer. In our above Bill and Jack example we were signalling with data. We had a pice of data (work) that left
	Bill's hands and went into Jack's hands, and then from Jack's to Bills'.
	But when we are doing signalling with data it's only a 1:1 signal.

	Question. What if we want to signal 1:many? What if we needed all bunch of go routines to know about
	an event.
	Answer. Well there's only one piece of data, we could signal individually to each Go routine, but that would
	take some time. So we have the ability to signal without data, we do that by closing a channel
	Idea is that we just "turn the lights off". If we had auditorium of people and we turned the lights off,
	everybody would see the events. The audience wouldn't have anything to take home with them (no data shared)
	but they would see the event.
	This is what closing the channel does. We don't need to close the channel for memory leaks or issues like that.
	It is literally a state change from open to close.
	A channel could be in one of 3 states, A channel could be in
	1. Open state, we do that by making the channel. "make" creates and puts the channel in open state.
		Then the receives and sends can happen in the ways we described above.
	2. Zero value state or nil state. That mean sends and receives are going to be blocked.
		We can use that for short term shortages like rate limiting, and things like that. May be in an event loops.
	3. Closed state - We will use that mainly for cancellations and shut down.


	Signal semantics we will be using going forward.

	1. Does the sender need a guarantee that their signal has been received.
		How does that work?
		The receive happens before the send. Both sender and receiver have to come together at that moment in time.
		When we say that receive happens before the send, it's going to happen nanoseconds before (we are still in multi
			threaded environment).

	If we take the guarantee then we have to deal with the unknown latency. If we don't take a guarantee that
	the signal being sent has been received we can reduce or eliminate that latency (but there's a risk)
	And more the risk more the problems.

	We need to know the state of the channel so that we know how it is going to operate and
	we also want to know do we need to deliver the signal with or without the data.

*/