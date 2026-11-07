package message

import (
	"encoding/json"
	"fmt"
	"strconv"
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

type FlexBool bool

func (b *FlexBool) UnmarshalJSON(data []byte) error {
	var v bool
	if err := json.Unmarshal(data, &v); err == nil {
		*b = FlexBool(v)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("cannot unmarshal %s into FlexBool", data)
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("cannot unmarshal %q into FlexBool", s)
	}
	*b = FlexBool(v)
	return nil
}

type FlexInt int

func (i *FlexInt) UnmarshalJSON(data []byte) error {
	var v int
	if err := json.Unmarshal(data, &v); err == nil {
		*i = FlexInt(v)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("cannot unmarshal %s into FlexInt", data)
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("cannot unmarshal %q into FlexInt", s)
	}
	*i = FlexInt(v)
	return nil
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
