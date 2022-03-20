package bnws

import (
	"encoding/json"
)

type UpdateMsg struct {
	// Time    int64           `json:"time"`
	// Id      *int64          `json:"id,omitempty"`
	// Channel string          `json:"channel"`
	Channel string `json:"e"`
	// Error   *ServiceError   `json:"error,omitempty"`
	// Result json.RawMessage `json:"result"`
}

type UpdateMsgRaw json.RawMessage

type ServiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ServiceError) Error() string {
	return e.Message
}

type WSEvent struct {
	UpdateMsg
}

type ChannelEvent struct {
	Event  string
	Market []string
}

type WebsocketRequest struct {
	Market []string
}

type Request struct {
	App     string   `json:"app,omitempty"`
	Time    int64    `json:"time"`
	Id      *int64   `json:"id,omitempty"`
	Channel string   `json:"channel"`
	Event   string   `json:"event"`
	Auth    Auth     `json:"auth"`
	Payload []string `json:"payload"`
}

type Auth struct {
	Method string `json:"method"`
	Key    string `json:"KEY"`
	Secret string `json:"SIGN"`
}

type requestHistory struct {
	Channel string   `json:"channel"`
	Event   string   `json:"event"`
	Payload []string `json:"payload"`
	op      *SubscribeOptions
}
