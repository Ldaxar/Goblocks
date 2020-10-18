package main

import (
	"goblocks/util"
	"log"
	"os"
	"strings"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

var (
	blocks    []string
	channels  []chan bool
	signalMap map[string]int = make(map[string]int)
)

func main() {
	// establish connection to X server
	x, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer x.Close()
	root := xproto.Setup(x).DefaultScreen(x).Root

	config := util.ReadConfig("goblocks.json")
	channels = make([]chan bool, len(config.Actions))
	//recChannel is common for gothreads contributing to status bar
	recChannel := make(chan util.Change)
	for i, action := range config.Actions {
		//Assign a cell for each separator/prefix/action/suffix
		if config.Separator != "" {
			blocks = append(blocks, config.Separator)
		}
		if value, ok := action["prefix"]; ok {
			blocks = append(blocks, value.(string))
		}
		blocks = append(blocks, "action")
		actionId := len(blocks) - 1
		if value, ok := action["suffix"]; ok {
			blocks = append(blocks, value.(string))
		}
		//Create an unique channel for each action
		channels[i] = make(chan bool)
		signalMap["signal "+action["updateSignal"].(string)] = i
		if (action["command"].(string))[0] == '#' {
			go util.FunctionMap[action["command"].(string)](actionId, recChannel, channels[i], action)
		} else {
			go util.RunCmd(actionId, recChannel, channels[i], action)
		}
		timer := action["timer"].(string)
		if timer != "0" {
			go util.Schedule(channels[i], timer)
		}
	}
	go handleSignals(util.GetSIGRTchannel())
	//start event loop
	for {
		//Block until some gothread has an update
		res := <-recChannel
		if res.Success {
			blocks[res.BlockId] = res.Data
		} else {
			log.Println(res.Data)
			blocks[res.BlockId] = "ERROR"
		}
		updateBar(x, root)
	}
}
//Goroutine that pings a channel according to received signal
func handleSignals(rec chan os.Signal) {
	for {
		sig := <-rec
		if index, ok := signalMap[sig.String()]; ok {
			channels[index] <- true
		}
	}
}
//Craft status text out of blocks data
func updateBar(conn *xgb.Conn, rootWindow xproto.Window) {
	var builder strings.Builder
	for _, s := range blocks {
		builder.WriteString(s)
	}
	xproto.ChangeProperty(conn, xproto.PropModeReplace, rootWindow, xproto.AtomWmName, xproto.AtomString, 8, uint32(builder.Len()), []byte(builder.String()))
}
