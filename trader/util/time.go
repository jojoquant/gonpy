package util

import (
	"log"
	"time"
)

func FuncExecDuration(funcName string, f func()) time.Duration{
	start := time.Now()
	f()
	t := time.Since(start)
	log.Printf("Function [%s] exec time cost : %s", funcName, t)
	return t
}