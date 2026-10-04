package ipguard

import (
	"encoding/json"
	"log"

	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

const version = "3.4.6"

type Socket interface {
	SendMessage([][]byte) error
	RecvMessage() ([][]byte, error)
	Identity() string
	SetIdentity(string)
	Connect(string) error
	SetCurveServerkey(string)
	Bind(string) (int, error)
	SetCurveServer(int)
}

type ZeroMQ struct {
	name        string
	socket      Socket
	appHandlers map[string]Handler
}

type Status int

const (
	StatusNotConnected Status = iota
	StatusConnected    Status = iota
)

func (api *ZeroMQ) send(msg [][]byte) error {
	err := api.socket.SendMessage(msg)
	if len(msg) == 1 {
		log.Printf("[%s]-<- %s", api.name, msg[0])
	}
	if len(msg) == 2 {
		log.Printf("[%s]-<-[%s] %s", api.name, msg[0], msg[1])
	}
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (api *ZeroMQ) handler() {
	for true {
		var data []byte
		var target string
		msg, err := api.socket.RecvMessage()
		if err != nil {
			log.Fatal(err)
		}

		switch len(msg) {
		case 1:
			target = ""
			data = msg[0]
		case 2:
			target = string(msg[0])
			data = msg[1]
		default:
		}
		var resp message.GenericResponse
		json.Unmarshal(data, &resp)
		log.Printf("[%s]->- %s", api.name, string(data))
		if h, ok := api.appHandlers[resp.Header.Event]; ok {
			if err := h(target, data); err != nil {
				log.Printf("[AUTH] %s: %v", resp.Header.Event, err)
			}
		}
	}
}

type Handler func(target string, data []byte) error

type handlerRegistrar interface {
	handlers() map[string]Handler
}

func (api *ZeroMQ) handlers() map[string]Handler {
	if api.appHandlers == nil {
		api.appHandlers = make(map[string]Handler)
	}
	return api.appHandlers
}

func Bind[T any](api handlerRegistrar, event string, fn func(message.Message[T]) error) {
	api.handlers()[event] = func(_ string, data []byte) error {
		msg, err := message.Decode[T](data)
		if err != nil {
			return err
		}
		return fn(*msg)
	}
}

func BindWithTarget[T any](api handlerRegistrar, event string, fn func(string, message.Message[T]) error) {
	api.handlers()[event] = func(target string, data []byte) error {
		msg, err := message.Decode[T](data)
		if err != nil {
			return err
		}
		return fn(target, *msg)
	}
}
