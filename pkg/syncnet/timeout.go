package syncnet

import (
	"net"
	"time"
)

type IdleTimeoutConn struct {
	net.TCPConn
	IdleTimeout time.Duration
}

func (c *IdleTimeoutConn) Read(b []byte) (int, error) {
	err := c.TCPConn.SetReadDeadline(time.Now().Add(c.IdleTimeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Read(b)
}

func (c *IdleTimeoutConn) Write(b []byte) (int, error) {
	err := c.TCPConn.SetWriteDeadline(time.Now().Add(c.IdleTimeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Write(b)
}
