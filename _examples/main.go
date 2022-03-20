package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	bnws "github.com/lvdou-bing/binance-ws"
)

func main() {
	// create WsService with ConnConf, this is recommended, key and secret will be needed by some channels
	// ctx and logger could be nil, they'll be initialized by default
	// ws, err := bnws.NewWsService(nil, nil, bnws.NewConnConf("",
	// 	"YOUR_API_KEY", "YOUR_API_SECRET", 10))
	// RECOMMEND this way to get a ConnConf
	ws, err := bnws.NewWsService(nil, nil, bnws.NewConnConfFromOption(&bnws.ConfOptions{
		Key: "YOUR_API_KEY", Secret: "YOUR_API_SECRET", MaxRetryConn: 10, SkipTlsVerify: false,
	}))
	// we can also do nothing to get a WsService, all parameters will be initialized by default and default url is spot
	// but some channels need key and secret for auth, we can also use set function to set key and secret
	// ws, err := bnws.NewWsService(nil, nil, nil)
	// ws.SetKey("YOUR_API_KEY")
	// ws.SetSecret("YOUR_API_SECRET")
	if err != nil {
		log.Printf("NewWsService err:%s", err.Error())
		return
	}

	// create callback functions for receive messages
	callOrder := bnws.NewCallBack(func(msg *bnws.UpdateMsg) {
		// parse the message to struct we need
		var order []bnws.SpotOrderMsg
		if err := json.Unmarshal(msg.Result, &order); err != nil {
			log.Printf("order Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", order)
	})
	callTrade := bnws.NewCallBack(func(msg *bnws.UpdateMsg) {
		var trade bnws.SpotTradeMsg
		if err := json.Unmarshal(msg.Result, &trade); err != nil {
			log.Printf("trade Unmarshal err:%s", err.Error())
		}
		log.Printf("%+v", trade)
	})
	// first, we need set callback function
	ws.SetCallBack(bnws.ChannelSpotOrder, callOrder)
	ws.SetCallBack(bnws.ChannelSpotPublicTrade, callTrade)
	// second, after set callback function, subscribe to any channel you are interested into
	if err := ws.Subscribe(bnws.ChannelSpotPublicTrade, []string{"BCH_USDT"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
		return
	}
	if err := ws.Subscribe(bnws.ChannelSpotOrder, []string{"BCH_USDT"}); err != nil {
		log.Printf("Subscribe err:%s", err.Error())
		return
	}

	// example for maintaining local order book
	// LocalOrderBook(context.Background(), ws, []string{"BTC_USDT"})

	ch := make(chan os.Signal)
	signal.Ignore(syscall.SIGPIPE, syscall.SIGALRM)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT, syscall.SIGKILL)
	<-ch
}
