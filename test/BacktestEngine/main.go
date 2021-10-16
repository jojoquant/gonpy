package main

import (
	"gonpy/trader"
	"gonpy/trader/database"
	"gonpy/trader/engine/BacktestEngine"
	"gonpy/trader/strategy"
	"gonpy/trader/util"
	"log"
	"math"
	"time"
)

func main() {
	// xx :=fmt.Sprintf("%s.%d", trader.STOP, 123)
	// fmt.Println(xx)

	// r := BacktestEngine.NewDailyResult("dsad", 12.3)
	// fmt.Println(r)

	start, _ := time.Parse("2006-01-02 15:04:05", "2021-04-25 12:00:33")
	end, _ := time.Parse("2006-01-02 15:04:05", "2021-10-31 15:00:00")

	// , m, d := time.Now().Date()
	// fmt.Println(y, m, d, time.Now().Format("2006-01-02"))

	p := BacktestEngine.NewParameters(
		"ETH_BTC", "BINANCE", start, end, 0, 0, 1, math.Pow(10, -6),
		math.Pow(10, 6), 0,
		trader.BarMode, trader.MINUTE, false)
	
	// 创建数据源和展示数据库
	host := "localhost"
	dk := util.GetDockerComposeYml("../../docker-compose.yml")
	
	// db := database.NewMongoDB("192.168.0.113", 27017)
	srcDB := database.NewMongoDB(
		host, dk.Services.Mongo.PortSrc,
		dk.Services.Mongo.Env.Username, dk.Services.Mongo.Env.Password)
	
	displayDB:= database.NewInfluxDB(
		host, dk.Services.Influxdb.PortSrc,
		dk.Services.Influxdb.Env.Username, dk.Services.Influxdb.Env.Password,
		dk.Services.Influxdb.Env.Org, dk.Services.Influxdb.Env.Bucket,
		false,
	)

	// s := &strategy.Strategy{}
	s := strategy.NewDualMA()
	b := BacktestEngine.NewBacktestEngine(p, srcDB, displayDB, s)
	q := &database.QueryParam{
		Db:         "binance",
		Collection: "db_bar_data",
	}
	util.FuncExecDuration("LoadData", func() { b.LoadData(q) })
	// b.LoadData()

	// b.AddStrategy()

	util.FuncExecDuration("RunBacktest", func(){b.RunBacktest(true, true)})
	// b.RunBacktest()
	sMap := b.CalculateResult()
	
	srcDB.Close()
	displayDB.Close()
	
	log.Println(sMap)
}
