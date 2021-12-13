package main

import (
	"fmt"
	"syscall/js"
	"time"
)

var Const_timers map[string]*timer = make(map[string]*timer)

type timer struct {
	millisecondsPassed time.Duration
	millisecondGoal    time.Duration
	startTime          time.Time
	endTime            time.Time
	id                 string
	//event queue. Church-Turing Thesis: 1 channel can do anything that parallel channels can do, but possibly (and often) slower.
	//also, this is a reference type, treating like pointer
	signals chan byte
}

func startTimerGo(t *timer) (errcode int) {
	println("startTimerGo()")
	var i int
	var usedTimer timer
	//reference type
	var synValue byte
	var processedString string
	var mapTimerPtr *timer

	// var dur time.Duration

	fmt.Println(t.endTime)

	usedTimer = timer{
		id:                 t.id,
		millisecondsPassed: time.Since(t.startTime),
		millisecondGoal:    time.Until(t.endTime) * time.Millisecond,
		startTime:          t.startTime,
		endTime:            t.startTime.Add(time.Millisecond * time.Until(t.endTime)),
		signals:            make(chan byte),
	}

	mapTimerPtr = &usedTimer

	Const_timers[t.id] = mapTimerPtr

	for ; time.Now().Before(t.endTime); time.Sleep(time.Millisecond) {
		// dur = time.Until(t.endTime)
		JS_GLOBAL.Get("timer_display").Set("value", js.ValueOf(time.Until(t.endTime).String()))
		// t.mut.Lock()
		// defer t.mut.Unlock()
		println(len(mapTimerPtr.signals))
		if len(mapTimerPtr.signals) != 0 {
			processedString = fmt.Sprintf("%v (stopped)", time.Until(t.endTime).String())
			JS_GLOBAL.Get("timer_display").Set("value", js.ValueOf(processedString))
			synValue = <-mapTimerPtr.signals
			switch synValue {
			//stop funcall triggered.
			case 0:
				//should finish all other functions involving variables used by the timer and the timer itself, then clear the map association, then return.
				//Hinge safety on language design. That's the point of Go, right?
				//LOGIC: Suppose a function is already running and needs a pointer, If it has the pointer, then it's golden.
				//	If not, then it should nil check and this cleanup should stop this timer's future functionality until the next start API call.
				//	We should already be done executing the function, so this should be trivial.

				//in an ideal world, this assignment triggers a drop function. Here, we are left at the mercy of the GC about when to deallocate the actual struct.
				Const_timers[t.id] = nil
				return
			default:
				break
			}
		}
	}

	JS_GLOBAL.Get("timer_display").Set("value", js.ValueOf("Perfection"))

	for i = 0; i < 3; i++ {
		JS_GLOBAL.Get("timer_display").Set("value", js.ValueOf("Perfection🔥"))
		time.Sleep(time.Millisecond * time.Duration(333))
		JS_GLOBAL.Get("timer_display").Set("value", js.ValueOf("Perfection🔥🔥"))
		time.Sleep(time.Millisecond * time.Duration(333))
		JS_GLOBAL.Get("timer_display").Set("value", js.ValueOf("Perfection🔥🔥🔥"))
		time.Sleep(time.Millisecond * time.Duration(334))
	}
	return
}
func stopTimerGo(timr *timer, ecc chan byte) (errcode int) {
	println("stopTimerGo()")
	var callTime time.Time

	callTime = time.Now()

	if timr == nil {
		fmt.Println("timr == nil")
		errcode = 1
		return
	}

	if timr.endTime.Before(callTime) {
		fmt.Println("timr.endTime.Before(callTime)")
		errcode = 2
		return
	}

	//timer should be handled via the channel in this stop event.
	timr.signals <- 0
	fmt.Printf("passed stop token (0) into channel")
	return
}

func StopwatchGo(timr *timer) {

}
func updateTimerGo(timr *timer) (errcode int) {
	println("updateTimerGo()")
	var duration_ms uint64

	JS_GLOBAL.Get(timr.id).Set("value", fmt.Sprintf("%f seconds", float64(duration_ms)/1000))
	return
}
