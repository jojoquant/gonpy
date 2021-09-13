package trader

import (
	"fmt"
	"time"
)

type BaseData struct {
	Gateway string
	Symbol   string
	Exchange Exchange
	VtSymbol  string // "symbol.exchange"	
}

type LogData struct {
	BaseData
	Msg   string
	level LogLevel
	time  time.Time
}

func NewLogData(msg, gatewayName string) *LogData {
	l := &LogData{Msg: msg, level: INFO}
	l.Gateway = gatewayName
	l.time = time.Now()
	return l
}

type OrderData struct {
	BaseData
	VtOrderId string // "gateway.OrderId"

	OrderId   string
	OrderType OrderType
	Direction Direction
	Offset    Offset

	Price    float64
	Volume   float64
	Traded   float64
	Status   Status
	Datetime time.Time

	IsActive bool
}

type TraderData struct {
	BaseData
	VtOrderId string // "gateway.OrderId"
	VtTradeId string // "gateway.TradeId"

	OrderId string
	TradeId string

	Direction Direction
	Offset    Offset

	Price    float64
	Volume   float64
	Datetime time.Time
}

func NewTradeData(
	gateway, symbol, orderId, tradeId string,
	exchange Exchange, direction Direction, offset Offset,
	price, volume float64,
	datetime time.Time,
) *TraderData {

	trade := &TraderData{
		OrderId:   orderId,
		TradeId:   tradeId,
		Direction: direction,
		Offset:    offset,
		Price:     price,
		Volume:    volume,
		Datetime:  datetime,
	}
	trade.Symbol = symbol
	trade.Exchange = exchange
	trade.VtSymbol= fmt.Sprintf("%s.%s", symbol, exchange)
	
	trade.Gateway = gateway
	trade.VtOrderId = fmt.Sprintf("%s.%s", gateway, orderId)
	trade.VtTradeId = fmt.Sprintf("%s.%s", gateway, tradeId)

	return trade
}
