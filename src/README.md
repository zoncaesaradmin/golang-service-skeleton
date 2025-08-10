# Golang Component Skeleton

This project provides a comprehensive Go component skeleton with a REST API component and a separate test runner using Go workspace for multi-module management.

## Project Structure

```
katharos/                     # Root workspace
├── go.work                   # Go workspace file
├── Makefile                  # Root orchestration
├── README.md                 # Documentation
│
├── component/               # Main REST API component
│   ├── go.mod               # Module: katharos/component
│   ├── go.sum               # Go module checksums
│   ├── Makefile             # Component build automation
│   ├── cmd/                 # Application entry point
│   │   └── main.go          # Component main function
│   ├── internal/            # Private application code
│   │   ├── api/             # HTTP handlers and routes
│   │   │   └── handlers.go  # REST API handlers
│   │   ├── config/          # Configuration management
│   │   │   └── config.go    # Config structures
│   │   ├── models/          # Data models
│   │   │   └── models.go    # User and API models
│   │   └── service/         # Business logic
│   │       ├── user_service.go      # User service implementation
│   │       └── user_service_test.go # Unit tests
│   ├── pkg/                 # Public library code (if any)
│   └── bin/                 # Build output
│
└── testrunner/              # Test runner service
    ├── go.mod               # Module: katharos/testrunner
    ├── go.sum               # Go module checksums
    ├── Makefile             # Test runner build automation
    ├── cmd/                 # Test runner entry point
    │   └── main.go          # Test runner main function
    ├── internal/            # Private test code
    │   ├── client/          # API client for testing
    │   │   └── client.go    # HTTP client implementation
    │   └── tests/           # Integration tests
    │       └── integration_tests.go # Test implementations
    ├── pkg/                 # Public test utilities (if any)
    └── bin/                 # Build output
```

## Features

### Main Service
- **REST API** with Go standard library (net/http)
- **User Management** (CRUD operations)
- **Health Check** endpoint
- **Configuration Management** via environment variables
- **Graceful Shutdown**
- **CORS Support**
- **Thread-safe in-memory user storage**
- **Unit Tests** with testify
- **JSON API responses** with proper status codes

### Test Runner
- **Integration Tests** for all API endpoints
- **Performance Tests** with configurable concurrency
- **Load Testing** capabilities
- **Comprehensive Test Reports**
- **Configurable Test Modes**
- **Independent test binary** for service boundary testing

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Make

### Running the Service

1. Build and run the service from root:
```bash
make run-service
```

Or navigate to service directory:
```bash
cd service
make run
```

Or use individual commands:
```bash
cd service
make build
./bin/service
```

The service will start on `localhost:8080` by default.

### Running Tests

1. Start the service (in another terminal):
```bash
make run-service
# OR
cd service && make run
```

2. Run the test runner:
```bash
make run-tests
# OR
cd testrunner && make run
```

### Using Go Workspace

This project uses Go 1.18+ workspace feature:
```bash
go work sync    # Sync workspace modules
go work use ./service ./testrunner  # Add modules to workspace
```

## API Endpoints

### Health Check
- `GET /health` - Service health status

### User Management
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `GET /api/v1/users/search?q=query` - Search users

### Statistics
- `GET /api/v1/stats` - Get service statistics

## Configuration

The service can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_HOST | localhost | Server host |
| SERVER_PORT | 8080 | Server port |
| SERVER_READ_TIMEOUT | 10 | Read timeout in seconds |
| SERVER_WRITE_TIMEOUT | 10 | Write timeout in seconds |
| LOG_LEVEL | info | Log level |
| LOG_FORMAT | json | Log format |

## Testing

### Unit Tests
```bash
# From root
make test

# Individual service
cd service
make test
```

### Integration Tests
```bash
# From root
make run-tests

# From testrunner directory
cd testrunner
make run
```

### Performance Tests
```bash
cd testrunner
make run-performance-tests
```

### Test Coverage
```bash
cd service
make test-coverage
```

### End-to-End Testing
```bash
# From root - starts service, runs tests, stops service
make e2e-test
```

## Cleanup Commands

### Quick Clean (removes build artifacts)
```bash
# Clean both projects from root
make clean

# Clean individual projects
cd service && make clean
cd testrunner && make clean
```

### Deep Clean (removes build artifacts, vendor, and module cache)
```bash
# Deep clean both projects from root
make deep-clean

# Deep clean individual projects
cd service && make deep-clean
cd testrunner && make deep-clean
```

The clean commands remove:
- `bin/` directories and binaries
- `coverage.out` and `coverage.html` files
- `*.log` files
- `tmp/` and `temp/` directories
- `*.test` and `*.out` files

The deep-clean commands additionally remove:
- `vendor/` directories
- Go module cache

## Build Commands

### Root Level Commands
```bash
make help           # Show all available targets
make build          # Build both service and testrunner
make test           # Run all tests
make clean          # Clean build artifacts from both projects
make deep-clean     # Deep clean including vendor and module cache
make run-service    # Build and run the main service
make run-tests      # Build and run integration tests
make dev-setup      # Complete development setup
make ci             # Full CI pipeline
make e2e-test       # End-to-end testing
```

### Service
```bash
cd service
make help           # Show available targets
make build          # Build the binary
make test           # Run unit tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts and generated files
make deep-clean     # Deep clean including vendor and module cache
make fmt            # Format code
make vet            # Vet code
make lint           # Lint code (requires golangci-lint)
```

### Test Runner
```bash
cd testrunner
make help                    # Show available targets
make build                   # Build the binary
make test                    # Run unit tests
make clean                   # Clean build artifacts and generated files
make deep-clean              # Deep clean including vendor and module cache
make run-integration-tests   # Run integration tests
make run-performance-tests   # Run performance tests
```

## Development

### Prerequisites
- Go 1.21 or higher
- Make
- golangci-lint (optional, for linting)

### Project Architecture
- **Multi-module workspace**: Uses Go 1.18+ workspace feature
- **Standard HTTP**: Uses Go standard library instead of frameworks
- **Thread-safe storage**: In-memory user storage with mutex protection
- **Clean separation**: Service and test runner are completely independent

### Adding New Features
1. Add models in `service/internal/models/`
2. Implement business logic in `service/internal/service/`
3. Add HTTP handlers in `service/internal/api/`
4. Add routes in the `SetupRoutes` function
5. Write unit tests
6. Add integration tests in `testrunner/internal/tests/`

### Module Management
```bash
# Working with workspace
go work sync                 # Sync all modules in workspace
go work use ./service        # Add service module
go work use ./testrunner     # Add testrunner module

# Working with individual modules
cd service && go mod tidy    # Manage service dependencies
cd testrunner && go mod tidy # Manage testrunner dependencies
```

### Testing Strategy
1. **Unit Tests**: Test individual functions and methods in `service/internal/service/`
2. **Integration Tests**: Test API endpoints end-to-end via HTTP client
3. **Performance Tests**: Test system under load with configurable concurrency
4. **Load Tests**: Test system scalability and limits
5. **Service Boundary Tests**: External testing via independent test runner

## Examples

### Creating a User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Getting All Users
```bash
curl http://localhost:8080/api/v1/users
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Go Workspace Benefits

This project uses Go 1.18+ workspace feature which provides:

- **Multi-module management**: Work with multiple related modules in a single workspace
- **Cross-module development**: Make changes across service and testrunner simultaneously
- **Shared dependencies**: Efficient dependency management across modules
- **IDE support**: Better code navigation and refactoring across modules

### Workspace Commands
```bash
go work init                 # Initialize workspace
go work use ./service        # Add service module to workspace
go work use ./testrunner     # Add testrunner module to workspace
go work sync                 # Sync workspace modules
go work edit                 # Edit go.work file
```

## Docker Support (Optional)

You can add Docker support by creating a Dockerfile in each service directory:

```dockerfile
# service/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/service cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/service .
EXPOSE 8080
CMD ["./service"]
```

### Docker Compose Example
```yaml
# docker-compose.yml
version: '3.8'
services:
  service:
    build: ./service
    ports:
      - "8080:8080"
    environment:
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
    
  testrunner:
    build: ./testrunner
    depends_on:
      - service
    environment:
      - SERVICE_URL=http://service:8080
```

## License

MIT License - see LICENSE file for details.
