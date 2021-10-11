package converter

import (
	"gonpy/trader"
	"gonpy/trader/object"
	"strings"
)

type PositionHolding struct {
	VtSymbol string
	Exchange trader.Exchange

	ActiveOrders map[string]*object.OrderData

	LongPos float64
	LongYd  float64
	LongTd  float64

	ShortPos float64
	ShortYd  float64
	ShortTd  float64

	LongPosFrozen float64
	LongYdFrozen  float64
	LongTdFrozen  float64

	ShortPosFrozen float64
	ShortYdFrozen  float64
	ShortTdFrozen  float64
}

func NewPositionHolding(contract object.ContractData)*PositionHolding{
	p := &PositionHolding{
		VtSymbol: contract.VtSymbol,
		Exchange: contract.Exchange,
		ActiveOrders: make(map[string]*object.OrderData),
	}
	return p
}

func (p *PositionHolding) UpdatePosition(position *object.PositionData) {
	if position.Direction == trader.LONG {
		p.LongPos = position.Volume
		p.LongYd = position.YdVolume
		p.LongTd = p.LongPos - p.LongYd
	} else if position.Direction == trader.SHORT {
		p.ShortPos = position.Volume
		p.ShortYd = position.YdVolume
		p.ShortTd = p.ShortPos - p.ShortYd
	}
}

func (p *PositionHolding) UpdateOrder(order *object.OrderData) {
	if order.IsActive {
		p.ActiveOrders[order.VtOrderId] = order
	} else {
		delete(p.ActiveOrders, order.VtOrderId)
	}

	p.CalculateFrozen()
}

func (p *PositionHolding) UpdateOrderRequest(req *object.OrderRequest, vtOrderId string) {
	vtOrderIdSlice := strings.Split(vtOrderId, ".")
	gateway, orderId := vtOrderIdSlice[0], vtOrderIdSlice[1]
	p.UpdateOrder(req.CreateOrderData(orderId, gateway))
}

func (p *PositionHolding) UpdateTrade(trade *object.TradeData) {
	if trade.Direction == trader.LONG {

		switch trade.Offset {
		case trader.OPEN:
			p.LongTd += trade.Volume
		case trader.CLOSETODAY:
			p.ShortTd -= trade.Volume
		case trader.CLOSEYESTERDAY:
			p.ShortYd -= trade.Volume
		case trader.CLOSE:
			if (trade.Exchange == trader.SHFE) || (trade.Exchange == trader.INE) {
				p.ShortYd -= trade.Volume
			} else {
				p.ShortTd -= trade.Volume

				if p.ShortTd < 0 {
					p.ShortYd += p.ShortTd
					p.ShortTd = 0
				}
			}
		}

	} else {
		switch trade.Offset {
		case trader.OPEN:
			p.ShortTd += trade.Volume
		case trader.CLOSETODAY:
			p.LongTd -= trade.Volume
		case trader.CLOSEYESTERDAY:
			p.LongYd -= trade.Volume
		case trader.CLOSE:
			if (trade.Exchange == trader.SHFE) || (trade.Exchange == trader.INE) {
				p.LongYd -= trade.Volume
			} else {
				p.LongTd -= trade.Volume

				if p.LongTd < 0 {
					p.LongYd += p.LongTd
					p.LongTd = 0
				}
			}
		}
	}

	p.LongPos = p.LongTd + p.LongYd
	p.ShortPos = p.ShortTd + p.ShortYd
}

func (p *PositionHolding) CalculateFrozen() {
	p.LongPosFrozen, p.LongYdFrozen, p.LongTdFrozen = 0, 0, 0
	p.ShortPosFrozen, p.ShortYdFrozen, p.ShortTdFrozen = 0, 0, 0

	for _, order := range p.ActiveOrders {
		if order.Offset == trader.OPEN {
			continue
		}

		frozen := order.Volume - order.Traded

		if order.Direction == trader.LONG {
			switch order.Offset {
			case trader.CLOSETODAY:
				p.ShortTdFrozen += frozen
			case trader.CLOSEYESTERDAY:
				p.ShortYdFrozen += frozen
			case trader.CLOSE:
				p.ShortTdFrozen += frozen
				if p.ShortTdFrozen > p.ShortTd {
					p.ShortYdFrozen += (p.ShortTdFrozen - p.ShortTd)
					p.ShortTdFrozen = p.ShortTd
				}
			}
		} else if order.Direction == trader.SHORT {
			switch order.Offset {
			case trader.CLOSETODAY:
				p.LongTdFrozen += frozen
			case trader.CLOSEYESTERDAY:
				p.LongYdFrozen += frozen
			case trader.CLOSE:
				p.LongTdFrozen += frozen
				if p.LongTdFrozen > p.LongTd {
					p.LongYdFrozen += (p.LongTdFrozen - p.LongTd)
					p.LongTdFrozen = p.LongTd
				}
			}
		}

		p.LongPosFrozen = p.LongTdFrozen + p.LongYdFrozen
		p.ShortPosFrozen = p.ShortTdFrozen + p.ShortYdFrozen
	}
}

func (p *PositionHolding) ConvertOrderRequestSHFE(req *object.OrderRequest)[]*object.OrderRequest{
	if req.Offset == trader.OPEN{
		return []*object.OrderRequest{req,}
	}

	return nil
}

func (p *PositionHolding) ConvertOrderRequestLock(req *object.OrderRequest)[]*object.OrderRequest{
	return nil
}

func (p *PositionHolding) ConvertOrderRequestNet(req *object.OrderRequest)[]*object.OrderRequest{
	return nil
}

