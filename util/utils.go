package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Change struct {
	BlockId int
	Data    string
	Success bool
}

//Based on "timer" prorty from config file
//Schedule gothread that will ping other gothreads via send channel
func Schedule(send chan bool, duration string) {
	u, err := time.ParseDuration(duration)
	if err == nil {
		for {
			send <- true
			time.Sleep(u)
		}
	} else {
		fmt.Println("Couldn't set a scheduler due to improper time format: " + duration)
	}
}

//Run any arbitrary POSIX shell command
func RunCmd(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	cmdStr := action["command"].(string)
	run := true

	for run {
		out, err := exec.Command("sh", "-c", cmdStr).Output()
		if err == nil {
			send <- Change{blockId, strings.TrimSuffix(string(out), "\n"), true}
		} else {
			send <- Change{blockId, err.Error(), false}
		}
		//Block untill other thread will ping you
		run = <-rec
	}
}

//Create a channel that captures all 34-64 signals
func GetSIGRTchannel() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	sigArr := make([]os.Signal, 31)
	for i := range sigArr {
		sigArr[i] = syscall.Signal(i + 0x22)
	}
	signal.Notify(sigChan, sigArr...)
	return sigChan
}
