package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "gateway_service/models/crypro"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// Client represents a WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub maintains active WebSocket connections
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

var hub *Hub

type server struct {
	pb.UnimplementedMessageStreamerServer
}

func (s *server) StreamMessages(stream pb.MessageStreamer_StreamMessagesServer) error {
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
		for _, update := range updateList.Updates {
			fmt.Printf("  - %s (%s): $%.4f\n",
				update.Name,
				update.Symbol,
				update.CurrentPrice)
		}
		fmt.Println("---")

		// Broadcast updates to WebSocket clients
		if hub != nil {
			jsonData, err := json.Marshal(updateList.Updates)
			if err == nil {
				hub.broadcast <- jsonData
			}
		}
	}
}

// WebSocket client handler
func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// WebSocket HTTP handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func main() {
	// Initialize the WebSocket hub
	hub = newHub()
	go hub.run()

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterMessageStreamerServer(s, &server{})

		log.Println("Gateway gRPC server listening on port 50051...")
		log.Println("Waiting for crypto updates from fetcher...")

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Set up HTTP server for WebSocket connections
	http.HandleFunc("/ws", handleWebSocket)
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
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Crypto WebSocket Stream</title>
</head>
<body>
    <h1>Crypto Price Updates</h1>
    <div id="updates"></div>
    <script>
        const ws = new WebSocket('ws://localhost:8080/ws');
        const updatesDiv = document.getElementById('updates');
        
        ws.onmessage = function(event) {
            const updates = JSON.parse(event.data);
            const updateHtml = updates.map(crypto => 
                ` + "`<p>${crypto.name} (${crypto.symbol}): $${crypto.current_price.toFixed(4)}</p>`" + `
            ).join('');
            updatesDiv.innerHTML = '<h3>Latest Updates:</h3>' + updateHtml;
        };
        
        ws.onopen = function(event) {
            console.log('WebSocket connected');
        };
        
        ws.onclose = function(event) {
            console.log('WebSocket disconnected');
        };
    </script>
</body>
</html>`))
	})

	log.Println("WebSocket server listening on port 8080...")
	log.Println("Visit http://localhost:8080 to see live crypto updates")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
