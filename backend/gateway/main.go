package main

import (
	"log"
	"net"
	"net/http"

	"gateway_service/managers"
	pb "gateway_service/models/crypto"

	"google.golang.org/grpc"
)

func main() {
	hub := managers.NewHub()
	managers.SetGlobalHub(hub)
	go hub.Run()

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterMessageStreamerServer(s, managers.NewServer(hub))

		log.Println("Gateway gRPC server listening on port 50051...")

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	http.HandleFunc("/ws", managers.HandleWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("WebSocket Gateway Server"))
	})

	log.Println("WebSocket server listening on port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
