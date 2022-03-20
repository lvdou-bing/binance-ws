package main

import (
	"encoding/json"
	"log"
	"time"

	bnws "github.com/lvdou-bing/binance-ws"
)

func main2() {
	spotWs, err := bnws.NewWsService(nil, nil, bnws.NewConnConfFromOption(&bnws.ConfOptions{
		URL:          bnws.BaseUrl,
		Key:          "YOUR_API_KEY",
		Secret:       "YOUR_API_SECRET",
		MaxRetryConn: 10,
	}))
	if err != nil {
		log.Printf("new spot wsService err:%s", err.Error())
		return
	}

	futureWs, err := bnws.NewWsService(nil, nil, bnws.NewConnConfFromOption(&bnws.ConfOptions{
		URL:          bnws.FuturesUsdtUrl,
		Key:          "YOUR_API_KEY",
		Secret:       "YOUR_API_SECRET",
		MaxRetryConn: 10,
	}))
	if err != nil {
		log.Printf("new futures wsService err:%s", err.Error())
		return
	}

	// create callback functions for receive messages
	// spot order book update
	callOrderBookUpdate := bnws.NewCallBack(func(msg *bnws.UpdateMsg) {
		// parse the message to struct we need
		var update bnws.SpotUpdateDepthMsg
		if err := json.Unmarshal(msg.Result, &update); err != nil {
			log.Printf("order book update Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", update)
	})

	// futures trade
	callTrade := bnws.NewCallBack(func(msg *bnws.UpdateMsg) {
		var trades []bnws.FuturesTrade
		if err := json.Unmarshal(msg.Result, &trades); err != nil {
			log.Printf("trade Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", trades)
	})

	// first, set callback
	spotWs.SetCallBack(bnws.ChannelSpotOrderBookUpdate, callOrderBookUpdate)
	futureWs.SetCallBack(bnws.ChannelFutureTrade, callTrade)
	if err := spotWs.Subscribe(bnws.ChannelSpotOrderBookUpdate, []string{"BTC_USDT", "100ms"}); err != nil {
		log.Printf("spotWs Subscribe err:%s", err.Error())
		return
	}

	if err := futureWs.Subscribe(bnws.ChannelFutureTrade, []string{"BTC_USDT"}); err != nil {
		log.Printf("futureWs Subscribe err:%s", err.Error())
		return
	}

	ch := make(chan bool)
	defer close(ch)

	for {
		select {
		case <-ch:
			log.Printf("manual done")
		case <-time.After(time.Second * 1000):
			log.Printf("auto done")
			return
		}
	}
}
