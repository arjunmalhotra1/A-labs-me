/*
	Golang comes with built in testing tool that we can call from the command line called "go test".
	There are a lot fo ways we can write unit testing in Go.
	Bill does not recoomend bringing in thrid party test packages to do testing
	primarily because, it becomes a dependency for other people to downloadn and run your tests.

	Bill recommends that as a team we should have some level of consistency on how everybody write tests.

	"Given When Should"
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	We will do some basic unit testing around the http package and using the testing package.

	First core thing about a test is that it has to be in a file that is named
	"_test.go"
	Go is about convention over configuration.
	And the convention of the file name is what is driving the testing tool to know what a test file is.
	We also have to have functions that start with the word "Test" and
	these functions have to be exported.

	Testing tool will be scanning this file for the functions that start with "Test" and also these
	Testing functions will also take the testing 'T' pointer as a parameter that's going to
	indicate that this is a test function.
	"func TestDownload(t *testing.T) {"

	Note:
	We can use '_' after 'Test' in the function name which Bill said he is not a fan off because
	we don't really do that in Go. "Test_Download"
	OR
	the next letter has to be Capitalized like - "TestDownload"
	If we use "Testdownload" the testing tool will not find this function.

	Here we are testing our "TestDownload" this will help us test the ability to download some content
	off the internet. We are trying to download some blog content.
	We are also expecting to see the status code 200.

	It's good to sometimes set your intput and output on top of the unit test. That's what
	we try to do here.

	"url := "https://www.ardanlabs.com/blog/index.xml"
	statusCode := 200"

	"Given when should"
	So bil is here doing some special formatting here. The "t" variable has
	an API of log, error and fatal.
	Log is for verbose output anything we want to comment or trace about in test.
	The "Error" function is going to say that this function ahs failed but we will continue to execute
	the code inside the function and
	"Fatal" means this test has failed we are done and move on to another function.
	We will use all the 3 APis in this code.

	This is our given comment. Given comment here is why does this test exist.
	"t.Log("Given the need to test downloading content.")"

	"Given the need to test downloading content."

	We also use extra code blocks to separate the "Given", "Whens" and the "Should"
	Just so we can visually read all the three.

	Next we have "When",
	"t.Logf("\tTest %d:\tWhen checking %q for status code %d", testID, url, statusCode)"

	"When" is showing us the data we are using. We are saying here,
	"When we are cheking this url for this particular status code."
	We can clearly look at the when very clearly at the output and know what our input was and
	what my potential output is. Sometimes we can just show our input as well.
	We are marking this "Test 0:" as a part of this unit test called "Download".

	Inside the code block for "When" we are now showing the "Should" aspect.
	"Should" is all the things we are testing and what should occur with that test.

	First thing we are doing is the
	"resp, err := http.Get(url)" get call, on the URL and then we are cheking the error.
	It is very very important that we check errors all the way through our unit tests
	as if this was production code.

	If there is an error, we are now in a fatal situation, there is no reason to execute
	anything else if we can't even get that URL.
	So we will "Fatalf". Fatal is also like an exist so there's no return when we call
	t.Fatalf()
	Or if not an error we will go into the Log.

	Not the Fatal case and the Error case are almost identical, in terms of the output.

	But Bill says he likes the idea of being able to scan the output regardless of the case
	whether there is failure or success.
	So the more consistency there is with small marker changes, it will help you scan the output
	to know what's going on.

	Finally we check the status code and log if status is 200 or not.

	We don't need an else clause on fatal because, that's an automatic exit out.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	Next we run this test. We can do just do "go test" and it will find all the
	"_test" files and it will run them.
	When he did that, we got 1.png output.
	It just said "Pass" and we saw the time taken.

	If we want to see all the verbose output, we do
	"go test -v" and now, we get all the verbose output that we were putting together.
	This is why we should write tests like this because it becomes very clear here
	On the given,
		All the different tests (Test 0) we are executing, may be all the different inputs/
		outputs we are using.
	And then the idea that everything has succeeded like we can see with those check marks.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	then he made some edits and forced the test to fail. See example_2.go.












*/

// Sample test to show how to write a basic unit test.
package example1

import (
	"net/http"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

// TestDownload validates the http Get function can download content.
func TestDownload(t *testing.T) {
	url := "https://www.ardanlabs.com/blog/index.xml"
	statusCode := 200

	t.Log("Given the need to test downloading content.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen checking %q for status code %d", testID, url, statusCode)
		{
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to make the Get call : %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to make the Get call.", succeed, testID)

			defer resp.Body.Close()

			if resp.StatusCode == statusCode {
				t.Logf("\t%s\tTest %d:\tShould receive a %d status code.", succeed, testID, statusCode)
			} else {
				t.Errorf("\t%s\tTest %d:\tShould receive a %d status code : %d", failed, testID, statusCode, resp.StatusCode)
			}
		}
	}
}
