package database

import (
	. "gonpy/trader/object"
	"time"
)

type TickData struct {
	BaseData

	Name string

	Open     float64
	High     float64
	Low      float64
	PreClose float64

	LastPrice  float64
	LastVolume float64
	LimitUp    float64
	LimitDown  float64

	OpenInterest float64
	Volume       float64
	Turnover     float64
	Datetime     time.Time
	LocalTime    time.Time

	BidPrice1 float64
	BidPrice2 float64
	BidPrice3 float64
	BidPrice4 float64
	BidPrice5 float64

	AskPrice1 float64
	AskPrice2 float64
	AskPrice3 float64
	AskPrice4 float64
	AskPrice5 float64

	BidVolume1 float64
	BidVolume2 float64
	BidVolume3 float64
	BidVolume4 float64
	BidVolume5 float64

	AskVolume1 float64
	AskVolume2 float64
	AskVolume3 float64
	AskVolume4 float64
	AskVolume5 float64
}
