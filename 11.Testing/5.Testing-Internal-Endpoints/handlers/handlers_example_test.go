/*
	Here we can see that Go has another type of test that we can run, called
	example tests.

	What's  nice is not only this becomes an example inside Go doc.
	It also now, acts as a test based on the standard out.

	If we want this example to not just be a test but also to show up in go documentation,
	then we define the test function as
	"ExampleSendJSON"
	Also not that this function doesn't take the testing "t".
	After the word "Example" what we would want to do is make sure that we are binding it to something
	that is truly exported from the API.
	Like we have "SendJSON" here exported from "handlers.go". If we don't bind it to
	something that is truly exported from the API then it will not show up in the Go documentation.

	In this example test, we will use our request and recorder again and run everything through our mux.

	r := httptest.NewRequest("GET", "/sendjson", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)

	But this time after we have run through the mux,

	"if err := json.NewDecoder(w.Body).Decode(&u); err != nil {"

	we will take the response body, the json that we expected to get back, and we will decode it
	and then we will write it to standard out.

	Example tests are about testing what's coming out of standard out.

	This comment in the example test
	// Output:
	// {Bill bill@ardanlabs.com}
	is saying this "{Bill bill@ardanlabs.com}" is the expected output that we expect to get when this
	example finished running.

	We run
	"go text ExampleSendJSON -v"
	see 3.png as output.

	If we change the output, remove the '}' closing bracket.
	// Output:
	// {Bill bill@ardanlabs.com

	and run again we set 4.png.
	It says now that the test failed.

	It says we got "{Bill bill@ardanlabs.com}" output but we were expecting "{Bill bill@ardanlabs.com" output.



*/

// go test -run ExampleSendJSON

// Sample to show how to write a basic example.
package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

// ExampleSendJSON provides a basic example example.
func ExampleSendJSON() {
	r := httptest.NewRequest("GET", "/sendjson", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)

	var u struct {
		Name  string
		Email string
	}

	if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
		log.Println("ERROR:", err)
	}

	fmt.Println(u)
	// Output:
	// {Bill bill@ardanlabs.com}
}
