package engine

import (
	"fmt"
	"gonpy/trader"
	"gonpy/trader/database"
	"gonpy/trader/object"
	"log"
)

type OmsEngine struct {
	// BaseEnginer
	BaseEngine
	Ticks     map[string]*database.TickData
	Orders    map[string]*object.OrderData
	Trades    map[string]*object.TradeData
	Positions map[string]*object.PositionData
	Accounts  map[string]*object.AccountData
	Contracts map[string]*object.ContractData
	Quotes    map[string]*object.QuoteData

	ActiveOrders map[string]*object.OrderData
	ActiveQuotes map[string]*object.QuoteData
}

func NewOmsEngine(e *EventEngine) *OmsEngine {
	o := &OmsEngine{
		Ticks:        make(map[string]*database.TickData),
		Orders:       make(map[string]*object.OrderData),
		Trades:       make(map[string]*object.TradeData),
		Positions:    make(map[string]*object.PositionData),
		Accounts:     make(map[string]*object.AccountData),
		Contracts:    make(map[string]*object.ContractData),
		Quotes:       make(map[string]*object.QuoteData),
		ActiveOrders: make(map[string]*object.OrderData),
		ActiveQuotes: make(map[string]*object.QuoteData),
	}

	o.Name = "oms"
	o.EventEngine = e

	o.EventEngine.Register(trader.EVENT_TICK, o.ProcessTickEvent)

	return o
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

func (o *OmsEngine) ProcessTickEvent(event *trader.Event) {
	tick, ok := event.Data.(*database.TickData)
	if !ok {
		log.Fatalln("TickData type assertion fail!")
		return
	}
	o.Ticks[tick.VtSymbol] = tick
}

func (o *OmsEngine) ProcessOrderEvent(event *trader.Event) {
	order, ok := event.Data.(*object.OrderData)
	if !ok {
		log.Fatalln("OrderData type assertion fail!")
		return
	}

	o.Orders[order.VtOrderId] = order

	// if order is active, then update data in map
	if order.IsActive {
		o.ActiveOrders[order.OrderId] = order
	} else {
		delete(o.ActiveOrders, order.VtOrderId)
	}
}

func (o *OmsEngine) ProcessTradeEvent(event *trader.Event) {
	trade, ok := event.Data.(*object.TradeData)
	if !ok {
		log.Fatalln("TradeData type assertion fail!")
		return
	}
	o.Trades[trade.VtTradeId] = trade
}

func (o *OmsEngine) ProcessPositionEvent(event *trader.Event) {
	position, ok := event.Data.(*object.PositionData)
	if !ok {
		log.Fatalln("PositionData type assertion fail!")
		return
	}
	o.Positions[position.VtPositionId] = position
}

func (o *OmsEngine) ProcessAccountEvent(event *trader.Event) {
	account, ok := event.Data.(*object.AccountData)
	if !ok {
		log.Fatalln("AccountData type assertion fail!")
		return
	}
	o.Accounts[account.VtAccountId] = account
}

func (o *OmsEngine) ProcessContractEvent(event *trader.Event) {
	contract, ok := event.Data.(*object.ContractData)
	if !ok {
		log.Fatalln("ContractData type assertion fail!")
		return
	}
	o.Contracts[contract.VtSymbol] = contract
}

func (o *OmsEngine) ProcessQuoteEvent(event *trader.Event) {
	quote, ok := event.Data.(*object.QuoteData)
	if !ok {
		log.Fatalln("QuoteData type assertion fail!")
		return
	}
	o.Quotes[quote.VtQuoteId] = quote

	// if quote is active, then update data in map
	if quote.IsActive {
		o.ActiveQuotes[quote.VtQuoteId] = quote
	} else {
		delete(o.ActiveQuotes, quote.VtQuoteId)
	}
}

func (o *OmsEngine) GetTick(vtSymbol string) *database.TickData {
	if t, ok := o.Ticks[vtSymbol]; ok {
		return t
	}
	return nil
}

func (o *OmsEngine) GetOrder(vtOrderId string) *object.OrderData {
	if v, ok := o.Orders[vtOrderId]; ok {
		return v
	}
	return nil
}

func (o *OmsEngine) GetTrade(vtTradeId string) *object.TradeData {
	if v, ok := o.Trades[vtTradeId]; ok {
		return v
	}
	return nil
}

func (o *OmsEngine) GetPosition(vtPositionId string) *object.PositionData {
	if v, ok := o.Positions[vtPositionId]; ok {
		return v
	}
	return nil
}

func (o *OmsEngine) GetAccount(vtAccountId string) *object.AccountData {
	if v, ok := o.Accounts[vtAccountId]; ok {
		return v
	}
	return nil
}

func (o *OmsEngine) GetContract(vtSymbol string) *object.ContractData {
	if v, ok := o.Contracts[vtSymbol]; ok {
		return v
	}
	return nil
}

func (o *OmsEngine) GetQuote(VtQuoteId string) *object.QuoteData {
	if v, ok := o.Quotes[VtQuoteId]; ok {
		return v
	}
	return nil
}

func (o *OmsEngine) GetAllTicks() []*database.TickData {
	vs := make([]*database.TickData, 0, len(o.Ticks))
	for _, v := range o.Ticks {
		vs = append(vs, v)
	}
	return vs
}

func (o *OmsEngine) GetAllOrders() []*object.OrderData {
	vs := make([]*object.OrderData, 0, len(o.Orders))
	for _, v := range o.Orders {
		vs = append(vs, v)
	}
	return vs
}

func (o *OmsEngine) GetAllTrades() []*object.TradeData {
	vs := make([]*object.TradeData, 0, len(o.Trades))
	for _, v := range o.Trades {
		vs = append(vs, v)
	}
	return vs
}

func (o *OmsEngine) GetAllPositions() []*object.PositionData {
	vs := make([]*object.PositionData, 0, len(o.Positions))
	for _, v := range o.Positions {
		vs = append(vs, v)
	}
	return vs
}

func (o *OmsEngine) GetAllAccounts() []*object.AccountData {
	vs := make([]*object.AccountData, 0, len(o.Accounts))
	for _, v := range o.Accounts {
		vs = append(vs, v)
	}
	return vs
}

func (o *OmsEngine) GetAllContracts() []*object.ContractData {
	vs := make([]*object.ContractData, 0, len(o.Contracts))
	for _, v := range o.Contracts {
		vs = append(vs, v)
	}
	return vs
}

func (o *OmsEngine) GetAllQuotes() []*object.QuoteData {
	vs := make([]*object.QuoteData, 0, len(o.Quotes))
	for _, v := range o.Quotes {
		vs = append(vs, v)
	}
	return vs
}

func (o *OmsEngine) GetAllActiveOrders() {}
func (o *OmsEngine) GetAllActiveQuotes() {}
