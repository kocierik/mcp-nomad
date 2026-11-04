# Development configuration for mcp-nomad

## Environment Variables

### Nomad Server
```bash
export NOMAD_ADDR=http://localhost:4646
export NOMAD_TOKEN=<your-nomad-acl-token>
```

### MCP Server
The binary does not read `MCP_TRANSPORT` from `.env`; use CLI flags (default **`stdio`**). HTTP example:

```bash
export MCP_TOKEN="<set-a-static-ingress-token>"
export MCP_ALLOWED_ORIGINS="http://localhost:6274,http://127.0.0.1:6274"
./mcp-nomad -transport=streamable-http -listen-address=127.0.0.1 -port=8080
```

PowerShell equivalent:

```powershell
$env:MCP_TOKEN = "<set-a-static-ingress-token>"
$env:MCP_ALLOWED_ORIGINS = "http://localhost:6274,http://127.0.0.1:6274"
./mcp-nomad -transport=streamable-http -listen-address=127.0.0.1 -port=8080
```

Notes:
- `MCP_TOKEN` is required for `sse` and `streamable-http`.
- `NOMAD_TOKEN` remains downstream-only for Nomad API calls.
- HTTP clients must send `Authorization: Bearer <MCP_TOKEN>` on every request.
- `MCP_ALLOWED_ORIGINS` values must be exact origins only (no path/query/fragment/userinfo).

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

Alternative launch command from source root:

```bash
go run . -transport=streamable-http -listen-address=127.0.0.1 -port=8080
```

`go run main.go` is not the supported project entrypoint for this repository layout.

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

`make docker-run` forwards `NOMAD_ADDR`, `NOMAD_TOKEN`, `MCP_TOKEN`, and `MCP_ALLOWED_ORIGINS` from your shell and binds `-listen-address=0.0.0.0` explicitly for container ingress.

If Nomad runs outside the container, avoid `NOMAD_ADDR=http://localhost:4646` inside containerized execution because container `localhost` does not point to the host Nomad process.

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

## Architecture notes

- The Nomad HTTP client follows the official `/v1/` API and is split across `utils/client_*.go`; MCP tools depend on narrow interfaces in `utils/nomad_tool_interfaces.go`.
- `NomadClient.MakeRequest` (a few cluster/legacy call sites) rejects paths outside an internal allow-list; prefer typed helpers such as `StopAllocation`.
- `tool_policy.go` holds the catalog of every tool. A tool missing from it is not exposed, and a test fails if the catalog and the registered tools diverge. Read tools get `readOnlyHint`, write tools `destructiveHint`; `-read-only` drops write tools from discovery and calls.
- README and LICENSE are embedded in the binary (`assets.go`) and served as `docs://` resources. `system://info` reports the build version and a process start time fixed at startup.

### Tool results

- JSON tools return `structuredContent` with an `outputSchema`, plus the original JSON as text for older clients. Objects are returned as-is; arrays are wrapped as `{ "items": [...] }`. Allocation logs stay text-only.
- Numbers are decoded with `UseNumber`, so large Nomad indexes keep their exact value.

### Variables and CAS

- `cas` is sent as Nomad's `cas` query parameter. Omitted: unconditional write or delete. `cas: 0` on create: only create if the variable does not exist yet. `cas: <ModifyIndex>`: only write or delete if nobody changed the variable since that index.
- `lock_operation` is not supported and is rejected before any request is sent.

### Code Mode

Clients that use Code Mode can orchestrate the existing tools; the structured outputs make that easier. The server itself does not run user code and has no `search`/`execute` tools.

### Live tests

`nomad_live_test.go` runs the full MCP server against a real Nomad API. It only runs when `NOMAD_INTEGRATION_ADDR` points to a disposable agent (it creates and deletes variables and a batch job): `NOMAD_INTEGRATION_ADDR=http://127.0.0.1:4646 make test-live`.

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
