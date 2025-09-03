# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files first (better Docker caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install swag and generate documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN make generate-docs

# Build the application
RUN go build -o fizzbuzz-api ./cmd/server/main.go

# Production stage  
FROM alpine:latest

# Install ca-certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/fizzbuzz-api .

# Copy the generated docs
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./fizzbuzz-api"]