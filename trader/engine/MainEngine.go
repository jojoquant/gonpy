package engine

import (
	"fmt"
	"gonpy/trader"
	"gonpy/trader/gateway"
	"gonpy/trader/object"
)

type MainEngine struct {
	EventEngine *EventEngine
	OmsEngine   *OmsEngine
	Gateways    map[string]gateway.BaseGateway
	Engines     map[string]BaseEnginer
	Exchanges   []trader.Exchange
}

func NewMainEngine(eventEngine *EventEngine) *MainEngine {
	m := &MainEngine{EventEngine: eventEngine}
	m.EventEngine.Start()

	m.Gateways = make(map[string]gateway.BaseGateway)
	m.Engines = make(map[string]BaseEnginer)
	m.Exchanges = make([]trader.Exchange, 4)

	m.OmsEngine = NewOmsEngine(eventEngine)
	return m
}

func (m *MainEngine) AddEngine(b BaseEnginer) BaseEnginer {
	b.SetEventEngine(m.EventEngine)
	key := b.GetName()
	m.Engines[key] = b
	return b
}

func (m *MainEngine) GetEngine(engineName string) BaseEnginer {
	if e, ok := m.Engines[engineName]; ok {
		return e
	}
	m.WriteLog(fmt.Sprintf("找不到引擎:%v", engineName), "")
	return nil
}

func (m *MainEngine) AddGateway()         {}
func (m *MainEngine) GetGateway()         {}
func (m *MainEngine) GetDefaultSetting()  {}
func (m *MainEngine) GetAllGatewayNames() {}
func (m *MainEngine) Connect()            {}
func (m *MainEngine) Subscribe()          {}
func (m *MainEngine) SendOrder()          {}
func (m *MainEngine) CancelOrder()        {}
func (m *MainEngine) SendOrders()         {}
func (m *MainEngine) CancelOrders()       {}
func (m *MainEngine) QueryHistory()       {}

func (m *MainEngine) GetAllExchanges() []trader.Exchange {
	return m.Exchanges
}

func (m *MainEngine) WriteLog(msg, source string) {
	m.EventEngine.Put(
		&trader.Event{
			Type: trader.EVENT_LOG,
			Data: object.NewLogData(msg, source),
		},
	)
}

func (m *MainEngine) Close() {
	m.EventEngine.Stop()
	for _, engine := range m.Engines {
		engine.Close()
	}
}
