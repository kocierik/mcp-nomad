# Builder stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /workspace

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (CGO_ENABLED=0 so we don't need gcc)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o mcp-nomad .

# Runtime stage - using Alpine for small size but with shell access
FROM alpine:3.19

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary from the builder stage
COPY --from=builder /workspace/mcp-nomad /mcp-nomad

# Set executable permissions
RUN chmod +x /mcp-nomad

# Switch to non-root user
USER appuser

# Set the entry point (Smithery will override this with environment variables)
ENTRYPOINT ["/mcp-nomad"]
