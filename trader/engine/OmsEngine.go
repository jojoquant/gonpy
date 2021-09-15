package engine

import (
	"fmt"
	. "gonpy/trader/object"
	"gonpy/trader/database"
)

type OmsEngine struct {
	// BaseEnginer
	BaseEngine
	Ticks map[string]*database.TickData
	Orders map[string]*OrderData
	Trades map[string]*TradeData
}

func (o *OmsEngine) Close() {
	fmt.Println("Oms engine close")
}

func (o *OmsEngine) GetName() string {
	return o.Name
}

func (o *OmsEngine) SetEventEngine(eventEngine *EventEngine) {
	o.EventEngine = eventEngine
}