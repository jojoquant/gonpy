package strategy

import (
	"gonpy/trader"
)

type TradeEnginer interface {
	SendOrder(
		strategy *Strategy, 
		direction trader.Direction,
		offset trader.Offset, 
		price, volume float64,
		stop, lock bool,
	)string
	CancelOrder(strategy *Strategy, vtOrderId string)
	CancelAll(strategy *Strategy)

	SendStopOrder(
		strategy *Strategy, contract *trader.ContractData, 
		direction trader.Direction, offset trader.Offset,
		price, volume float64, lock bool,
	)string

	SendLimitOrder(
		strategy *Strategy, contract *trader.ContractData, 
		direction trader.Direction, offset trader.Offset,
		price, volume float64, lock bool,
	)string
}


type Strategy struct {
	Author   string

	// vnpy中的 cta_engine
	TradeEngine TradeEnginer
	Name     string
	VtSymbol string

	Inited  bool
	Trading bool
	Pos     float64
}

func (s *Strategy) OnInit()      {}
func (s *Strategy) OnStart()     {}
func (s *Strategy) OnStop()      {}
func (s *Strategy) OnTick()      {}
func (s *Strategy) OnBar()       {}
func (s *Strategy) OnTrade()     {}
func (s *Strategy) OnOrder()     {}
func (s *Strategy) OnStopOrder() {}

func (s *Strategy) Buy()   {}
func (s *Strategy) Sell()  {}
func (s *Strategy) Short() {}
func (s *Strategy) Cover() {}

func (s *Strategy) SendOrder(
	direction trader.Direction, offset trader.Offset,
	price, volume float64,
	stop, lock bool)string {
	
		if s.Trading{
		vtOrderId := s.TradeEngine.SendOrder(s, direction, offset,price,volume,stop,lock)
		return vtOrderId
	}

	return ""
}

func (s *Strategy) CancelOrder(vtOrderId string) {
	if s.Trading{
		s.TradeEngine.CancelOrder(s, vtOrderId)
	}
}
func (s *Strategy) CancelAll()   {
	if s.Trading{
		s.TradeEngine.CancelAll(s)
	}
}

func (s *Strategy) WriteLog() {}

func (s *Strategy) GetEngineType() {}
func (s *Strategy) GetPriceTick()  {}

func (s *Strategy) LoadBar()  {}
func (s *Strategy) LoadTick() {}

func (s *Strategy) PutEvent() {}

func (s *Strategy) SendEmail() {}
func (s *Strategy) SyncData()  {}
