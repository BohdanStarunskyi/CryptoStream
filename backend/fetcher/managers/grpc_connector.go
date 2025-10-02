package managers

import (
	"context"

	"fetceher_service/models/crypto"

	"google.golang.org/grpc"
)

type GRPCConnector struct {
	conn   *grpc.ClientConn
	client crypto.MessageStreamerClient
	stream crypto.MessageStreamer_StreamMessagesClient
}

func NewGRPCConnector(address string) (*GRPCConnector, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := crypto.NewMessageStreamerClient(conn)

	stream, err := client.StreamMessages(context.Background())
	if err != nil {
		return nil, err
	}

	return &GRPCConnector{
		conn:   conn,
		client: client,
		stream: stream,
	}, nil
}

func (g *GRPCConnector) SendUpdates(updates []*crypto.CryptoUpdate) error {
	updateList := &crypto.CryptoUpdateList{
		Updates: updates,
	}
	return g.stream.Send(updateList)
}

func (g *GRPCConnector) Close() error {
	return g.conn.Close()
}
