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

	Handlers        map[string][]HandlerFunc
	GeneralHandlers []HandlerFunc
}

type HandlerFunc func(event *trader.Event)

func NewEventEngine() *EventEngine {
	// 这里真特殊, 不能直接return, 必须有个变量接然后return
	engine := &EventEngine{
		Name:     "EventEngine",
		Interval: 1,
		Queue:    make(chan *trader.Event, 10),
		Active:   false,

		Handlers: map[string][]HandlerFunc{},
	}
	return engine 
}

func (e *EventEngine) run() {
	ticker := time.NewTicker(time.Duration(e.Interval) * time.Second)
	for e.Active {
		select {
		case event := <-e.Queue:
			e.Process(event)
		case <-ticker.C:
			e.Put(&trader.Event{Type: trader.EVENT_TIMER, Data: time.Now()})
		}
	}
}

func (e *EventEngine) Process(event *trader.Event) {
	handlers, handlersExist := e.Handlers[event.Type]
	if handlersExist {
		for i, handler := range handlers {
			fmt.Println(i, handler)
			wg.Add(1)
			go handler(event)
		}
	}

	if len(e.GeneralHandlers) != 0 {
		for _, generalHandler := range e.GeneralHandlers {
			wg.Add(1)
			go generalHandler(event)
		}
	}

}

func (e *EventEngine) Start() {
	e.Active = true
	wg.Add(1)
	go e.run()
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
