package object

import (
	"fmt"
	"time"
	. "gonpy/trader"
)

type TradeData struct {
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
) *TradeData {

	trade := &TradeData{
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
	trade.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)

	trade.Gateway = gateway
	trade.VtOrderId = fmt.Sprintf("%s.%s", gateway, orderId)
	trade.VtTradeId = fmt.Sprintf("%s.%s", gateway, tradeId)

	return trade
}