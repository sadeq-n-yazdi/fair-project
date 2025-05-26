# Build stage
FROM golang:1.24-alpine AS builder

# Install git
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Get the current git branch hash
RUN BRANCH_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    # Build the application with the branch hash
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/sadeq/fair-project-go/pkg/version.BranchHash=$BRANCH_HASH'" -o /fair-project-server ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /fair-project-server .

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create data directory and set permissions
RUN mkdir -p /app/data && chown -R appuser:appgroup /app/data

# Set environment variables
ENV PORT=8080
ENV DATA_DIR=/app/data

# Switch to non-root user
USER appuser

# Expose the port
EXPOSE 8080

# Create a volume for the data directory
VOLUME ["/app/data"]

# Run the application
ENTRYPOINT ["./fair-project-server"]
