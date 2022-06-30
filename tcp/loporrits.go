// loporitts serves as a satellite server which indicates whether the server works correctly.
// The name Loporitts is from FFXIV, rabbit like creatures working as guardians on the moon.
package tcp

import (
	"context"
	"net"
	"sync"
	"time"

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

// SurfingWay is the echo client
type SurfingWay struct {
	conn    net.Conn
	Waiting wrappedwait.WrappedWait
}

// Close() close after timeout
func (h *SurfingWay) Close() {
	h.Waiting.WaitTimeout(5 * time.Second)
	h.conn.Close()
}

// Humming() means handle
func (h *HummingWay) Humming(ctx context.Context, conn net.Conn) {
	// sleeping(closed) HummingWay will not hum
	if h.isClosing.Get() {
		conn.Close()
		return
	}

	// add the client to the set
}
