#!/bin/bash
echo "Starting Docker services..."
docker-compose -f deployments/docker-compose.yaml up -d

echo "Waiting for dependencies..."
sleep 5

echo "Starting main application..."
go run cmd/server/main.go