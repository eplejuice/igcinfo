package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var startTime = time.Now()

type metaData struct {
	Uptime  string
	Info    string
	Version string
}

var Files map[int]igcFile

type igcFile struct {
	Url string `json: "url"`
}

var globalCount int

func main() {

	Files = make(map[int]igcFile)

	globalCount = 0
	http.HandleFunc("/igcinfo", handleJustIgcinfo)
	http.HandleFunc("/", handleRubbish)
	http.HandleFunc("/igcinfo/api/", handleApiSlashRubbish)
	http.HandleFunc("/igcinfo/api", handleMetaData)
	http.HandleFunc("/igcinfo/api/igc", handlePost)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
		fmt.Println("Could not handle url")
	}
}

func calcTime() string {

	year := 31536000
	month := 2592000
	week := 604800
	day := 86400
	hour := 3600
	minute := 60

	years := 0
	months := 0
	weeks := 0
	days := 0
	hours := 0
	minutes := 0

	uptime := time.Since(startTime)
	uptimee := uptime.String()
	tmp := strings.Split(uptimee, ".")
	tmp2 := tmp[0]
	uptimeS, _ := strconv.Atoi(tmp2)

	fmt.Println(uptimeS)
	fmt.Println(year)
	fmt.Println(uptimeS % year)

	years = uptimeS % year

	months = uptimeS % month

	weeks = uptimeS % week

	days = uptimeS % day

	hours = uptimeS % hour

	minutes = uptimeS % minute

	fmt.Println(years)
	fmt.Println(months)
	fmt.Println(weeks)
	fmt.Println(days)
	fmt.Println(hours)
	fmt.Println(minutes)
	fmt.Println(uptimeS)

	returnVal := strings.Join([]string{"P", strconv.Itoa(years), "Y", strconv.Itoa(months), "M", strconv.Itoa(weeks), "W", strconv.Itoa(days), "D", "T", strconv.Itoa(hours), "H", strconv.Itoa(minutes), "M", strconv.Itoa(uptimeS), "S"}, "")
	return returnVal
}
