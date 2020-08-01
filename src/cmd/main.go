package main

import (
	handler "godis/src/redis"
	"godis/src/server"
	"time"
)

func main() {
	conf := server.Config{
		Address:    "127.0.0.1:80",
		MaxConnect: 12,
		Timeout:    time.Second * 2,
	}

	server.ListenAndServer(&conf, handler.MakeHandler())
}
