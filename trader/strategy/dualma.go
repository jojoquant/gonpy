package strategy

import (
	"github.com/markcheno/go-talib"
	"gonpy/trader/database"
	// "log"
)

type DualMA struct {
	Strategy

	fastWindow int
	slowWindow int

	fastMa  []float64
	// fastMa0 float64
	// fastMa1 float64

	slowMa  []float64
	// slowMa0 float64
	// slowMa1 float64

	lenCount    int
	arrayLenMax int
	closeArray  []float64
}

func NewDualMA() *DualMA {
	arrayLenMax := 380000
	d := &DualMA{
		lenCount:    0,
		fastWindow: 10,
		slowWindow: 20,
		arrayLenMax: arrayLenMax,
		closeArray:  make([]float64, 0, arrayLenMax+1),
	}
	d.Name = "DualMA"

	return d
}

func (d *DualMA) OnBar(bar *database.BarData) {
	// log.Println("Dual MA on bar:", bar)

	d.lenCount++
	d.closeArray = append(d.closeArray, bar.Close)
	if d.lenCount < d.slowWindow {
		return
	}
	// fmt.Println("d.closeArray load ok")
	d.fastMa = talib.Sma(d.closeArray, d.fastWindow)
	d.slowMa = talib.Sma(d.closeArray, d.slowWindow)

	crossOver := (d.fastMa[len(d.fastMa)-1] > d.slowMa[len(d.slowMa)-1]) &&
		(d.fastMa[len(d.fastMa)-2] < d.slowMa[len(d.slowMa)-2])

	crossBelow := (d.fastMa[len(d.fastMa)-1] < d.slowMa[len(d.slowMa)-1]) &&
		(d.fastMa[len(d.fastMa)-2] > d.slowMa[len(d.slowMa)-2])

	if crossOver{
		if d.Pos == 0{
			d.Buy(bar.Close, 1, false, false)
		}else if d.Pos<0{
			d.Cover(bar.Close, 1, false, false)
			d.Buy(bar.Close, 1, false, false)
		}
	}else if crossBelow{
		if d.Pos == 0{
			d.Short(bar.Close, 1, false, false)
		}else if d.Pos<0{
			d.Sell(bar.Close, 1, false, false)
			d.Short(bar.Close, 1, false, false)
		}
	}
}
