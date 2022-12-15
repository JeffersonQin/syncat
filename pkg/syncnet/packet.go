package syncnet

type PacketType uint8

const TypeLength = 1

const SizeLength = 8

const (
	ACK PacketType = iota
	AUTH
	REPLY
	PING
	PONG
	FILE
	SYNC
	META
	BYE
)
