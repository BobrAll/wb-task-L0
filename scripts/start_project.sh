#!/bin/bash
echo "Starting Docker services..."
docker-compose -f deployments/docker-compose.yaml up -d

echo "Waiting for dependencies..."
sleep 5

echo "Starting orders services (ports 8081 and 8082)..."
PORT=8081 go run cmd/server/main.go &
PORT=8082 go run cmd/server/main.go