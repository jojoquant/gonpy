package main

import (
	"fmt"
	"gonpy/trader"
	"gonpy/trader/engine"
)

func main() {
	e := engine.NewEventEngine()
	e.Handlers = map[string][]engine.HandlerFunc{
		trader.EVENT_TIMER:  {func(event *trader.Event) {fmt.Println(trader.EVENT_TIMER, event)}},
		trader.EVENT_LOG:  {func(event *trader.Event) {fmt.Println(trader.EVENT_LOG,event, event.Data)}},
	}
	// event1 := &trader.Event{Type: "eTimer", Data: "123"}
	// event2 := &trader.Event{Type: "eTimer2", Data: "12333"}

	m := engine.NewMainEngine(e)
	m.WriteLog("hhh", "test")
	m.WriteLog("hhhdddd", "test2")
	fmt.Println(m)
	cta := engine.NewCtaEngine(e)
	m.AddEngine(cta)
	select{}
}