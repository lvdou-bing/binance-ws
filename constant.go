package bnws

const (
	BaseUrl          = "wss://stream.binance.com:9443/stream"
	AuthMethodApiKey = "api_key"
	MaxRetryConn     = 10
)

// spot channels
const (
	// ChannelSpotBalance         = "spot.balances"
	// ChannelSpotCandleStick     = "spot.candlesticks"
	// ChannelSpotOrder           = "spot.orders"
	// ChannelSpotOrderBook       = "spot.order_book"
	// ChannelSpotBookTicker      = "spot.book_ticker"
	// ChannelSpotOrderBookUpdate = "spot.order_book_update"
	// ChannelSpotTicker          = "spot.tickers"
	// ChannelSpotUserTrade       = "spot.usertrades"
	// ChannelSpotPublicTrade     = "spot.trades"
	// ChannelSpotFundingBalance  = "spot.funding_balances"
	// ChannelSpotMarginBalance   = "spot.margin_balances"
	// ChannelSpotCrossBalance    = "spot.cross_balances"
	ChannelSpotAggTrade       = "aggTrade"
	ChannelSpotTrade          = "trade"
	ChannelSpotKline          = "kline"
	ChannelSpot24hrMiniTicker = "24hrMiniTicker"
	ChannelSpot24hrTicker     = "24hrTicker"
)

// future channels
const (
	ChannelFutureTicker           = "futures.tickers"
	ChannelFutureTrade            = "futures.trades"
	ChannelFutureOrderBook        = "futures.order_book"
	ChannelFutureBookTicker       = "futures.book_ticker"
	ChannelFutureOrderBookUpdate  = "futures.order_book_update"
	ChannelFutureCandleStick      = "futures.candlesticks"
	ChannelFutureOrder            = "futures.orders"
	ChannelFutureUserTrade        = "futures.usertrades"
	ChannelFutureLiquidates       = "futures.liquidates"
	ChannelFutureAutoDeleverages  = "futures.auto_deleverages"
	ChannelFuturePositionCloses   = "futures.position_closes"
	ChannelFutureBalance          = "futures.balances"
	ChannelFutureReduceRiskLimits = "futures.reduce_risk_limits"
	ChannelFuturePositions        = "futures.positions"
	ChannelFutureAutoOrders       = "futures.autoorders"
)

// var (
// 	authChannel = map[string]bool{
// 		// spot
// 		ChannelSpotBalance:        true,
// 		ChannelSpotFundingBalance: true,
// 		ChannelSpotMarginBalance:  true,
// 		ChannelSpotOrder:          true,
// 		ChannelSpotUserTrade:      true,

// 		// future
// 		ChannelFutureOrder:            true,
// 		ChannelFutureUserTrade:        true,
// 		ChannelFutureLiquidates:       true,
// 		ChannelFutureAutoDeleverages:  true,
// 		ChannelFuturePositionCloses:   true,
// 		ChannelFutureReduceRiskLimits: true,
// 		ChannelFuturePositions:        true,
// 		ChannelFutureAutoOrders:       true,
// 		ChannelFutureBalance:          true,
// 	}
// )

const (
	Subscribe   = "SUBSCRIBE"
	UnSubscribe = "UNSUBSCRIBE"

	ServiceTypeSpot    = 1
	ServiceTypeFutures = 2
)
