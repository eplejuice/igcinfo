package main

import (
	"net/http"
	"time"
)

// Starts to count the time for the uptime function
var startTime = time.Now()

// A struct that cointains the return data of the metaData function
// I chose to make it global for easier accessability
type metaData struct {
	Uptime  string
	Info    string
	Version string
}

// Files is a map connecting a id(int) to a struct containing the Url to the IGC file.
var Files map[int]igcFile

// Stores the url in a struct to easier encode and decode as a json object
type igcFile struct {
	URL string `json:"url"`
}

// A global variable used as the unique ID in the map containing the IGC file structs
var globalCount int

func main() {
	// Creates the map
	Files = make(map[int]igcFile)

	globalCount = 0
	// Sends every request to the router function with Regex.
	http.HandleFunc("/", handleRouter)

	//Listens to the Url given by heroku
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// If the Url is wrong the program shuts down immediately.
		panic(err)
	}
}
