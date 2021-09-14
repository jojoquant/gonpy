package main

import (
	"fmt"
	"gonpy/trader/util"
)

func main() {
	fmt.Println(util.RoundTo(1000.2, 0.1))
	fmt.Println(util.RoundTo(1000.23, 0.1))
	fmt.Println(util.RoundTo(1000.0, 1))
}