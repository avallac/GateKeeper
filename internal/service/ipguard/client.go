package ipguard

import (
	"log"
	"slices"
	"time"

	"github.com/avallac/gatekeeper/internal/service"
	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

type AppStatus int

const (
	AppStatusStarting       AppStatus = iota
	AppStatusInputPhone     AppStatus = iota
	AppStatusCheckingAuth   AppStatus = iota
	AppStatusInputCode      AppStatus = iota
	AppStatusValidatingAuth AppStatus = iota
	AppStatusWaitingKey     AppStatus = iota
	AppStatusReady          AppStatus = iota
	AppStatusHelloSent      AppStatus = iota
	AppStatusConnected      AppStatus = iota
	AppStatusError          AppStatus = iota
)

const maxLogSize = 50

type RelaysUpdatedFunc func(sites []AppSite)
type CallLogFunc func(boardID string, log AppLog, body message.BrdSpCallLogBody)

type Client struct {
	UserConfig      *service.UserConfig
	AppConfig       *service.AppConfig
	Sites           []AppSite
	Log             []AppLog
	loginApi        *GuardAPILogin
	mainApi         *GuardAPIMain
	status          AppStatus
	checkingCode    string
	onRelaysUpdated []RelaysUpdatedFunc
	onCallLog       []CallLogFunc
}

type AppSite struct {
	Name   string
	Boards []AppBoard
}

type AppLog struct {
	Title string
	Time  string
	Photo string
	Open  bool
}

type AppBoard struct {
	Identity string
	Name     string
}

func Run(authSocket, mainSocket Socket, userConfig *service.UserConfig, appConfig *service.AppConfig) *Client {
	app := Client{
		status: AppStatusStarting,
	}
	app.AppConfig = appConfig
	app.UserConfig = userConfig
	app.loginApi = InitLoginApi(authSocket, app.UserConfig.Sender, app.AppConfig.AuthAddress)
	app.mainApi = InitApi(mainSocket, app.UserConfig.Sender, app.AppConfig.MainAddress)
	Bind(app.loginApi, message.SrvSpCheckAuthResultHeader, app.handleSrvSpCheckAuthResult)
	Bind(app.loginApi, message.SrvSpValidateAuthResultHeder, app.handleSrvSpValidateAuthResultMessage)
	Bind(app.loginApi, message.SrvSpAskKeyResultHeader, app.handleSrvSpAskKeyResultMessage)
	Bind(app.mainApi, message.BcSpRelaysHeader, app.handleBcSpRelays)
	Bind(app.mainApi, message.BrdSpOpenDoorAckHeader, app.handleSpOpenDoorAck)
	Bind(app.mainApi, message.BrdSpCallLogHeader, app.handleBrdSpCallLog)
	Bind(app.mainApi, message.PsSpSdpOfferHeader, app.handlePsSpSdpOffer)
	Bind(app.mainApi, message.SrvSpDeleteHeader, app.handleSrvSpDelete)
	Bind(app.mainApi, message.SrvSpCheckUsernameResultHeader, app.handleSrvSpCheckUsernameResult)
	Bind(app.mainApi, message.SrvSpUpdateHeader, app.handleSrvSpUpdate)
	if app.AppConfig.FakeServer {
		app.UserConfig.Key = ""
	}
	if app.UserConfig.Name == "" || app.UserConfig.Phone == "" {
		app.UpdateStatus(AppStatusInputPhone)
	} else if app.UserConfig.Code == "" {
		app.UpdateStatus(AppStatusInputPhone)
	} else if app.UserConfig.Key == "" {
		app.UpdateStatus(AppStatusWaitingKey)
		app.loginApi.AskKey(app.UserConfig.Phone, app.UserConfig.Code)
	} else {
		app.UpdateStatus(AppStatusReady)
	}

	return &app
}

func (app *Client) handleSrvSpCheckAuthResult(msg message.SrvSpAskKeyResultMessage) error {
	if msg.Body.Result == 1 {
		app.UpdateStatus(AppStatusInputCode)
	} else {
		app.UpdateStatus(AppStatusError)
	}
	return nil
}

func (app *Client) handleSrvSpAskKeyResultMessage(msg message.SrvSpAskKeyResultMessage) error {
	if msg.Body.Result == 1 {
		app.UserConfig.UpdateKey(msg.Body.Key)
		app.UpdateStatus(AppStatusReady)
	} else {
		app.UpdateStatus(AppStatusError)
	}
	return nil
}

func (app *Client) handleSrvSpValidateAuthResultMessage(msg message.SrvSpAskKeyResultMessage) error {
	if msg.Body.Result == 1 {
		app.UserConfig.UpdateCode(app.checkingCode)
		app.loginApi.AskKey(app.UserConfig.Phone, app.UserConfig.Code)
		app.UpdateStatus(AppStatusWaitingKey)
	} else {
		app.UpdateStatus(AppStatusError)
	}
	return nil
}

func (app *Client) handleBcSpRelays(msg message.BcSpRelaysMessage) error {
	var appSites []AppSite
	for _, site := range msg.Body.Sites {
		appSite := AppSite{
			Name: site.Name,
		}
		var appBoards []AppBoard
		for _, board := range site.Boards {
			appBoard := AppBoard{
				Identity: board.Identity,
				Name:     board.Name,
			}
			appBoards = append(appBoards, appBoard)
		}
		appSite.Boards = appBoards
		appSites = append(appSites, appSite)
	}
	app.Sites = appSites
	app.UpdateStatus(AppStatusConnected)
	for _, fn := range app.onRelaysUpdated {
		fn(app.Sites)
	}
	return nil
}

func (app *Client) handleSpOpenDoorAck(msg message.BrdSpOpenDoorAckMessage) error {
	logTime := time.Now()
	event := AppLog{}
	event.Title = "Door \"" + app.getDoorNameById(msg.Header.Sender) + "\" was open"
	event.Time = logTime.Format("2006-01-02 15:04:05")
	app.appendLog(event)
	return nil
}

func (app *Client) handleBrdSpCallLog(msg message.BrdSpCallLogMessage) error {
	logTime := time.Unix(msg.Body.Timestamp, 0)
	event := AppLog{}
	event.Title = "Incoming call from \"" + app.getDoorNameById(msg.Header.Sender) + "\""
	event.Photo = msg.Body.PhotoData
	event.Time = logTime.Format("2006-01-02 15:04:05")
	app.appendLog(event)
	for _, fn := range app.onCallLog {
		fn(msg.Header.Sender, event, msg.Body)
	}
	return app.mainApi.CallLogAck(msg.Body.Token)
}

func (app *Client) handlePsSpSdpOffer(msg message.PsSpSdpOfferMessage) error {
	app.OpenDoor(msg.Body.RemoteIdentity)
	return nil
}

func (app *Client) handleSrvSpDelete(msg message.SrvSpDeleteMessage) error {
	return app.mainApi.CheckUser(app.UserConfig.Phone, app.UserConfig.Code, app.UserConfig.Name)
}

func (app *Client) handleSrvSpCheckUsernameResult(msg message.SrvSpCheckUsernameResultMessage) error {
	if msg.Body.Result == 0 {
		return app.Override()
	}
	return nil
}

func (app *Client) handleSrvSpUpdate(msg message.SrvSpUpdateMessage) error {
	return app.mainApi.UpdateUsernameAck(msg.Body.Token)
}

func (app *Client) CheckUser() error {
	return app.mainApi.CheckUser(app.UserConfig.Phone, app.UserConfig.Code, app.UserConfig.Name)
}

func (app *Client) Override() error {
	return app.mainApi.OverrideUser(app.UserConfig.Phone, app.UserConfig.Code, app.UserConfig.Name)
}

func (app *Client) UpdateStatus(status AppStatus) {
	log.Printf("Status %d -> %d", app.status, status)
	app.status = status
	if status == AppStatusReady {
		app.Connect()
	}
}

func (app *Client) Connect() {
	if app.status == AppStatusReady {
		app.mainApi.Connect(app.UserConfig.Key, app.UserConfig.Name, app.UserConfig.Phone, app.UserConfig.Code)
		app.UpdateStatus(AppStatusHelloSent)
	}
}

func (app *Client) ValidateCode(code string) {
	app.checkingCode = code
	app.loginApi.ValidateCode(app.UserConfig.Phone, code)
	app.UpdateStatus(AppStatusValidatingAuth)
}

func (app *Client) IsConnected() bool {
	return app.status == AppStatusConnected
}

func (app *Client) IsConnecting() bool {
	return app.status == AppStatusReady
}

func (app *Client) UpdateUserInfo(name, phone string) {
	app.UserConfig.UpdateUserInfo(name, phone)
}

func (app *Client) AskCode() {
	app.loginApi.AskCode(app.UserConfig.Phone)
	app.UpdateStatus(AppStatusCheckingAuth)
}

func (app *Client) OpenDoor(id string) {
	now := time.Now()
	event := AppLog{}
	event.Title = "Request sent to open door \"" + app.getDoorNameById(id) + "\""
	event.Time = now.Format("2006-01-02 15:04:05")
	event.Open = true
	app.appendLog(event)
	app.mainApi.OpenDoor(id)
}

func (app *Client) IsNeededPhone() bool {
	return app.status == AppStatusInputPhone
}

func (app *Client) IsNeededCode() bool {
	return app.status == AppStatusInputCode
}

func (app *Client) IsWebWaiting() bool {
	waitingStatuses := []AppStatus{AppStatusCheckingAuth, AppStatusValidatingAuth, AppStatusWaitingKey, AppStatusReady, AppStatusHelloSent}
	return slices.Contains(waitingStatuses, app.status)
}

func (app *Client) OnRelaysUpdated(fn RelaysUpdatedFunc) {
	app.onRelaysUpdated = append(app.onRelaysUpdated, fn)
}

func (app *Client) OnCallLog(fn CallLogFunc) {
	app.onCallLog = append(app.onCallLog, fn)
}

func (app *Client) appendLog(event AppLog) {
	app.Log = append(app.Log, event)
	if len(app.Log) > maxLogSize {
		app.Log = app.Log[len(app.Log)-maxLogSize:]
	}
}

func (app *Client) getDoorNameById(id string) string {
	for _, site := range app.Sites {
		for _, board := range site.Boards {
			if board.Identity == id {
				return board.Name
			}
		}
	}

	return ""
}
