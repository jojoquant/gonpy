package trader

import (
	"time"
)

type BaseData struct {
	GatewayName string
}

type LogData struct {
	BaseData
	Msg   string
	level LogLevel
	time  time.Time
}

func NewLogData(msg, gatewayName string) *LogData {
	l := &LogData{Msg: msg, level: INFO}
	l.GatewayName = gatewayName
	l.time = time.Now()
	return l
}
