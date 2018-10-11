package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleMetaData(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handleMetaData)

	handler.ServeHTTP(response, req)

	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Wrong status code: got %v , expected %v", status, http.StatusOK)
	}

	expected := `{"Uptime":"P0Y0M0DT0H0M0S","Info":"Info about igc api reader","Version":"1"}`

	got := response.Body.String()

	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}
}

func Test_handleGet(t *testing.T) {
	Files = make(map[int]igcFile)
	Files[1] = igcFile{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}

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

	expected := "[1]\n"

	got := response.Body.String()

	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}

}

func Test_handlePost(t *testing.T) {

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

	/*
		type returnTmp struct {
			id string `json:"id"`
		}

		rtrnTmp := returnTmp{`"id":1`}
		expect, err := json.Marshal(rtrnTmp)
		expected := bytes.NewReader(expect)
	*/

	got := response.Body.String()
	expected := "{\"id\":1}\n"
	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}

}

func Test_handleAPIIgcId(t *testing.T) {

	Files = make(map[int]igcFile)
	Files[1] = igcFile{"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}

	req, err := http.NewRequest("GET", "/api/igc/id", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handleAPIIgcID("/", "api/igc/1"))

	handler.ServeHTTP(response, req)

	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Wrong status code: got %v , expected %v", status, http.StatusOK)
	}

	expected := "[1]\n"

	got := response.Body.String()

	if expected != got {
		t.Errorf("Test failed: got %v, wanted %v", got, expected)
	}

}
