package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	PORT string = ":8080"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I should list commands here\n"))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/start/"):]
	fmt.Fprintf(w, "Received a start request on: %s\n", title)
}

func main() {
	log.Printf("Listening on port %s\n", PORT)
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/start/", startHandler)
	err := http.ListenAndServe(PORT, nil)
	log.Fatal("Listen and serve error: ", err)
}
