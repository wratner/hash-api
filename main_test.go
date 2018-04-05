package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/wratner/hash-api/logger"
)

// TestHashHandler makes sure that the hashHandler can be called successfully,
// and returns the expected base64 encode SHA512 hash.
func TestHashHandler(t *testing.T) {
	logger.Init(os.Stdout, os.Stdout)

	app := App{}

	data := url.Values{}
	data.Set("password", "angryMonkey")

	req, err := http.NewRequest("POST", "/hash", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.hashHandler)

	handler.ServeHTTP(rr, req)

	encodedHash, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	actualResult := string(encodedHash)
	expectedResult := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

// TestEmptyPasswordHashHandler ensures that when a client makes a request
// to the /hash endpoint without a password field value set that it is handled
// appropriately and returns the appropriate 400 response code.
func TestEmptyPasswordHashHandler(t *testing.T) {
	logger.Init(os.Stdout, os.Stdout)

	app := App{}

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.hashHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected %d but got %d", http.StatusBadRequest, rr.Code)
	}
}

// TestBadShutdownHandler makes sure that if a user makes a request
// besides GET to the /shutdown endpoint that it is handled appropriately
// and returns a 400 response code.
func TestBadShutdownHandler(t *testing.T) {
	app := App{}

	req, err := http.NewRequest("POST", "/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.shutdownHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected %d but got %d", http.StatusBadGateway, rr.Code)
	}
}

// TestEmptyStatsHandler checks that the /stats endpoint can handle
// being called when no requests to the /hash endpoint have been made.
func TestEmptyStatsHandler(t *testing.T) {
	app := App{}
	statistics := Statistics{}

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.statsHandler)
	handler.ServeHTTP(rr, req)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(bodyBytes, &statistics)
	if err != nil {
		t.Fatal(err)
	}

	expectedTotal := 0
	expectedAverage := float64(0)

	if statistics.Total != expectedTotal && statistics.Average != expectedAverage {
		t.Fatalf("Expected %d total and %f average but got %d and %f", expectedTotal, expectedAverage, statistics.Total, statistics.Average)
	}
}

// TestStatsHandler ensures that the total number of requests and
// average response time is calculated properly and successfully
// returned to the client in the proper JSON format.
func TestStatsHandler(t *testing.T) {
	app := App{Total: 1, ResponseTime: float64(5001)}
	statistics := Statistics{}

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.statsHandler)
	handler.ServeHTTP(rr, req)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(bodyBytes, &statistics)
	if err != nil {
		t.Fatal(err)
	}

	expectedTotal := 1
	expectedAverage := float64(5001)

	if statistics.Total != expectedTotal && statistics.Average != expectedAverage {
		t.Fatalf("Expected %d total and %f average but got %d and %f", expectedTotal, expectedAverage, statistics.Total, statistics.Average)
	}
}
