package message

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Header struct {
	Event     string `json:"event"`
	Sender    string `json:"sender"`
	Timestamp int64  `json:"timestamp"`
}

type Message[T any] struct {
	Header Header `json:"header"`
	Body   T      `json:"body"`
}

type BaseBody struct {
	Token string `json:"token"`
}

func NewBaseBody() BaseBody {
	return BaseBody{Token: uuid.NewString()}
}

func New[T any](event, sender string, body T) Message[T] {
	return Message[T]{
		Header: Header{
			Event:     event,
			Sender:    sender,
			Timestamp: time.Now().UnixMilli(),
		},
		Body: body,
	}
}

func Encode[T any](msg Message[T]) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Decode[T any](data []byte) (*Message[T], error) {
	var msg Message[T]
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
