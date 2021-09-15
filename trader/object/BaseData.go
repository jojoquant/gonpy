package object

import (
	. "gonpy/trader"
)

type BaseData struct {
	Gateway  string
	Symbol   string
	Exchange Exchange
	VtSymbol string // "symbol.exchange"
}




