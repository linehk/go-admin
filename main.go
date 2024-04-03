package main

import (
	"log"
	"net/http"
	"time"

	"github.com/linehk/go-admin/config"
	"github.com/linehk/go-admin/controller"
)

func main() {
	config.Setup()
	handler := controller.Setup()

	server := &http.Server{
		Addr:           config.Raw.String("ADDR"),
		Handler:        handler,
		ReadTimeout:    time.Duration(config.Raw.Int("READ_TIMEOUT") * int(time.Second)),
		WriteTimeout:   time.Duration(config.Raw.Int("WRITE_TIMEOUT") * int(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
