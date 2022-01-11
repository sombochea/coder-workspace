// Code generated by protoc-gen-go-drpc. DO NOT EDIT.
// protoc-gen-go-drpc version: v0.0.26
// source: peerbroker.proto

package proto

import (
	context "context"
	errors "errors"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	drpc "storj.io/drpc"
	drpcerr "storj.io/drpc/drpcerr"
)

type drpcEncoding_File_peerbroker_proto struct{}

func (drpcEncoding_File_peerbroker_proto) Marshal(msg drpc.Message) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_peerbroker_proto) MarshalAppend(buf []byte, msg drpc.Message) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(buf, msg.(proto.Message))
}

func (drpcEncoding_File_peerbroker_proto) Unmarshal(buf []byte, msg drpc.Message) error {
	return proto.Unmarshal(buf, msg.(proto.Message))
}

func (drpcEncoding_File_peerbroker_proto) JSONMarshal(msg drpc.Message) ([]byte, error) {
	return protojson.Marshal(msg.(proto.Message))
}

func (drpcEncoding_File_peerbroker_proto) JSONUnmarshal(buf []byte, msg drpc.Message) error {
	return protojson.Unmarshal(buf, msg.(proto.Message))
}

type DRPCPeerBrokerClient interface {
	DRPCConn() drpc.Conn

	NegotiateConnection(ctx context.Context) (DRPCPeerBroker_NegotiateConnectionClient, error)
}

type drpcPeerBrokerClient struct {
	cc drpc.Conn
}

func NewDRPCPeerBrokerClient(cc drpc.Conn) DRPCPeerBrokerClient {
	return &drpcPeerBrokerClient{cc}
}

func (c *drpcPeerBrokerClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcPeerBrokerClient) NegotiateConnection(ctx context.Context) (DRPCPeerBroker_NegotiateConnectionClient, error) {
	stream, err := c.cc.NewStream(ctx, "/peerbroker.PeerBroker/NegotiateConnection", drpcEncoding_File_peerbroker_proto{})
	if err != nil {
		return nil, err
	}
	x := &drpcPeerBroker_NegotiateConnectionClient{stream}
	return x, nil
}

type DRPCPeerBroker_NegotiateConnectionClient interface {
	drpc.Stream
	Send(*NegotiateConnection_ClientToServer) error
	Recv() (*NegotiateConnection_ServerToClient, error)
}

type drpcPeerBroker_NegotiateConnectionClient struct {
	drpc.Stream
}

func (x *drpcPeerBroker_NegotiateConnectionClient) Send(m *NegotiateConnection_ClientToServer) error {
	return x.MsgSend(m, drpcEncoding_File_peerbroker_proto{})
}

func (x *drpcPeerBroker_NegotiateConnectionClient) Recv() (*NegotiateConnection_ServerToClient, error) {
	m := new(NegotiateConnection_ServerToClient)
	if err := x.MsgRecv(m, drpcEncoding_File_peerbroker_proto{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (x *drpcPeerBroker_NegotiateConnectionClient) RecvMsg(m *NegotiateConnection_ServerToClient) error {
	return x.MsgRecv(m, drpcEncoding_File_peerbroker_proto{})
}

type DRPCPeerBrokerServer interface {
	NegotiateConnection(DRPCPeerBroker_NegotiateConnectionStream) error
}

type DRPCPeerBrokerUnimplementedServer struct{}

func (s *DRPCPeerBrokerUnimplementedServer) NegotiateConnection(DRPCPeerBroker_NegotiateConnectionStream) error {
	return drpcerr.WithCode(errors.New("Unimplemented"), drpcerr.Unimplemented)
}

type DRPCPeerBrokerDescription struct{}

func (DRPCPeerBrokerDescription) NumMethods() int { return 1 }

func (DRPCPeerBrokerDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/peerbroker.PeerBroker/NegotiateConnection", drpcEncoding_File_peerbroker_proto{},
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return nil, srv.(DRPCPeerBrokerServer).
					NegotiateConnection(
						&drpcPeerBroker_NegotiateConnectionStream{in1.(drpc.Stream)},
					)
			}, DRPCPeerBrokerServer.NegotiateConnection, true
	default:
		return "", nil, nil, nil, false
	}
}

func DRPCRegisterPeerBroker(mux drpc.Mux, impl DRPCPeerBrokerServer) error {
	return mux.Register(impl, DRPCPeerBrokerDescription{})
}

type DRPCPeerBroker_NegotiateConnectionStream interface {
	drpc.Stream
	Send(*NegotiateConnection_ServerToClient) error
	Recv() (*NegotiateConnection_ClientToServer, error)
}

type drpcPeerBroker_NegotiateConnectionStream struct {
	drpc.Stream
}

func (x *drpcPeerBroker_NegotiateConnectionStream) Send(m *NegotiateConnection_ServerToClient) error {
	return x.MsgSend(m, drpcEncoding_File_peerbroker_proto{})
}

func (x *drpcPeerBroker_NegotiateConnectionStream) Recv() (*NegotiateConnection_ClientToServer, error) {
	m := new(NegotiateConnection_ClientToServer)
	if err := x.MsgRecv(m, drpcEncoding_File_peerbroker_proto{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (x *drpcPeerBroker_NegotiateConnectionStream) RecvMsg(m *NegotiateConnection_ClientToServer) error {
	return x.MsgRecv(m, drpcEncoding_File_peerbroker_proto{})
}