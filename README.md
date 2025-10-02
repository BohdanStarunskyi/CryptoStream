# CryptoStream

Real-time cryptocurrency price tracker with microservices architecture. Features live WebSocket updates, gRPC inter-service communication, and a beautiful responsive web interface.

## Features

- Real-time crypto price updates via WebSocket
- Microservices architecture with gRPC communication
- Responsive web interface with live connection status
- Auto-reconnection with exponential backoff
- Docker containerization with service orchestration
- Clean, animated UI with price change indicators
- Search and sort functionality

## Tech

- **Backend**: Go + Gorilla WebSocket + gRPC
- **Frontend**: Vanilla JavaScript + CSS animations
- **Protocol**: Protocol Buffers for service communication
- **Containerization**: Docker + Docker Compose
- **Data Source**: CoinGecko API

## Architecture

```
┌─────────────┐    gRPC     ┌─────────────┐   WebSocket   ┌─────────────┐
│   Fetcher   │ ──────────► │   Gateway   │ ────────────► │     Web     │
│  Service    │             │   Service   │               │   Client    │
└─────────────┘             └─────────────┘               └─────────────┘
      │                           │                             │
      │                           │                             │
   Fetches                    Broadcasts                   Displays
 crypto data                 to WebSocket                real-time
from CoinGecko                  clients                    prices
```

## Getting started

1) Clone and navigate to project

```bash
git clone <repository>
cd test_sockets
```

2) Start all services with Docker

```bash
make run
```

Or manually:

```bash
docker compose up --build
```

3) Access the application

- **Web Interface**: http://localhost:3000
- **WebSocket Endpoint**: ws://localhost:8080/ws
- **Gateway Health**: http://localhost:8080

## Development

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Make (optional, for convenience commands)

### Local development

1) Generate Protocol Buffer files

```bash
./generate_proto.sh
```

2) Run services individually

**Terminal 1 - Gateway:**
```bash
cd backend/gateway
go run main.go
```

**Terminal 2 - Fetcher:**
```bash
cd backend/fetcher
go run main.go
```

**Terminal 3 - Web Server:**
```bash
cd web
python3 -m http.server 3000
```

### Make commands

```bash
make build    # Build all Docker images
make up       # Start all services
make down     # Stop all services
make logs     # View logs for all services
make restart  # Restart all services
make clean    # Clean up containers, images, volumes
```

## Services overview

### Fetcher Service
**Port**: Internal gRPC client  
**Role**: Fetches cryptocurrency data from CoinGecko API every 10 seconds and streams to Gateway via gRPC

### Gateway Service
**Ports**: 8080 (WebSocket), 50051 (gRPC)  
**Role**: Receives crypto updates via gRPC and broadcasts to WebSocket clients with connection management

### Web Service
**Port**: 3000  
**Role**: Serves static frontend with real-time WebSocket connection to display live cryptocurrency prices

## API Reference

### WebSocket Connection

Connect to `ws://localhost:8080/ws` to receive real-time cryptocurrency updates.

**Message Format:**
```json
[
  {
    "id": "bitcoin",
    "symbol": "btc", 
    "name": "Bitcoin",
    "image": "https://assets.coingecko.com/coins/images/1/large/bitcoin.png",
    "current_price": 43250.75,
    "price_change_24h": 2.45
  }
]
```

### gRPC Service (Internal)

**Service**: `MessageStreamer`  
**Method**: `StreamMessages`  
**Proto**: [crypto.proto](proto/crypto.proto)

## Project layout

```
├── backend/
│   ├── fetcher/           # Data fetching service
│   │   ├── dto/          # API response structures
│   │   ├── managers/     # gRPC and HTTP clients
│   │   └── models/       # Protocol Buffer models
│   └── gateway/          # WebSocket gateway service
│       ├── managers/     # WebSocket hub and gRPC server
│       └── models/       # Protocol Buffer models
├── proto/                # Protocol Buffer definitions
├── web/                  # Frontend application
│   ├── index.html       # Main HTML structure
│   ├── app.js           # WebSocket client and UI logic
│   └── styles.css       # Responsive styling with animations
├── docker-compose.yml    # Service orchestration
└── Makefile             # Development commands
```

## Features in detail

### Real-time Updates
- WebSocket connections with automatic reconnection
- Live price change indicators (green/red)
- Visual feedback on data updates
- Connection status monitoring

### Microservices Communication
- gRPC streaming for efficient data transfer
- Protocol Buffers for type-safe messaging
- Service health checks and dependencies

### Web Interface
- Responsive grid layout for crypto cards
- Search and sort functionality
- Mobile-friendly design
- Smooth animations and transitions
- Connection status indicator

### DevOps Ready
- Multi-stage Docker builds
- Service orchestration with Docker Compose
- Health checks and restart policies
- Environment-based configuration