package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"panda-data/src/interface/tcp"
	"panda-data/src/lib/logger"
	"panda-data/src/lib/sync/atomic"
	"syscall"
	"time"
)

type Config struct {
	Address    string
	MaxConnect uint32
	Timeout    time.Duration
}

func ListenAndServer(cfg *Config, handler tcp.Handler) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Error()
	}
	var closing atomic.BoolAtomic
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.Info("shutdown ing ...")
			closing.Set(true)
			er := listener.Close()
			logger.Fatal(er)
		}
	}()
	logger.Info(fmt.Sprintf("start lisrening bind: %s", cfg.Address))
	defer func() {
		err1 := handler.Close()
		if err1 != nil {
			logger.Info(err1)
		}
	}()
	defer listener.Close()

	ctx, _ := context.WithCancel(context.Background())
	for {
		conn, err := listener.Accept()
		if err != nil {
			if closing.Get() {
				return
			}
			logger.Error(fmt.Sprintf("accept err: %v", err))
			continue
		}
		logger.Info("accept success")
		go handler.Handler(ctx, conn)
	}

}
