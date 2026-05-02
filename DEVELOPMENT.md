# Development configuration for mcp-nomad

## Environment Variables

### Nomad Server
```bash
export NOMAD_ADDR=http://localhost:4646
export NOMAD_TOKEN=your-token-here
```

### MCP Server
The binary does not read `MCP_TRANSPORT` from `.env`; use CLI flags (default **`stdio`**). HTTP example:

```bash
./mcp-nomad -transport=streamable-http -port=8080
```

### Development
```bash
export DEBUG=true
export LOG_LEVEL=debug
export SKIP_INTEGRATION=false
```

## Quick Start

### 1. Start Nomad Server (with Docker)
```bash
make start-nomad
```

### 2. Check Nomad Status
```bash
make nomad-status
```

### 3. Run MCP Server
```bash
# Default transport is stdio (same as plain `make run`)
make run

# Explicit stdio
make run-stdio

# Legacy SSE transport
make run-sse

# Streamable HTTP (Inspector / `:8080/mcp`)
make run-http
```

### 4. Development Mode (with hot reload)
```bash
make dev
```

## Testing

### Run Tests
```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests only
make test-integration

# With coverage
make test-coverage
```

### Start Nomad for Integration Tests
```bash
make start-nomad
make test-integration
make stop-nomad
```

## Building

### Quick Build
```bash
make quick-build
```

### Release Build
```bash
make build-all
```

### Docker Build
```bash
make docker-build
make docker-run
```

## Development Tools

### Install Tools
```bash
make install-tools
```

### Code Quality
```bash
make lint
make format
make security
```

## Troubleshooting

### Check Status
```bash
make status
make version
```

### Clean Everything
```bash
make clean-all
```

### Reset Dependencies
```bash
make deps
```
