package main

import (
	"reflect"
	"strings"
	"testing"
)

const TEST_CONF_FILE string = `[
    {
        "name": "test",
        "cmd": "ls -lR /",
        "numprocs": 1,
        "workingdir": "/tmp",
        "stdout": "/tmp/test_foo",
        "stderr": "/tmp/test_bar"
    },
    {
        "name": "foo",
        "cmd": "echo foo",
        "numprocs": 1,
        "autostart": true,
        "autorestart": "never",
        "exitcodes": [0, 2],
        "startretries": 3,
        "starttime": 5,
        "stopsignal": "KILL",
        "stoptime": 10,
        "stdout": "/tmp/foo",
        "stderr": "/tmp/bar",
        "env" : ["mykey=myvalue"]
    }
]`

var TEST_PARSED_FILE = []TaskSettings{
	TaskSettings{
		Name:        "test",
		Cmd:         "ls -lR /",
		Numprocs:    1,
		Autostart:   true,
		Autorestart: "NEVER",
		Stopsignal:  "TSTP",
		Workingdir:  "/tmp",
		Stdout:      "/tmp/test_foo",
		Stderr:      "/tmp/test_bar",
	},
	TaskSettings{
		Name:         "foo",
		Cmd:          "echo foo",
		Numprocs:     1,
		Autostart:    true,
		Autorestart:  "never",
		Exitcodes:    []int{0, 2},
		Startretries: 3,
		Starttime:    5,
		Stopsignal:   "KILL",
		Stoptime:     10,
		Stdout:       "/tmp/foo",
		Stderr:       "/tmp/bar",
		Env:          []string{"mykey=myvalue"},
	},
}

func TestReadConfFile(t *testing.T) {
	r := strings.NewReader(TEST_CONF_FILE)

	settings, err := ReadConfiguration(r)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(settings, TEST_PARSED_FILE) {
		t.Fatalf("Parsed conf doesn't match test conf\nTEST CONF: %v\nPARSED CONF: %v\n", TEST_PARSED_FILE, settings)
	}
}
