/*
	A lot of times we would be building endpoints and webservices, web APIs in Go.
	When doing that we would want to test those endpoints as well.
	http package gives us the ability to test our endpoints in our unit tests without having to stand up a server.

	To see this, we have set up a basic webserver like a web application.
	We will use our "ListenAndServe" and we will listen on port 4000
	To bind some route we decide to create a separate package called handlers.
	Look at handlers.go inside the handlers folder.
	xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	When we build the main and we run it, we are listening on port 4000.
	In the browser we then type "localhost:4000/sendjson"
	And we see the output
	See 1.png.
	Now we want to test the route without having to start the server and then
	open a browser and then typing the url.

	We will add "handler_test.go" file and we will add a unit test.

	Note: Whenever you are testing these routes it is very important to bind these routes to our mux.
*/

// Sample program that implements a simple web service.
package main

import (
	"log"
	"net/http"

	"github.com/ardanlabs/gotraining/topics/go/testing/tests/example4/handlers"
)

func main() {
	handlers.Routes()

	log.Println("listener : Started : Listening on: http://localhost:4000")
	http.ListenAndServe(":4000", nil)
}
