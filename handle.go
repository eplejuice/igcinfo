package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/marni/goigc"
)

//This is the router which takes all the http requests and sends them further to the right handleFunc based on the matching Regular expression
func handleRouter(w http.ResponseWriter, r *http.Request) {
	// This handles the GET /api
	regHandleAPI, err := regexp.Compile("^/igcinfo/api/?$")
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}

	// This handles the POST/GET /api/igc
	regHandleAPIIgc, err := regexp.Compile("^/igcinfo/api/igc/?$")
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}

	// This handles the	GET /api/igc/id
	regHandleAPIIgcID, err := regexp.Compile("^/igcinfo/api/igc/[0-9]+/?$")
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}

	// This handles the GET /api/igc/id/field
	regHandleAPIIgcIDField, err := regexp.Compile("^/igcinfo/api/igc/[0-9]+/(pilot|glider|glider_id|track_lenght|H_date)$")
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}

	// This is a switch that always runs routes the http request to the right handlefunc
	// Otherwise the dafault gives the user a httpBadRequest response
	switch {
	case regHandleAPI.MatchString(r.URL.Path):
		handleMetaData(w, r)
	case regHandleAPIIgc.MatchString(r.URL.Path):
		handleAPIIgc(w, r)
	case regHandleAPIIgcID.MatchString(r.URL.Path):
		handleAPIIgcID(w, r)
	case regHandleAPIIgcIDField.MatchString(r.URL.Path):
		handleAPIIgcIDField(w, r)
	default:
		handleError(w, r, nil, http.StatusBadRequest)
	}
}

// This function handles all the errors and writes them as a reponse to the user
// with the right error code based on the parameter recieved
func handleError(w http.ResponseWriter, r *http.Request, err error, status int) {

	http.Error(w, fmt.Sprintf("%s/t%s", http.StatusText(status), err), status)
}

// This is the function which gives the user information about the api
func handleMetaData(w http.ResponseWriter, r *http.Request) {
	// Using a struct to easily encode to a json
	metaInfo := metaData{
		Uptime:  calcTime(startTime),
		Info:    "Info about igc api reader",
		Version: "1",
	}

	// Using Marshal instead of Endoce, because i believe Marshal is used to encode strings
	// and the struct mainly consist of strings.
	metaResp, _ := json.Marshal(metaInfo)
	// Sets the header to json, and returns a json object as the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(metaResp)
}

// This fuction checks wether the user send a POST or a GET request
// and routes to the right handleFunc, otherwise returns a httpBadRequest
func handleAPIIgc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	default:
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
	}
}

// This function handles the GET, when a user wants to get a array containing all IDs of all
// existing IGC files currently stored in the in-memory database.
func handleGet(w http.ResponseWriter, r *http.Request) {
	// Initializes the struct containint the array for the return object
	type rID struct {
		ReturnID []int `json:"id"`
	}
	w.Header().Set("Content-Type", "application/json")
	returnObj := rID{make([]int, 0)}
	// If the in-memory database is empty, returns a empty array
	if len(Files) == 0 {
		err := json.NewEncoder(w).Encode(returnObj.ReturnID)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}
	} else {
		// Loops through the map/database
		for i := range Files {
			// Appends the IDs ( [int] field of the map) to the array.
			returnObj.ReturnID = append(returnObj.ReturnID, i)
		}
		// Encodes the array to json and sends it as Response
		err := json.NewEncoder(w).Encode(returnObj.ReturnID)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}
	}

}

// This function handles the POST request when a user wants to post a new Url to an IGC file
func handlePost(w http.ResponseWriter, r *http.Request) {
	// Gets the Url stored in the body
	postAPI := r.Body
	defer r.Body.Close()
	// Creates a new struct, and decodes the json into it
	var tmp igcFile
	err := json.NewDecoder(postAPI).Decode(&tmp)
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}
	// Checks to see if the given Url actually is real, using the marni/goigc library
	s := tmp.URL
	_, err = igc.ParseLocation(s)
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}
	// Maps the encoded struct to the map, with a unique ID from the global variable
	// this variable can assign the same ID twice, if two POST requests are sent exactly
	// at the same time, however this seems very unlikely.
	globalCount++
	Files[globalCount] = tmp

	// Creates a json object containing the ID given to the Url as a Response to the user.
	type ReturnVal struct {
		ID int `json:"id"`
	}

	RV := ReturnVal{globalCount}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(RV)
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// This funtion handles the GET ID call, which is supposed to return the data of an Igc file
// based on the given ID.
func handleAPIIgcID(w http.ResponseWriter, r *http.Request) {
	// Base lets us get the last value of the Url, which in this case is the ID
	tmp := path.Base(r.URL.Path)
	// The value is given as a string, so it has to be converted to an int before use.
	tmpInt, err := strconv.Atoi(tmp)
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}
	// Checks if the [int] field of the map contains the given ID,
	// if yes, elem returns as the struct connected to the int matching the given ID in the map
	elem, ok := Files[tmpInt]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		// Here we use the marni/goigc library functions to get the data of the Igc file
		// based on a given Url.
		s := elem.URL
		track, err := igc.ParseLocation(s)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}

		// Here a struct is made to contains the values the Igcfile/track returns
		// to lates be converted to a json for the response
		type returnValues struct {
			HDate       time.Time `json:"H_date"`
			Pilot       string    `json:"pilot"`
			Glider      string    `json:"glider"`
			GliderID    string    `json:"glider:id"`
			TrackLenght float64   `json:"track_lenght"`
		}
		// The values is put into the struct
		returnObject := returnValues{
			HDate:       track.Header.Date,
			Pilot:       track.Pilot,
			Glider:      track.GliderType,
			GliderID:    track.GliderID,
			TrackLenght: getTrackLenght(track),
		}
		// Return the struct as a json with the right values (hopefully).
		err = json.NewEncoder(w).Encode(returnObject)
	}
}

// This functions handles the GET /api/igc/id/field call, which is supposed to returns a single
// value of a Igc file/ track, based on a given ID and the value to be returned.
func handleAPIIgcIDField(w http.ResponseWriter, r *http.Request) {
	// First use Base to get the last value of the Url which is field.
	field := path.Base(r.URL.Path)
	// Dir returns everything in the URL, but the last value.
	// this way we can use Base again, but this time ID is the last value of the Url
	tmp := path.Dir(r.URL.Path)
	nummer := path.Base(tmp)

	index, err := strconv.Atoi(nummer)
	if err != nil {
		handleError(w, r, err, http.StatusBadRequest)
		return
	}

	// Checks if the databse contains the requested track.
	elem, ok := Files[index]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
	} else {
		// If the track is found use the marni/goigc library functions to get values from it
		s := elem.URL
		track, err := igc.ParseLocation(s)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}
		// A switch which checks what field the request is asking for
		// and returns the right value by Writing it as a plain text
		switch field {
		case "H_date":
			text, err := track.Header.Date.MarshalText()
			if err != nil {
				handleError(w, r, err, http.StatusBadRequest)
				return
			}
			w.Write(text)
		case "pilot":
			w.Write([]byte(track.Pilot))
		case "glider":
			w.Write([]byte(track.GliderType))
		case "glider_id":
			w.Write([]byte(track.GliderID))
		case "track_lenght":
			w.Write([]byte(strconv.Itoa(int(getTrackLenght(track)))))
		}
	}
}

// This is a funtion to calculate the track_lenght variable based on a set of given coordinates
func getTrackLenght(s igc.Track) float64 {
	totalDistance := 0.0
	// Loops through all given coordinates and adds to the total distance variable
	for i := 0; i < len(s.Points)-1; i++ {
		totalDistance += s.Points[i].Distance(s.Points[i+1])
	}
	return totalDistance
}
