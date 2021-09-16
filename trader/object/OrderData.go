package object

import (
	"fmt"
	"time"
	. "gonpy/trader"
)

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
	Reference string

	IsActive bool
}

// OrderType 默认为 LIMIT 
// Traded 默认为 0.0
// Reference 默认为 ""
func NewOrderData(
	gateway, symbol string, exchange Exchange,
	orderId string, direction Direction, offset Offset,
	price, volume float64, status Status,
	datetime time.Time,
) *OrderData {
	order := &OrderData{
		OrderId:   orderId,
		OrderType: LIMIT,
		Direction: direction,
		Offset:    offset,

		Price:    price,
		Volume:   volume,
		Status:   status,
		Datetime: datetime,
	}
	order.Symbol = symbol
	order.Exchange = exchange
	order.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)

	order.Gateway = gateway
	order.VtOrderId = fmt.Sprintf("%s.%s", gateway, orderId)

	for _, s := range ACTIVE_STATUSES {
		if s == status {
			order.IsActive = true
			break
		}
		order.IsActive = false
	}

	return order
}