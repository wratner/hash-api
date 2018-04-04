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

// Statistics ...
type Statistics struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

// App ...
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
		hashedPassword := utils.HashPassword([]byte(password))
		encodedHashedPassword := utils.Base64(hashedPassword)

		<-timer.C
		fmt.Fprintf(w, "%s", encodedHashedPassword)

		app.ResponseTime += float64(time.Since(start)) / float64(time.Millisecond)
		app.Total++

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (app *App) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		fmt.Fprint(w, "Gracefully shutting down server...")
		app.Shutdown <- true
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBody)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
