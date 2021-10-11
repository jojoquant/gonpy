package object

import (
	"fmt"
	. "gonpy/trader"
	"time"
)

type OrderRequest struct {
	BaseData
	Direction Direction
	Type      OrderType
	Volume    float64
	Price     float64
	Offset    Offset
	Reference string
}

func NewOrderRequest(
	gateway, symbol string, exchange Exchange,
	direction Direction, offset Offset,
	orderType OrderType,
	price, volume float64, reference string) *OrderRequest {

	o := &OrderRequest{
		Direction: direction,
		Type: orderType,
		Volume: volume,
		Price: price,
		Offset: offset,
		Reference: reference,
	}
	o.Gateway = gateway
	o.Symbol = symbol
	o.Exchange = exchange
	o.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)

	return o
}

func (o *OrderRequest)CreateOrderData(orderId string, gateway string)*OrderData{
	od := NewOrderData(gateway, o.Symbol,o.Exchange, orderId,o.Direction,o.Offset,o.Price,o.Volume,SUBMITTING,time.Now())
	return od
}
