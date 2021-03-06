package bnws

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsService struct {
	Logger *log.Logger
	Ctx    context.Context
	Client *websocket.Conn
	once   *sync.Once
	msgChs *sync.Map // business chan
	calls  *sync.Map
	conf   *ConnConf
}

// ConnConf default URL is spot websocket
type ConnConf struct {
	subscribeMsg  *sync.Map
	URL           string
	MaxRetryConn  int
	SkipTlsVerify bool
}

// option is used to create conf
type ConfOptions struct {
	URL           string
	MaxRetryConn  int
	SkipTlsVerify bool
}

func NewWsService(ctx context.Context, logger *log.Logger, conf *ConnConf) (*WsService, error) {
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var cfg *ConnConf
	if conf != nil {
		cfg = conf
	} else {
		cfg = getInitConnConf()
	}

	stop := false
	retry := 0
	var conn *websocket.Conn
	for !stop {
		dialer := websocket.DefaultDialer
		if cfg.SkipTlsVerify {
			dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		c, _, err := dialer.Dial(cfg.URL, nil)
		if err != nil {
			if retry >= cfg.MaxRetryConn {
				logger.Printf("max reconnect time %d reached, give it up", cfg.MaxRetryConn)
				return nil, err
			}
			retry++
			logger.Printf("failed to connect to server for the %d time, try again later", retry)
			time.Sleep(time.Millisecond * (time.Duration(retry) * 500))
			continue
		} else {
			stop = true
			conn = c
		}
	}

	if retry > 0 {
		logger.Printf("reconnect succeeded after retrying %d times", retry)
	}

	ws := &WsService{
		conf:   cfg,
		Logger: logger,
		Ctx:    ctx,
		Client: conn,
		calls:  new(sync.Map),
		msgChs: new(sync.Map),
		once:   new(sync.Once),
	}

	return ws, nil
}

func getInitConnConf() *ConnConf {
	return &ConnConf{
		subscribeMsg:  new(sync.Map),
		MaxRetryConn:  MaxRetryConn,
		URL:           BaseUrl,
		SkipTlsVerify: false,
	}
}

func NewConnConf(url string, maxRetry int, skipTlsVerify bool) *ConnConf {
	if url == "" {
		url = BaseUrl
	}
	if maxRetry == 0 {
		maxRetry = MaxRetryConn
	}
	return &ConnConf{
		subscribeMsg:  new(sync.Map),
		MaxRetryConn:  maxRetry,
		URL:           url,
		SkipTlsVerify: skipTlsVerify,
	}
}

// NewConnConfFromOption conf from options, recommend using this
func NewConnConfFromOption(op *ConfOptions) *ConnConf {
	if op.URL == "" {
		op.URL = BaseUrl
	}
	if op.MaxRetryConn == 0 {
		op.MaxRetryConn = MaxRetryConn
	}

	return &ConnConf{
		subscribeMsg:  new(sync.Map),
		MaxRetryConn:  op.MaxRetryConn,
		URL:           op.URL,
		SkipTlsVerify: op.SkipTlsVerify,
	}
}

func (ws *WsService) GetConnConf() *ConnConf {
	return ws.conf
}

func (ws *WsService) reconnect() error {
	stop := false
	retry := 0
	for !stop {
		c, _, err := websocket.DefaultDialer.Dial(ws.conf.URL, nil)
		if err != nil {
			if retry >= ws.conf.MaxRetryConn {
				ws.Logger.Printf("max reconnect time %d reached, give it up", ws.conf.MaxRetryConn)
				return err
			}
			retry++
			ws.Logger.Printf("failed to connect to server for the %d time, try again later", retry)
			time.Sleep(time.Millisecond * (time.Duration(retry) * 500))
			continue
		} else {
			stop = true
			ws.Client = c
		}
	}

	// resubscribe after reconnect
	ws.conf.subscribeMsg.Range(func(key, value interface{}) bool {
		// key is channel, value is []requestHistory
		for _, req := range value.([]requestHistory) {
			if req.op == nil {
				req.op = &SubscribeOptions{
					IsReConnect: true,
				}
			} else {
				req.op.IsReConnect = true
			}
			if err := ws.baseSubscribe(req.Event, req.Stream, req.op); err != nil {
				ws.Logger.Printf("after reconnect, subscribe stream[%s] err:%s", key.(string), err.Error())
			} else {
				ws.Logger.Printf("reconnect stream[%s] success", key.(string))
			}
		}
		return true
	})

	return nil
}

func (ws *WsService) SetMaxRetryConn(max int) {
	ws.conf.MaxRetryConn = max
}

func (ws *WsService) GetMaxRetryConn() int {
	return ws.conf.MaxRetryConn
}

// func (ws *WsService) GetChannelMarkets(channel string) []string {
// 	var markets []string
// 	set := mapset.NewSet()
// 	if v, ok := ws.conf.subscribeMsg.Load(channel); ok {
// 		for _, req := range v.([]requestHistory) {
// 			if req.Event == Subscribe {
// 				set.Add(req.Stream)
// 			} else {
// 				set.Remove(req.Stream)
// 			}
// 		}

// 		for _, v := range set.ToSlice() {
// 			markets = append(markets, v.(string))
// 		}
// 		return markets
// 	}
// 	return markets
// }

func (ws *WsService) GetChannels() []string {
	var channels []string
	ws.calls.Range(func(key, value interface{}) bool {
		channels = append(channels, key.(string))
		return true
	})
	return channels
}

func (ws *WsService) GetConnection() *websocket.Conn {
	return ws.Client
}
