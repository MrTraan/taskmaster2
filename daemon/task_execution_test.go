package main

import (
	"fmt"
	"testing"
)

var TEST_CONF = TaskSettings{
	Name: "test command",
	Cmd:  "sleep 5",
}

func TestNewTask(t *testing.T) {
	task, err := NewTask(TEST_CONF)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(task)
	if err = task.Start(); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", <-task.Log)
}
