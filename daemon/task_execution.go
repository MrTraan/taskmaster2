package main

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

const (
	STATUS_STOPPED = iota
	STATUS_RUNNING
	STATUS_STARTING
)

var (
	ERR_TASK_ALREADY_RUNNING = errors.New("Task is already running")
)

type Task struct {
	TaskSettings
	Command *exec.Cmd
	Pid     int
	Stdout  io.ReadCloser
	Stderr  io.ReadCloser
	Status  int
	Uptime  time.Time
	Log     chan string
}

//Status formats task status and returns it as a string
func (t *Task) GetStatus() string {
	switch t.Status {
	case STATUS_STOPPED:
		return "Stopped"
	case STATUS_RUNNING:
		return "Running"
	case STATUS_STARTING:
		return "Starting"
	default:
		return "Unkown status"
	}
}

//String formats a task status and returns it as a string
//Useful to log task to user
func (t *Task) String() string {
	if t.Status == STATUS_RUNNING {
		return fmt.Sprintf("%-20s%-20s%-20d%-20v\n",
			t.Name, t.GetStatus(), t.Pid, time.Since(t.Uptime))
	} else {
		return fmt.Sprintf("%-20s%-20s%-20d%-20v\n",
			t.Name, t.GetStatus(), t.Pid, 0)
	}
}

//Start launch a task
//It returns an error if the task is already running or if an error occured while starting the task
func (t *Task) Start() error {
	if t.Status != STATUS_STOPPED {
		return ERR_TASK_ALREADY_RUNNING
	}

	t.Status = STATUS_STARTING
	if err := t.Command.Start(); err != nil {
		t.Status = STATUS_STOPPED
		return err
	}

	t.Status = STATUS_RUNNING
	t.Uptime = time.Now()

	go func() {
		if err := t.Command.Wait(); err != nil {
			t.Log <- "Error"
			t.Status = STATUS_STOPPED
			return
		}
		t.Log <- "Done"
		t.Status = STATUS_STOPPED
	}()
	return nil
}

func NewTask(settings TaskSettings) (task *Task, err error) {
	task = new(Task)
	task.TaskSettings = settings

	task.Log = make(chan string)
	args := strings.Split(task.Cmd, " ")
	task.Command = exec.Command(args[0], args[1:]...)
	err = nil
	return task, err
}

func TaskPs(tasks []*Task) string {
	str := ""
	str += fmt.Sprintf("%-20s%-20s%-20s%-20s\n", "NAME", "STATUS", "PID", "UPTIME")
	for _, t := range tasks {
		str += fmt.Sprint(t)
	}
	return str
}
