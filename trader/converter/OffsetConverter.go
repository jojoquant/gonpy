package converter

import (
	"gonpy/trader"
	"gonpy/trader/engine"
	"gonpy/trader/object"
)

type OffsetConverter struct {
	MainEngine *engine.MainEngine
	Holdings   map[string]*PositionHolding
}

func (o *OffsetConverter) IsConvertRequired(vtSymbol string) bool {
	if contract := o.MainEngine.OmsEngine.GetContract(vtSymbol); contract != nil {
		if contract.NetPosition {
			return false
		} else {
			return true
		}
	}
	return false
}

func (o *OffsetConverter) UpdatePosition(position *object.PositionData) {
	if !o.IsConvertRequired(position.VtSymbol) {
		return
	}

	holding := o.GetPositionHolding(position.VtSymbol)
	holding.UpdatePosition(position)
}

func (o *OffsetConverter) UpdateTrade(trade *object.TradeData) {
	if !o.IsConvertRequired(trade.VtSymbol) {
		return
	}

	holding := o.GetPositionHolding(trade.VtSymbol)
	holding.UpdateTrade(trade)
}

func (o *OffsetConverter) UpdateOrder(order *object.OrderData) {
	if !o.IsConvertRequired(order.VtSymbol) {
		return
	}

	holding := o.GetPositionHolding(order.VtSymbol)
	holding.UpdateOrder(order)
}

func (o *OffsetConverter) UpdateOrderRequest(req *object.OrderRequest, vtOrderId string) {
	if !o.IsConvertRequired(req.VtSymbol) {
		return
	}

	holding := o.GetPositionHolding(req.VtSymbol)
	holding.UpdateOrderRequest(req, vtOrderId)
}

func (o *OffsetConverter) GetPositionHolding(vtSymbol string) *PositionHolding {
	holding, ok := o.Holdings[vtSymbol]
	if !ok {
		contract := o.MainEngine.OmsEngine.GetContract(vtSymbol)
		holding = NewPositionHolding(contract)
		o.Holdings[vtSymbol] = holding
	}

	return holding
}

func (o *OffsetConverter) ConvertOrderRequest(req *object.OrderRequest, lock bool, net bool)[]*object.OrderRequest{
	if !o.IsConvertRequired(req.VtSymbol) {
		return []*object.OrderRequest{req,}
	}

	holding := o.GetPositionHolding(req.VtSymbol)

	if lock{
		return holding.ConvertOrderRequestLock(req)
	}else if net{
		return holding.ConvertOrderRequestNet(req)
	}else if (req.Exchange == trader.SHFE) || (req.Exchange == trader.INE){
		return holding.ConvertOrderRequestSHFE(req)
	}else{
		return []*object.OrderRequest{req,}
	}
}
