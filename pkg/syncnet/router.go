package syncnet

import (
	"encoding/binary"
	pb "github.com/JeffersonQin/syncat/pkg/proto"
	"golang.org/x/exp/slices"
)

func Wait(conn *IdleTimeoutConn, typeList []PacketType) (SyncatRequest, error) {
	request, err := RouteConn(conn)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(typeList, request.GetType()) {
		return nil, ErrUnexpectedPacketType{byte(request.GetType())}
	}
	return request, nil
}

func RouteConn(conn *IdleTimeoutConn) (SyncatRequest, error) {
	headData := make([]byte, TypeLength+SizeLength)
	count, err := conn.Read(headData)
	if err != nil {
		return nil, err
	}
	if count != TypeLength+SizeLength {
		return nil, ErrInvalidPacket{count}
	}
	length := binary.BigEndian.Uint64(headData[TypeLength:])
	switch headData[0] {
	case byte(ACK):
		return &SyncatAckRequest{
			SyncatRequestHeader{
				PacketType: ACK,
				Length:     length,
			},
		}, nil
	case byte(AUTH):
		return &SyncatAuthRequest{
			SyncatRequestHeader{
				PacketType: AUTH,
				Length:     length,
			},
			pb.SyncatAuthRequestBody{},
		}, nil
	case byte(REPLY):
		return &SyncatReplyRequest{
			SyncatRequestHeader{
				PacketType: REPLY,
				Length:     length,
			},
			pb.SyncatReplyRequestBody{},
		}, nil
	}
	return nil, ErrInvalidPacketType{headData[0]}
}
