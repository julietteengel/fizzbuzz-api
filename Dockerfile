# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files first (better Docker caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o fizzbuzz-api ./cmd/server/main.go

# Production stage  
FROM alpine:latest

# Install ca-certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/fizzbuzz-api .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./fizzbuzz-api"]