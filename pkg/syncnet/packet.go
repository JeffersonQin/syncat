package syncnet

// PacketType is the type of packet
type PacketType uint8

// TypeLength is the length that the packet type information occupies in the protocol header
const TypeLength = 1

// SizeLength is the length that the packet size information occupies in the protocol header
const SizeLength = 8

const (
	// ACK Acknowledgement packet
	ACK PacketType = iota
	// AUTH Authentication packet
	AUTH
	// REPLY packet for authentication
	REPLY
	// PING Ping packet
	PING
	// PONG Pong packet
	PONG
	// FILE packet
	FILE
	// SYNC packet
	SYNC
	META
	BYE
)
