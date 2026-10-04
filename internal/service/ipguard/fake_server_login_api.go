package ipguard

import (
	"log"

	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

type FakeServerAPILogin struct {
	ZeroMQ
}

func (api *FakeServerAPILogin) SendSpSrvValidateAuth(target []byte, result int) error {
	jsonStr, err := message.CreateSrvSpValidateAuthResultMessage(result)
	if err != nil {
		return err
	}
	return api.send([][]byte{target, jsonStr})
}

func (api *FakeServerAPILogin) SendSrvSpCheckAuthResult(target []byte, result int) error {
	jsonStr, err := message.CreateSrvSpCheckAuthResultMessage(result)
	if err != nil {
		return err
	}
	return api.send([][]byte{target, jsonStr})
}

func (api *FakeServerAPILogin) SendSrvSpAskKeyResult(target []byte, key string) error {
	jsonStr, err := message.CreateSrvSpAskKeyResultMessage(1, key)
	if err != nil {
		return err
	}
	return api.send([][]byte{target, jsonStr})
}

func runFakeServerAPILogin(socket Socket) *FakeServerAPILogin {
	api := FakeServerAPILogin{}
	api.socket = socket
	api.socket.SetIdentity("auth")
	port, err := api.socket.Bind("tcp://*:5579")
	if err != nil {
		log.Fatalf("bind %s: %v", "tcp://*:5579", err)
	}
	log.Printf("login API port %d", port)
	go api.handler()
	return &api
}
