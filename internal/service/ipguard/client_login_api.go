package ipguard

import (
	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

type GuardAPILogin struct {
	ZeroMQ
	status Status
	url    string
}

func (api *GuardAPILogin) AskCode(phone string) {
	if api.status == StatusNotConnected {
		api.connect()
	}

	jsonStr, err := message.CreateSpSrvCheckAuthMessage(api.socket.Identity(), phone)
	if err != nil {
		panic(err)
	}
	api.send([][]byte{jsonStr})
}

func (api *GuardAPILogin) ValidateCode(phone, code string) {
	if api.status == StatusNotConnected {
		api.connect()
	}

	jsonStr, err := message.CreateSpSrvValidateAuthMessage(api.socket.Identity(), phone, code)
	if err != nil {
		panic(err)
	}
	api.send([][]byte{jsonStr})
}

func (api *GuardAPILogin) AskKey(phone, code string) {
	if api.status == StatusNotConnected {
		api.connect()
	}

	jsonStr, err := message.CreateSpSrvAskKeyMessage(api.socket.Identity(), phone, code)
	if err != nil {
		panic(err)
	}
	api.send([][]byte{jsonStr})
}

func (api *GuardAPILogin) connect() {
	api.socket.Connect(api.url)
	jsonStr, err := message.CreateAuthHelloMessage(api.socket.Identity(), version)
	if err != nil {
		panic(err)
	}
	api.send([][]byte{jsonStr})
	api.status = StatusConnected
}

func InitLoginApi(socket Socket, sender string, url string) *GuardAPILogin {
	api := GuardAPILogin{
		status: StatusNotConnected,
		url:    url,
	}
	api.socket = socket
	api.name = "AUTH"
	api.socket.SetIdentity(sender)

	go api.handler()
	return &api
}
