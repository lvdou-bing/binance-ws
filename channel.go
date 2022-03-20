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
@param channel string: the stream type
@param payload []string: the list of detail stream message, e.g. <symbol>@type_<option>
*/
func (ws *WsService) Subscribe(channel string, payload []string) error {
	// if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[channel] {
	// 	return newAuthEmptyErr()
	// }

	msgCh, ok := ws.msgChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsgRaw, 1)
		go ws.receiveCallMsg(channel, msgCh.(chan *UpdateMsgRaw))
	}

	err := ws.newBaseChannel(channel, payload, msgCh.(chan *UpdateMsgRaw), nil)
	if err != nil {
		return err
	}

	return nil
}

func (ws *WsService) SubscribeWithOption(channel string, payload []string, op *SubscribeOptions) error {
	// if (ws.conf.Key == "" || ws.conf.Secret == "") && authChannel[channel] {
	// 	return newAuthEmptyErr()
	// }

	msgCh, ok := ws.msgChs.Load(channel)
	if !ok {
		msgCh = make(chan *UpdateMsgRaw, 1)
		go ws.receiveCallMsg(channel, msgCh.(chan *UpdateMsgRaw))
	}

	err := ws.newBaseChannel(channel, payload, msgCh.(chan *UpdateMsgRaw), op)
	if err != nil {
		return err
	}

	return nil
}

func (ws *WsService) UnSubscribe(channel string, payload []string) error {
	return ws.baseSubscribe(UnSubscribe, channel, payload, nil)
}

func (ws *WsService) newBaseChannel(channel string, payload []string, bch chan *UpdateMsgRaw, op *SubscribeOptions) error {
	err := ws.baseSubscribe(Subscribe, channel, payload, op)
	if err != nil {
		return err
	}

	if _, ok := ws.msgChs.Load(channel); !ok {
		ws.msgChs.Store(channel, bch)
	}

	ws.readMsg()

	return nil
}

func (ws *WsService) baseSubscribe(event string, channel string, payload []string, op *SubscribeOptions) error {
	req := make(map[string]interface{})
	req["method"] = event
	req["params"] = payload
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

	if v, ok := ws.conf.subscribeMsg.Load(channel); ok {
		if op != nil && op.IsReConnect {
			return nil
		}
		reqs := v.([]requestHistory)
		reqs = append(reqs, requestHistory{
			Channel: channel,
			Event:   event,
			Payload: payload,
		})
		ws.conf.subscribeMsg.Store(channel, reqs)
	} else {
		ws.conf.subscribeMsg.Store(channel, []requestHistory{{
			Channel: channel,
			Event:   event,
			Payload: payload,
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

					if bch, ok := ws.msgChs.Load(rawTrade.Channel); ok {
						bch.(chan *UpdateMsgRaw) <- (*UpdateMsgRaw)(&message)
					}
				}
			}
		}()
	})
}

type callBack func(*UpdateMsgRaw)

func NewCallBack(f func(*UpdateMsgRaw)) func(*UpdateMsgRaw) {
	return f
}

func (ws *WsService) SetCallBack(channel string, call callBack) {
	if call == nil {
		return
	}
	ws.calls.Store(channel, call)
}

func (ws *WsService) receiveCallMsg(channel string, msgCh chan *UpdateMsgRaw) {
	defer close(msgCh)
	for {
		select {
		case <-ws.Ctx.Done():
			ws.Logger.Printf("received parent context exit")
			return
		case msg := <-msgCh:
			if call, ok := ws.calls.Load(channel); ok {
				call.(callBack)(msg)
			}
		}
	}
}
