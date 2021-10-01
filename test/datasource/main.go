package main

import (
	"context"
	"gonpy/trader"
	"gonpy/trader/database"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/binance"
	"go.mongodb.org/mongo-driver/bson"
)

func MapGoExKlinePeriodToVnpyInterval(kp goex.KlinePeriod)trader.Interval{
	var i trader.Interval
	switch kp{
	case goex.KLINE_PERIOD_1MIN:
		i = trader.MINUTE
	case goex.KLINE_PERIOD_1H:
		i = trader.HOUR
	case goex.KLINE_PERIOD_1DAY:
		i = trader.DAILY
	case goex.KLINE_PERIOD_1WEEK:
		i = trader.WEEKLY
	default:
		log.Fatal("未知 Interval 类型")
		return ""
	}
	return i
}

func main() {
	beginTime := time.Date(2017, 12, 18, 0, 0, 0, 0, time.Local) //开始时间2017年12月18日,需自行修改
	var klinePeriod goex.KlinePeriod = goex.KLINE_PERIOD_1MIN    //see: github.com/nntaoli-project/goex/Const.go
	dbInterval := MapGoExKlinePeriodToVnpyInterval(klinePeriod)
	currencyPair := goex.LTC_USDT
	proxyUrl := "http://127.0.0.1:7890"  // 国内目前不挂代理无法请求数据

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c
		cancel()
	}()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	if proxyUrl != "" {
		log.Println("proxy:", proxyUrl)
		httpClient.Transport = &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return url.Parse(proxyUrl) //ss proxy
			},
		}
	}

	ba := binance.NewWithConfig(&goex.APIConfig{
		HttpClient: httpClient,
	})


	since := map[string]interface{}{
		"startTime": int(beginTime.Unix()) * 1000,
		// "endTime": int(endTime.Unix()) * 1000,
	}
	interval := time.NewTimer(200 * time.Millisecond)

	db := database.NewMongoDB("192.168.0.113", 27017)
	insertParam := &database.InsertParam{
		Db: "binance",
		Collection: "db_bar_data",
		Ordered: false,
	}

	dataNum := 0
	var startTime int
	for {
		select {
		case <-ctx.Done():
			return
		case <-interval.C:
			klines, err := ba.GetKlineRecords(currencyPair, klinePeriod, 1000, since)
			if err != nil {
				log.Println(err)
				interval.Reset(200 * time.Millisecond)
				continue
			}
			insertParam.Doc = make([]interface{}, 0)
			for _, k := range klines {
				insertParam.Doc = append(insertParam.Doc, bson.M{
					"datetime":time.Unix(k.Timestamp,0),
					"symbol":k.Pair.String(),
					"open_price":k.Open,
					"high_price":k.High,
					"low_price":k.Low,
					"close_price":k.Close,
					"volume":k.Vol,
					"interval": dbInterval,
				})
			}
			db.InsertMany(insertParam)
			dataNum += len(klines)

			startTime = int(klines[len(klines)-1].Timestamp)*1000 + 1
			since["startTime"] = startTime
			if len(klines) < 1000 {
				cancel()
			}

			interval.Reset(200 * time.Millisecond)
		}

		log.Println("当前数据量为: ",dataNum, "条.", "最后数据日期:", time.Unix(int64(startTime/1000),0))
	}
}
