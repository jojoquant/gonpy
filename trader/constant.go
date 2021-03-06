package trader

type Exchange string

const (
	// Chinese
	CFFEX Exchange = "CFFEX" // China Financial Futures Exchange
	SHFE  Exchange = "SHFE"  // Shanghai Futures Exchange
	CZCE  Exchange = "CZCE"  // Zhengzhou Commodity Exchange
	DCE   Exchange = "DCE"   // Dalian Commodity Exchange
	INE   Exchange = "INE"   // Shanghai International Energy Exchange
	SSE   Exchange = "SSE"   // Shanghai Stock Exchange
	SZSE  Exchange = "SZSE"  // Shenzhen Stock Exchange
	SGE   Exchange = "SGE"   // Shanghai Gold Exchange
	WXE   Exchange = "WXE"   // Wuxi Steel Exchange
	CFETS Exchange = "CFETS" // China Foreign Exchange Trade System

	// Global
	SMART    Exchange = "SMART"    // Smart Router for US stocks
	NYSE     Exchange = "NYSE"     // New York Stock Exchnage
	NASDAQ   Exchange = "NASDAQ"   // Nasdaq Exchange
	ARCA     Exchange = "ARCA"     // ARCA Exchange
	EDGEA    Exchange = "EDGEA"    // Direct Edge Exchange
	ISLAND   Exchange = "ISLAND"   // Nasdaq Island ECN
	BATS     Exchange = "BATS"     // Bats Global Markets
	IEX      Exchange = "IEX"      // The Investors Exchange
	NYMEX    Exchange = "NYMEX"    // New York Mercantile Exchange
	COMEX    Exchange = "COMEX"    // COMEX of CME
	GLOBEX   Exchange = "GLOBEX"   // Globex of CME
	IDEALPRO Exchange = "IDEALPRO" // Forex ECN of Interactive Brokers
	CME      Exchange = "CME"      // Chicago Mercantile Exchange
	ICE      Exchange = "ICE"      // Intercontinental Exchange
	SEHK     Exchange = "SEHK"     // Stock Exchange of Hong Kong
	HKFE     Exchange = "HKFE"     // Hong Kong Futures Exchange
	HKSE     Exchange = "HKSE"     // Hong Kong Stock Exchange
	SGX      Exchange = "SGX"      // Singapore Global Exchange
	CBOT     Exchange = "CBT"      // Chicago Board of Trade
	CBOE     Exchange = "CBOE"     // Chicago Board Options Exchange
	CFE      Exchange = "CFE"      // CBOE Futures Exchange
	DME      Exchange = "DME"      // Dubai Mercantile Exchange
	EUREX    Exchange = "EUX"      // Eurex Exchange
	APEX     Exchange = "APEX"     // Asia Pacific Exchange
	LME      Exchange = "LME"      // London Metal Exchange
	BMD      Exchange = "BMD"      // Bursa Malaysia Derivatives
	TOCOM    Exchange = "TOCOM"    // Tokyo Commodity Exchange
	EUNX     Exchange = "EUNX"     // Euronext Exchange
	KRX      Exchange = "KRX"      // Korean Exchange
	OTC      Exchange = "OTC"      // OTC Product (Forex/CFD/Pink Sheet Equity)
	IBKRATS  Exchange = "IBKRATS"  // Paper Trading Exchange of IB

	// CryptoCurrency
	BITMEX   Exchange = "BITMEX"
	OKEX     Exchange = "OKEX"
	HUOBI    Exchange = "HUOBI"
	BITFINEX Exchange = "BITFINEX"
	BINANCE  Exchange = "BINANCE"
	BYBIT    Exchange = "BYBIT" // bybit.com
	COINBASE Exchange = "COINBASE"
	DERIBIT  Exchange = "DERIBIT"
	GATEIO   Exchange = "GATEIO"
	BITSTAMP Exchange = "BITSTAMP"

	// Special Function
	LOCAL Exchange = "LOCAL" // For local generated data
)

type LogLevel int

const (
	CRITICAL LogLevel = 50
	FATAL    LogLevel = CRITICAL
	ERROR    LogLevel = 40
	WARNING  LogLevel = 30
	WARN     LogLevel = WARNING
	INFO     LogLevel = 20
	DEBUG    LogLevel = 10
	NOTSET   LogLevel = 0
)

type BacktestMode int

const (
	BarMode  BacktestMode = 1
	TickMode BacktestMode = 2
)

type OrderType string

const (
	LIMIT  OrderType = "??????"
	MARKET OrderType = "??????"
	STOP   OrderType = "STOP"
	FAK    OrderType = "FAK"
	FOK    OrderType = "FOK"
	RFQ    OrderType = "??????"
)

type OptionType string
const(
	CALL OptionType = "????????????"
	PUT OptionType = "????????????"
)

type Direction string

const (
	// Direction of order/trade/position.
	LONG  Direction = "???"
	SHORT Direction = "???"
	NET   Direction = "???"
)

type Offset string

const (
	//Offset of order/trade.
	NONE           Offset = ""
	OPEN           Offset = "???"
	CLOSE          Offset = "???"
	CLOSETODAY     Offset = "??????"
	CLOSEYESTERDAY Offset = "??????"
)

type Status string

const (
	// Order status.
	SUBMITTING Status = "?????????"
	NOTTRADED  Status = "?????????"
	PARTTRADED Status = "????????????"
	ALLTRADED  Status = "????????????"
	CANCELLED  Status = "?????????"
	REJECTED   Status = "??????"

	// stop order ????????????
	WAITING Status = "?????????"
	TRIGGERED Status = "?????????"
)

var ACTIVE_STATUSES [4]Status = [4]Status{SUBMITTING, NOTTRADED, PARTTRADED, WAITING}

type Product string

const (
	EQUITY  Product = "??????"
	FUTURES Product = "??????"
	OPTION  Product = "??????"
	INDEX   Product = "??????"
	FOREX   Product = "??????"
	SPOT    Product = "??????"
	ETF     Product = "ETF"
	BOND    Product = "??????"
	WARRANT Product = "??????"
	SPREAD  Product = "??????"
	FUND    Product = "??????"
)

type Interval string

const (
	// Interval of bar data.
	MINUTE Interval = "1m"
	HOUR   Interval = "1h"
	DAILY  Interval = "d"
	WEEKLY Interval = "w"
	TICK   Interval = "tick"
)
