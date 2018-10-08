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

func handleRouter(w http.ResponseWriter, r *http.Request) {
	regHandleApi, err := regexp.Compile("^/igcinfo/api/$")
	if err != nil {
		panic(err)
	}

	regHandleApiIgc, err := regexp.Compile("^/igcinfo/api/igc/$")
	if err != nil {
		panic(err)
	}

	regHandleApiIgcId, err := regexp.Compile("^/igcinfo/api/igc/[0-9]+$")
	if err != nil {
		panic(err)
	}

	regHandleApiIgcIdField, err := regexp.Compile("^/igcinfo/api/igc/[0-9]+/(pilot|glider|glider_id|track_lenght|H_date)$")
	if err != nil {
		panic(err)
	}

	switch {
	case regHandleApi.MatchString(r.URL.Path):
		fmt.Println("/api")
		handleMetaData(w, r)
	case regHandleApiIgc.MatchString(r.URL.Path):
		fmt.Println("/api/igc")
		handleApiIgc(w, r)
	case regHandleApiIgcId.MatchString(r.URL.Path):
		handleApiIgcId(w, r)
		fmt.Println("/api/igc/id")
	case regHandleApiIgcIdField.MatchString(r.URL.Path):
		handleApiIgcIdField(w, r)
		fmt.Println("/api/igc/id/field")
	default:
		handleError(w, r)
		fmt.Println("Default")
	}
}

func handleError(w http.ResponseWriter, r *http.Request) {

	status := http.StatusBadRequest
	http.Error(w, http.StatusText(status), status)
}

func handleApiSlashRubbish(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest
	http.Error(w, http.StatusText(status), status)
}

func handleMetaData(w http.ResponseWriter, r *http.Request) {
	metaInfo := metaData{
		Uptime:  calcTime(),
		Info:    "Info about igc api reader",
		Version: "1",
	}

	metaResp, _ := json.Marshal(metaInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(metaResp)
}

func handleApiIgc(w http.ResponseWriter, r *http.Request) {
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
func handleGet(w http.ResponseWriter, r *http.Request) {
	type rId struct {
		returnId []int `json:"id"`
	}
	w.Header().Set("Content-Type", "application/json")
	returnObj := rId{make([]int, 0)}
	if len(Files) == 0 {
		err := json.NewEncoder(w).Encode(returnObj.returnId)
		if err != nil {
			panic(err)
		}
	} else {
		for i := range Files {
			returnObj.returnId = append(returnObj.returnId, i)
		}
		err := json.NewEncoder(w).Encode(returnObj.returnId)
		if err != nil {
			panic(err)
		}
	}

}

func handlePost(w http.ResponseWriter, r *http.Request) {
	postApi := r.Body
	defer r.Body.Close()
	var tmp igcFile
	err := json.NewDecoder(postApi).Decode(&tmp)
	if err != nil {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}
	s := tmp.Url
	_, err = igc.ParseLocation(s)
	if err != nil {
		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		fmt.Errorf("Body does not contain proper URL", err)
	}

	globalCount++
	Files[globalCount] = tmp
	fmt.Println(Files)

	type ReturnVal struct {
		Id int `json: "id"`
	}

	RV := ReturnVal{globalCount}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(RV)
	if err != nil {
		fmt.Println("Error with reponse in api/POST")
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func handleApiIgcId(w http.ResponseWriter, r *http.Request) {
	tmp := path.Base(r.URL.Path)
	tmpInt, err := strconv.Atoi(tmp)
	if err != nil {
		panic(err)
	}
	elem, ok := Files[tmpInt]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		s := elem.Url
		track, err := igc.ParseLocation(s)
		if err != nil {
			handleError(w, r)
		}

		type returnValues struct {
			H_date       time.Time `json:"H_date"`
			Pilot        string    `json:"pilot"`
			Glider       string    `json:"glider"`
			Glider_id    string    `json:"glider:id"`
			Track_lenght float64   `json:"track_lenght"`
		}

		returnObject := returnValues{
			H_date:       track.Header.Date,
			Pilot:        track.Pilot,
			Glider:       track.GliderType,
			Glider_id:    track.GliderID,
			Track_lenght: getTrackLenght(track),
		}
		err = json.NewEncoder(w).Encode(returnObject)
	}
}

func handleApiIgcIdField(w http.ResponseWriter, r *http.Request) {
	field := path.Base(r.URL.Path)
	tmp := path.Dir(r.URL.Path)
	nummer := path.Base(tmp)

	fmt.Println(field)
	fmt.Println(nummer)
	index, err := strconv.Atoi(nummer)
	if err != nil {
		panic(err)
	}

	elem, ok := Files[index]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
	} else {
		s := elem.Url
		track, err := igc.ParseLocation(s)
		if err != nil {
			panic(err)
		}
		switch field {
		case "H_date":
			text, err := track.Header.Date.MarshalText()
			if err != nil {
				panic(err)
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

func getTrackLenght(s igc.Track) float64 {
	totalDistance := 0.0
	for i := 0; i < len(s.Points)-1; i++ {
		totalDistance += s.Points[i].Distance(s.Points[i+1])
	}
	return totalDistance
}
