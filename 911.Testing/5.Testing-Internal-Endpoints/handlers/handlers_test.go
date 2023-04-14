/*
	VVIMP
	Note: Whenever you are testing these routes it is very important to bind these routes to our mux.

	It is one of the reasons why in this handlers_test example Bill has taken the routes and
	put them into separate package
	"handlers.Routes()"
	We call this function within "init()" and make sure that the "SendJson" route is bounded into this test
	before we test it.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	"TestSendJSON()" function takes "t *testing.T"
	url - "/sendjson" and we
	expect status code 200.

	Within the httptest package we have "NewRequest" function.
	This "NewRequest" function actually returns a concrete pointer to the request value that we need to do
	any sought of request.
	"httptest.NewRequest("GET", url, nil)"
	"GET" call on the "url" and no post data.

	r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, r)

	Now next we need a concrete type that knows how to be a recorder.
	"httptest.NewRecorder()" returns a pointer to the response recorder type.
	Which actually implements a response writer interface.
	Now that we have a request and the recorder that implements the response writer interface.
	We can directly pass these values into the mux itself.
	"serveHTTP" is how we process a request. So we are by passing the network we are telling the mux
	take this request and write to this recorder for the response writer and then we can look at the recorder
	to see how well we did processing this request from the mux down.

	We then check the recorder to see if we got 200 and then we can also decode and then validate that we
	have the right name & email.

	To run this we do in the handlers folder and write
	"go test -run SendJson -v"
	See 2.png

	Next we take a look at "handlers_example_test.go"




*/

// go test -run TestSendJSON -race -cpu 16

// Sample test to show how to test the execution of an internal endpoint.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/gotraining/topics/go/testing/tests/example4/handlers"
)

const succeed = "\u2713"
const failed = "\u2717"

func init() {
	handlers.Routes()
}

// TestSendJSON testing the sendjson internal endpoint.
func TestSendJSON(t *testing.T) {
	url := "/sendjson"
	statusCode := 200

	t.Log("Given the need to test the SendJSON endpoint.")
	{
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, r)

		testID := 0
		t.Logf("\tTest %d:\tWhen checking %q for status code %d", testID, url, statusCode)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of %d for the response. Received[%d].", failed, testID, statusCode, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of %d for the response.", succeed, testID, statusCode)

			var u struct {
				Name  string
				Email string
			}

			if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to decode the response.", failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to decode the response.", succeed, testID)

			if u.Name == "Bill" {
				t.Logf("\t%s\tTest %d:\tShould have \"Bill\" for Name in the response.", succeed, testID)
			} else {
				t.Errorf("\t%s\tTest %d:\tShould have \"Bill\" for Name in the response : %q", failed, testID, u.Name)
			}

			if u.Email == "bill@ardanlabs.com" {
				t.Logf("\t%s\tTest %d:\tShould have \"bill@ardanlabs.com\" for Email in the response.", succeed, testID)
			} else {
				t.Errorf("\t%s\tTest %d:\tShould have \"bill@ardanlabs.com\" for Email in the response : %q", failed, testID, u.Email)
			}
		}
	}
}
