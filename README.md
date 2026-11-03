<h4 align="center">Golang-based MCP server connecting to Nomad</h4>

<h1 align="center">
  <img src="https://github.com/user-attachments/assets/77e291ef-11ae-4b12-94b1-3409f4356ceb" alt="nomad-futuristic-logo" style="width:200px;"/>
   <br/>
   MCP Nomad Go
</h1>

<p align="center">
  <a href="#use-with-claude">Use with Claude</a> ⚙
  <a href="#server-options">Server options</a> ⚙
  <a href="#browse-with-mcp-inspector">MCP Inspector (testing)</a> ⚙
  <a href="https://github.com/kocierik/mcp-nomad/blob/main/CONTRIBUTING.md">Contributing ↗</a> ⚙
  <a href="https://modelcontextprotocol.io">About MCP ↗</a>
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/kocierik/mcp-nomad"><img src="https://goreportcard.com/badge/github.com/kocierik/mcp-nomad" alt="Go Report"></a>
  <a href="https://github.com/kocierik/mcp-nomad/releases/latest"><img src="https://img.shields.io/github/v/release/kocierik/mcp-nomad?logo=github&color=22ff22" alt="latest release badge"></a>
  <a href="https://github.com/kocierik/mcp-nomad/blob/main/LICENSE"><img src="https://img.shields.io/github/license/kocierik/mcp-nomad" alt="license badge"></a>
</p>

## Use With Claude

https://github.com/user-attachments/assets/731621d7-0acf-4045-bacc-7b34a7d83648


### Installation Options

|              | <a href="#prebuilt-from-npm">Pre-built NPM</a> | <a href="#from-github-releases">Pre-built in Github</a> | <a href="#building-from-source">From sources</a> |
| ------------ | ---------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------ |
| Claude Setup | Manual                                         | Manual                                                  | Manual                                           |
| Prerequisite | Node.js                                        | None                                                    | Golang                                           |

### Prebuilt from npm

The package publishes a **`mcp-nomad`** CLI. Easiest zero-install option (downloads to npm’s cache; needs Node/npm):

```bash
npx -y @kocierik/mcp-nomad
```

Or install globally so `mcp-nomad` is on your `PATH`:

```bash
npm install -g @kocierik/mcp-nomad
```

`claude_desktop_config.json` with **`npx`** (recommended):

```json
{
  "mcpServers": {
    "mcp_nomad": {
      "command": "npx",
      "args": ["-y", "@kocierik/mcp-nomad"],
      "env": {
        "NOMAD_TOKEN": "${NOMAD_TOKEN}",
        "NOMAD_ADDR": "${NOMAD_ADDR}"
      }
    }
  }
}
```

If you used **`npm install -g`**, keep `command` / `args` as the binary directly:

```json
{
  "mcpServers": {
    "mcp_nomad": {
      "command": "mcp-nomad",
      "args": [],
      "env": {
        "NOMAD_TOKEN": "${NOMAD_TOKEN}",
        "NOMAD_ADDR": "${NOMAD_ADDR}"
      }
    }
  }
}
```

### From GitHub Releases

Download the binary and configure Claude Desktop like so:

```json
{
  "mcpServers": {
    "mcp_nomad": {
      "command": "mcp-nomad",
      "args": [],
      "env": {
        "NOMAD_TOKEN": "${NOMAD_TOKEN}",
        "NOMAD_ADDR": "${NOMAD_ADDR}"
      }
    }
  }
}
```

### Building from Source

```bash
go get github.com/kocierik/mcp-nomad
go install github.com/kocierik/mcp-nomad
```

## Server options

Command-line flags (also relevant when pairing with MCP Inspector against a manually started binary):

```
  -transport string
	    Transport type (stdio, sse, or streamable-http) (default "stdio")
  -port string
	    Port for HTTP server (default "8080")
  -listen-address string
	    Address for HTTP server listener (default "127.0.0.1")
  -read-only
	    Register only read tools (default "false")
```

### Environment variables

- `NOMAD_ADDR`: Nomad HTTP API address (default: http://127.0.0.1:4646)
- `NOMAD_TOKEN`: Nomad ACL token (optional)
- `MCP_TOKEN`: required for `sse` and `streamable-http` transports; this is separate from `NOMAD_TOKEN`
- `MCP_ALLOWED_ORIGINS`: optional comma-separated exact HTTP(S) origins for HTTP transports (strict validation)
- `NOMAD_REGION`: forwarded as the REST `region` query parameter when callers do not override it (multi-region clusters)
- `NOMAD_NAMESPACE`: default namespace for tools/resources that accept an optional namespace when the request omits it (an explicit namespace argument/query still overrides this default)
- TLS: `NOMAD_CACERT`, `NOMAD_SKIP_VERIFY`, `NOMAD_TLS_SERVER_NAME` (see `utils/client.go` / `buildTLSConfig`)

The HTTP client follows the official `/v1/` API and is split across `utils/client_*.go`; MCP tools depend on narrow interfaces in `utils/nomad_tool_interfaces.go`.

`NomadClient.MakeRequest` (used only for a few cluster/legacy call sites) rejects paths outside an internal allow-list — prefer typed helpers such as `StopAllocation`.

### HTTP authentication and origin policy

- HTTP transports require `Authorization: Bearer <MCP_TOKEN>` on every request.
- Authentication identity is a static shared token; it is not OAuth and does not provide per-user delegation.
- Authentication and downstream Nomad authorization are separate: `MCP_TOKEN` secures MCP ingress; `NOMAD_TOKEN` is used only for Nomad API calls.
- Requests with no `Origin` header are allowed (after successful auth).
- Duplicate `Origin` headers are rejected.
- `MCP_ALLOWED_ORIGINS` entries must be exact origins only: `scheme://host[:port]`.
- Paths, query strings, fragments, and userinfo are rejected.
- Empty entries in `MCP_ALLOWED_ORIGINS` are rejected.
- Canonical default ports are applied, so `https://example.com` and `https://example.com:443` are treated identically (same for HTTP :80).
- Local defaults are always allowed for the configured port: `http://localhost:<port>`, `http://127.0.0.1:<port>`, `http://[::1]:<port>`.

Set HTTP auth/origin environment variables using placeholders:

```bash
export MCP_TOKEN="<set-a-static-ingress-token>"
export MCP_ALLOWED_ORIGINS="http://localhost:6274,http://127.0.0.1:6274"
```

```powershell
$env:MCP_TOKEN = "<set-a-static-ingress-token>"
$env:MCP_ALLOWED_ORIGINS = "http://localhost:6274,http://127.0.0.1:6274"
```

The MCP Inspector browser UI runs on port `6274`; include those origins explicitly when testing HTTP transports.

### Tool behavior notes

- `-read-only=true` removes all write tools from discovery and from callable operations.
- Structured tool outputs expose objects directly and wrap arrays as `{ "items": [...] }`.
- Tool text output is preserved as fallback content.
- Allocation log responses remain text-only.
- Embedded README/license resources are served from compiled assets, not from the current working directory.
- System resource `start_time` is initialized once at startup and remains stable for the running process.
- Variable CAS semantics:
  - Create: omitted `cas` means no CAS query; explicit `cas: 0` creates only if the variable does not already exist.
  - Delete: omitted `cas` is unconditional delete; explicit `cas: 0` is a conditional CAS=0 delete.
  - Large numeric values are decoded with `UseNumber` and preserved as raw numeric values in structured output.
- `lock_operation` for variables is deliberately unsupported and rejected.

### Code Mode scope

This server is compatible with client-side Code Mode orchestration through the MCP tool catalog. It does not implement a runtime gateway, arbitrary code execution endpoint, or server-side search/execute runtime.

## Browse with MCP Inspector

Use this for **local testing and debugging** — not required for Claude Desktop daily use.

To run the latest published npm build under the MCP Inspector:

```bash
npx @modelcontextprotocol/inspector npx @kocierik/mcp-nomad
```

### Inspector with a local HTTP server (optional)

Default transport is **stdio**. To attach the Inspector as **Streamable HTTP**, start the binary in another terminal first:

```bash
go run . -transport=streamable-http -listen-address=127.0.0.1 -port=8080
```

Then open **`http://localhost:8080/mcp`** in the Inspector. For `-transport=sse`, use **`http://localhost:8080/sse`**.

For container or remote-host ingress, use `-listen-address=0.0.0.0` intentionally and put TLS/auth hardening at a trusted reverse proxy boundary.

If running in Docker, `localhost` inside the container is the container itself. Set `NOMAD_ADDR` to a reachable host/network address for the Nomad cluster.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
