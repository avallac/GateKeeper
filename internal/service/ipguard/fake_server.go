package ipguard

import (
	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

type FakeServer struct {
	loginApi              *FakeServerAPILogin
	mainApi               *FakeServerAPIMain
	pubKey                string
	SpSrvCheckAuthResult  int
	SrvValidateAuthResult int
}

func RunFakeServer(socketAuth, socketMain Socket, pubKey string) *FakeServer {
	server := FakeServer{
		pubKey:                pubKey,
		SpSrvCheckAuthResult:  1,
		SrvValidateAuthResult: 1,
	}
	server.loginApi = runFakeServerAPILogin(socketAuth)
	server.mainApi = runFakeServerAPIMain(socketMain)
	Bind(server.loginApi, message.SpSrvCheckAuthHeader, server.handleSpSrvCheckAuth)
	Bind(server.loginApi, message.SpSrvValidateAuthHeader, server.handleSpSrvValidateAuth)
	Bind(server.loginApi, message.SpSrvAskKeyHeader, server.handleSpSrvAskKey)
	Bind(server.mainApi, message.HelloHeader, server.handleHello)
	BindWithTarget(server.mainApi, message.SpBrdOpenDoorHeader, server.handleSpBrdOpenDoor)
	return &server
}

func (s *FakeServer) handleSpSrvCheckAuth(msg message.Message[message.SpSrvCheckAuthBody]) error {
	return s.loginApi.SendSrvSpCheckAuthResult([]byte(msg.Header.Sender), s.SpSrvCheckAuthResult)
}

func (s *FakeServer) handleSpSrvValidateAuth(msg message.SpSrvAskKeyMessage) error {
	return s.loginApi.SendSpSrvValidateAuth([]byte(msg.Header.Sender), s.SrvValidateAuthResult)
}

func (s *FakeServer) handleSpSrvAskKey(msg message.SpSrvAskKeyMessage) error {
	return s.loginApi.SendSrvSpAskKeyResult([]byte(msg.Header.Sender), s.pubKey)
}

func (s *FakeServer) handleHello(msg message.SpSrvAskKeyMessage) error {
	return s.mainApi.SendBcSpRelays([]byte(msg.Header.Sender))
}

func (s *FakeServer) handleSpBrdOpenDoor(target string, msg message.SpSrvAskKeyMessage) error {
	return s.mainApi.SendBrdSpOpenDoorAck([]byte(msg.Header.Sender), target)
}
