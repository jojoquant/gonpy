package object

import (
	"fmt"
	"gonpy/trader"
	"time"
)

type QuoteData struct {
	BaseData
	QuoteId string
	VtQuoteId string

	BidPrice  float64
	BidVolume int
	AskPrice  float64
	AskVolume int
	BidOffset trader.Offset
	AskOffset trader.Offset
	Status trader.Status
	Datetime time.Time
	Reference string

	IsActive bool
}

// is_active 和 create_cancel_request 方法日后实现
// Traded 默认为 0.0
// Reference 默认为 ""
// status 如无特殊情况, 传入SUBMITTING
func NewQuoteData(
	gateway, symbol string, exchange trader.Exchange,
	quoteId string, direction trader.Direction, offset trader.Offset,
	price, volume float64, status trader.Status,
	datetime time.Time,
) *QuoteData {
	quote := &QuoteData{
		QuoteId:   quoteId,
		// OrderType: LIMIT,
		// Direction: direction,
		// Offset:    offset,

		// Price:    price,
		// Volume:   volume,
		Status:   status,
		Datetime: datetime,
	}
	quote.Symbol = symbol
	quote.Exchange = exchange
	quote.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)

	quote.Gateway = gateway
	quote.VtQuoteId = fmt.Sprintf("%s.%s", gateway, quoteId)

	for _, s := range trader.ACTIVE_STATUSES {
		if s == status {
			quote.IsActive = true
			break
		}
		quote.IsActive = false
	}

	return quote
}