# FizzBuzz REST API

A production-ready FizzBuzz REST API server built with Go and Echo framework following clean architecture principles.

## Overview

This API provides a FizzBuzz endpoint that accepts parameters and returns a customized FizzBuzz sequence, along with statistics tracking for the most frequently requested parameters.

## Features

- **FizzBuzz Endpoint**: Generate FizzBuzz sequences with custom parameters
- **Statistics Tracking**: Track and retrieve most frequently used parameters
- **Health Check**: Monitor API health and status
- **Clean Architecture**: Separated layers for maintainability and testability
- **Comprehensive Testing**: Unit and integration tests with mocking
- **Production Ready**: Includes logging, error handling, graceful shutdown

## API Endpoints

### POST /fizzbuzz
Generate a FizzBuzz sequence with custom parameters.

**Parameters:**
- `int1` (integer): First divisor
- `int2` (integer): Second divisor  
- `limit` (integer): Upper limit for the sequence
- `str1` (string): Replacement string for multiples of int1
- `str2` (string): Replacement string for multiples of int2

### GET /stats
Get statistics about the most frequently requested parameters.

### GET /health
Health check endpoint.

## Tech Stack

- **Framework**: Echo v4
- **Architecture**: Clean Architecture (Controller → Service → Repository)
- **Database**: PostgreSQL with GORM ORM
- **Dependency Injection**: Uber FX
- **Configuration**: Viper for environment-based config
- **Testing**: Testify for assertions and mocking
- **Mock Generation**: Mockery v3 for interface mocks
- **Containerization**: Docker with multi-stage builds

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── controller/          # HTTP handlers and tests
│   ├── service/             # Business logic and tests
│   ├── repository/          # Data access layer and tests
│   ├── model/               # Domain models
│   ├── mocks/               # Generated interface mocks
│   └── middleware/          # HTTP middleware
├── pkg/
│   └── validator/
├── .mockery.yaml            # Mock generation configuration
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
└── README.md
```

## Development

### Prerequisites
- Go 1.21+
- Docker (optional)

### Running the Application

```bash
# Install dependencies
go mod download

# Run the server
make run

# Or run directly
go run cmd/server/main.go
```

### Testing

```bash
# Run all tests
make test

# Generate mocks from interfaces
make generate-mocks

# Run tests with coverage
make test-coverage

# Run specific test suite
go test ./internal/service/...
```

### Building

```bash
# Build binary
make build

# Build Docker image
make docker-build
```

## Configuration

The application can be configured via environment variables:

- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment (development/production, default: development)
- `MAX_LIMIT`: Maximum limit for FizzBuzz sequences (default: 10000)
- `DATABASE_URL`: PostgreSQL connection string (default: postgres://user:password@localhost:5432/fizzbuzz_db?sslmode=disable)
- `STATS_STORAGE`: Storage backend for statistics (memory/postgres, default: memory)

## Database Setup

### PostgreSQL

The application supports PostgreSQL for persistent statistics storage:

1. **Using existing PostgreSQL instance:**
   ```bash
   # Create database
   createdb fizzbuzz_db
   
   # Set environment variables
   export DATABASE_URL="postgres://user:password@localhost:5432/fizzbuzz_db?sslmode=disable"
   export STATS_STORAGE="postgres"
   ```

2. **Using Docker Compose:**
   ```bash
   # Start PostgreSQL
   docker compose up -d postgres
   
   # Database will be automatically created with credentials:
   # - Host: localhost:5432
   # - User: postgres
   # - Password: postgres
   # - Database: fizzbuzz_db
   ```

3. **TablePlus Connection:**
   - Host: `localhost`
   - Port: `5432`
   - User: `user`
   - Password: `password`
   - Database: `fizzbuzz_db`

## Contributing

1. Follow Go best practices and idioms
2. Write tests for new functionality
3. Ensure all tests pass before committing
4. Follow the existing code style and architecture

## License

MIT License