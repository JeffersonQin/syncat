package syncnet

import "log"

func (c *IdleTimeoutConn) Log(v ...any) {
	log.Println("["+c.RemoteAddr().String()+"]", v)
}
