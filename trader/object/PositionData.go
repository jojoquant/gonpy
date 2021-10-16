package object

import (
	"fmt"
	"gonpy/trader"
)

type PositionData struct {
	BaseData
	VtPositionId string // "gateway.OrderId"

	Direction trader.Direction

	Volume   float64
	Frozen   float64
	Price    float64
	Pnl      float64
	YdVolume float64
}

func NewPositionData(
	gateway, symbol string, exchange trader.Exchange,
	direction trader.Direction,
	price, volume, traded float64,
) *PositionData {
	p := &PositionData{
		Direction: direction,

		Price:    price,
		Volume:   volume,

	}
	p.Symbol = symbol
	p.Exchange = exchange
	p.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)

	p.Gateway = gateway
	p.VtPositionId = fmt.Sprintf("%s.%s", p.VtSymbol, direction)

	return p
}
