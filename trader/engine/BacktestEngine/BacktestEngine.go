package BacktestEngine

import (
	"fmt"
	. "gonpy/trader"
	"gonpy/trader/database"
	"gonpy/trader/engine"
	. "gonpy/trader/object"
	"gonpy/trader/strategy"
	"gonpy/trader/util"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Parameters struct {
	Symbol     string
	VtSymbol   string
	Exchange   Exchange
	Start      time.Time
	End        time.Time
	Rate       float64
	Slippage   float64
	Size       float64
	PriceTick  float64
	Capital    float64
	RiskFree   float64
	AnnualDays int

	Mode     BacktestMode
	Interval Interval
	Inverse  bool
}

type BacktestEngine struct {
	engine.BaseEngine
	Parameters

	Gateway string

	Strategy *strategy.Strategy
	Bar      *database.BarData
	Tick     *database.TickData
	Datetime time.Time

	Database *database.MongoDB

	Days            int
	Callback        string
	HistoryBarData  []*database.BarData
	HistoryTickData []*database.TickData

	ActiveLimitOrders map[string]*OrderData
	LimitOrders       map[string]*OrderData
	LimitOrderCount   int

	ActiveStopOrders map[string]*StopOrderData
	StopOrders       map[string]*StopOrderData
	StopOrderCount   int

	TradeCount int
	Trades     map[string]*TradeData

	DailyResults map[string]*DailyResult
}

func NewBacktestEngine(param Parameters) *BacktestEngine {
	b := &BacktestEngine{}
	b.Parameters = param
	b.Gateway = "BacktestEngine"
	return b
}

func (b *BacktestEngine) Close() {
	fmt.Println("backtest engine close")
}

func (b *BacktestEngine) GetName() string {
	return b.Name
}

func (b *BacktestEngine) SetEventEngine(eventEngine *engine.EventEngine) {
	b.EventEngine = eventEngine
}

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

func (b *BacktestEngine) NewBar(bar *database.BarData) {
	b.Bar = bar
	b.Datetime = bar.Datetime

	b.CrossLimitOrder()
	b.CrossStopOrder()
	// b.Strategy.OnBar(bar)

	b.UpdateDailyClose(bar.Close)
}

func (b *BacktestEngine) CrossLimitOrder() {
	var longCrossPrice float64
	var shortCrossPrice float64
	var longBestPrice float64
	var shortBestPrice float64

	if b.Mode == BarMode {
		longCrossPrice = b.Bar.Low
		shortCrossPrice = b.Bar.High
		longBestPrice = b.Bar.Open
		shortBestPrice = b.Bar.Open
	} else if b.Mode == TickMode {
		log.Println("Tick mode TODO")
	}

	for _, limitOrder := range b.ActiveLimitOrders {
		if limitOrder.Status == SUBMITTING {
			limitOrder.Status = NOTTRADED
			// TODO 传入策略中的响应函数中
			// b.strategy.OnOrder(order)
		}

		// Check whether limit order can be filled
		longCross := (limitOrder.Direction == LONG && limitOrder.Price >= longCrossPrice && longCrossPrice > 0)
		shortCross := (limitOrder.Direction == SHORT && limitOrder.Price <= shortCrossPrice && shortCrossPrice > 0)

		if !longCross && !shortCross {
			continue
		}

		// Push order update with status "all traded" (filled)
		limitOrder.Traded = limitOrder.Volume
		limitOrder.Status = ALLTRADED
		// TODO 传入策略中的响应函数中
		// b.strategy.OnOrder(order)
		delete(b.ActiveLimitOrders, limitOrder.VtOrderId)

		b.TradeCount++

		var tradePrice, posChange float64
		if longCross {
			tradePrice = math.Min(limitOrder.Price, longBestPrice)
			posChange = limitOrder.Volume
		} else if shortCross {
			tradePrice = math.Max(limitOrder.Price, shortBestPrice)
			posChange = -limitOrder.Volume
		}

		trade := NewTradeData(
			b.Gateway, limitOrder.Symbol, limitOrder.OrderId,
			fmt.Sprintf("%d", b.TradeCount),
			limitOrder.Exchange, limitOrder.Direction, limitOrder.Offset,
			tradePrice, limitOrder.Volume, b.Datetime,
		)

		//TODO strategy.pos+=poschange
		log.Println(posChange)
		//TODO strategy.OnTrade(trade)
		b.Trades[trade.VtTradeId] = trade
	}
}

func (b *BacktestEngine) CrossStopOrder() {
	var longCrossPrice float64
	var shortCrossPrice float64
	var longBestPrice float64
	var shortBestPrice float64

	if b.Mode == BarMode {
		longCrossPrice = b.Bar.High
		shortCrossPrice = b.Bar.Low
		longBestPrice = b.Bar.Open
		shortBestPrice = b.Bar.Open
	} else if b.Mode == TickMode {
		log.Println("Tick mode TODO")
	}

	for _, stopOrder := range b.ActiveStopOrders {
		// if order.Status == SUBMITTING {
		// 	order.Status = NOTTRADED
		// 	// TODO 传入策略中的响应函数中
		// 	// b.strategy.OnOrder(order)
		// }

		// Check whether limit order can be filled
		longCross := (stopOrder.Direction == LONG && stopOrder.Price <= longCrossPrice && longCrossPrice > 0)
		shortCross := (stopOrder.Direction == SHORT && stopOrder.Price >= shortCrossPrice && shortCrossPrice > 0)

		if !longCross && !shortCross {
			continue
		}

		// Push order update with status "all traded" (filled)
		// order.Traded = order.Volume
		// order.Status = ALLTRADED
		// TODO 传入策略中的响应函数中
		// b.strategy.OnOrder(order)
		// delete(b.ActiveLimitOrders, order.VtOrderId)

		// turn stop order into limit order
		b.LimitOrderCount++
		// limit orderId +1
		fromStopToLimitOrder := NewOrderData(
			stopOrder.Gateway, stopOrder.Symbol, stopOrder.Exchange,
			fmt.Sprint(b.LimitOrderCount), stopOrder.Direction,
			stopOrder.Offset, stopOrder.Price, stopOrder.Volume,
			ALLTRADED, b.Datetime,
		)

		b.LimitOrders[fromStopToLimitOrder.VtOrderId] = fromStopToLimitOrder

		// create trade data
		var tradePrice, posChange float64
		if longCross {
			tradePrice = math.Max(fromStopToLimitOrder.Price, longBestPrice)
			posChange = fromStopToLimitOrder.Volume
		} else if shortCross {
			tradePrice = math.Min(fromStopToLimitOrder.Price, shortBestPrice)
			posChange = -fromStopToLimitOrder.Volume
		}

		b.TradeCount++

		trade := NewTradeData(
			b.Gateway, fromStopToLimitOrder.Symbol, fromStopToLimitOrder.OrderId,
			fmt.Sprintf("%d", b.TradeCount),
			fromStopToLimitOrder.Exchange, fromStopToLimitOrder.Direction, fromStopToLimitOrder.Offset,
			tradePrice, fromStopToLimitOrder.Volume, b.Datetime,
		)
		b.Trades[trade.VtTradeId] = trade

		// Update stop order
		// stop_order.vt_orderids 这里没有按照vnpy写, 感觉没什么用
		stopOrder.Status = TRIGGERED
		// stop order 状态变为 triggered 存回StopOrders中
		b.StopOrders[stopOrder.VtOrderId] = stopOrder
		delete(b.ActiveStopOrders, stopOrder.OrderId)

		// TODO strategy.OnStopOrder(stopOrder)
		// TODO strategy.OnOrder(fromStopToLimitOrder)

		//TODO strategy.pos+=poschange
		log.Println(posChange)
		//TODO strategy.OnTrade(trade)
	}
}

func (b *BacktestEngine) UpdateDailyClose(close float64) {
	date := b.Datetime.Format("2006-01-02")
	if dailyResult, ok := b.DailyResults[date]; ok {
		dailyResult.ClosePrice = close
	} else {
		b.DailyResults[date] = NewDailyResult(date, close)
	}
}

func (b *BacktestEngine) SendOrder(
	strategy *strategy.Strategy,
	direction Direction,
	offset Offset,
	price, volume float64,
	stop, lock, net bool) string {

	var vtOrderId string
	var contract *ContractData
	price = util.RoundTo(price, b.PriceTick)

	if stop {
		vtOrderId = b.SendStopOrder(strategy, contract, direction, offset, price, volume, false)
	} else {
		vtOrderId = b.SendLimitOrder(strategy, contract, direction, offset, price, volume, false)
	}

	return vtOrderId
}

func (b *BacktestEngine) SendStopOrder(
	strategy *strategy.Strategy, contract *ContractData,
	direction Direction, offset Offset,
	price, volume float64, lock bool,
) string {

	b.StopOrderCount++
	stopOrder := NewStopOrderData(
		b.Gateway, b.Symbol, Exchange(b.Parameters.Exchange),
		direction, offset, price, volume, strategy.Name,
		fmt.Sprintf("STOP.%d", b.StopOrderCount), b.Datetime)

	b.ActiveStopOrders[stopOrder.StopOrderId] = stopOrder
	b.StopOrders[stopOrder.StopOrderId] = stopOrder

	return stopOrder.StopOrderId
}

func (b *BacktestEngine) SendLimitOrder(
	strategy *strategy.Strategy, contract *ContractData,
	direction Direction, offset Offset,
	price, volume float64, lock bool,
) string {

	b.LimitOrderCount++
	order := NewOrderData(
		b.Gateway, b.Symbol, b.Exchange, fmt.Sprintf("%d", b.LimitOrderCount),
		direction, offset, price, volume, SUBMITTING, b.Datetime)

	b.ActiveLimitOrders[order.VtOrderId] = order
	b.LimitOrders[order.VtOrderId] = order

	return order.VtOrderId
}
