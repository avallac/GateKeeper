package ipguard

import (
	"testing"
	"time"

	"github.com/avallac/gatekeeper/internal/service"
	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
	"github.com/google/uuid"
)

func getClient(userConfig *service.UserConfig, appConfig *service.AppConfig) (*Client, *FakeServer) {
	authClient, authServer := createPipe()
	mainClient, mainServer := createPipe()
	fakeServer := RunFakeServer(authServer, mainClient, "pubKey")
	client := Run(authClient, mainServer, userConfig, appConfig)
	return client, fakeServer
}

func TestEmptyConfig(t *testing.T) {
	userConfig := &service.UserConfig{
		Sender: uuid.NewString(),
	}
	client, _ := getClient(userConfig, &service.AppConfig{})
	if client.status != AppStatusInputPhone {
		t.Errorf("Default status = %d, want %d", client.status, AppStatusInputPhone)
	}
}

func TestCheckingAuth(t *testing.T) {
	tests := []struct {
		name       string
		result     int
		wantStatus AppStatus
	}{
		{"Successful", 1, AppStatusInputCode},
		{"Failed", 0, AppStatusError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userConfig := &service.UserConfig{
				Sender: uuid.NewString(),
				Name:   "Test",
				Phone:  "123",
			}
			client, server := getClient(userConfig, &service.AppConfig{})
			server.SpSrvCheckAuthResult = tt.result
			client.AskCode()
			time.Sleep(100 * time.Millisecond)
			if client.status != tt.wantStatus {
				t.Errorf("status = %d, want %d", client.status, tt.wantStatus)
			}
		})
	}
}

func TestValidateAuth(t *testing.T) {
	tests := []struct {
		name       string
		result     int
		wantStatus AppStatus
	}{
		{"Successful", 1, AppStatusConnected},
		{"Failed", 0, AppStatusError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userConfig := &service.UserConfig{
				Sender: uuid.NewString(),
				Name:   "Test",
				Phone:  "123",
			}
			client, server := getClient(userConfig, &service.AppConfig{})
			server.SrvValidateAuthResult = tt.result
			client.ValidateCode("123456")
			time.Sleep(100 * time.Millisecond)
			if client.status != tt.wantStatus {
				t.Errorf("status = %d, want %d", client.status, tt.wantStatus)
			}
		})
	}
}

func TestOnRelaysUpdatedHook(t *testing.T) {
	userConfig := &service.UserConfig{
		Sender: uuid.NewString(),
		Name:   "Test",
		Phone:  "123",
	}
	client, server := getClient(userConfig, &service.AppConfig{})
	server.SrvValidateAuthResult = 1

	var gotSites []AppSite
	client.OnRelaysUpdated(func(sites []AppSite) {
		gotSites = sites
	})

	client.ValidateCode("123456")
	time.Sleep(100 * time.Millisecond)

	if len(gotSites) != 1 || gotSites[0].Name != "Test Place" {
		t.Fatalf("OnRelaysUpdated sites = %+v, want one site named \"Test Place\"", gotSites)
	}
	if len(gotSites[0].Boards) != 1 || gotSites[0].Boards[0].Identity != "00000000-1111-2222-3333-444444444444" {
		t.Fatalf("OnRelaysUpdated boards = %+v, want the fixture board", gotSites[0].Boards)
	}
}

func TestOnCallLogHook(t *testing.T) {
	userConfig := &service.UserConfig{
		Sender: uuid.NewString(),
		Name:   "Test",
		Phone:  "123",
	}
	client, server := getClient(userConfig, &service.AppConfig{})
	server.SrvValidateAuthResult = 1
	client.ValidateCode("123456")
	time.Sleep(100 * time.Millisecond)
	if client.status != AppStatusConnected {
		t.Fatalf("precondition failed: status = %d, want %d (Connected)", client.status, AppStatusConnected)
	}

	var gotBoardID string
	var gotBody message.BrdSpCallLogBody
	client.OnCallLog(func(boardID string, _ AppLog, body message.BrdSpCallLogBody) {
		gotBoardID = boardID
		gotBody = body
	})

	msg := message.New(message.BrdSpCallLogHeader, "00000000-1111-2222-3333-444444444444", message.BrdSpCallLogBody{
		BaseBody:      message.NewBaseBody(),
		ResponderName: "Front Door",
		Timestring:    "2026-10-02 14:32:00",
	})
	if err := client.handleBrdSpCallLog(msg); err != nil {
		t.Fatalf("handleBrdSpCallLog() error = %v", err)
	}

	if gotBoardID != "00000000-1111-2222-3333-444444444444" {
		t.Errorf("OnCallLog boardID = %q, want fixture board id", gotBoardID)
	}
	if gotBody.ResponderName != "Front Door" {
		t.Errorf("OnCallLog body.ResponderName = %q, want %q", gotBody.ResponderName, "Front Door")
	}
}
