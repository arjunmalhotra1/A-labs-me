/*
	In the routes function we add a route called "/sendjson"

	"http.HandleFunc("/sendjson", SendJSON)"
	Where anytime we hit the sendjson route, we will call the "SendJSON" function.

	Inside the SendJSON function we are mocking a response using a struct literal.
	Then we set the content type, the header, and then we will send it back down under "Json".

*/
// Package handlers provides the endpoints for the web service.
package handlers

import (
	"encoding/json"
	"net/http"
)

// Routes sets the routes for the web service.
func Routes() {
	http.HandleFunc("/sendjson", SendJSON)
}

// SendJSON returns a simple JSON document.
func SendJSON(rw http.ResponseWriter, r *http.Request) {
	u := struct {
		Name  string
		Email string
	}{
		Name:  "Bill",
		Email: "bill@ardanlabs.com",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(&u)
}
