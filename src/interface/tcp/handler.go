package tcp

import (
	"context"
	"net"
)

type HandlerFunc func(ctx context.Context, conn net.Conn)

type Handler interface {
	Handler(ctx context.Context, conn net.Conn)
	Close() error
}
