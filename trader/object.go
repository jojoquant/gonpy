package trader

import (
	"fmt"
	"time"
)

type BaseData struct {
	Gateway  string
	Symbol   string
	Exchange Exchange
	VtSymbol string // "symbol.exchange"
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

func NewOrderData(
	gateway, symbol string, exchange Exchange,
	orderId string, orderType OrderType, direction Direction, offset Offset,
	price, volume, traded float64, status Status,
	datetime time.Time,
) *OrderData {
	order := &OrderData{
		OrderId:   orderId,
		OrderType: orderType,
		Direction: direction,
		Offset:    offset,

		Price:    price,
		Volume:   volume,
		Traded:   traded,
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

type ContractData struct {
	BaseData
	Name      string
	Product   Product
	Size      float64
	PriceTick float64

	MinVolume   float64
	StopSupport bool
	NetPosition bool
	HistoryData bool

	OptionStrike     float64
	OptionUnderlying string
	OptionType       OptionType
	OptionExpiry     time.Time
	OptionPortfolio  string
	OptionIndex      string
}

func NewContractData(gateway, symbol string,
	exchange Exchange, direction Direction, offset Offset,
	price, volume float64,
) *ContractData {

	contract := &ContractData{}
	contract.Symbol = symbol
	contract.Exchange = exchange
	contract.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)
	contract.Gateway = gateway
	return contract
}
