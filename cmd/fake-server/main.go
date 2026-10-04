package main

import (
	"log"
	"net/http"

	"github.com/avallac/gatekeeper/internal/handlers"
	"github.com/avallac/gatekeeper/internal/service/ipguard"
	"github.com/zeromq/goczmq"
)

func main() {
	addr := ":8081"
	auth := goczmq.NewAuth()
	auth.Curve(goczmq.CurveAllowAny)

	socket := goczmq.NewSock(goczmq.Dealer)
	socket2 := goczmq.NewSock(goczmq.Dealer)
	cert := goczmq.NewCert()
	cert.SetMeta("name", "fake-server")
	cert.Apply(socket2)

	ipguard.RunFakeServer(socket, socket2, cert.PublicText())

	h := handlers.IpGuardFakeHandlers{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Index)

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Printf("ipguard fake server started on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
