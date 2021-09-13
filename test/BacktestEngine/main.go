package main

import (
	"gonpy/trader/engine"
	"time"
)

func main() {

	start, _ := time.Parse("2006-01-02 15:04:05","2019-10-13 12:00:33")
	end, _ := time.Parse("2006-01-02","2020-10-13")
	
	p := engine.Parameters{
		Start: start,
		End: end,
	}

	b := engine.NewBacktestEngine(p)
	b.LoadData()
}