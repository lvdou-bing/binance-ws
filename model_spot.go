package bnws

type SpotKlineMsgPayload struct {
	StartTime            int64   `json:"t"` // 这根K线的起始时间
	EndTime              int64   `json:"T"` // 这根K线的结束时间
	Symbol               string  `json:"s"` // 交易对 大写
	Interval             string  `json:"i"` // K线间隔
	FirstTradeId         int64   `json:"f"` // 这根K线期间第一笔成交ID
	LastTradeId          int64   `json:"L"` // 这根K线期间末一笔成交ID
	OpenPrice            float64 `json:"o"` // 这根K线期间第一笔成交价
	ClosePrice           float64 `json:"c"` // 这根K线期间末一笔成交价
	HighestPrice         float64 `json:"h"` // 这根K线期间最高成交价
	LowestPrice          float64 `json:"l"` // 这根K线期间最低成交价
	VolumeOfAllTrades    float64 `json:"v"` // 这根K线期间成交量
	NumberOfAllTrades    float64 `json:"n"` // 这根K线期间成交笔数
	QuoteOfAllTrades     float64 `json:"q"` // 这根K线期间成交额
	VolumeOfActiveTrades float64 `json:"V"` // 主动买入的成交量
	QuoteOfActiveTrades  float64 `json:"Q"` // 主动买入的成交额
	IsFinished           bool    `json:"x"` // 这根K线是否完结(是否已经开始下一根K线)
	//   "B": "123456"   // 忽略此参数
}

type SpotKlineMsg struct {
	Event     string               `json:"e"` // 事件类型
	EventTime int64                `json:"E"` // 事件时间
	Symbol    string               `json:"s"` // 交易对 大写
	Kline     *SpotKlineMsgPayload `json:"k"` // payload
}

type SpotBalancesMsg struct {
	Timestamp        string `json:"timestamp"`
	TimestampInMilli string `json:"timestamp_ms"`
	User             string `json:"user"`
	Asset            string `json:"currency"`
	Change           string `json:"change"`
	Total            string `json:"total"`
	Available        string `json:"available"`
}

type SpotCandleUpdateMsg struct {
	Time   string `json:"t"`
	Volume string `json:"v"`
	Close  string `json:"c"`
	High   string `json:"h"`
	Low    string `json:"l"`
	Open   string `json:"o"`
	Name   string `json:"n"`
}

// SpotUpdateDepthMsg update order book
type SpotUpdateDepthMsg struct {
	TimeInMilli  int64      `json:"t"`
	Event        string     `json:"e"`
	ETime        int64      `json:"E"`
	CurrencyPair string     `json:"s"`
	FirstId      int64      `json:"U"`
	LastId       int64      `json:"u"`
	Bid          [][]string `json:"b"`
	Ask          [][]string `json:"a"`
}

// SpotUpdateAllDepthMsg all order book
type SpotUpdateAllDepthMsg struct {
	TimeInMilli  int64       `json:"t"`
	LastUpdateId int64       `json:"lastUpdateId"`
	CurrencyPair string      `json:"s"`
	Bid          [][2]string `json:"bids"`
	Ask          [][2]string `json:"asks"`
}

type SpotFundingBalancesMsg struct {
	Timestamp        string `json:"timestamp"`
	TimestampInMilli string `json:"timestamp_ms"`
	User             string `json:"user"`
	Asset            string `json:"currency"`
	Change           string `json:"change"`
	Freeze           string `json:"freeze"`
	Lent             string `json:"lent"`
}

type SpotMarginBalancesMsg struct {
	Timestamp        string `json:"timestamp"`
	TimestampInMilli string `json:"timestamp_ms"`
	User             string `json:"user"`
	Market           string `json:"currency_pair"`
	Asset            string `json:"currency"`
	Change           string `json:"change"`
	Available        string `json:"available"`
	Freeze           string `json:"freeze"`
	Borrowed         string `json:"borrowed"`
	Interest         string `json:"interest"`
}

type SpotBookTickerMsg struct {
	TimeInMilli  int64  `json:"t"`
	LastId       int64  `json:"u"`
	CurrencyPair string `json:"s"`
	Bid          string `json:"b"`
	BidSize      string `json:"B"`
	Ask          string `json:"a"`
	AskSize      string `json:"A"`
}

type SpotTickerMsg struct {
	// Currency pair
	CurrencyPair string `json:"currency_pair,omitempty"`
	// Last trading price
	Last string `json:"last,omitempty"`
	// Lowest ask
	LowestAsk string `json:"lowest_ask,omitempty"`
	// Highest bid
	HighestBid string `json:"highest_bid,omitempty"`
	// Change percentage.
	ChangePercentage string `json:"change_percentage,omitempty"`
	// Base currency trade volume
	BaseVolume string `json:"base_volume,omitempty"`
	// Quote currency trade volume
	QuoteVolume string `json:"quote_volume,omitempty"`
	// Highest price in 24h
	High24h string `json:"high_24h,omitempty"`
	// Lowest price in 24h
	Low24h string `json:"low_24h,omitempty"`
}

type SpotUserTradesMsg struct {
	Id           uint64 `json:"id"`
	UserId       uint64 `json:"user_id"`
	OrderId      string `json:"order_id"`
	CurrencyPair string `json:"currency_pair"`
	CreateTime   int64  `json:"create_time"`
	CreateTimeMs string `json:"create_time_ms"`
	Side         string `json:"side"`
	Amount       string `json:"amount"`
	Role         string `json:"role"`
	Price        string `json:"price"`
	Fee          string `json:"fee"`
	FeeCurrency  string `json:"fee_currency"`
	PointFee     string `json:"point_fee"`
	GtFee        string `json:"gt_fee"`
	Text         string `json:"text"`
}

type SpotTradeMsg struct {
	Id           uint64 `json:"id"`
	CreateTime   int64  `json:"create_time"`
	CreateTimeMs string `json:"create_time_ms"`
	Side         string `json:"side"`
	CurrencyPair string `json:"currency_pair"`
	Amount       string `json:"amount"`
	Price        string `json:"price"`
}

type OrderMsg struct {
	// SpotOrderMsg ID
	Id string `json:"id,omitempty"`
	// User defined information. If not empty, must follow the rules below:  1. prefixed with `t-` 2. no longer than 28 bytes without `t-` prefix 3. can only include 0-9, A-Z, a-z, underscore(_), hyphen(-) or dot(.)
	Text string `json:"text,omitempty"`
	// SpotOrderMsg creation time
	CreateTime string `json:"create_time,omitempty"`
	// SpotOrderMsg last modification time
	UpdateTime string `json:"update_time,omitempty"`
	// SpotOrderMsg status  - `open`: to be filled - `closed`: filled - `cancelled`: cancelled
	Status string `json:"status,omitempty"`
	// Currency pair
	CurrencyPair string `json:"currency_pair"`
	// SpotOrderMsg type. limit - limit order
	Type string `json:"type,omitempty"`
	// Account type. spot - use spot account; margin - use margin account
	Account string `json:"account,omitempty"`
	// SpotOrderMsg side
	Side string `json:"side"`
	// SpotTradeMsg amount
	Amount string `json:"amount"`
	// SpotOrderMsg price
	Price string `json:"price"`
	// Time in force  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled, taker only - poc: PendingOrCancelled, makes a post-only order that always enjoys a maker fee
	TimeInForce string `json:"time_in_force,omitempty"`
	// Amount to display for the iceberg order. Null or 0 for normal orders
	Iceberg string `json:"iceberg,omitempty"`
	// Used in margin trading(i.e. `account` is `margin`) to allow automatic loan of insufficient part if balance is not enough.
	AutoBorrow bool `json:"auto_borrow,omitempty"`
	// Amount left to fill
	Left string `json:"left,omitempty"`
	// Total filled in quote currency. Deprecated in favor of `filled_total`
	FillPrice string `json:"fill_price,omitempty"`
	// Total filled in quote currency
	FilledTotal string `json:"filled_total,omitempty"`
	// Fee deducted
	Fee string `json:"fee,omitempty"`
	// Fee currency unit
	FeeCurrency string `json:"fee_currency,omitempty"`
	// Point used to deduct fee
	PointFee string `json:"point_fee,omitempty"`
	// GT used to deduct fee
	GtFee string `json:"gt_fee,omitempty"`
	// Whether GT fee discount is used
	GtDiscount bool `json:"gt_discount,omitempty"`
	// Rebated fee
	RebatedFee string `json:"rebated_fee,omitempty"`
	// Rebated fee currency unit
	RebatedFeeCurrency string `json:"rebated_fee_currency,omitempty"`
}

type SpotOrderMsg struct {
	OrderMsg
	CreateTimeMs string `json:"create_time_ms,omitempty"`
	UpdateTimeMs string `json:"update_time_ms,omitempty"`
	User         int64  `json:"user"`
	Event        string `json:"event"`
}
