package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/avallac/gatekeeper/internal/handlers"
	"github.com/avallac/gatekeeper/internal/service"
	"github.com/avallac/gatekeeper/internal/service/ipguard"
	"github.com/avallac/gatekeeper/internal/service/mqtt"
	"github.com/zeromq/goczmq"
)

func main() {
	appConfig, err := service.LoadAppConfig("configs/app.yaml")
	appConfig = service.OverwriteAppConfig(appConfig, "configs/app.local.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var userConfig *service.UserConfig
	if os.Getenv("SUPERVISOR_TOKEN") != "" {
		appConfig = service.LoadEnvOverrides(appConfig)
		appConfig = service.LoadAddonOptions(appConfig, "/data/options.json")
		appConfig = service.LoadSupervisorMQTT(appConfig, os.Getenv("SUPERVISOR_TOKEN"))
		userConfig = service.InitUserConfig("/data/config.yaml")
	} else {
		userConfig = service.InitUserConfig("configs/user.yaml")
	}

	authSocket, err := goczmq.NewDealer("")
	cert := goczmq.NewCert()
	mainSocket, err := goczmq.NewDealer("")
	cert.Apply(mainSocket)
	configureSocket(authSocket)
	configureSocket(mainSocket)

	app := ipguard.Run(authSocket, mainSocket, userConfig, appConfig)

	if _, err := mqtt.NewClient(appConfig, app); err != nil {
		log.Printf("mqtt disabled: %v", err)
	}

	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	h := handlers.Handlers{
		App:  app,
		Tmpl: tmpl,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/ask_code", h.AskCode)
	mux.HandleFunc("/send_code", h.SendCode)
	mux.HandleFunc("/open_door", h.OpenDoor)
	mux.HandleFunc("/check", h.Check)
	mux.HandleFunc("/override", h.Override)

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Printf("server started on %s", appConfig.Addr)
	err = http.ListenAndServe(appConfig.Addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func configureSocket(sock *goczmq.Sock) {
	sock.SetTcpKeepalive(1)
	sock.SetTcpKeepaliveIdle(60)
	sock.SetTcpKeepaliveIntvl(60)

	sock.SetHeartbeatIvl(15000)
	sock.SetHeartbeatTimeout(30000)

	sock.SetReconnectIvl(1000)
	sock.SetReconnectIvlMax(30000)
}
