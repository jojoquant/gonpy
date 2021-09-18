package database

import (
	"gonpy/trader"
	. "gonpy/trader/object"
	"time"
)

type BarData struct {
	BaseData
	Datetime     time.Time `bson:"datetime"`
	Interval     trader.Interval `bson:"interval"`
	
	Open         float64   `bson:"open_price"`
	High         float64   `bson:"high_price"`
	Low          float64   `bson:"low_price"`
	Close        float64   `bson:"close_price"`
	OpenInterest float64   `bson:"open_interest"`
	Volume       float64   `bson:"volume"`
	
	// 以下字段在 BaseData 中
	// Symbol       string    
	// Exchange    trader.Exchange `bson:"exchange"`
}
