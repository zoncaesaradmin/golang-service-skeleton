#!/bin/bash

# Quick test script for the refactored component
echo "Testing the refactored REST API component..."

# Build and start the component in background
cd src/component
echo "Building component..."
go build -o test-component ./cmd

echo "Starting component..."
./test-component &
COMPONENT_PID=$!

# Wait a moment for component to start
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
kill $COMPONENT_PID 2>/dev/null
rm -f test-component
echo "Test completed!"
