# Build stage
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git and other build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o oolio ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1001 -S oolio && \
    adduser -u 1001 -S oolio -G oolio

# Set the working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/oolio .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Change ownership to oolio user
RUN chown -R oolio:oolio /app

# Switch to non-root user
USER oolio

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./oolio"]