package builtins

import (
	"os/exec"
	"time"
	"strings"
	"fmt"
)

type Change struct {
	BlockId int
	Data string
	Success bool
}

var FunctionMap = map[string]func(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	"#Date":Date,
}
//Based on "timer" prorty from config file
//Schedule gothread that will ping other gothreads via send channel
func Schedule(send chan bool, duration string){
	u, err := time.ParseDuration(duration)
	if err == nil {
		for {
			send <- true
			time.Sleep(u)
		}
	} else {
		fmt.Println("Couldn't set a scheduler due to improper time format: "+ duration)
	}
}
//Run any arbitrary POSIX shell command
//Use "send" channel to update blocks data
func RunCmd(blockId int, send chan Change, rec chan bool, action map[string]interface{} ) {
	cmdStr := action["command"].(string)
	run := true
	for run {
		out, err := exec.Command("sh","-c",cmdStr).Output()
		send <- Change{blockId, strings.TrimSuffix(string(out), "\n"), err == nil}
		//Block untill other thread will ping you
		run = <- rec
	}
}
//Update time according to "format" property
func Date(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		send <- Change{blockId, time.Now().Format(action["format"].(string)), true}
		//Block untill other thread will ping you
		run = <- rec
	}
}
