#!/bin/bash

# Quick test script for the refactored service
echo "Testing the refactored REST API service..."

# Build and start the service in background
cd src/service
echo "Building service..."
go build -o test-service ./cmd

echo "Starting service..."
./test-service &
SERVICE_PID=$!

# Wait a moment for service to start
sleep 2

echo "Testing health endpoint..."
curl -s http://localhost:8080/health | jq '.' || echo "Health check failed"

echo ""
echo "Testing user creation..."
curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","first_name":"Test","last_name":"User"}' | jq '.' || echo "User creation failed"

echo ""
echo "Testing user retrieval..."
curl -s http://localhost:8080/api/v1/users/1 | jq '.' || echo "User retrieval failed"

echo ""
echo "Testing stats endpoint..."
curl -s http://localhost:8080/api/v1/stats | jq '.' || echo "Stats failed"

# Clean up
echo ""
echo "Cleaning up..."
kill $SERVICE_PID 2>/dev/null
rm -f test-service
echo "Test completed!"
