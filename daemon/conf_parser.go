package main

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	BUF_SIZE int = 1024
)

type TaskSettings struct {
	Name         string   `json:Name`
	Cmd          string   `json:cmd`
	Umask        int      `json:umask`
	Numprocs     int      `json:numprocs`
	Workingdir   string   `json:workingdir`
	Autostart    bool     `json:autostart`
	Autorestart  string   `json:autorestart`
	Exitcodes    []int    `json:exitcodes`
	Startretries int      `json:startretries`
	Starttime    int      `json:starttime`
	Stopsignal   string   `json:stopsignal`
	Stoptime     int      `json:stoptime`
	Stdout       string   `json:stdout`
	Stderr       string   `json:stderr`
	Env          []string `json:env`
}

//String returns a string to display a TaskSettings structure
func (s TaskSettings) String() string {
	return fmt.Sprintf(`TaskSetting:	
Name: %s\n
Cmd: %s\n
Numprocs: %d\n
Autostart: %v
Autorestart: %s
Exitcode: %v
Startretries: %d
Starttime: %d
Stopsignal: %s
Stoptime: %d
Stdout: %s
Stderr: %s
Env: %v
`,
		s.Name, s.Cmd, s.Numprocs, s.Autostart, s.Autorestart, s.Exitcodes, s.Startretries,
		s.Starttime, s.Stopsignal, s.Stoptime, s.Stdout, s.Stderr, s.Env)
}

//ReadConfiguration reads a configuration stored in a reader and returns a slice of TaskSettings
//Or an error if the file could not be found or badly formatted
func ReadConfiguration(reader io.Reader) (settings []TaskSettings, err error) {
	dec := json.NewDecoder(reader)

	_, err = dec.Token()
	if err != nil {
		return nil, err
	}
	for dec.More() {
		newEntry := TaskSettings{
			Numprocs:    1,
			Autostart:   false,
			Autorestart: "NEVER",
			Stopsignal:  "KILL",
			Stdout:      "/dev/null",
			Stderr:      "/dev/null",
			Env:         nil,
		}

		err := dec.Decode(&newEntry)
		if err != nil {
			return nil, err
		}
		settings = append(settings, newEntry)
		for i := 2; i <= newEntry.Numprocs; i++ {
			tmp := newEntry
			tmp.Name = tmp.Name + fmt.Sprintf(":%d", i)
			settings = append(settings, tmp)
		}
	}

	return settings, err
}

func CheckConfiguration(settings []TaskSettings) error {
	for _, entry := range settings {
		duplicates := 0
		for _, e := range settings {
			if e.Name == entry.Name {
				duplicates++
			}
		}
		if duplicates > 1 {
			return fmt.Errorf("Duplicate task name: %s\n", entry.Name)
		}
	}
	return nil
}
