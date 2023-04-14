/*
	He changed
	"statusCode := 200" to "statusCode := 204"
	and say we are expecting 204 for some reason.
	Now when we run the test even without "-v".
	We will see the verbose mode, see 3.png.' We will see the failure.
	Ans we can see the "x" on front of the test case.

	We can see that, we have to come up with some sought of formatting.
	We don't need third party packages to do this.
	We need something that lets us scan and helps us identify quickly what is going on.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	We can run a specific test if we have multiple,
	command wasn't clear
	"<xxxxx> -run Down"
	Then all the methods with "Down" word will run.

	We can also do that with,
	"<xxxxx> -run Down -v"
	See 4.png.

	In this test we had one single input expecting one single output.


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
	// statusCode := 200
	statusCode := 204

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
