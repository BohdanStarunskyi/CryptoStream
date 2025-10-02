#!/bin/bash

# Script to regenerate proto models for fetcher and gateway services

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔄 Regenerating proto models...${NC}"

# Get the script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Add Go bin directory to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}❌ protoc is not installed. Please install Protocol Buffers compiler.${NC}"
    echo "On macOS: brew install protobuf"
    exit 1
fi

# Check if Go protoc plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${YELLOW}⚠️  protoc-gen-go not found. Installing...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${YELLOW}⚠️  protoc-gen-go-grpc not found. Installing...${NC}"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directories if they don't exist
mkdir -p backend/fetcher/models/crypto
mkdir -p backend/gateway/models/crypto

echo -e "${YELLOW}📦 Generating models for fetcher service...${NC}"
protoc \
    --go_out=backend/fetcher/ \
    --go-grpc_out=backend/fetcher/ \
    --proto_path=proto \
    proto/crypto.proto

echo -e "${YELLOW}📦 Generating models for gateway service...${NC}"
protoc \
    --go_out=backend/gateway/ \
    --go-grpc_out=backend/gateway/ \
    --proto_path=proto \
    proto/crypto.proto

echo -e "${GREEN}✅ Proto models regenerated successfully!${NC}"
echo -e "${GREEN}📁 Generated files:${NC}"
echo "  - backend/fetcher/models/crypto/crypto.pb.go"
echo "  - backend/fetcher/models/crypto/crypto_grpc.pb.go"
echo "  - backend/gateway/models/crypto/crypto.pb.go"
echo "  - backend/gateway/models/crypto/crypto_grpc.pb.go"