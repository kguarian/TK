package main

import (
	"strconv"
	"syscall/js"
	"time"
)

var JS_GLOBAL js.Value = js.Global()
var display js.Value = JS_GLOBAL.Get("timer_display")

func main() {
	var lifetime chan byte = make(chan byte)

	JS_GLOBAL.Set("wasm_startTimer", js.FuncOf(startTimerAPI))
	JS_GLOBAL.Set("wasm_stopTimer", js.FuncOf(stopTimerAPI))

	updateTimerGo(&timer{id: "timer_display"})
	display.Set("value",
		"stop uses a *token* now. Should allow for increasing timer functionality for sleep-tracking tab.\n"+
			"TODO: Sleep for at least 8 hours a day and set a proper circadian rhythm :)\n\n"+
			"Implement a mutex for the token so that only one start/runner function can run at once.\n\n"+
			"Allow for whole app to be implemented react-style, with a bitwise constant index.html (index HTML) file\n\n"+
			"Implement timer under new GO_LOOT.js framework",
	)
	//TODO: The display is your instruction set
	_ = <-lifetime
}

//I wanted to place the private API functions in their own files, so these API
//functions convert, parse, and store input for organized function calls
func startTimerAPI(this js.Value, val []js.Value) (retVal interface{}) {
	var id string
	var duration int
	var refTime time.Time
	var inputDuration time.Duration
	var timr timer
	var err error
	if len(val) == 0 {
		retVal = js.ValueOf(1)
		return
	}
	id = val[0].String()
	duration, err = strconv.Atoi(val[1].String())
	if err != nil {
		retVal = js.ValueOf(2)
		return
	}
	refTime = time.Now()
	inputDuration = time.Duration(duration) * time.Millisecond

	//"Everything with a space, everything in its space"
	timr = timer{
		id:                 id,
		millisecondsPassed: 0,
		millisecondGoal:    inputDuration * time.Millisecond,
		startTime:          refTime,
		endTime:            refTime.Add(inputDuration),
	}

	Const_timers[timr.id] = &timr

	go startTimerGo(&timr)
	return
}

//convert, parse, function call.
func stopTimerAPI(this js.Value, val []js.Value) (retVal interface{}) {
	var timr *timer
	var id string
	var errorCodeChannel chan byte
	if len(val) == 0 {
		retVal = js.ValueOf(1)
		return
	}
	errorCodeChannel = make(chan int)
	id = val[0].String()
	timr = Const_timers[id]
	if timr == nil {
		return js.ValueOf(1)
	}
	go stopTimerGo(timr, errorCodeChannel)
	println(id)
	return
}
