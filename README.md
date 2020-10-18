# Goblocks
Goblocks is a [dwm](https://dwm.suckless.org/) status bar program partially inspired by [dwmblocks](https://github.com/torrinfail/dwmblocks).

![Goblocks in action](https://i.imgur.com/kH1G8tr.png)

It is lightweight, fast, multithreaded, well(-ish) documented and easily customizable. It is also one of very few status bars done correctly. Most of them "naively"
scan for changes every 1s or so, whereas goblocks only updates when there are updates to be done.

## Installation
Pull this repo and run:
```
make install
```
This will build this project and generate default config file at $HOME/.config/goblocks.json. After that you can simply run it with
```
goblocks &!
```

If you want to uninstall it, just run.
```
make uninstall
```


## Customization
You can add your own blocks in [config file](https://github.com/Stargarth/Goblocks/blob/master/goblocks.json). Just remember to modify config at ~/.config/
goblocks.json.
After you have applied your changes, just re run the aplication.

In order to add new block to the bar add new action object to actions array in goblocks.json:
```
"actions":
	[
		{
			"prefix": "🗓 ",
			"updateSignal": "35",
			"command": "#Date",
			"format": "Monday 02/01/2006 15:04:05",
			"timer": "1s"
		},
		{
			"prefix": "🔊 ",
			"updateSignal": "36",
			"command": "amixer get Master | awk ' \/%\/ {print $5 $6}' | head -1 | sed -e 's\/]\\[\/ \/g'",
			"timer": "0"
		},
		{
			"prefix": "Mem: ",
			"updateSignal": "37",
			"command": "#Memory",
			"suffix": "%",
			"format": "%.2f",	
			"timer": "2s"
		},
		{
			"prefix": "CPU: ",
			"updateSignal": "38",
			"command": "#Cpu",
			"suffix": "%",
			"format": "%.2f",	
			"timer": "2s"
		}
	]
```
### Config properties
1. **prefix(optional)**:text that will be displayed before the output of command.
2. **updateSignal**: signal that can be sent to goblocks to re-run the command and refresh output in bar. This value must be in range 34-64 inclusive.
3. **command**: functionality that will produce the text for status bar. If it starts with hash-tag goblocks will try to invoke corresponding built-in function.
If it doesn't then goblocks will run command in shell(sh) and put it's output in status bar.
4. **suffix(optional)**: Text that will be displayed after the output of command
5. **Timer**: Time period after which command is being re-run and text in status bar refreshed. s for seconds(1s), m for minutes(1m), h for hours(1h).
When set to 0 the block will only update at startup and when **updateSignal** is sent to goblocks.

### Signaling Blocks
Blocks can be updated immediately by sending the corresponding **updateSignal** to the goblocks process.
For example, `kill -35 $(pidof goblocks)` will update the block with **updateSignal** 35.
Alternatively, you can signal a block with `pkill -RTMIN+<x> goblocks` where `x` is the **updateSignal** minus 34.
So `pkill -RTMIN+1 goblocks` will update blocks with **updateSignal** 35.

> Note if multiple blocks are assigned the same **updateSignal** only the first will be updated.

## Miscellaneous
### Formating date with #Date built in program
Format a reference layout parameter (**Mon Jan 2 15:04:05 MST 2006**) in a way that you would like to have your time formatted.
If you have trouble check out this [guideline](https://www.golangprograms.com/get-current-date-and-time-in-various-format-in-golang.html).
### Escape sh commands
Some sh commands have characters (e.g. \\) which have a special meaning in JSON format. To easily escape them just dump your sh commands to
[json formatter](https://www.freeformatter.com/json-escape.html).
### TODO
1. ~~Create a block that displays weather condition.~~ Done
2. ~~Add implementation documentation~~ Code is pretty readable and has comments. This should be enough
3. Add support for colour text in blocks.

## Feedback
Your feedback is more than welcome and if you have some questions/issues/feature ideas just post them on
[issues page](https://github.com/Stargarth/Goblocks/issues).
