# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY databus/go.mod databus/go.sum ./
RUN go mod download

# Copy source code
COPY databus/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/databus-server ./cmd/entrypoint

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/databus-server .

# Copy the documents directory from build context
COPY documents/ /documents/

# Expose the API port
EXPOSE 8080

# Run the application
CMD ["./databus-server"]

