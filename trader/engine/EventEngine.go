package engine

import (
	"fmt"
	"gonpy/trader"
	"reflect"
	"sync"
	"time"
)

var wg sync.WaitGroup

type EventEngine struct {
	Name     string
	Interval int
	Queue    chan *trader.Event
	Active   bool

	Thread          string
	Timer           string
	Handlers        map[string][]HandlerFunc
	GeneralHandlers []HandlerFunc
}

type HandlerFunc func(event *trader.Event)

func (e *EventEngine) run() {
	for {
		if e.Active {
			event := <-e.Queue
			e.Process(event)
		}
	}
}

func (e *EventEngine) runTimer() {
	for e.Active {
		time.Sleep(time.Duration(e.Interval) * time.Second)
		e.Put(&trader.Event{trader.EVENT_TIMER, time.Now()})
	}
}

func (e *EventEngine) Process(event *trader.Event) {
	handlers, handlersExist := e.Handlers[event.Type]
	if handlersExist {
		for i, handler := range handlers {
			fmt.Println(i, handler)
			handler(event)
		}
	}

	if len(e.GeneralHandlers) != 0 {
		for _, generalHandler := range e.GeneralHandlers {
			generalHandler(event)
		}
	}

}

func (e *EventEngine) Start() {
	e.Active = true
	wg.Add(2)
	go e.run()
	go e.runTimer()

}

func (e *EventEngine) Stop() {
	e.Active = false
	wg.Done()
}

func (e *EventEngine) Put(event *trader.Event) {
	e.Queue <- event
}

func (e *EventEngine) Register(Type string, handler HandlerFunc) {
	if value, ok := e.Handlers[Type]; ok {
		for _, item := range value {
			if reflect.ValueOf(item).Pointer() == reflect.ValueOf(handler).Pointer() {
				return
			}
		}
		e.Handlers[Type] = append(value, handler)
	}
}

func (e *EventEngine) Unregister(Type string, handler HandlerFunc) {
	if value, ok := e.Handlers[Type]; ok {
		for index, item := range value {
			if reflect.ValueOf(item).Pointer() == reflect.ValueOf(handler).Pointer() {

				if index == len(value)-1 {
					e.Handlers[Type] = value[:index]
					return
				}

				e.Handlers[Type] = append(value[:index], value[index+1:]...)
			}
		}

	}
}
