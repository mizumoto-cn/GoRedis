package tcp

import (
	"bufio"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test ListenAndServe tests the ListenAndServe function
func TestListenAndServe(t *testing.T) {
	var (
		err       error
		l         net.Listener
		closeChan = make(chan struct{})
	)
	l, err = net.Listen("tcp", ":0")
	assert.NoError(t, err)

	addr := l.Addr().String()
	go func() {
		ListenAndServe(l, NewHummingWay(), closeChan)
	}()

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		randomString := strconv.Itoa(rand.Int())
		_, err = conn.Write([]byte(randomString + "\n"))
		assert.NoError(t, err)

		bufReader := bufio.NewReader(conn)
		line, _, err := bufReader.ReadLine()
		assert.NoError(t, err)
		assert.Equal(t, randomString, string(line))
	}

	conn.Close()
	for i := 0; i < 5; i++ {
		// idle connection shall be closed after 5 seconds
		net.Dial("tcp", addr)
	}
	closeChan <- struct{}{}
	time.Sleep(time.Second)
}
