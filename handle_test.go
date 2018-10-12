package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Tests the functuon which returns metadata about the api
func Test_handleMetaData(t *testing.T) {
	// Makes a request for the handler
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	// This repsonse is created in place of a http.ResponseWriter to handle the response
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handleMetaData)

	// Send the request and the response to handler with serveHTTP,
	// which acts as a http.Handler
	handler.ServeHTTP(response, req)

	// Checks if the function returned the right status code
	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Wrong status code: got %v , expected %v", status, http.StatusOK)
	}

	// Expected string from the function
	expected := `{"Uptime":"P0Y0M0DT0H0M0S","Info":"Info about igc api reader","Version":"1"}`

	// Gets the response as a string for easier comparison
	got := response.Body.String()

	// Checks if the response matches the expected
	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}
}

// Tests the funtion which is supposed to return an array with the IDs of all tracks
func Test_handleGet(t *testing.T) {
	// First had to hardcode values for the function test
	Files = make(map[int]igcFile)
	Files[1] = igcFile{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}

	// The rest of the function is very similar to the one above
	req, err := http.NewRequest("GET", "/api/igc", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGet)

	handler.ServeHTTP(response, req)

	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Wrong status code: got %v , expected %v", status, http.StatusOK)
	}

	// Had to add \n because the response body contains a newline
	expected := "[1]\n"

	got := response.Body.String()

	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}

}

// Tests the function which is supposed to post a track,
// and return the ID assigned to the track
func Test_handlePost(t *testing.T) {

	// Hardcoded values, with a legit Url which points to a IGC file.
	// Coverts a struct to json sends it.
	testTmp := igcFile{`http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc`}
	testTmpJSON, err := json.Marshal(testTmp)
	req, err := http.NewRequest("POST", "/api/igc", bytes.NewReader(testTmpJSON))
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePost)

	handler.ServeHTTP(response, req)

	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Wrong status code: got %v , expected %v", status, http.StatusOK)
	}

	// Again i tried to make a json object to compare with the response,
	// But found out its easier to just compare two strings
	got := response.Body.String()
	expected := "{\"id\":1}\n"
	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}

}

//Tests the function which gets the data of a stored track based on gived ID
func Test_handleAPIIgcID(t *testing.T) {

	Files = make(map[int]igcFile)
	Files[1] = igcFile{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}

	// Hardcoded the value 1 to send
	req, err := http.NewRequest("GET", "/api/igc/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	hh := httptest.NewRecorder()

	// Mux is a library which contains a router based on regular expressions
	h := mux.NewRouter()

	// Routes based on a regular expression, so any sequnce of numbers are accepted.
	h.HandleFunc("/api/igc/{id:[0-9]+}", handleAPIIgcID).Methods("GET")

	h.ServeHTTP(hh, req)

	status := hh.Code
	if status != http.StatusOK {
		t.Errorf("Wrong status code: got %v , expected %v", status, http.StatusOK)
	}

	expected := `{"H_date":"2016-02-19T00:00:00Z","pilot":"Miguel Angel Gordillo","glider":"RV8","glider:id":"EC-XLL","track_lenght":443.2573603705269}
`

	got := hh.Body.String()

	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}

}

// Tests the function which gets a single value from the track based on ID and field
func Test_handleAPIIgcIDField(t *testing.T) {

	Files = make(map[int]igcFile)
	Files[1] = igcFile{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}

	req, err := http.NewRequest("GET", "/api/igc/1/pilot", nil)
	if err != nil {
		t.Fatal(err)
	}

	hh := httptest.NewRecorder()

	h := mux.NewRouter()

	h.HandleFunc("/api/igc/{id:[0-9]+}/{pilot}", handleAPIIgcIDField).Methods("GET")

	h.ServeHTTP(hh, req)

	status := hh.Code
	if status != http.StatusOK {
		t.Errorf("Wrong status code: got %v , expected %v", status, http.StatusOK)
	}

	expected := "Miguel Angel Gordillo"

	got := hh.Body.String()

	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}

}
