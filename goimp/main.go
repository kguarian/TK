package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var JS_GLOBAL js.Value = js.Global()
var display js.Value = JS_GLOBAL.Get("timer_display")

func main() {
	var lifetime chan JsFunction

	var startTimerAPI_Hook, stopTimerAPI_Hook JsFunction

	lifetime = make(chan JsFunction)

	startTimerAPI_Hook.Init("wasm_startTimer", startTimerAPI)
	stopTimerAPI_Hook.Init("wasm_stopTimer", stopTimerAPI)

	startTimerAPI_Hook.Expose(JS_GLOBAL)
	stopTimerAPI_Hook.Expose(JS_GLOBAL)

	display.Set("value", "11.9 - fixing stop")
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
		signals:            make(chan byte),
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
	errorCodeChannel = make(chan byte)
	id = val[0].String()
	timr = Const_timers[id]
	if timr == nil {
		return js.ValueOf(1)
	}
	fmt.Println(timr)
	go stopTimerGo(timr, errorCodeChannel)
	println(id)
	return
}
