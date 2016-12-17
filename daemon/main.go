package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	PORT              string = ":8080"
	DEFAULT_CONF_FILE string = "./conf.json"
)

var (
	taskHolder  []*Task
	gLogChannel chan string
	gConfFile   string
)

func main() {
	gConfFile = DEFAULT_CONF_FILE
	gLogChannel = make(chan string)

	sighupChannel := make(chan os.Signal, 1)
	signal.Notify(sighupChannel, syscall.Signal(0x1))

	if err := reloadConf(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			log.Println(<-gLogChannel)
		}
	}()

	go func() {
		for {
			<-sighupChannel
			if err := reloadConf(); err != nil {
				log.Println(err)
			} else {
				log.Println("Configuration file reloaded")
			}
		}
	}()

	HandleRoutes()
	log.Printf("Listening on port %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	log.Fatal("Listen and serve error: ", err)
}

func isEnvEqual(e1 []string, e2 []string) bool {
	if len(e1) != len(e2) {
		return false
	}

	for i, v := range e1 {
		if v != e2[i] {
			return false
		}
	}
	return true
}

func reloadConf() error {
	confFile, err := os.Open(gConfFile)
	if err != nil {
		return fmt.Errorf("Error: could not open configuration file\n")
	}
	defer confFile.Close()

	settings, err := ReadConfiguration(confFile)
	if err != nil {
		return fmt.Errorf("Error: an error occured while reading configuration file: %s\n", err)
	}

	err = CheckConfiguration(settings)
	if err != nil {
		return fmt.Errorf("Error: invalid configuration file: %s\n", err)
	}

	persisting := []*Task{}
	for _, t := range taskHolder {
		persists := false
		for _, s := range settings {
			if t.Name == s.Name {
				persists = true
			}
		}

		if !persists {
			t.Kill()
		} else {
			persisting = append(persisting, t)
		}
	}
	taskHolder = persisting

	for _, s := range settings {
		duplicate := false
		for i, t := range taskHolder {
			if t.Name == s.Name {
				duplicate = true

				if t.Cmd != s.Cmd || t.Umask != s.Umask || t.Workingdir != s.Workingdir ||
					t.Stdout != s.Stdout || t.Stderr != s.Stderr || !isEnvEqual(t.Env, s.Env) {
					t.Kill()
					taskHolder = append(taskHolder[:i], taskHolder[i+1:]...)
					newTask, err := NewTask(s, gLogChannel)
					if err != nil {
						log.Printf("Error while creating task: %v", err)
					} else {
						taskHolder = append(taskHolder, newTask)
						if newTask.Autostart {
							if err = newTask.Start(); err != nil {
								log.Fatal(err)
							}
						}
					}
				} else {
					t.Autostart = s.Autostart
					t.Autorestart = s.Autorestart
					t.Exitcodes = s.Exitcodes
					t.Startretries = s.Startretries
					t.Starttime = s.Starttime
					t.Stopsignal = s.Stopsignal
					t.Stoptime = s.Stoptime
				}
				break
			}
		}

		if !duplicate {
			newTask, err := NewTask(s, gLogChannel)
			if err != nil {
				log.Printf("Error while creating task: %v", err)
			} else {
				taskHolder = append(taskHolder, newTask)
				if newTask.Autostart {
					if err = newTask.Start(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
	return nil
}
