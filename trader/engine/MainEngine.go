package engine

import (
	"gonpy/trader/gateway"
)

type MainEngine struct {
	event_engine EventEngine
	MapGateways map[string]gateway.BaseGateway
	MapEngines map[string]*BaseEngine
}

func NewMainEngine(event_engine EventEngine) *MainEngine {
	m := MainEngine{event_engine: event_engine}
	return &m
}

func(m *MainEngine)AddEngine(b *BaseEngine)*BaseEngine{
	m.MapEngines[b.name] = b
	return b
}

