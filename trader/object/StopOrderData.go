package object

import (
	"fmt"
	. "gonpy/trader"
	"time"
)

//     Gateway  string
//     Symbol   string
//     Exchange Exchange
//     VtSymbol string   // "symbol.exchange"
//     VtOrderId string  // "gateway.OrderId"
//     StopOrder string  // "STOP.stopOrderCount"
type StopOrderData struct {
	OrderData
	StrategyName string
	StopOrderId  string
	Lock         bool
	Net          bool
}

// Lock, Net 默认为 false, IsActive 默认为 true
// VtOrderId 默认为"." 
// OrderId, Reference 默认为 ""
// traded 默认为 0
// Status 默认为 WAITING = "等待中"
func NewStopOrderData(
	gateway, symbol string, exchange Exchange,
	direction Direction, offset Offset,
	price, volume float64,
	strategyName, stopOrderId string, 
	datetime time.Time,
) *StopOrderData {
	order := &StopOrderData{
		StopOrderId: stopOrderId,
		StrategyName: strategyName,
	}

	order.Direction = direction
	order.Offset = offset
	order.Price = price
	order.Volume = volume
	order.Status = WAITING
	order.IsActive = true
	order.Datetime = datetime

	order.Symbol = symbol
	order.Exchange = exchange
	order.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)

	order.Gateway = gateway
	// order.VtOrderId = fmt.Sprintf("%s.%s", gateway, orderId)
	// order.Reference = reference
	
	// for _, s := range ACTIVE_STATUSES {
	// 	if s == status {
	// 		order.IsActive = true
	// 		break
	// 	}
	// 	order.IsActive = false
	// }

	return order
}
