package server

import (
	"github.com/JeffersonQin/syncat/pkg/config"
	"github.com/JeffersonQin/syncat/pkg/syncnet"
	"log"
	"net"
	"strconv"
	"time"
)

func handleConnection(conn *syncnet.IdleTimeoutConn) {
	defer func(conn *syncnet.IdleTimeoutConn) {
		_ = conn.Close()
		conn.Log("Connection closed")
	}(conn)
	// wait for AUTH
	req, err := syncnet.Wait(conn, []syncnet.PacketType{syncnet.AUTH})
	if err != nil {
		conn.Log("Failed to wait for auth packet", err)
		return
	}
	err = req.Handle(conn)
	if err != nil {
		conn.Log("Failed to handle auth packet", err)
		return
	}
	// server should wait for 3 different packets: PING, SYNC, and BYE
	// PING for maintaining the connection in case of timeout
	// SYNC for starting a sync session
	// BYE for closing the connection
	// the detailed implementations are handled in requests.go
	for {
		req, err = syncnet.Wait(conn, []syncnet.PacketType{syncnet.PING, syncnet.SYNC, syncnet.BYE})
		if err != nil {
			conn.Log("Failed to wait for packet", err)
			return
		}
		// if the packet is BYE, then the server should close the connection
		if req.GetType() == syncnet.BYE {
			break
		}
		// otherwise handle the packet
		err = req.Handle(conn)
		if err != nil {
			conn.Log("Failed to handle packet", err)
			return
		}
	}
}

func StartSyncatServer() error {
	serverConfig := GetConfig()
	addr := serverConfig.Host + ":" + strconv.Itoa(serverConfig.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	log.Println("Syncat server started at " + addr + "...")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		idleTimeoutConn := &syncnet.IdleTimeoutConn{
			TCPConn:     conn,
			IdleTimeout: time.Duration(config.GetConfig().Protocol.Timeout) * time.Second,
		}
		idleTimeoutConn.Log("Connection established")
		go handleConnection(idleTimeoutConn)
	}
}
