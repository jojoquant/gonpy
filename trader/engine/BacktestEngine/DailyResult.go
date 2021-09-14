package BacktestEngine

import (
	"gonpy/trader"
)

type DailyResult struct {
	Date       string
	ClosePrice float64
	PreClose   float64

	Trades     []*trader.TraderData
	TradeCount int

	StartPos float64
	EndPos   float64

	Turnover   float64
	Commission float64
	Slippage   float64

	TradingPnl float64
	HoldingPnl float64
	TotalPnl   float64
	NetPnl     float64
}

func NewDailyResult(date string, closePrice float64)*DailyResult{
	d := &DailyResult{
		Date: date,
		ClosePrice: closePrice,
	}
	return d
}

func (d *DailyResult) AddTrade(trade *trader.TraderData) {
	d.Trades = append(d.Trades, trade)
}

func (d *DailyResult) CalculatePnl(
	preClose, startPos, size, rate, slippage float64,
	inverse bool) {
	// if no preClose provided on the first day
	// use value 1 to avoid zero division error
	if preClose > 0 {
		d.PreClose = preClose
	} else {
		d.PreClose = 1
	}

	// Holding pnl is the pnl from holding position at day start
	d.StartPos = startPos
	d.EndPos = startPos

	if !inverse {
		d.HoldingPnl = d.StartPos * (d.ClosePrice - d.PreClose) * size
	} else {
		d.HoldingPnl = d.StartPos * (1/d.PreClose - 1/d.ClosePrice) * size
	}

	d.TradeCount = len(d.Trades)

	var posChange float64
	var turnover float64
	for _, trade := range d.Trades {

		if trade.Direction == trader.LONG {
			posChange = trade.Volume
		} else {
			// 这里包含 SHORT 和 NET
			posChange = -trade.Volume
		}
		d.EndPos += posChange

		if !inverse {
			turnover = trade.Volume * size * trade.Price
			d.TradingPnl += posChange * (d.ClosePrice - trade.Price) * size
			d.Slippage += trade.Volume * size * slippage
		} else {
			turnover = trade.Volume * size / trade.Price
			d.TradingPnl += posChange * (1/trade.Price - 1/d.ClosePrice) * size
			d.Slippage += trade.Volume * size * slippage / (trade.Price * trade.Price)
		}
		d.Turnover += turnover
		d.Commission += turnover*rate
	}

	// Net pnl takes account of commission and slippage cost
	d.TotalPnl = d.TradingPnl + d.HoldingPnl
	d.NetPnl = d.TotalPnl - d.Commission - d.Slippage
}
