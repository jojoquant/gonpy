package main

import (
	"fmt"
	"gonpy/trader/engine/BacktestEngine"
	"time"
)

func main() {

	r := BacktestEngine.NewDailyResult("dsad", 12.3)
	fmt.Println(r)

	start, _ := time.Parse("2006-01-02 15:04:05", "2019-10-13 12:00:33")
	end, _ := time.Parse("2006-01-02", "2020-10-13")
	y, m, d := time.Now().Date()
	
	fmt.Println(y, m, d, time.Now().Format("2006-01-02"))

	p := BacktestEngine.Parameters{
		Start: start,
		End:   end,
	}

	b := BacktestEngine.NewBacktestEngine(p)
	b.LoadData()
}
