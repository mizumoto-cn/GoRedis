// loporitts serves as a satellite server which indicates whether the server works correctly.
// The name Loporitts is from FFXIV, rabbit like creatures working as guardians on the moon.
package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/mizumoto-cn/goredis/lib/logs"
	"github.com/mizumoto-cn/goredis/lib/sync/atomic"
	"github.com/mizumoto-cn/goredis/lib/sync/wrappedwait"
)

// HummingWay is the echo handler
// echos the received message back to the client
type HummingWay struct {
	// map as a set, storing all working clients
	// safe to use concurrently
	activeConns sync.Map
	isClosing   atomic.Boolean
}

// NewHummingWay creates a new HummingWay
func NewHummingWay() *HummingWay {
	return &HummingWay{}
}

// Close() close the HummingWay
func (h *HummingWay) Close() error {
	return h.shut()
}

func (h *HummingWay) shut() error {
	logs.Info("shutting down")
	h.isClosing.Set(true)
	// set client in the set to close
	h.activeConns.Range(func(key, value interface{}) bool {
		client := key.(*SurfingWay)
		client.Close()
		return true
	})
	return nil
}

// SurfingWay is the echo client
type SurfingWay struct {
	conn    net.Conn
	Waiting wrappedwait.WrappedWait
}

// Close() close after timeout
func (c *SurfingWay) Close() error {
	c.Waiting.WaitTimeout(5 * time.Second)
	c.conn.Close()
	return nil
}

// Humming() means handle
func (h *HummingWay) Handle(ctx context.Context, conn net.Conn) {
	h.humming(ctx, conn)
}

func (h *HummingWay) humming(ctx context.Context, conn net.Conn) {
	// sleeping(closed) HummingWay will not hum
	if h.isClosing.Get() {
		conn.Close()
		return
	}

	// add the client to the set
	client := &SurfingWay{conn: conn}
	h.activeConns.Store(client, true)

	// handle the client
	reader := bufio.NewReader(conn)
	for {
		// read the message
		message, err := reader.ReadString('\n')
		if err != nil {
			// remove the client from the set if read EOF
			if err == io.EOF {
				h.activeConns.Delete(client)
				logs.Info("connection closed")
			} else {
				logs.Warn("read error:", err)
			}
			return
		}
		// Set SurfingWay to waiting state, preventing itself from being closed
		client.Waiting.AddOne()

		// Simulate sending not completed when client is closed
		// logs.Info("sleep 10s")
		// time.Sleep(10 * time.Second)

		b := []byte(message)
		conn.Write(b)

		// release the waiting state
		client.Waiting.Done()
	}
}
