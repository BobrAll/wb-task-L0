#!/bin/bash
echo "Running unit tests..."
go test test/unit/ -v

echo "Running integration tests..."
go test test/integration/ -v -timeout 30s
