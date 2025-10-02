package managers

import (
	"encoding/json"
	"fmt"
	"gateway_service/models/crypto"
	"io"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	crypto.UnimplementedMessageStreamerServer
	hub *Hub
}

func NewServer(hub *Hub) *server {
	return &server{hub: hub}
}

func (s *server) StreamMessages(stream crypto.MessageStreamer_StreamMessagesServer) error {
	log.Println("Client connected, starting to receive crypto updates...")

	for {
		updateList, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client finished sending updates")
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			log.Printf("Error receiving updates: %v", err)
			return err
		}

		fmt.Printf("Received %d crypto updates:\n", len(updateList.Updates))

		if s.hub != nil {
			jsonData, err := json.Marshal(updateList.Updates)
			if err == nil {
				s.hub.SetLatestData(jsonData)
				s.hub.GetBroadcastChannel() <- jsonData
			}
		}
	}
}
