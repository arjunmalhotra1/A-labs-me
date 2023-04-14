/*
	We just saw 2 differnet types of unit tests. A simple unit test and then a data driven tests.
	In both cases the tests were out there hitting the internet to do the Get call.
	It is fine that here we have internet access. But we shouldn't assume that we wil always have the internet access.
	Also it could be dangerous to hit the liver server when we are running the tests.
	Hence we should mock these calls every once in a while, not just to be able to validate success but a lot of times
	we want to validate the failure cases.
	If we are hitting a live server, we can't assume that the request is going to fail. But we know that
	it will may be one day doesn't return something bad comes we want to test all these cases.
	So mocking is a great way to do that.

	Go has httptest package in - "net/http/httptest"
	this will help us do the mocking at the http level.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	We have a raw string,
	'var feed' withing `` this tell the compiler to include the carriage returns and the other spaces.
	With this raw string we are basically documenting what an rss feed looks like.
	Header section of rss feed document this is going to be the mock that we will use to pretend
	that we are hitting the Ardan labs blog. Idea is that instead of hitting the Ardan labs rss feed document what we are
	doing is simulate as if we actually did and in a good case scenario.
	This string "feed" is the document that we will return form our Mock.

	Next, we do some marshalling.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	func mockServer() *httptest.Server {
		f := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/xml")
			fmt.Fprintln(w, feed)
		}

		return httptest.NewServer(http.HandlerFunc(f))
	}

	Mock server will help us set up a server that we can run our mocks against.
	So we are not hitting the live server we are hitting this server.
	"mockServer()" returns a pointer to the "*httptest.Server"
	We are actually going to run a server on localhost on some port.
	and what we will do is intercept those get calls nad then return the mock data back out.
	"f := func(w http.ResponseWriter, r *http.Request) {"
	We define a literal function and that function will be a web handler in Go.
	Takes in a ResponseWriter and pointer to a request.

	Next we set the writeHeader to 200 and content-type to xml.
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/xml")
	Next we send the feed data out,
			fmt.Fprintln(w, feed)

	Finally we take this function and bind it into our new mock server
			"return httptest.NewServer(http.HandlerFunc(f))"
	What this wil do is tell the mock server that anytime a url of any description hits the mock server you
	execute this mock function 'f' and we will return the feed document.
				"fmt.Fprintln(w, feed)"

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	We turn our mock server on,
		"server := mockServer()"
	Next we defer the close the server.
		"defer server.Close()"
	Notice we are doing the same exact Get call,
	"resp, err := http.Get(server.URL)"
	But this time we are using the URL that the server is going to provide.
	Since the server is the only one that knows what port it started up on.

	Now we will have to do one more thing, unmarshall the response body out.
		"xml.NewDecoder(resp.Body).Decode(&d)"

	And then validate that we have at least one item.
		"if len(d.Channel.Items) == 1 {"

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	Finally we run it.
	"Go test -v"
	See 1.png for the output.

	If you notice that when we ran the previous test we were running it for about 0.5 seconds
	to go back to out internet and come back.
	This one is running in 16 micro seconds or 0.016s.
	So mocking has increased the speed of the testing up,
	And we can see the port number on which our mock server stood up on.
	When we run it again we can see that we are now running on a different port.

	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	We have to be careful with the mock,
	"A unit test is a test of behavior whose success or failure wholly determined by the correctness of
	the test and the correctness of unit under test"
	When we are mocking something, we are still mockign it. It's not like hitting the live system.
	Mocks are valueable specially for validating the error cases that the live system probably aren't producing.


*/

// Sample test to show how to mock an HTTP GET call internally.
package example3

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

// feed is mocking the XML document we expect to receive.
var feed = `<?xml version="1.0" encoding="UTF-8"?>
<rss>
<channel>
    <title>Going Go Programming</title>
    <description>Golang : https://github.com/goinggo</description>
    <link>http://www.goinggo.net/</link>
    <item>
        <pubDate>Sun, 15 Mar 2015 15:04:00 +0000</pubDate>
        <title>Object Oriented Programming Mechanics</title>
        <description>Go is an object oriented language.</description>
        <link>http://www.goinggo.net/2015/03/object-oriented</link>
    </item>
</channel>
</rss>`

// Item defines the fields associated with the item tag in
// the buoy RSS document.
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
}

// Channel defines the fields associated with the channel tag in
// the buoy RSS document.
type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Items       []Item   `xml:"item"`
}

// Document defines the fields associated with the buoy RSS document.
type Document struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
	URI     string
}

// mockServer returns a pointer to a server to handle the mock get call.
func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, feed)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

// TestDownload validates the http Get function can download content and
// the content can be unmarshaled and clean.
func TestDownload(t *testing.T) {
	statusCode := http.StatusOK

	server := mockServer()
	defer server.Close()

	t.Log("Given the need to test downloading content.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen checking %q for status code %d", testID, server.URL, statusCode)
		{
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to make the Get call : %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to make the Get call.", succeed, testID)

			defer resp.Body.Close()

			if resp.StatusCode != statusCode {
				t.Fatalf("\t%s\tTest %d:\tShould receive a %d status code : %v", failed, testID, statusCode, resp.StatusCode)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a %d status code.", succeed, testID, statusCode)

			var d Document
			if err := xml.NewDecoder(resp.Body).Decode(&d); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response.", succeed, testID)

			if len(d.Channel.Items) == 1 {
				t.Logf("\t%s\tTest %d:\tShould have 1 item in the feed.", succeed, testID)
			} else {
				t.Errorf("\t%s\tTest %d:\tShould have 1 item in the feed : %d", failed, testID, len(d.Channel.Items))
			}
		}
	}
}
