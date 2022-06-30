package tcp

import (
	"context"
	"net"
)

// Handler is a TCP server handler.
// represents application logic over tcp.
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

// HandlerFunc represents a handler function
// that will be called when a connection is established.
type HandlerFunc func(ctx context.Context, conn net.Conn)
