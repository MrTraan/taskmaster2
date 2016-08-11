package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Route struct {
	Name    string
	Path    string
	Handler http.HandlerFunc
}

var Routes = []Route{
	Route{
		Name:    "Root",
		Path:    "/",
		Handler: mainHandler,
	},
	Route{
		Name:    "Status",
		Path:    "/status",
		Handler: statusHandler,
	},
	Route{
		Name:    "Start",
		Path:    "/start/",
		Handler: startHandler,
	},
	Route{
		Name:    "Stop",
		Path:    "/stop/",
		Handler: stopHandler,
	},
	Route{
		Name:    "Restart",
		Path:    "/restart/",
		Handler: restartHandler,
	},
	Route{
		Name:    "Shutdown",
		Path:    "/shutdown",
		Handler: shutdownHandler,
	},
}

func HandleRoutes() {
	for _, r := range Routes {
		http.HandleFunc(r.Path, r.Handler)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I should list commands here\n"))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(TaskPs(taskHolder)))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/start/"):]
	log.Printf("Received a start request on: %s\n", target)
	for _, t := range taskHolder {
		if t.Name == target {
			if err := t.Start(); err != nil {
				fmt.Fprintf(w, "Error while starting task %s: %v\n", target, err)
			} else {
				fmt.Fprintf(w, "Started task %s\n", target)
			}
			return
		}
	}
	fmt.Fprintf(w, "Error: unknown task %s\n", target)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/stop/"):]
	log.Printf("Received a stop request on: %s\n", target)
	for _, t := range taskHolder {
		if t.Name == target {
			if err := t.Stop(); err != nil {
				fmt.Fprintf(w, "Error while stopping task %s: %v\n", target, err)
			} else {
				fmt.Fprintf(w, "Stoped task %s\n", target)
			}
		}
		return
	}
	fmt.Fprintf(w, "Error: unknown task %s\n", target)
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/restart/"):]
	log.Printf("Received a restart request on: %s\n", target)
	for _, t := range taskHolder {
		if t.Name == target {
			if err := t.Stop(); err != nil {
				fmt.Fprintf(w, "Error while stopping task %s: %v\n", target, err)
				return
			}
			if err := t.Start(); err != nil {
				fmt.Fprintf(w, "Error while starting task %s: %v\n", target, err)
			} else {
				fmt.Fprintf(w, "Restarted task %s\n", target)
			}
			return
		}
	}
	fmt.Fprintf(w, "Error: unknown task %s\n", target)
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	for _, t := range taskHolder {
		if err := t.Kill(); err != nil {
			log.Printf("Error while shutting down task %s\n", t.Name)
		}
	}
	fmt.Fprintf(w, "Daemon is shutting down\n")
	defer os.Exit(0)
}
