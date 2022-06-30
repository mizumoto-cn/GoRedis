package tcp

// server.go is a TCP server

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mizumoto-cn/goredis/interface/tcp"
	"github.com/mizumoto-cn/goredis/lib/logs"
)

type Config struct {
	Address string        `json:"address"`
	Timeout time.Duration `json:"timeout"`
	MaxConn uint32        `json:"max_conn"`
}

// ListenAndServe handles requests, blocks until close signal is received
func ListenAndServe(listener net.Listener, handler tcp.Handler, close <-chan struct{}) error {
	// listen close signal
	go func() {
		<-close
		logs.Info("TCP server is closing")
		listener.Close()
		handler.Close()
	}()
	defer func() {
		// when unexpected errors occur, close the listener and handler
		listener.Close()
		handler.Close()
	}()

	ctx := context.Background()
	var waitGroup sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logs.Error("TCP request is canceled %v", err)
				break
			}
			logs.Error("TCP server accept error: %v", err)
			break
		}
		logs.Info("TCP server accept a connection from %s", conn.RemoteAddr().String())
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			handler.Handle(ctx, conn)
		}()
	}
	waitGroup.Wait()
}

// ListenAndServeWithConfig handles requests, blocks until close signal is received
func ListenAndServeWithConfig(config Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	// buffered channel to avoid leaking goroutines
	signalChan := make(chan os.Signal, 1)
	// Notify passing signal from system-call to the signal channel
	// https://go.googlesource.com/proposal/+/refs/heads/master/design/freeze-syscall.md
	// For more posix signals, see: https://dsa.cs.tsinghua.edu.cn/oj/static/unix_signal.html
	// Here it passes Interrupt, Termination, Quit from Keyboard, Hangup signals.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// listen close signal and pass it to the close channel in a separate goroutine
	go func() {
		signal := <-signalChan
		switch signal {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, SIGHUP:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}

	logs.Info("TCP server is listening on %s", config.Address)
	ListenAndServe(listener, handler, closeChan)
	return nil
}
