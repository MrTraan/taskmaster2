package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	PORT string = ":8080"
)

var (
	taskHolder []*Task
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I should list commands here\n"))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/start/"):]
	fmt.Fprintf(w, "Received a start request on: %s\n", title)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(TaskPs(taskHolder)))
}

func main() {
	confFile, err := os.Open("./conf.json")
	if err != nil {
		log.Fatal("Error while opening conf file: ", err)
	}
	settings, err := ReadConfiguration(confFile)
	if err != nil {
		log.Fatal("Error while parsing configuration: ", err)
	}
	for _, s := range settings {
		t, err := NewTask(s)
		if err != nil {
			log.Fatal("Error while creating task: ", err)
		}
		taskHolder = append(taskHolder, t)
	}

	log.Printf("Listening on port %s\n", PORT)
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/start/", startHandler)
	http.HandleFunc("/status", statusHandler)
	err = http.ListenAndServe(PORT, nil)
	log.Fatal("Listen and serve error: ", err)
}
