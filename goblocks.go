package main
import (
	"encoding/json"
	"fmt"
	"goblocks/builtins"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type configStruct struct {
	Separator string
	ConfigReloadSignal int //currently unused
	Actions []map[string]interface{}
}

var blocks []string
var channels []chan bool
var signalMap map[string]int = make(map[string]int)

func main() {
	config, err := readConfig(os.Getenv("HOME")+"/.config/goblocks.json")
	if  err == nil {
		channels = make([]chan bool, len(config.Actions))
		//recChannel is common for gothreads contributing to status bar
		recChannel := make(chan builtins.Change)
		for i, action := range config.Actions {
			//Assign a cell for each separator/prefix/action/suffix
			if config.Separator != "" {
				blocks = append(blocks, config.Separator)
			}
			if value, ok := action["prefix"]; ok {
				blocks = append(blocks, value.(string))
			}
			blocks = append(blocks, "action")
			actionId := len(blocks)-1
			if value, ok := action["suffix"]; ok {
				blocks = append(blocks, value.(string))
			}
			//Create an unique channel for each action
			channels[i] = make(chan bool)
			signalMap["signal "+action["updateSignal"].(string)] = i
			if (action["command"].(string))[0] == '#' {
				go builtins.FunctionMap[action["command"].(string)](actionId, recChannel, channels[i], action)
			} else {
				go builtins.RunCmd(actionId, recChannel, channels[i], action)
			}
			timer := action["timer"].(string)
			if timer != "0" {
				go builtins.Schedule(channels[i], timer)
			}
		}
		go handleSignals(getSIGRTchannel())
		//start event loop
		for {
			//Block untill some gothread has an update
			res := <- recChannel
			if res.Success {
				blocks[res.BlockId] = res.Data
			} else {
				fmt.Println(res.Data)
				blocks[res.BlockId] = "ERROR"
			}
			updateStatusBar()
		}
	} else {
		fmt.Println(err)
	}
}
//Read config and map it to configStruct
func readConfig(path string) ( config configStruct, err error) {
	var file *os.File
	file, err = os.Open(path)
	defer file.Close()
	if err == nil {
		var byteValue []byte
		byteValue, err = ioutil.ReadAll(file)
		if err == nil {
			err = json.Unmarshal([]byte(byteValue), &config)
		}
	}
	return config, err
}
//Create a channel that captures all 34-64 signals
func getSIGRTchannel() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	sigArr := make([]os.Signal, 31)
	for i := range sigArr {
		sigArr[i] = syscall.Signal(i + 0x22)
	}
	signal.Notify(sigChan, sigArr...)
	return sigChan
}
//Goroutine that pings a channel according to received signal
func handleSignals(rec chan os.Signal) {
	for {
		sig := <- rec
		if index, ok := signalMap[sig.String()]; ok {
			channels[index] <- true
		}
	}
}
//Craft status text out of blocks data
func updateStatusBar() {
	var builder strings.Builder
	for _, s := range blocks {
		builder.WriteString(s)
	}
	//	fmt.Println(builder.String())
	//	set dwm status text
	exec.Command("xsetroot","-name",builder.String()).Run()
}
