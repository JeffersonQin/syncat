package syncnet

import "fmt"

type ErrInvalidPacket struct {
	length int
}

func (e ErrInvalidPacket) Error() string {
	return fmt.Sprintf("invalid packet length: %d", e.length)
}

type ErrInvalidPacketType struct {
	packetType byte
}

func (e ErrInvalidPacketType) Error() string {
	return fmt.Sprintf("invalid packet type: %d", e.packetType)
}

type ErrUnexpectedPacketType struct {
	packetType byte
}

func (e ErrUnexpectedPacketType) Error() string {
	return fmt.Sprintf("unexpected packet type: %d", e.packetType)
}

type ErrAuthFailed struct {
	message string
}

func (e ErrAuthFailed) Error() string {
	return fmt.Sprintf("auth failed: %s", e.message)
}
