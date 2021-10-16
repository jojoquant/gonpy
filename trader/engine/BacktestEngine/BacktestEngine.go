package BacktestEngine

import (
	"encoding/json"
	"fmt"
	"gonpy/trader"
	"gonpy/trader/database"
	"gonpy/trader/engine"
	"gonpy/trader/object"
	"gonpy/trader/strategy"
	"gonpy/trader/util"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/montanaflynn/stats"
	"go.mongodb.org/mongo-driver/bson"
)

type Parameters struct {
	Symbol   string
	Exchange trader.Exchange
	VtSymbol string

	Start time.Time
	End   time.Time

	Rate      float64
	Slippage  float64
	Size      float64
	PriceTick float64
	Capital   float64
	RiskFree  float64

	AnnualDays int

	Mode     trader.BacktestMode
	Interval trader.Interval
	Inverse  bool
}

func NewParameters(
	symbol string,
	exchange trader.Exchange,
	start, end time.Time,
	rate, slippage, size, priceTick, capital, riskFree float64,
	mode trader.BacktestMode, interval trader.Interval, inverse bool,
) Parameters {
	p := Parameters{
		Symbol:     symbol,
		Exchange:   exchange,
		VtSymbol:   fmt.Sprintf("%s.%s", symbol, exchange),
		Start:      start,
		End:        end,
		Rate:       rate,
		Slippage:   slippage,
		Size:       size,
		PriceTick:  priceTick,
		Capital:    capital,
		RiskFree:   riskFree,
		AnnualDays: 240,
		Mode:       mode,
		Interval:   interval,
		Inverse:    inverse,
	}
	return p
}

type BacktestEngine struct {
	engine.BaseEngine
	Parameters

	Gateway string

	Strategy strategy.Strategyer
	Bar      *database.BarData
	Tick     *database.TickData
	Datetime time.Time

	Database *database.MongoDB
	DisplayDB *database.InfluxDB

	Days            int
	BarCallback     strategy.BarCallback
	TickCallback    strategy.TickCallback
	HistoryBarData  []*database.BarData
	HistoryTickData []*database.TickData

	ActiveLimitOrders map[string]*object.OrderData
	LimitOrders       map[string]*object.OrderData
	LimitOrderCount   int

	ActiveStopOrders map[string]*object.StopOrderData
	StopOrders       map[string]*object.StopOrderData
	StopOrderCount   int

	TradeCount int
	Trades     map[string]*object.TradeData

	DailyResults     map[string]*DailyResult
	DailyResultsKeys []string
	Dailydf          dataframe.DataFrame

	HasCalculated bool  // 用于标记回测结果是否统计过
}

func NewBacktestEngine(param Parameters, database *database.MongoDB, strategy strategy.Strategyer) *BacktestEngine {

	b := &BacktestEngine{
		Database:          database,
		Parameters:        param,
		Gateway:           "BacktestEngine",
		Strategy:          strategy,
		ActiveLimitOrders: make(map[string]*object.OrderData),
		ActiveStopOrders:  make(map[string]*object.StopOrderData),
		LimitOrders:       make(map[string]*object.OrderData),
		StopOrders:        make(map[string]*object.StopOrderData),
		Trades:            make(map[string]*object.TradeData),
		DailyResults:      make(map[string]*DailyResult),
		DailyResultsKeys:  make([]string, 0, 50),
	}
	// b.Parameters = param
	// b.Gateway = "BacktestEngine"

	// b.Strategy.VtSymbol = b.VtSymbol
	b.Strategy.SetVtSymbol(b.VtSymbol)
	strategy.SetTradeEngine(b)

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

func (b *BacktestEngine) LoadData(q *database.QueryParam) {
	if b.Start.After(b.End) {
		log.Fatalln("起始日期必须小于结束日期")
		return
	}

	log.Println("加载数据: ", b.Start, " -> ", b.End)

	q.Filter = bson.M{
		"symbol": b.Symbol, "exchange": b.Exchange, "interval": b.Interval,
		"datetime": bson.M{"$gte": b.Start, "$lte": b.End},
	}

	b.HistoryBarData = b.Database.Query(q)

	b.Database.Close()
}

func (b *BacktestEngine) LoadBar(
	vtSymbol string, days int,
	interval trader.Interval, callback strategy.BarCallback,
	useDatabase bool) {
	b.Days = days
	b.BarCallback = callback
}

func (b *BacktestEngine) LoadTick(
	vtSymbol string, days int,
	callback strategy.TickCallback,
	useDatabase bool) {
	b.Days = days
	b.TickCallback = callback
}

func (b *BacktestEngine) AddStrategy() {}

func (b *BacktestEngine) RunBacktest() {
	
	b.HasCalculated = false
	b.Strategy.OnInit()

	var index int
	dayCount := 0
	if b.Mode == trader.BarMode {
		for ix, data := range b.HistoryBarData {
			if !b.Datetime.IsZero() && (data.Datetime.Day() != b.Datetime.Day()) {
				dayCount++
				if dayCount >= b.Days {
					break
				}
			}

			b.Datetime = data.Datetime
			b.BarCallback(data)
			index = ix
		}

		// b.Strategy.Inited = true
		b.Strategy.SetInited(true)
		log.Println("策略初始化完成")

		b.Strategy.OnStart()
		// b.Strategy.Trading = true
		b.Strategy.SetTrading(true)

		log.Println("开始回放 Bar 历史数据")
		if len(b.HistoryBarData[index:]) <= 1 {
			log.Println("历史数据不足, 回测终止")
			return
		}

		for _, data := range b.HistoryBarData {
			b.NewBar(data)
			// log.Printf("当前回放进度: %d / %d \n", i, len(b.HistoryBarData[index:]))
		}

	} else if b.Mode == trader.TickMode {
		for _, data := range b.HistoryTickData {
			if !b.Datetime.IsZero() && (data.Datetime.Day() != b.Datetime.Day()) {
				dayCount++
				if dayCount >= b.Days {
					break
				}
			}

			b.Datetime = data.Datetime
			b.TickCallback(data)
			// index = ix
		}
		// function := b.NewTick
	}
}

func (b *BacktestEngine) NewBar(bar *database.BarData) {
	b.Bar = bar
	b.Datetime = bar.Datetime

	b.CrossLimitOrder(b.Strategy, "")
	b.CrossStopOrder(b.Strategy, "")
	b.Strategy.OnBar(bar)

	b.UpdateDailyClose(bar.Close)
}

func (b *BacktestEngine) NewTick(tick *database.TickData) {
	b.Tick = tick
	b.Datetime = tick.Datetime

	b.CrossLimitOrder(b.Strategy, "")
	b.CrossStopOrder(b.Strategy, "")
	b.Strategy.OnTick(tick)

	b.UpdateDailyClose(tick.LastPrice)
}

func (b *BacktestEngine) CrossLimitOrder(strategy strategy.Strategyer, vtOrderId string) {
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

	for _, limitOrder := range b.ActiveLimitOrders {
		if limitOrder.Status == trader.SUBMITTING {
			limitOrder.Status = trader.NOTTRADED
			// TODO 传入策略中的响应函数中
			// b.strategy.OnOrder(order)
		}

		// Check whether limit order can be filled
		longCross := (limitOrder.Direction == trader.LONG && limitOrder.Price >= longCrossPrice && longCrossPrice > 0)
		shortCross := (limitOrder.Direction == trader.SHORT && limitOrder.Price <= shortCrossPrice && shortCrossPrice > 0)

		if !longCross && !shortCross {
			continue
		}

		// Push order update with status "all traded" (filled)
		limitOrder.Traded = limitOrder.Volume
		limitOrder.Status = trader.ALLTRADED
		// TODO 传入策略中的响应函数中
		// b.strategy.OnOrder(order)
		delete(b.ActiveLimitOrders, limitOrder.VtOrderId)

		b.TradeCount++

		var tradePrice float64
		// var posChange float64
		if longCross {
			tradePrice = math.Min(limitOrder.Price, longBestPrice)
			// posChange = limitOrder.Volume
		} else if shortCross {
			tradePrice = math.Max(limitOrder.Price, shortBestPrice)
			// posChange = -limitOrder.Volume
		}

		trade := object.NewTradeData(
			b.Gateway, limitOrder.Symbol, limitOrder.OrderId,
			fmt.Sprintf("%d", b.TradeCount),
			limitOrder.Exchange, limitOrder.Direction, limitOrder.Offset,
			tradePrice, limitOrder.Volume, b.Datetime,
		)

		//TODO strategy.pos+=poschange
		// log.Println(posChange)
		//TODO strategy.OnTrade(trade)
		b.Trades[trade.VtTradeId] = trade
	}
}

func (b *BacktestEngine) CrossStopOrder(strategy strategy.Strategyer, vtOrderId string) {
	var longCrossPrice float64
	var shortCrossPrice float64
	var longBestPrice float64
	var shortBestPrice float64

	if b.Mode == trader.BarMode {
		longCrossPrice = b.Bar.High
		shortCrossPrice = b.Bar.Low
		longBestPrice = b.Bar.Open
		shortBestPrice = b.Bar.Open
	} else if b.Mode == trader.TickMode {
		log.Println("Tick mode TODO")
	}

	for _, stopOrder := range b.ActiveStopOrders {
		// if order.Status == SUBMITTING {
		// 	order.Status = NOTTRADED
		// 	// TODO 传入策略中的响应函数中
		// 	// b.strategy.OnOrder(order)
		// }

		// Check whether limit order can be filled
		longCross := (stopOrder.Direction == trader.LONG && stopOrder.Price <= longCrossPrice && longCrossPrice > 0)
		shortCross := (stopOrder.Direction == trader.SHORT && stopOrder.Price >= shortCrossPrice && shortCrossPrice > 0)

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
		fromStopToLimitOrder := object.NewOrderData(
			stopOrder.Gateway, stopOrder.Symbol, stopOrder.Exchange,
			fmt.Sprint(b.LimitOrderCount), stopOrder.Direction,
			stopOrder.Offset, stopOrder.Price, stopOrder.Volume,
			trader.ALLTRADED, b.Datetime,
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

		trade := object.NewTradeData(
			b.Gateway, fromStopToLimitOrder.Symbol, fromStopToLimitOrder.OrderId,
			fmt.Sprintf("%d", b.TradeCount),
			fromStopToLimitOrder.Exchange, fromStopToLimitOrder.Direction, fromStopToLimitOrder.Offset,
			tradePrice, fromStopToLimitOrder.Volume, b.Datetime,
		)
		b.Trades[trade.VtTradeId] = trade

		// Update stop order
		// stop_order.vt_orderids 这里没有按照vnpy写, 感觉没什么用
		stopOrder.Status = trader.TRIGGERED
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
		b.DailyResultsKeys = append(b.DailyResultsKeys, date)
	}
}

func (b *BacktestEngine) SendOrder(
	strategy strategy.Strategyer,
	direction trader.Direction,
	offset trader.Offset,
	price, volume float64,
	stop, lock, net bool) string {

	var vtOrderId string
	var contract *object.ContractData
	price = util.RoundTo(price, b.PriceTick)

	if stop {
		vtOrderId = b.SendStopOrder(strategy, contract, direction, offset, price, volume, lock, net)
	} else {
		vtOrderId = b.SendLimitOrder(strategy, contract, direction, offset, price, volume, lock, net)
	}

	return vtOrderId
}

func (b *BacktestEngine) SendStopOrder(
	strategy strategy.Strategyer, contract *object.ContractData,
	direction trader.Direction, offset trader.Offset,
	price, volume float64, lock, net bool,
) string {

	b.StopOrderCount++
	stopOrder := object.NewStopOrderData(
		b.Gateway, b.Symbol, trader.Exchange(b.Parameters.Exchange),
		direction, offset, price, volume, strategy.GetStrategyName(),
		fmt.Sprintf("%s.%d", trader.STOP, b.StopOrderCount), b.Datetime)

	b.ActiveStopOrders[stopOrder.StopOrderId] = stopOrder
	b.StopOrders[stopOrder.StopOrderId] = stopOrder

	return stopOrder.StopOrderId
}

func (b *BacktestEngine) SendLimitOrder(
	strategy strategy.Strategyer, contract *object.ContractData,
	direction trader.Direction, offset trader.Offset,
	price, volume float64, lock, net bool,
) string {

	b.LimitOrderCount++
	order := object.NewOrderData(
		b.Gateway, b.Symbol, b.Exchange, fmt.Sprintf("%d", b.LimitOrderCount),
		direction, offset, price, volume, trader.SUBMITTING, b.Datetime)

	b.ActiveLimitOrders[order.VtOrderId] = order
	b.LimitOrders[order.VtOrderId] = order

	return order.VtOrderId
}

func (b *BacktestEngine) CancelOrder(strategy strategy.Strategyer, vtOrderId string) {
	if strings.HasPrefix(vtOrderId, string(trader.STOP)) {
		b.CancelStopOrder(strategy, vtOrderId)
	} else {
		b.CancelLimitOrder(strategy, vtOrderId)
	}

}

func (b *BacktestEngine) CancelStopOrder(strategy strategy.Strategyer, vtOrderId string) {
	if order, ok := b.ActiveStopOrders[vtOrderId]; ok {
		order.Status = trader.CANCELLED
		b.Strategy.OnStopOrder(order)
		delete(b.ActiveStopOrders, vtOrderId)
	}
}

func (b *BacktestEngine) CancelLimitOrder(strategy strategy.Strategyer, vtOrderId string) {
	if order, ok := b.ActiveLimitOrders[vtOrderId]; ok {
		order.Status = trader.CANCELLED
		b.Strategy.OnOrder(order)
		delete(b.ActiveLimitOrders, vtOrderId)
	}
}

func (b *BacktestEngine) CancelAll(s strategy.Strategyer) {
	for vtOrderId := range b.ActiveLimitOrders {
		b.CancelLimitOrder(s, vtOrderId)
	}

	for stopOrderId := range b.ActiveStopOrders {
		b.CancelStopOrder(s, stopOrderId)
	}
}

func (b *BacktestEngine) CalculateResult() map[string]interface{} {
	log.Println("开始计算逐日盯市盈亏")

	if len(b.Trades) == 0 {
		log.Panicln("成交记录为空, 无法计算")
		return nil
	}

	for _, trade := range b.Trades {
		date := trade.Datetime.Format("2006-01-02")
		dailyResult := b.DailyResults[date]
		dailyResult.AddTrade(trade)
	}

	var preClose float64 = 0
	var startPos float64 = 0
	dr := []DailyResult{}
	for _, key := range b.DailyResultsKeys {
		b.DailyResults[key].CalculatePnl(preClose, startPos, b.Size, b.Rate, b.Slippage, b.Inverse)
		preClose = b.DailyResults[key].ClosePrice
		startPos = b.DailyResults[key].EndPos

		dr = append(dr, *b.DailyResults[key])
	}

	log.Println("逐日盯市盈亏计算完成")

	drJson, err := json.Marshal(dr)
	if err != nil {
		log.Println(err)
	}

	b.Dailydf = dataframe.ReadJSON(strings.NewReader(string(drJson)))
	DailydfLength := b.Dailydf.Nrow()
	balance, _ := stats.CumulativeSum(b.Dailydf.Col("NetPnl").Float())
	
	b.HasCalculated = true

	// 计算 return
	returnS := make([]float64, DailydfLength)
	highlevel := make([]float64, DailydfLength)
	drawdown := make([]float64, DailydfLength)
	ddpercent := make([]float64, DailydfLength)

	balance[0] = balance[0] + b.Capital
	highlevel[0] = balance[0]
	drawdown[0] = balance[0] - highlevel[0]
	ddpercent[0] = drawdown[0] / highlevel[0] * 100.0

	for i := 1; i < len(balance); i++ {
		balance[i] = balance[i] + b.Capital
		if i == 1 {
			returnS[i-1] = math.Log(balance[i] / b.Capital)
		} else {
			returnS[i-1] = math.Log(balance[i] / balance[i-1])
		}

		if math.IsNaN(returnS[i-1]) || math.IsInf(returnS[i-1], 0) {
			returnS[i-1] = 0.0
		}

		highlevel[i], err = stats.Max(balance[0 : i+1])
		if err != nil {
			log.Println(err)
		}
		drawdown[i] = balance[i] - highlevel[i]
		ddpercent[i] = drawdown[i] / highlevel[i] * 100.0
	}

	b.Dailydf = b.Dailydf.Mutate(series.New(balance, series.Float, "Balance"))
	b.Dailydf = b.Dailydf.Mutate(series.New(returnS, series.Float, "Return"))

	//TODO save HistoryBardata and dailydf into influxdb to display on grafana

	startDate := b.Dailydf.Select("Date").Records()[1][0]
	endDate := b.Dailydf.Select("Date").Records()[DailydfLength][0]

	totalDays := DailydfLength
	profitDays := b.Dailydf.Filter(
		dataframe.F{Colname: "NetPnl", Comparator: series.Greater, Comparando: 0.0},
	).Nrow()
	lossDays := b.Dailydf.Filter(
		dataframe.F{Colname: "NetPnl", Comparator: series.Less, Comparando: 0.0},
	).Nrow()

	endBalance := balance[len(balance)-1]
	maxDrawdown, _ := stats.Min(drawdown)
	maxDDpercent, _ := stats.Min(ddpercent)

	maxDrawdownEndIndex := util.SliceIndex(len(drawdown), func(i int) bool { return drawdown[i] == maxDrawdown })
	maxDrawdownStart, _ := stats.Max(balance[0:maxDrawdownEndIndex])
	maxDrawdownStartIndex := util.SliceIndex(len(drawdown), func(i int) bool { return balance[0:maxDrawdownEndIndex][i] == maxDrawdownStart })
	maxDrawdownDuration := maxDrawdownEndIndex - maxDrawdownStartIndex

	totalNetPnl, _ := stats.Sum(b.Dailydf.Col("NetPnl").Float())
	dailyNetPnl := totalNetPnl / float64(totalDays)

	totalCommission, _ := stats.Sum(b.Dailydf.Col("Commission").Float())
	dailyCommission := totalCommission / float64(totalDays)

	totalSlippage, _ := stats.Sum(b.Dailydf.Col("Slippage").Float())
	dailySlippage := totalSlippage / float64(totalDays)

	totalTurnover, _ := stats.Sum(b.Dailydf.Col("Turnover").Float())
	dailyTurnover := totalTurnover / float64(totalDays)

	totalTradeCount, _ := stats.Sum(b.Dailydf.Col("TradeCount").Float())
	dailyTradeCount := totalTradeCount / float64(totalDays)

	totalReturn := (endBalance/b.Capital - 1) * 100
	annualReturn := totalReturn / float64(totalDays) * float64(b.AnnualDays)

	returnMean, _ := stats.Mean(returnS)
	dailyReturn := returnMean * 100
	returnStd, _ := stats.StandardDeviation(returnS)
	returnStd = returnStd * 100

	var dailyRiskFree float64
	var sharpRatio float64
	if returnStd != 0 {
		sqrtAnnualDays := math.Sqrt(float64(b.AnnualDays))
		dailyRiskFree = b.RiskFree / sqrtAnnualDays
		sharpRatio = (dailyReturn - dailyRiskFree) / returnStd * sqrtAnnualDays
	} else {
		sharpRatio = 0
	}

	returnDarwdownRatio := -totalReturn / maxDDpercent

	log.Println("----------------------------------------------")
	log.Println("首个交易日:", startDate)
	log.Println("最后交易日:", endDate)
	log.Println("总交易日:", totalDays)
	log.Println("盈利交易日:", profitDays)
	log.Println("亏损交易日:", lossDays)

	log.Println("起始资金:", strconv.FormatFloat(b.Capital, 'f', 2, 64))
	log.Println("结束资金:", strconv.FormatFloat(endBalance, 'f', 2, 64))

	log.Println("总收益率:", strconv.FormatFloat(totalReturn, 'f', 2, 64), "%")
	log.Println("年化收益率:", strconv.FormatFloat(annualReturn, 'f', 2, 64), "%")
	log.Println("最大回撤:", strconv.FormatFloat(maxDrawdown, 'f', 2, 64))
	log.Println("百分比最大回撤:", strconv.FormatFloat(maxDDpercent, 'f', 2, 64), "%")
	log.Println("最长回撤天数:", maxDrawdownDuration)

	log.Println("总盈亏:", strconv.FormatFloat(totalNetPnl, 'f', 2, 64))
	log.Println("总手续费:", strconv.FormatFloat(totalCommission, 'f', 2, 64))
	log.Println("总滑点:", strconv.FormatFloat(totalSlippage, 'f', 2, 64))
	log.Println("总成交金额:", strconv.FormatFloat(totalTurnover, 'f', 2, 64))
	log.Println("总成交笔数:", totalTradeCount)
	log.Println("-----------------------")
	log.Println("日均盈亏:", strconv.FormatFloat(dailyNetPnl, 'f', 2, 64))
	log.Println("日均手续费:", strconv.FormatFloat(dailyCommission, 'f', 2, 64))
	log.Println("日均滑点:", strconv.FormatFloat(dailySlippage, 'f', 2, 64))
	log.Println("日均成交金额:", strconv.FormatFloat(dailyTurnover, 'f', 2, 64))
	log.Println("日均成交笔数:", strconv.FormatFloat(dailyTradeCount, 'f', 2, 64))
	log.Println("日均收益率:", strconv.FormatFloat(dailyReturn, 'f', 2, 64), "%")
	log.Println("收益标准差:", strconv.FormatFloat(returnStd, 'f', 2, 64), "%")
	log.Println("Sharp Ratio:", strconv.FormatFloat(sharpRatio, 'f', 2, 64))
	log.Println("收益回撤比:", strconv.FormatFloat(returnDarwdownRatio, 'f', 2, 64))
	
	log.Println("----------------------------------------------")
	log.Println("策略统计指标计算完成")
	
	statistics := map[string]interface{}{
		"startDate":startDate,
		"endDate":endDate,
		"totalDays":totalDays,
		"profitDays":profitDays,
		"lossDays":lossDays,
		"capital":b.Capital,
		"endBalance":endBalance,
		"maxDrawdown":maxDrawdown,
		"maxDDpercent":maxDDpercent,
		"maxDrawdownDuration":maxDrawdownDuration,
		"totalNetPnl":totalNetPnl,
		"dailyNetPnl":dailyNetPnl,
		"totalCommission":totalCommission,
		"dailyCommission":dailyCommission,
		"totalSlippage":totalSlippage,
		"dailySlippage":dailySlippage,
		"totalTurnover":totalTurnover,
		"dailyTurnover":dailyTurnover,
		"totalTradeCount":totalTradeCount,
		"dailyTradeCount":dailyTradeCount,
		"totalReturn":totalReturn,
		"annualReturn":annualReturn,
		"dailyReturn":dailyReturn,
		"returnStd":returnStd,
		"sharpRatio":sharpRatio,
		"returnDarwdownRatio":returnDarwdownRatio,
	}
	return statistics 
}

func(b *BacktestEngine) SaveBacktestResultToInfluxDB(){
	if !b.HasCalculated{
		log.Println("未完成统计，禁止保存数据到 InfluxDB")
		return
	}


}