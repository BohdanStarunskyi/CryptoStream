# CryptoStream Docker Commands

.PHONY: build up down logs clean restart

# Build all services
build:
	docker compose build

# Start all services
up:
	docker compose up -d

# Stop all services
down:
	docker compose down

# View logs for all services
logs:
	docker compose logs -f

# View logs for specific service
logs-gateway:
	docker compose logs -f gateway

logs-fetcher:
	docker compose logs -f fetcher

logs-web:
	docker compose logs -f web

# Clean up everything (containers, images, volumes)
clean:
	docker compose down -v --rmi all
	docker system prune -f

# Restart all services
restart: down up

# Build and start
run: build up

# Development mode - build and start with logs
dev: build
	docker compose up

# Check status
status:
	docker compose ps