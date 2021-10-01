package engine

import (
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

	sync.Mutex
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
	ticker.Stop()
}

func (e *EventEngine) Process(event *trader.Event) {

	if handlers, handlersExist := e.Handlers[event.Type]; handlersExist {
		for _, handler := range handlers {
			wg.Add(1)
			go func(handler HandlerFunc){
				handler(event)
				wg.Done()
			}(handler)
		}
	}

	if len(e.GeneralHandlers) != 0 {
		for _, generalHandler := range e.GeneralHandlers {
			wg.Add(1)
			go func(generalHandler HandlerFunc){
				generalHandler(event)
				wg.Done()
			}(generalHandler)
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
	close(e.Queue)
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
		e.Lock()
		e.Handlers[Type] = append(value, handler)
		e.Unlock()
	}

}

func (e *EventEngine) Unregister(Type string, handler HandlerFunc) {
	if value, ok := e.Handlers[Type]; ok {
		for index, item := range value {
			if reflect.ValueOf(item).Pointer() == reflect.ValueOf(handler).Pointer() {
				
				e.Lock()
				if index == len(value)-1 {
					e.Handlers[Type] = value[:index]
					return
				}
				e.Handlers[Type] = append(value[:index], value[index+1:]...)
				e.Unlock()
			}
		}
	}

}
