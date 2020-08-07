package main

import (
	handler "panda-data/src/redis"
	"panda-data/src/server"
	"time"
)

func main() {
	conf := server.Config{
		Address:    "127.0.0.1:6379",
		MaxConnect: 12,
		Timeout:    time.Second * 2,
	}

	server.ListenAndServer(&conf, handler.MakeHandler())
}
