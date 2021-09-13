package engine

import (
	"fmt"
	"gonpy/trader"
	"gonpy/trader/database"
	"gonpy/trader/util"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Parameters struct {
	Symbol     string
	VtSymbol   string
	Exchange   string
	Start      time.Time
	End        time.Time
	Rate       float64
	Slippage   float64
	Size       float64
	PriceTick  float64
	Capital    float64
	RiskFree   float64
	AnnualDays int

	Mode    trader.BacktestMode
	Inverse bool
}

type BacktestEngine struct {
	BaseEngine
	Parameters

	Gateway string

	Strategy        string
	Bar             *database.BarCollection
	Tick            *database.TickCollection
	Database        *database.MongoDB
	HistoryBarData  []*database.BarCollection
	HistoryTickData []*database.TickCollection

	ActiveLimitOrder map[string]*trader.OrderData
	LimitOrders      map[string]*trader.OrderData
	LimitOrderCount  int

	TradeCount int
	Trades     map[string]*trader.TraderData
}

func NewBacktestEngine(param Parameters) *BacktestEngine {
	b := &BacktestEngine{}
	b.Parameters = param
	b.Gateway = "BacktestEngine"
	return b
}

func (b *BacktestEngine) Close() {}

func (b *BacktestEngine) GetName() string {
	return b.Name
}

func (b *BacktestEngine) SetEventEngine(eventEngine *EventEngine) {}

func (b *BacktestEngine) LoadData() {
	if b.Start.After(b.End) {
		log.Fatalln("起始日期必须小于结束日期")
		return
	}

	log.Println("加载数据: ", b.Start, " -> ", b.End)
	util.FuncExecDuration(func() { time.Sleep(2 * time.Second) })

	m := database.NewMongoDB("192.168.0.113", 27017)
	b.HistoryBarData = m.Query(
		&database.QueryParam{
			Db:         "vnpy",
			Collection: "SHFE_d_AUL8",
			Filter:     bson.D{{}},
		},
	)
}

func (b *BacktestEngine) AddStrategy() {}

func (b *BacktestEngine) NewBar(bar *database.BarCollection) {
	b.Bar = bar
	b.CrossLimitOrder()
}

func (b *BacktestEngine) CrossLimitOrder() {
	var longCrossPrice float64
	var shortCrossPrice float64
	var longBestPrice float64
	var shortBestPrice float64

	if b.Mode == trader.BarMode {
		longCrossPrice = b.Bar.Low
		shortCrossPrice = b.Bar.High
		longBestPrice = b.Bar.Open
		shortBestPrice = b.Bar.Open
	} else if b.Mode == trader.TickMode {
		log.Println("Tick mode TODO")
	}

	for _, order := range b.ActiveLimitOrder {
		if order.Status == trader.SUBMITTING {
			order.Status = trader.NOTTRADED
			// TODO 传入策略中的响应函数中
			// b.strategy.OnOrder(order)
		}

		// Check whether limit order can be filled
		longCross := (order.Direction == trader.LONG && order.Price >= longCrossPrice && longCrossPrice > 0)
		shortCross := (order.Direction == trader.SHORT && order.Price <= shortCrossPrice && shortCrossPrice > 0)

		if !longCross && !shortCross {
			continue
		}

		// Push order update with status "all traded" (filled)
		order.Traded = order.Volume
		order.Status = trader.ALLTRADED
		// TODO 传入策略中的响应函数中
		// b.strategy.OnOrder(order)
		delete(b.ActiveLimitOrder, order.VtOrderId)

		b.TradeCount++

		var tradePrice, posChange float64
		if longCross {
			tradePrice = math.Min(order.Price, longBestPrice)
			posChange = order.Volume
		} else if shortCross {
			tradePrice = math.Max(order.Price, shortBestPrice)
			posChange = -order.Volume
		}

		trade := trader.NewTradeData(
			b.Gateway, order.Symbol, order.OrderId,
			fmt.Sprintf("%d", b.TradeCount),
			order.Exchange, order.Direction, order.Offset,
			tradePrice, order.Volume, b.Bar.Datetime,
		)

		//TODO strategy.pos+=poschange
		log.Println(posChange)
		//TODO strategy.OnTrade(trade)
		b.Trades[trade.VtTradeId] = trade
	}

}
