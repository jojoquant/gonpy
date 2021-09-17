package strategy

import (
	. "gonpy/trader"
	"gonpy/trader/database"
	. "gonpy/trader/object"
)

type TickCallback func(*database.TickData)
type BarCallback func(*database.BarData)

type TradeEnginer interface {
	LoadBar(vtSymbol string, days int, interval Interval, callback BarCallback, useDatabase bool)
	SendOrder(
		strategy *Strategy,
		direction Direction,
		offset Offset,
		price, volume float64,
		stop, lock bool,
	) string
	CancelOrder(strategy *Strategy, vtOrderId string)
	CancelLimitOrder(strategy *Strategy, vtOrderId string)
	CancelStopOrder(strategy *Strategy, vtOrderId string)
	CancelAll(strategy *Strategy)

	SendStopOrder(
		strategy *Strategy, contract *ContractData,
		direction Direction, offset Offset,
		price, volume float64, lock bool,
	) string

	SendLimitOrder(
		strategy *Strategy, contract *ContractData,
		direction Direction, offset Offset,
		price, volume float64, lock bool,
	) string
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

func (s *Strategy) OnInit()                              {}
func (s *Strategy) OnStart()                             {}
func (s *Strategy) OnStop()                              {}
func (s *Strategy) OnTick(*database.TickData)            {}
func (s *Strategy) OnBar(*database.BarData)              {}
func (s *Strategy) OnTrade(trade *TradeData)             {}
func (s *Strategy) OnOrder(order *OrderData)             {}
func (s *Strategy) OnStopOrder(stopOrder *StopOrderData) {}

func (s *Strategy) Buy(price, volume float64, stop, lock bool) string {
	return s.SendOrder(LONG, OPEN, price, volume, stop, lock)
}
func (s *Strategy) Sell(price, volume float64, stop, lock bool) string {
	return s.SendOrder(SHORT, CLOSE, price, volume, stop, lock)
}
func (s *Strategy) Short(price, volume float64, stop, lock bool) string {
	return s.SendOrder(SHORT, OPEN, price, volume, stop, lock)
}
func (s *Strategy) Cover(price, volume float64, stop, lock bool) string {
	return s.SendOrder(LONG, CLOSE, price, volume, stop, lock)
}

func (s *Strategy) SendOrder(
	direction Direction, offset Offset,
	price, volume float64,
	stop, lock bool) string {

	if s.Trading {
		vtOrderId := s.TradeEngine.SendOrder(s, direction, offset, price, volume, stop, lock)
		return vtOrderId
	}

	return ""
}

func (s *Strategy) CancelOrder(vtOrderId string) {
	if s.Trading {
		s.TradeEngine.CancelOrder(s, vtOrderId)
	}
}
func (s *Strategy) CancelAll() {
	if s.Trading {
		s.TradeEngine.CancelAll(s)
	}
}

func (s *Strategy) WriteLog() {}

func (s *Strategy) GetEngineType() {}
func (s *Strategy) GetPriceTick()  {}

func (s *Strategy) LoadBar(days int, interval Interval, callback BarCallback) {
	if callback == nil {
		callback = s.OnBar
	}
	s.TradeEngine.LoadBar(s.VtSymbol, days, interval, callback, false)
}
func (s *Strategy) LoadTick(days int, callback TickCallback) {}

func (s *Strategy) PutEvent() {}

func (s *Strategy) SendEmail() {}
func (s *Strategy) SyncData()  {}
