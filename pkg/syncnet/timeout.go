package syncnet

import (
	"net"
	"time"
)

// IdleTimeoutConn is the connection with idle timeout
type IdleTimeoutConn struct {
	// TCPConn is the underlying TCP connection
	*net.TCPConn
	// IdleTimeout is the timeout for idle connection
	IdleTimeout time.Duration
}

// Read reads data from the connection
// The timeout is set for each read operation
func (c *IdleTimeoutConn) Read(b []byte) (int, error) {
	err := c.TCPConn.SetReadDeadline(time.Now().Add(c.IdleTimeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Read(b)
}

// Write writes data to the connection
// The timeout is set for each write operation
func (c *IdleTimeoutConn) Write(b []byte) (int, error) {
	err := c.TCPConn.SetWriteDeadline(time.Now().Add(c.IdleTimeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Write(b)
}
