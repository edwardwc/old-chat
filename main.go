package main

import (
	"e2ee-chat/handlers"
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:              ":6969",
		Handler:           handlers.NewProxyMux(),
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	fmt.Println("Listening on :6969")
	server.ListenAndServe()
}
