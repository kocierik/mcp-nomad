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
  -nomad-addr string
    	Nomad server address (default "http://localhost:4646")
  -port string
    	Port for HTTP server (default "8080")
  -transport string
    	Transport type (stdio, sse, or streamable-http) (default "stdio")
```

### Environment variables

- `NOMAD_ADDR`: Nomad HTTP API address (default: http://localhost:4646)
- `NOMAD_TOKEN`: Nomad ACL token (optional)
- `NOMAD_REGION`: forwarded as the REST `region` query parameter when callers do not override it (multi-region clusters)
- `NOMAD_NAMESPACE`: default namespace for tools that accept an optional namespace when the tool omits it
- TLS: `NOMAD_CACERT`, `NOMAD_SKIP_VERIFY`, `NOMAD_TLS_SERVER_NAME` (see `utils/client.go` / `buildTLSConfig`)

The HTTP client follows the official `/v1/` API and is split across `utils/client_*.go`; MCP tools depend on narrow interfaces in `utils/nomad_tool_interfaces.go`.

`NomadClient.MakeRequest` (used only for a few cluster/legacy call sites) rejects paths outside an internal allow-list — prefer typed helpers such as `StopAllocation`.

## Browse with MCP Inspector

Use this for **local testing and debugging** — not required for Claude Desktop daily use.

To run the latest published npm build under the MCP Inspector:

```bash
npx @modelcontextprotocol/inspector npx @kocierik/mcp-nomad
```

### Inspector with a local HTTP server (optional)

Default transport is **stdio**. To attach the Inspector as **Streamable HTTP**, start the binary in another terminal first:

```bash
go run . -transport=streamable-http -port=8080
```

Then open **`http://localhost:8080/mcp`** in the Inspector. For `-transport=sse`, use **`http://localhost:8080/sse`**.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
