package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/marni/goigc"
)

func handleJustIgcinfo(w http.ResponseWriter, r *http.Request) {
	status := http.StatusBadRequest
	http.Error(w, http.StatusText(status), status)
}

func handleRubbish(w http.ResponseWriter, r *http.Request) {

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
