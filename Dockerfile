FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o server .

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /app

# Add CA certificates in case you need HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./server"]