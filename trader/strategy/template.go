package strategy

import (
	"fmt"
	"gonpy/trader"
	"gonpy/trader/database"
	"gonpy/trader/object"
)

type TickCallback func(*database.TickData)
type BarCallback func(*database.BarData)

type TradeEnginer interface {
	LoadBar(
		vtSymbol string, days int,
		interval trader.Interval,
		callback BarCallback, useDatabase bool)
	SendOrder(
		strategy Strategyer,
		direction trader.Direction,
		offset trader.Offset,
		price, volume float64,
		stop, lock, net bool,
	) string
	CancelOrder(strategy Strategyer, vtOrderId string)
	CancelLimitOrder(strategy Strategyer, vtOrderId string)
	CancelStopOrder(strategy Strategyer, vtOrderId string)
	CancelAll(strategy Strategyer)

	SendStopOrder(
		strategy Strategyer, contract *object.ContractData,
		direction trader.Direction, offset trader.Offset,
		price, volume float64, lock, net bool,
	) string

	SendLimitOrder(
		strategy Strategyer, contract *object.ContractData,
		direction trader.Direction, offset trader.Offset,
		price, volume float64, lock, net bool,
	) string
}

type Strategyer interface {
	SetVtSymbol(string)
	SetInited(bool)
	SetTrading(bool)
	SetTradeEngine(TradeEnginer)
	GetStrategyName() string

	OnInit()
	OnStart()
	OnStop()
	OnTick(*database.TickData)
	OnBar(*database.BarData)
	OnTrade(trade *object.TradeData)
	OnOrder(order *object.OrderData)
	OnStopOrder(stopOrder *object.StopOrderData)
}

type Strategy struct {
	Author string

	// vnpy中的 cta_engine
	TradeEngine TradeEnginer
	Name        string
	VtSymbol    string

	Inited  bool
	Trading bool
	Pos     float64
}

func (s *Strategy) SetVtSymbol(v string) {
	s.VtSymbol = v
}
func (s *Strategy) SetInited(v bool) {
	s.Inited = v
}
func (s *Strategy) SetTrading(v bool) {
	s.Trading = v
}
func (s *Strategy) SetTradeEngine(t TradeEnginer) {
	s.TradeEngine = t
}

func (s *Strategy) GetStrategyName() string {
	return s.Name
}

func (s *Strategy) OnInit() {
	s.LoadBar(1, trader.MINUTE, s.OnBar, false)
}
func (s *Strategy) OnStart()                  {}
func (s *Strategy) OnStop()                   {}
func (s *Strategy) OnTick(*database.TickData) {}
func (s *Strategy) OnBar(b *database.BarData) {
	fmt.Println("Strategy OnBar:", b)
}
func (s *Strategy) OnTrade(trade *object.TradeData)             {}
func (s *Strategy) OnOrder(order *object.OrderData)             {}
func (s *Strategy) OnStopOrder(stopOrder *object.StopOrderData) {}

func (s *Strategy) Buy(price, volume float64, stop, lock bool) string {
	return s.SendOrder(trader.LONG, trader.OPEN, price, volume, stop, lock, false)
}
func (s *Strategy) Sell(price, volume float64, stop, lock bool) string {
	return s.SendOrder(trader.SHORT, trader.CLOSE, price, volume, stop, lock, false)
}
func (s *Strategy) Short(price, volume float64, stop, lock bool) string {
	return s.SendOrder(trader.SHORT, trader.OPEN, price, volume, stop, lock, false)
}
func (s *Strategy) Cover(price, volume float64, stop, lock bool) string {
	return s.SendOrder(trader.LONG, trader.CLOSE, price, volume, stop, lock, false)
}

func (s *Strategy) SendOrder(
	direction trader.Direction, offset trader.Offset,
	price, volume float64,
	stop, lock, net bool) string {

	if s.Trading {
		vtOrderId := s.TradeEngine.SendOrder(s, direction, offset, price, volume, stop, lock, net)
		return vtOrderId
	}

	return ""
}

func (s *Strategy) CancelOrder(vtOrderId string) {
	if s.Trading {
		s.TradeEngine.CancelOrder(s, vtOrderId)
	}
}
func (s *Strategy) CancelAll(strategy Strategyer) {
	if s.Trading {
		s.TradeEngine.CancelAll(s)
	}
}

func (s *Strategy) WriteLog() {}

func (s *Strategy) GetEngineType() {}
func (s *Strategy) GetPriceTick()  {}

func (s *Strategy) LoadBar(days int, interval trader.Interval, callback BarCallback, useDatabase bool) {
	s.TradeEngine.LoadBar(s.VtSymbol, days, interval, callback, useDatabase)
}

func (s *Strategy) LoadTick(days int, callback TickCallback) {}

func (s *Strategy) PutEvent() {}

func (s *Strategy) SendEmail() {}
func (s *Strategy) SyncData()  {}
