package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/wratner/hash-api/logger"
	"github.com/wratner/hash-api/utils"
)

// Statistics is the return format for the /stats endpoint.
type Statistics struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

// App is the class that holds the total number of requests,
// the total response time of all the requests, and the shutdown
// channel that is used to gracefully shutdown the server.
type App struct {
	Total        int
	ResponseTime float64
	Shutdown     chan bool
}

func main() {
	logger.Init(os.Stdout, os.Stdout)

	app := App{Shutdown: make(chan bool)}

	mux := http.NewServeMux()
	mux.Handle("/hash", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.hashHandler(w, r)
	}))
	mux.Handle("/shutdown", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.shutdownHandler(w, r)
	}))
	mux.Handle("/stats", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.statsHandler(w, r)
	}))

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error.Printf("Server stopped: %s\n", err)
		}
	}()

	<-app.Shutdown
	logger.Info.Println("Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	logger.Info.Println("Server gracefully shutdown")

}

// hashHandler reads the field value "password" and performs a
// SHA512 hash and then it is base64 encoded. It returns the response
// to the client after 5 seconds. It also increments the total request number
// everytime it is called as well as adds to the total response time.
func (app *App) hashHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	timer := time.NewTimer(time.Second * 5)
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println("Failed to parse form: ", err)
	}

	password := r.FormValue("password")
	if password == "" {
		logger.Error.Println("Password not provided")
		http.Error(w, "Password not provided", http.StatusBadRequest)
		return
	}
	switch {
	case r.Method == "POST":
		hashedPassword, err := utils.HashPassword([]byte(password))
		if err != nil {
			logger.Error.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encodedHashedPassword := utils.Base64(hashedPassword)

		<-timer.C
		fmt.Fprintf(w, "%s", encodedHashedPassword)

		app.ResponseTime += float64(time.Since(start)) / float64(time.Millisecond)
		app.Total++

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

// shutdownHandler will only accept a GET request and sends a response back
// to the client letting it know the server will be shut down. The shutdown
// channel is set to true which will unblock the code in the main function
// which will begin the graceful shutdown.
func (app *App) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		fmt.Fprint(w, "Gracefully shutting down server...")
		app.Shutdown <- true
	default:
		logger.Error.Println("Invalid HTTP Method")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

// statsHandler will only accept a GET request and constructs a JSON response
// which contains the total number of requests since the handler was called
// as well as the average response time of all of the requests.
func (app *App) statsHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		statistics := Statistics{}
		if app.Total == 0 {
			statistics = Statistics{app.Total, app.ResponseTime}
		} else {
			statistics = Statistics{app.Total, app.ResponseTime / float64(app.Total)}
		}

		jsonBody, err := json.Marshal(statistics)
		if err != nil {
			logger.Error.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBody)
	default:
		logger.Error.Println("Invalid HTTP Method")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
