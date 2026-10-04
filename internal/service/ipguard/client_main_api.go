package ipguard

import (
	"log"

	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

type GuardAPIMain struct {
	ZeroMQ
	status Status
	url    string
}

func (api *GuardAPIMain) OpenDoor(id string) {
	if api.status == StatusNotConnected {
		log.Panicln("Main api isn't connected. Message isn't sent.")
		return
	}
	msg, err := message.CreateOpenDoorMessage(api.socket.Identity())
	if err != nil {
		log.Fatal(err)
	}
	api.send([][]byte{[]byte(id), msg})
}

func (api *GuardAPIMain) CheckUser(phone, code, name string) error {
	if api.status == StatusNotConnected {
		log.Panicln("Main api isn't connected. Message isn't sent.")
		return nil
	}
	msg, err := message.CreateSpSrvCheckUsernameMessage(api.socket.Identity(), phone, code, name)
	if err != nil {
		log.Fatal(err)
	}
	return api.send([][]byte{msg})
}

func (api *GuardAPIMain) UpdateUsernameAck(token string) error {
	if api.status == StatusNotConnected {
		log.Panicln("Main api isn't connected. Message isn't sent.")
		return nil
	}
	msg, err := message.CreateSrvUpdateAckMessage(api.socket.Identity(), token)
	if err != nil {
		log.Fatal(err)
	}
	return api.send([][]byte{msg})
}

func (api *GuardAPIMain) CallLogAck(token string) error {
	if api.status == StatusNotConnected {
		log.Panicln("Main api isn't connected. Message isn't sent.")
		return nil
	}
	msg, err := message.CreateSpSrvCallLogAckMessage(api.socket.Identity(), token)
	if err != nil {
		log.Fatal(err)
	}
	return api.send([][]byte{msg})
}

func (api *GuardAPIMain) OverrideUser(phone, code, name string) error {
	if api.status == StatusNotConnected {
		log.Panicln("Main api isn't connected. Message isn't sent.")
		return nil
	}
	msg, err := message.CreateSpSrvOverrideUsernameBodyMessage(api.socket.Identity(), phone, code, name)
	if err != nil {
		log.Fatal(err)
	}
	return api.send([][]byte{msg})
}

func (api *GuardAPIMain) Connect(key, name, phone, code string) {
	api.socket.SetCurveServerkey(key)
	err := api.socket.Connect(api.url)
	if err != nil {
		panic(err)
	}
	jsonStr, err := message.CreateHelloMessage(
		api.socket.Identity(),
		version,
		name,
		phone,
		code,
	)
	if err != nil {
		panic(err)
	}
	api.send([][]byte{jsonStr})
	api.status = StatusConnected
}

func InitApi(socket Socket, sender, url string) *GuardAPIMain {
	api := GuardAPIMain{
		status: StatusNotConnected,
		url:    url,
	}
	api.socket = socket
	api.name = "MAIN"

	api.socket.SetIdentity(sender)
	go api.handler()
	return &api
}
