package syncnet

import (
	"encoding/binary"
	"github.com/JeffersonQin/syncat/pkg/config"
	"github.com/JeffersonQin/syncat/pkg/database"
	pb "github.com/JeffersonQin/syncat/pkg/proto"
	"github.com/golang/protobuf/proto"
)

type SyncatRequest interface {
	GetType() PacketType
	GetLength() uint64
	Handle(conn *IdleTimeoutConn) error
	Send(conn *IdleTimeoutConn) error
}

type SyncatRequestHeader struct {
	PacketType PacketType
	Length     uint64
}

func (r *SyncatRequestHeader) GetType() PacketType {
	return r.PacketType
}

func (r *SyncatRequestHeader) GetLength() uint64 {
	return r.Length
}

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

type SyncatAckRequest struct {
	SyncatRequestHeader
}

func (r *SyncatAckRequest) Handle(_ *IdleTimeoutConn) error {
	return nil
}

func NewSyncatAckRequest() *SyncatAckRequest {
	return &SyncatAckRequest{
		SyncatRequestHeader{
			PacketType: ACK,
			Length:     0,
		},
	}
}

type SyncatAuthRequest struct {
	SyncatRequestHeader
	pb.SyncatAuthRequestBody
}

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

type SyncatReplyRequest struct {
	SyncatRequestHeader
	pb.SyncatReplyRequestBody
}

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
