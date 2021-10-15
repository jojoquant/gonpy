package object

import (
	"fmt"
	"gonpy/trader"
)

type AccountData struct{
	BaseData
	AccountId string
	VtAccountId string

	Balance float64
	Frozen float64

	Available float64
}

func NewAccountData(
	gateway, symbol string, exchange trader.Exchange,
	accountId string, balance, frozen float64)*AccountData{
	a := &AccountData{
		AccountId: accountId,
		Balance: balance,
		Frozen: frozen,
		Available: balance-frozen,
	}

	a.Symbol = symbol
	a.Exchange = exchange
	a.VtSymbol = fmt.Sprintf("%s.%s", symbol, exchange)
	a.Gateway = gateway
	a.VtAccountId = fmt.Sprintf("%s.%s", gateway, accountId)

	return a
}