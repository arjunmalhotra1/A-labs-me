/*
	A new feature that cam ein golang probably 1.8 were sub tests.
	Sub tests are a really brilliant mechanisms specially when you are dealing with the data tests.
	We saw data tests earlier where we had different pieces of data in the collection.
	One of the problems with data test is that if we want to isolate just one piece of data on a given test run.
	You have to comment out all the other data.
	Sub tests, lets us treat each piece of data like it's own test without having to write
	extra test functions.

	We have added "name" field ot the tests.
	Now we will give a unique name to each of the data, the input and the expected output.
	Earlier we just had the "url" and the "status code".

	We go to our "Given" in "Given" we range over our table.
	But now for every piece of data we define a literal function.
	"func(t *testing.T) {"
	We are going to define a literal test function for every piece of data.
	Our "When" & "Should" are inside these individual test functions.
	Every piece of data is within it's on function "tf".

	Next we bind the literal test function to the name that was in the table test. We set it up as a sub test.
	"t.Run(test.name, tf)"
	These tests will now run in series.
	WE can run all these tests or can even filter by name.

	"go test -run TestDownload -v" See 1.png
	We can see that it ran both the tests "statusok" and "statusnotfound"
	Now we have this extra name, that's associated with each of these data test.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	Real power of sub test is that we can now run
	"go test -run TestDownload/statusok -v"
	See 2.png.
	We no longer have to comment the other data out.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	Now because we have separate test functions here, Go has also provided runing tests in parallel.
	So, normally when we run a set of unit tests inside of one package they run in series.
	We can run tests in parallel across packages using these sub tests.

	But now we added
	"t.Parallel()"
	This will tell the testing tool to not to wait for this test to finish before starting on the next test.
	Run all the sub tests in parallel.

	// When tests run in series.
	"go run -v TestDownload -v" Total time is 0.640 seconds.

	// when we run these test sin parallel.
	"go run TestParallelize -v" Total time is 0.238 seconds. Because the sub-tests were runing in parallel.




*/

// go test -v
// go test -run TestDownload/statusok -v
// go test -run TestDownload/statusnotfound -v
// go test -run TestParallelize -v

// Sample test to show how to write a basic sub unit table test.
package example2

import (
	"net/http"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

// TestDownload validates the http Get function can download content and
// handles different status conditions properly.
func TestDownload(t *testing.T) {
	tt := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"statusok", "https://www.ardanlabs.com/blog/index.xml", http.StatusOK},
		{"statusnotfound", "http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
	}

	t.Log("Given the need to test downloading different content.")
	{
		for testID, test := range tt {
			testID, test := testID, test
			tf := func(t *testing.T) {
				t.Logf("\tTest %d:\tWhen checking %q for status code %d", testID, test.url, test.statusCode)
				{
					resp, err := http.Get(test.url)
					if err != nil {
						t.Fatalf("\t%s\tTest %d:\tShould be able to make the Get call : %v", failed, testID, err)
					}
					t.Logf("\t%s\tTest %d:\tShould be able to make the Get call.", succeed, testID)

					defer resp.Body.Close()

					if resp.StatusCode == test.statusCode {
						t.Logf("\t%s\tTest %d:\tShould receive a %d status code.", succeed, testID, test.statusCode)
					} else {
						t.Errorf("\t%s\tTest %d:\tShould receive a %d status code : %v", failed, testID, test.statusCode, resp.StatusCode)
					}
				}
			}

			t.Run(test.name, tf)
		}
	}
}

// TestParallelize validates the http Get function can download content and
// handles different status conditions properly but runs the tests in parallel.
func TestParallelize(t *testing.T) {
	type tableTest struct {
		name       string
		url        string
		statusCode int
	}

	tt := []tableTest{
		{"statusok", "https://www.ardanlabs.com/blog/index.xml", http.StatusOK},
		{"statusnotfound", "http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
	}

	t.Log("Given the need to test downloading different content.")
	{
		for testID, test := range tt {
			testID, test := testID, test
			tf := func(t *testing.T) {
				t.Parallel()

				t.Logf("\tTest %d:\tWhen checking %q for status code %d", testID, test.url, test.statusCode)
				{
					resp, err := http.Get(test.url)
					if err != nil {
						t.Fatalf("\t%s\tTest %d:\tShould be able to make the Get call : %v", failed, testID, err)
					}
					t.Logf("\t%s\tTest %d:\tShould be able to make the Get call.", succeed, testID)

					defer resp.Body.Close()

					if resp.StatusCode == test.statusCode {
						t.Logf("\t%s\tTest %d:\tShould receive a %d status code.", succeed, testID, test.statusCode)
					} else {
						t.Errorf("\t%s\tTest %d:\tShould receive a %d status code : %v", failed, testID, test.statusCode, resp.StatusCode)
					}
				}
			}
			t.Run(test.name, tf)
		}
	}
}
