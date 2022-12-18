package syncnet

import (
	"encoding/binary"
	"github.com/JeffersonQin/syncat/pkg/config"
	"github.com/JeffersonQin/syncat/pkg/database"
	pb "github.com/JeffersonQin/syncat/pkg/proto"
	"github.com/golang/protobuf/proto"
)

// SyncatRequest is the interface for all syncat request
type SyncatRequest interface {
	// GetType Get the type of the request
	GetType() PacketType
	// GetLength Get the length of the request
	GetLength() uint64
	// Handle the request
	Handle(conn *IdleTimeoutConn) error
	// Send the request
	Send(conn *IdleTimeoutConn) error
}

// SyncatRequestHeader is the header of all syncat request
type SyncatRequestHeader struct {
	// PacketType is the type of the packet
	PacketType PacketType
	// Length is the length of the packet
	Length uint64
}

// GetType Get the type of the request from header info, default method for all request
func (r *SyncatRequestHeader) GetType() PacketType {
	return r.PacketType
}

// GetLength Get the length of the request from header info, default method for all request
func (r *SyncatRequestHeader) GetLength() uint64 {
	return r.Length
}

// Send the request based on configured header info
func (r *SyncatRequestHeader) Send(conn *IdleTimeoutConn) error {
	data := make([]byte, TypeLength+SizeLength)
	data[0] = byte(r.PacketType)
	binary.BigEndian.PutUint64(data[TypeLength:], r.Length)
	count, err := conn.Write(data)
	if err != nil {
		return err
	}
	if count != TypeLength+SizeLength {
		return ErrInvalidPacket{count}
	}
	return nil
}

// SyncatAckRequest is the request for ACK packet
type SyncatAckRequest struct {
	SyncatRequestHeader
}

// Handle ACK request does not need to be handled, the function is empty
func (r *SyncatAckRequest) Handle(_ *IdleTimeoutConn) error {
	return nil
}

// NewSyncatAckRequest Create a new SyncatAckRequest
func NewSyncatAckRequest() *SyncatAckRequest {
	return &SyncatAckRequest{
		SyncatRequestHeader{
			PacketType: ACK,
			Length:     0,
		},
	}
}

// SyncatAuthRequest is the request for AUTH packet
type SyncatAuthRequest struct {
	SyncatRequestHeader
	pb.SyncatAuthRequestBody
}

// Handle AUTH request
// Authenticate the token and check whether the client is registered
// If the client is not registered, auth will also fail
// If the client field is empty, the client will be registered
// REPLY packet will be sent back as response
// AUTH request will only be sent by the client to the server when the connection is established
func (r *SyncatAuthRequest) Handle(conn *IdleTimeoutConn) error {
	data := make([]byte, r.Length)
	_, err := conn.Read(data)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data, &r.SyncatAuthRequestBody)
	if err != nil {
		return err
	}
	if r.SyncatAuthRequestBody.Token != config.GetConfig().Auth.Token {
		err = NewSyncatReplyRequest(false, r.SyncatAuthRequestBody.ClientUuid,
			"Invalid token").Send(conn)
		return err
	}
	uuid := r.SyncatAuthRequestBody.ClientUuid
	// Allocate uuid for a new client if is a new client
	if r.SyncatAuthRequestBody.ClientUuid == "" {
		uuid, err = database.AllocateNewClientUuid()
		if err != nil {
			_ = NewSyncatReplyRequest(false, r.SyncatAuthRequestBody.ClientUuid,
				"Failed to allocate new uuid").Send(conn)
			return err
		}
	}
	// check whether the uuid exists in database
	exists, err := database.QueryExistsClientUuid(uuid)
	if err != nil {
		_ = NewSyncatReplyRequest(false, r.SyncatAuthRequestBody.ClientUuid,
			"Failed to query uuid").Send
		return err
	}
	if !exists {
		err = NewSyncatReplyRequest(false, r.SyncatAuthRequestBody.ClientUuid,
			"Invalid uuid").Send(conn)
		return err
	}
	// success
	err = NewSyncatReplyRequest(true, uuid, "OK").Send(conn)
	return err
}

// Send the AUTH request
func (r *SyncatAuthRequest) Send(conn *IdleTimeoutConn) error {
	data, err := proto.Marshal(&r.SyncatAuthRequestBody)
	if err != nil {
		return err
	}
	r.Length = uint64(len(data))
	err = r.SyncatRequestHeader.Send(conn)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

// NewSyncatAuthRequest Create a new SyncatAuthRequest
func NewSyncatAuthRequest() (*SyncatAuthRequest, error) {
	clientUuid, err := database.QueryClientUuid()
	if err != nil {
		return nil, err
	}
	return &SyncatAuthRequest{
		SyncatRequestHeader{
			PacketType: AUTH,
			Length:     0,
		},
		pb.SyncatAuthRequestBody{
			ClientUuid: clientUuid,
			Token:      config.GetConfig().Auth.Token,
		},
	}, nil
}

// SyncatReplyRequest is the request for REPLY packet
type SyncatReplyRequest struct {
	SyncatRequestHeader
	pb.SyncatReplyRequestBody
}

// Handle REPLY request
// Check whether the auth is successful, and also update the client's uuid when newly registered
// REPLY request will only be sent by the server to the client when the connection is established
func (r *SyncatReplyRequest) Handle(conn *IdleTimeoutConn) error {
	data := make([]byte, r.Length)
	_, err := conn.Read(data)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data, &r.SyncatReplyRequestBody)
	if err != nil {
		return err
	}
	if r.SyncatReplyRequestBody.Success {
		err = database.UpdateClientUuid(r.SyncatReplyRequestBody.ClientUuid)
		return err
	} else {
		return ErrAuthFailed{r.SyncatReplyRequestBody.Message}
	}
}

// Send the REPLY request
func (r *SyncatReplyRequest) Send(conn *IdleTimeoutConn) error {
	data, err := proto.Marshal(&r.SyncatReplyRequestBody)
	if err != nil {
		return err
	}
	r.Length = uint64(len(data))
	err = r.SyncatRequestHeader.Send(conn)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

// NewSyncatReplyRequest Create a new SyncatReplyRequest
func NewSyncatReplyRequest(success bool, clientUUid string, message string) *SyncatReplyRequest {
	return &SyncatReplyRequest{
		SyncatRequestHeader{
			PacketType: REPLY,
			Length:     0,
		},
		pb.SyncatReplyRequestBody{
			Success:    success,
			ClientUuid: clientUUid,
			Message:    message,
		},
	}
}
