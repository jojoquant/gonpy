package main

import (
	"fmt"
	"gonpy/trader"
	"gonpy/trader/engine"
	"time"
)

func ETimerHandlerFunc(event *trader.Event) {
	fmt.Println("ETimerHandlerFunc", event, event.Type, event.Data)
}
func ETimerHandlerFunc2(event *trader.Event) {
	fmt.Println("ETimerHandlerFunc2", event, event.Type, event.Data)
}

func main() {
	e := &engine.EventEngine{
		Name:     "case1EventEngine",
		Interval: 1,
		Queue:    make(chan *trader.Event, 0),
		Active:   false,

		Handlers: map[string][]engine.HandlerFunc{
			"eTimer":  {},
			"eTimer2": {ETimerHandlerFunc, ETimerHandlerFunc2},
		},
	}
	event1 := &trader.Event{Type: "eTimer", Data: "123"}
	event2 := &trader.Event{Type: "eTimer2", Data: "12333"}
	
	e.Register("eTimer", ETimerHandlerFunc)
	e.Register("eTimer", ETimerHandlerFunc)
	e.Register("eTimer", ETimerHandlerFunc)
	fmt.Println(e)
	e.Unregister("eTimer", ETimerHandlerFunc)
	e.Unregister("eTimer", ETimerHandlerFunc)
	fmt.Println(e)

	e.Start()
	e.Put(event1)
	e.Put(event2)
	time.Sleep(2*time.Second)
	e.Stop()

	
	select{}
}
