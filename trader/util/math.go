package util

import "math"

func RoundTo(value, target float64) float64 {
	return math.Round(value/target) * target
}

func FloorTo(value, target float64) float64 {
	return math.Floor(value/target) * target
}

func CeilTo(value, target float64)float64{
	return math.Ceil(value/target) * target
}
