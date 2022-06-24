package tcp

import (
	"context"
	"net"
)

// Handler is a TCP server handler.
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

// HandlerFunc is a TCP server handler function.
type HandlerFunc func(ctx context.Context, conn net.Conn)
