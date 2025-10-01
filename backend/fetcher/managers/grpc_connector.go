package managers

import (
	"context"

	"fetceher_service/models/crypro"

	"google.golang.org/grpc"
)

type GRPCConnector struct {
	conn   *grpc.ClientConn
	client crypro.MessageStreamerClient
	stream crypro.MessageStreamer_StreamMessagesClient
}

func NewGRPCConnector(address string) (*GRPCConnector, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := crypro.NewMessageStreamerClient(conn)

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

func (g *GRPCConnector) SendUpdates(updates []*crypro.CryptoUpdate) error {
	updateList := &crypro.CryptoUpdateList{
		Updates: updates,
	}
	return g.stream.Send(updateList)
}

func (g *GRPCConnector) Close() error {
	return g.conn.Close()
}
