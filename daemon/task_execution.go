package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

const (
	STATUS_STOPPED  = "STOPPED"
	STATUS_RUNNING  = "RUNNING"
	STATUS_STARTING = "STARTING"
	STOP_TIMEOUT    = 10000 * time.Millisecond
)

var (
	ERR_TASK_ALREADY_RUNNING = errors.New("Task is already running")
	ERR_TASK_ALREADY_STOPPED = errors.New("Task is already stopped")
)

var signalMap = map[string]syscall.Signal{
	"ABRT":   syscall.Signal(0x6),
	"ALRM":   syscall.Signal(0xe),
	"BUS":    syscall.Signal(0x7),
	"CHLD":   syscall.Signal(0x11),
	"CLD":    syscall.Signal(0x11),
	"CONT":   syscall.Signal(0x12),
	"FPE":    syscall.Signal(0x8),
	"HUP":    syscall.Signal(0x1),
	"ILL":    syscall.Signal(0x4),
	"INT":    syscall.Signal(0x2),
	"IO":     syscall.Signal(0x1d),
	"IOT":    syscall.Signal(0x6),
	"KILL":   syscall.Signal(0x9),
	"PIPE":   syscall.Signal(0xd),
	"POLL":   syscall.Signal(0x1d),
	"PROF":   syscall.Signal(0x1b),
	"PWR":    syscall.Signal(0x1e),
	"QUIT":   syscall.Signal(0x3),
	"SEGV":   syscall.Signal(0xb),
	"STKFLT": syscall.Signal(0x10),
	"STOP":   syscall.Signal(0x13),
	"SYS":    syscall.Signal(0x1f),
	"TERM":   syscall.Signal(0xf),
	"TRAP":   syscall.Signal(0x5),
	"TSTP":   syscall.Signal(0x14),
	"TTIN":   syscall.Signal(0x15),
	"TTOU":   syscall.Signal(0x16),
	"UNUSED": syscall.Signal(0x1f),
	"URG":    syscall.Signal(0x17),
	"USR1":   syscall.Signal(0xa),
	"USR2":   syscall.Signal(0xc),
	"VTALRM": syscall.Signal(0x1a),
	"WINCH":  syscall.Signal(0x1c),
	"XCPU":   syscall.Signal(0x18),
	"XFSZ":   syscall.Signal(0x19),
}

type Task struct {
	TaskSettings
	Process    *os.Process
	Attributes *os.ProcAttr
	ExitState  *os.ProcessState
	Status     string
	Uptime     time.Time
	Log        chan string
}

//NewTask create a new task instance according to settings given in parameter
//Logs will be output to logChannel
func NewTask(settings TaskSettings, logChannel chan string) (task *Task, err error) {
	task = new(Task)
	task.TaskSettings = settings

	task.Log = logChannel
	stdFiles, err := openTaskStdout(task)
	if err != nil {
		return nil, err
	}

	task.Attributes = &os.ProcAttr{
		Dir:   settings.Workingdir,
		Env:   settings.Env,
		Files: stdFiles,
	}
	task.Status = STATUS_STOPPED
	return task, nil
}

func openTaskStdout(t *Task) ([]*os.File, error) {
	stdoutFile, err := os.OpenFile(t.Stdout, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	stderrFile, err := os.OpenFile(t.Stderr, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	return []*os.File{nil, stdoutFile, stderrFile}, nil
}

//Start launch a task
//It returns an error if the task is already running or if an error occured while starting the task
func (t *Task) Start() error {
	var err error
	if t.Status != STATUS_STOPPED {
		return ERR_TASK_ALREADY_RUNNING
	}

	t.Status = STATUS_STARTING
	args := strings.Split(t.Cmd, " ")
	t.Process, err = os.StartProcess(args[0], args, t.Attributes)
	if err != nil {
		t.Status = STATUS_STOPPED
		return err
	}
	log.Printf("Task %s starting...\n", t.Name)
	t.Uptime = time.Now()

	go func() {
		time.Sleep(time.Duration(t.Starttime) * time.Millisecond)
		if t.Status == STATUS_STARTING {
			t.Status = STATUS_RUNNING
			t.Uptime = time.Now()
			log.Printf("Task %s is now running\n", t.Name)
		}
	}()

	go func() {
		t.ExitState, err = t.Process.Wait()
		if err != nil {
			t.Status = STATUS_STOPPED
			t.Log <- fmt.Sprintf("Task %s ended with error: %v", t.Name, err)
			return
		}
		t.Status = STATUS_STOPPED
		t.Log <- fmt.Sprintf("Task %s ended graciously: %s", t.Name, t.ExitState)
	}()
	return nil
}

func (t *Task) Stop() error {
	stopChan := make(chan string)

	if t.Status == STATUS_STOPPED {
		return ERR_TASK_ALREADY_STOPPED
	}

	if signalMap[t.Stopsignal] == syscall.Signal(0) {
		return errors.New("Unknown stop signal " + t.Stopsignal)
	}

	go func() {
		timeoutChan := time.After(STOP_TIMEOUT)
		for {
			select {
			case <-timeoutChan:
				stopChan <- "Error while stopping task: timeout"
				return
			default:
				if t.Status == STATUS_STOPPED {
					stopChan <- "SUCCESS"
					return
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	t.Process.Signal(signalMap[t.Stopsignal])

	msg := <-stopChan
	if msg == "SUCCESS" {
		return nil
	} else {
		return errors.New(msg)
	}
}

func (t *Task) Kill() error {
	return t.Process.Kill()
}

func TaskPs(tasks []*Task) string {
	str := ""
	str += fmt.Sprintf("%-20s%-20s%-20s%-20s\n", "NAME", "STATUS", "PID", "UPTIME")
	for _, t := range tasks {
		str += fmt.Sprint(t)
	}
	return str
}

//String formats a task status and returns it as a string
//Useful to log task to user
func (t *Task) String() string {
	if t.Status != STATUS_STOPPED {
		return fmt.Sprintf("%-20s%-20s%-20d%-20v\n",
			t.Name, t.Status, t.Process.Pid, time.Since(t.Uptime))
	} else {
		return fmt.Sprintf("%-20s%-20s%-20d%-20v\n",
			t.Name, t.Status, 0, 0)
	}
}
