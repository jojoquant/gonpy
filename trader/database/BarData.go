package database

import (
	"time"
	"gonpy/trader"
)

type BarData struct {
	trader.BaseData
	Open         float64
	High         float64
	Low          float64
	Close        float64
	OpenInterest float64
	Volume       float64
	Datetime     time.Time
}