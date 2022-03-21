package bnws

import (
	"encoding/json"
	"io"
	"net"

	"github.com/gorilla/websocket"
)

type SubscribeOptions struct {
	// ID          int64 `json:"id"`
	IsReConnect bool `json:"-"`
}

/*
@param stream string: the stream type
@param payload []string: the list of detail stream message, e.g. <symbol>@type_<option>
*/
func (ws *WsService) Subscribe(stream string) error {
	// if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[stream] {
	// 	return newAuthEmptyErr()
	// }

	msgCh, ok := ws.msgChs.Load(stream)
	if !ok {
		msgCh = make(chan *UpdateMsg, 1)
		go ws.receiveCallMsg(stream, msgCh.(chan *UpdateMsg))
	}

	err := ws.newBaseChannel(stream, msgCh.(chan *UpdateMsg), nil)
	if err != nil {
		return err
	}

	return nil
}

func (ws *WsService) SubscribeWithOption(stream string, op *SubscribeOptions) error {
	// if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[stream] {
	// 	return newAuthEmptyErr()
	// }

	msgCh, ok := ws.msgChs.Load(stream)
	if !ok {
		msgCh = make(chan *UpdateMsg, 1)
		go ws.receiveCallMsg(stream, msgCh.(chan *UpdateMsg))
	}

	err := ws.newBaseChannel(stream, msgCh.(chan *UpdateMsg), op)
	if err != nil {
		return err
	}

	return nil
}

func (ws *WsService) UnSubscribe(stream string) error {
	return ws.baseSubscribe(UnSubscribe, stream, nil)
}

func (ws *WsService) newBaseChannel(stream string, bch chan *UpdateMsg, op *SubscribeOptions) error {
	err := ws.baseSubscribe(Subscribe, stream, op)
	if err != nil {
		return err
	}

	if _, ok := ws.msgChs.Load(stream); !ok {
		ws.msgChs.Store(stream, bch)
	}

	ws.readMsg()

	return nil
}

func (ws *WsService) baseSubscribe(event string, stream string, op *SubscribeOptions) error {
	req := make(map[string]interface{})
	req["method"] = event
	req["params"] = []interface{}{stream}
	req["id"] = 1 // TODO, need a unique value

	byteReq, err := json.Marshal(req)
	if err != nil {
		ws.Logger.Printf("req Marshal err:%s", err.Error())
		return err
	}

	err = ws.Client.WriteMessage(websocket.TextMessage, byteReq)
	if err != nil {
		ws.Logger.Printf("wsWrite err:%s", err.Error())
		return err
	}

	if v, ok := ws.conf.subscribeMsg.Load(stream); ok {
		if op != nil && op.IsReConnect {
			return nil
		}
		reqs := v.([]requestHistory)
		reqs = append(reqs, requestHistory{
			Stream: stream,
			Event:  event,
		})
		ws.conf.subscribeMsg.Store(stream, reqs)
	} else {
		ws.conf.subscribeMsg.Store(stream, []requestHistory{{
			Stream: stream,
			Event:  event,
		}})
	}

	return nil
}

// readMsg only run once to read message
func (ws *WsService) readMsg() {
	ws.once.Do(func() {
		go func() {
			defer ws.Client.Close()

			for {
				select {
				case <-ws.Ctx.Done():
					ws.Logger.Printf("closing reader")
					return
				default:
					_, message, err := ws.Client.ReadMessage()
					if err != nil {
						ne, isNetErr := err.(net.Error)
						noe, isNetOpErr := err.(*net.OpError)
						if websocket.IsUnexpectedCloseError(err) || (isNetErr && ne.Timeout()) || (isNetOpErr && noe != nil) ||
							websocket.IsCloseError(err) || io.ErrUnexpectedEOF == err {
							ws.Logger.Printf("websocket err:%s", err.Error())
							if e := ws.reconnect(); e != nil {
								ws.Logger.Printf("reconnect err:%s", err.Error())
								return
							} else {
								ws.Logger.Printf("reconnect success, continue read message")
								continue
							}
						} else {
							ws.Logger.Printf("wsRead err:%s, type:%T", err.Error(), err)
							return
						}
					}
					var rawTrade UpdateMsg
					if err := json.Unmarshal(message, &rawTrade); err != nil {
						ws.Logger.Printf("Unmarshal err:%s, body:%s", err.Error(), string(message))
						continue
					}

					if bch, ok := ws.msgChs.Load(rawTrade.Stream); ok {
						bch.(chan *UpdateMsg) <- &rawTrade
					}
				}
			}
		}()
	})
}

type callBack func(*UpdateMsg)

func NewCallBack(f func(*UpdateMsg)) func(*UpdateMsg) {
	return f
}

func (ws *WsService) SetCallBack(stream string, call callBack) {
	if call == nil {
		return
	}
	ws.calls.Store(stream, call)
}

func (ws *WsService) receiveCallMsg(stream string, msgCh chan *UpdateMsg) {
	defer close(msgCh)
	for {
		select {
		case <-ws.Ctx.Done():
			ws.Logger.Printf("received parent context exit")
			return
		case msg := <-msgCh:
			if call, ok := ws.calls.Load(stream); ok {
				call.(callBack)(msg)
			}
		}
	}
}
