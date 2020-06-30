package builtins

import (
	"os/exec"
	"time"
	"strings"
	"fmt"
    "github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"math"
)

type Change struct {
	BlockId int
	Data string
	Success bool
}
//blockId is automatically allocated
//send channel is used to update blocks data
//Gothreads are waked up by messages on rec channels
//action is a map of whatever was in action json object for corressponding action
var FunctionMap = map[string]func(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	"#Date":Date,
	"#Memory":Memory,
	"#Cpu":Cpu,
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
//Get current % of used memory
func Memory(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		v, _ := mem.VirtualMemory()

		send <- Change{blockId, fmt.Sprintf("%.2f", math.Round(v.UsedPercent*100)/100), true}
		//Block untill other thread will ping you
		run = <- rec
	}
}
//Get current % of used CPU
func Cpu(blockId int, send chan Change, rec chan bool, action map[string]interface{}) {
	run := true
	for run {
		val, _ := cpu.Percent(time.Second, false)
		send <- Change{blockId, fmt.Sprintf("%.2f", math.Round(val[0]*100)/100), true}
		//Block untill other thread will ping you
		run = <- rec
	}
}
