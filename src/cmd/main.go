package main

import (
	handlerEcho "godis/src/echo/handler"
	"godis/src/server"
	"time"
)

func main() {
	conf := server.Config{
		Address:    "127.0.0.1:23",
		MaxConnect: 12,
		Timeout:    time.Second * 2,
	}

	server.ListenAndServer(&conf, handlerEcho.MakeEchoHandler())
}
