package object

import (
	"time"
	"gonpy/trader"
)

type LogData struct {
	BaseData
	Msg   string
	level trader.LogLevel
	time  time.Time
}

func NewLogData(msg, gatewayName string) *LogData {
	l := &LogData{Msg: msg, level: trader.INFO}
	l.Gateway = gatewayName
	l.time = time.Now()
	return l
}