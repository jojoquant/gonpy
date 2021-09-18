package object

import (
	. "gonpy/trader"
)

type BaseData struct {
	Gateway  string
	Symbol   string `bson:"symbol"`
	Exchange Exchange `bson:"exchange"`
	VtSymbol string // "symbol.exchange"
}




