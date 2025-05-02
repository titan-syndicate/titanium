# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/ti ./cmd/titanium/main.go

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/ti /app/ti

# Copy config template
COPY config.yaml /etc/titanium/config.yaml

# Set ownership
RUN chown -R appuser:appuser /app /etc/titanium

# Switch to non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["/app/ti"]

# Default command
CMD ["--help"]