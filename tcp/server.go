package tcp

// server.go is a TCP server

import (
	"net"
	"time"

	"github.com/mizumoto-cn/goredis/interface/tcp"
)

type Config struct {
	Address string        `json:"address"`
	Timeout time.Duration `json:"timeout"`
	MaxConn uint32        `json:"max_conn"`
}

// ListenAndServe starts a TCP server.
func ListenAndServe(lis net.Listener, han tcp.Handler, close <-chan struct{}) error {
	// listen close signal
	go func() {
		<-close
		log.info("TCP server is closing")
		lis.Close()
		han.Close()
	}()
	defer func() {
		// when unexpected errors occur, close the listener and handler
		lis.Close()
		han.Close()
	}()
	// Todo context & waitgroup
	// Todo for listen conn handle
}
