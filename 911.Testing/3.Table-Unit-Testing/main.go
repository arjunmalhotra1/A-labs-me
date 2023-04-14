/*
	Table test means that with just one test function we can, go ahead and test separate inputs and outputs
	all in one place. There could be times where the test logic doesn't have to change the data.

	Same example as the last one but now we construct our table,

	tt := []struct {
		url        string
		statusCode int
	}{
		{"https://www.ardanlabs.com/blog/index.xml", http.StatusOK},
		{"http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
	}

	A lot of times a table is just a slice of the struct values and this slice must not be a named type
	and could be literal type like here above.
	In this table we define both the input and the output.

	What's great is that we can come to this table and add more/different URLs for maybe all of the different status
	codes that we care about the 500s, 204s etc.
	Whatever it is that we want to test, we don't need a whole new separate function, we just add data to the table.

	Next, we go back to our given when should.
	"Given" "When" what should happen and what really happens.
	After "given" we go into "when" and before when we have a for range loop for our table.
	We use the variable 'i' for the test number.We use 'tt' to get the input and the expected output.

	We run this
	"go test -v"
	we will see the two tests that ran in the scope of this complete unit test.

	Now if in future we want to handle a new case we don't have to go and write a new test function,
	we will just add another test case in the table.



*/

// Sample test to show how to write a basic unit table test.
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
		url        string
		statusCode int
	}{
		{"https://www.ardanlabs.com/blog/index.xml", http.StatusOK},
		{"http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
	}

	t.Log("Given the need to test downloading different content.")
	{
		for testID, test := range tt {
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
	}
}
