package syncnet

import "log"

// Log is the logger for the connection
func (c *IdleTimeoutConn) Log(v ...any) {
	log.Println("["+c.RemoteAddr().String()+"]", v)
}
