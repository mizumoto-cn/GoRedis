package wrappedwait

import (
	"sync"
	"time"
)

type WrappedWait struct {
	waitGroup sync.WaitGroup
}

// AddOne adds 1 to the wait group.
func (w *WrappedWait) AddOne() {
	w.waitGroup.Add(1)
}

// Add adds the given amount (delta) to the wait group.
func (w *WrappedWait) Add(delta int) {
	w.waitGroup.Add(delta)
}

// Done decrements the wait group.
func (w *WrappedWait) Done() {
	w.waitGroup.Done()
}

// Wait blocks until the wait group is zero.
func (w *WrappedWait) Wait() {
	w.waitGroup.Wait()
}

// WaitTimeout blocks until the wait group is zero or the timeout is reached.
// Returns true if the timeout is reached.
func (w *WrappedWait) WaitTimeout(timeout time.Duration) bool {
	// buffer the channel to avoid leaking goroutines
	done := make(chan struct{}, 1)
	go func() {
		w.waitGroup.Wait()
		done <- struct{}{}
		close(done)
	}()
	select {
	case <-done:
		return false
	case <-time.After(timeout):
		return true
	}
}
