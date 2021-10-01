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

	start, _ := time.Parse("2006-01-02 15:04:05", "2013-08-25 12:00:33")
	end, _ := time.Parse("2006-01-02 15:04:05", "2021-02-25 15:00:00")
	
	// , m, d := time.Now().Date()
	// fmt.Println(y, m, d, time.Now().Format("2006-01-02"))

	p := BacktestEngine.NewParameters(
		"LL8", "DCE", start, end, 0.3/10000, 1, 10, 1, 
		math.Pow(10, 6), 0, 
		trader.BarMode, trader.MINUTE, false)

	db := database.NewMongoDB("192.168.0.113", 27017)
	
	// s := &strategy.Strategy{}
	s := strategy.NewDualMA()
	b := BacktestEngine.NewBacktestEngine(p, db, s)

	util.FuncExecDuration("LoadData", b.LoadData)
	// b.LoadData()
	
	// b.AddStrategy()
	
	util.FuncExecDuration("RunBacktest", b.RunBacktest)
	// b.RunBacktest()
	sMap := b.CalculateResult()
	log.Println(sMap)
}
