package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := Load()

	mailer := NewMailer(cfg)
	handler := NewHandler(mailer, cfg.AuthToken)

	addr := ":" + cfg.ServerPort
	fmt.Printf("msg_forwarder listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}