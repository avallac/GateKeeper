package ipguard

import (
	"log"

	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

type FakeServerAPIMain struct {
	ZeroMQ
}

func runFakeServerAPIMain(socket Socket) *FakeServerAPIMain {
	api := FakeServerAPIMain{}
	api.socket = socket
	api.socket.SetIdentity("main")
	api.socket.SetCurveServer(1)

	port, err := api.socket.Bind("tcp://*:5577")
	if err != nil {
		log.Fatalf("bind %s: %v", "tcp://*:5577", err)
	}
	log.Printf("main API port %d", port)

	go api.handler()
	return &api
}

func (api *FakeServerAPIMain) SendBcSpRelays(target []byte) error {
	jsonStr, err := message.CreateBcSpRelaysMessage()
	if err != nil {
		return err
	}
	return api.send([][]byte{target, jsonStr})
}

func (api *FakeServerAPIMain) SendBrdSpOpenDoorAck(target []byte, source string) error {
	jsonStr, err := message.CreateBrdSpOpenDoorAckMessage(source)
	if err != nil {
		return err
	}
	return api.send([][]byte{target, jsonStr})
}
