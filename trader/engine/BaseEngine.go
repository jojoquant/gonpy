package engine

import (
	. "gonpy/trader/object"
	"gonpy/trader/strategy"
	. "gonpy/trader"
)

type BaseEnginer interface {
	Close()
	GetName() string
	SetEventEngine(eventEngine *EventEngine)
}

type TradeEnginer interface{
	SendStopOrder(
		strategy *strategy.Strategy, contract *ContractData,
		direction Direction, offset Offset,
		price, volume float64, lock bool,
	) string
}

type BaseEngine struct {
	Name        string
	EventEngine *EventEngine
}

