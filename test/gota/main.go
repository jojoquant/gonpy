package main

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func main() {
	df := dataframe.New(
		series.New([]string{"b","a"}, series.String, "col1"),
		series.New([]int{1,2}, series.Int, "col2"),
		series.New([]float64{3.0,4.0}, series.Float, "col3"),
	)

	fmt.Println(df)
}