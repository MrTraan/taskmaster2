package main

import (
	"log"
	"net/http"
	"os"
)

const (
	PORT              string = ":8080"
	DEFAULT_CONF_FILE string = "./conf.json"
)

var (
	taskHolder []*Task
)

func main() {
	confFile, err := os.Open(DEFAULT_CONF_FILE)
	if err != nil {
		log.Fatal("Error while opening conf file: ", err)
	}
	settings, err := ReadConfiguration(confFile)
	if err != nil {
		log.Fatal("Error while parsing configuration: ", err)
	}

	logChannel := make(chan string)

	for _, s := range settings {
		t, err := NewTask(s, logChannel)
		if err != nil {
			log.Printf("Error while creating task: %v", err)
		} else {
			taskHolder = append(taskHolder, t)
			if t.Autostart {
				if err = t.Start(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	go func() {
		for {
			log.Println(<-logChannel)
		}
	}()

	HandleRoutes()
	log.Printf("Listening on port %s\n", PORT)
	err = http.ListenAndServe(PORT, nil)
	log.Fatal("Listen and serve error: ", err)
}
