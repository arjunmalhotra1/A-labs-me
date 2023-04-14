/*
	Hidden gem that Golang has in testing tool is the code coverage.
	Bill tries to get to 70-80% code coverage with tests.

	He wrote a tcp package to help support writing any sought of TCP portocol.
	What we want is make sure that the test files are giving us enough code coverage to feel comfortable,
	that we have enough unit testing.

	"go test -cover"
	This will run all the tests and then gives us the total number of coverage.
	Here we have 58.3% which is short of 70%.
	Question now becomes what is being covered and what isn't being covered.
	What we can do here is get "cover profile"

	We ask the coverprofile to write the output to some file like "c.out".
	"go test -coverprofile c.out"

	How to view the c.out file?
	Go have a cover tool.
	We just ned to tell the cover tool that we want to view that profile in html mode.

	"go tool cover -html=c.out"
	This will open our browser and at a file level tell us where we at with test coverage with each file.
	We can see that the code is red in color. Which means coverage missing.
	Green color mean that we are testing that code.

	this is what we can do while validating PRS and doing code reviews. We can then run those cover profiles,
	and get a nice visual look at what our cover profile is.



*/