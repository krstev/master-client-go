package service

import (
	"encoding/json"
	"fmt"
	"master-client-go/src/model"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

const serviceAppID = "service"
const minikubeService = "minikube-service"
const timeoutMilliseconds = 300

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	status := make(map[string]interface{})
	status["status"] = "OK"
	if err := json.NewEncoder(w).Encode(status); err != nil {
		panic(err)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	health := make(map[string]interface{})
	health["health"] = "OK"
	if err := json.NewEncoder(w).Encode(health); err != nil {
		panic(err)
	}
}

func v1PrimeNumbers(w http.ResponseWriter, r *http.Request) {
	PrimeNumbers(w, r, getV1InstanceUrls())
}

func v2PrimeNumbers(w http.ResponseWriter, r *http.Request) {
	PrimeNumbers(w, r, getV2InstanceUrls())
}

func v1CountVowels(w http.ResponseWriter, r *http.Request) {
	CountVowels(w,r,getV1InstanceUrls())
}

func v2CountVowels(w http.ResponseWriter, r *http.Request) {
	CountVowels(w,r,getV2InstanceUrls())
}

func GoogleSearch(w http.ResponseWriter, r *http.Request) {
	startTime := logStart()
	c := make(chan string, 3)
	var request model.GoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		panic(err)
	}
	go searchWeb(c, request.Search, v1ApiVersion)
	go searchImage(c, request.Search, v1ApiVersion)
	go searchVideo(c, request.Search, v1ApiVersion)

	var results string

	for i := 0; i < 3; i++ {
		results += <-c
	}
	if err := json.NewEncoder(w).Encode(results); err != nil {
		panic(err)
	}
	logEnd(startTime)
}

func googleQueryTimeout(w http.ResponseWriter, r *http.Request) {
	startTime := logStart()
	c := make(chan string, 3)
	var request model.GoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		panic(err)
	}
	go searchWeb(c, request.Search, v1ApiVersion)
	go searchImage(c, request.Search, v1ApiVersion)
	go searchVideo(c, request.Search, v1ApiVersion)

	timeout := time.After(timeoutMilliseconds * time.Millisecond)

	var results string

	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results += result
		case <-timeout:
			results = "Timeout"
			return
		}
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		panic(err)
	}
	logEnd(startTime)
}

func googleQueryTimeoutReplica(w http.ResponseWriter, r *http.Request) {
	startTime := logStart()
	c := make(chan string, 3)
	var request model.GoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		panic(err)
	}
	go searchWebReplica(c, request.Search, 2)
	go searchImageReplica(c, request.Search, 2)
	go searchVideoReplica(c, request.Search, 2)

	timeout := time.After(timeoutMilliseconds * time.Millisecond)

	var results string

	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results += result
		case <-timeout:
			results = "Timeout"
		}
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		panic(err)
	}
	logEnd(startTime)
}

