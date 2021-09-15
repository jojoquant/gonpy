package object

import (
	"fmt"
	"time"
	. "gonpy/trader"
)


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