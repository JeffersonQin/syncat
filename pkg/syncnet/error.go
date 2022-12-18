package syncnet

import "fmt"

// ErrInvalidPacket is returned when the packet length is invalid
type ErrInvalidPacket struct {
	length int
}

// Error returns the error message
func (e ErrInvalidPacket) Error() string {
	return fmt.Sprintf("invalid packet length: %d", e.length)
}

// ErrInvalidPacketType is returned when the packet type is invalid
type ErrInvalidPacketType struct {
	packetType byte
}

// Error returns the error message
func (e ErrInvalidPacketType) Error() string {
	return fmt.Sprintf("invalid packet type: %d", e.packetType)
}

// ErrUnexpectedPacketType is returned when the packet type is unexpected
type ErrUnexpectedPacketType struct {
	packetType byte
}

// Error returns the error message
func (e ErrUnexpectedPacketType) Error() string {
	return fmt.Sprintf("unexpected packet type: %d", e.packetType)
}

// ErrAuthFailed is returned when the authentication fails
type ErrAuthFailed struct {
	message string
}

// Error returns the error message
func (e ErrAuthFailed) Error() string {
	return fmt.Sprintf("auth failed: %s", e.message)
}
