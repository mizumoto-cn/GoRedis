package wrappedwait

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// BooTest asserts 1==1 to be deleted
func BooTest(t *testing.T) {
	assert.Equal(t, 1, 1)
}

// func testWaitGroup tests the Wrapped WaitGroup struct's functions except timeout
func testWaitGroup(t *testing.T, wg1 *WrappedWait, wg2 *WrappedWait, delta int) {

	exited := make(chan bool, delta)

	wg1.AddOne()
	wg1.Done()
	wg1.Add(delta)
	wg2.Add(delta)
	for i := 0; i < delta; i++ {
		go func() {
			wg1.Done()
			wg2.Wait()
			exited <- true
		}()
	}
	wg1.Wait()
	for i := 0; i < delta; i++ {
		select {
		case <-exited:
			t.Fatal("WaitGroup released group too soon")
		default:
		}
		wg2.Done()
	}
	for i := 0; i != delta; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
	assert.Equal(t, 0, len(exited))
}

// More Concurrent TestWrappedWaitGroup
func TestMultiWait(t *testing.T) {
	var (
		wg1 = &WrappedWait{}
		wg2 = &WrappedWait{}
	)

	for i := 0; i < 8; i++ {
		testWaitGroup(t, wg1, wg2, 16+i)
	}
}

// TestTimeout tests the Wrapped WaitGroup's timeout function
func TestTimeout(t *testing.T) {
	var (
		wg1 = &WrappedWait{}
		wg2 = &WrappedWait{}
	)

	wg1.AddOne()
	wg1.Done()
	wg1.Add(1)
	wg2.Add(1)
	go func() {
		wg1.Wait()
		wg2.Wait()
	}()
	assert.True(t, wg1.WaitTimeout(time.Millisecond))
	wg2.Done()
	wg2.Wait()
}
