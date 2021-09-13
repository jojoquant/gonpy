package util

import (
	"log"
	"time"
)

func FuncExecDuration(f func()) time.Duration{
	start := time.Now()
	f()
	t := time.Since(start)
	log.Println("Function exec time cost :", t)
	return t
}