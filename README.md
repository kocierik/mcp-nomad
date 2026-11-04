<h4 align="center">Golang-based MCP server connecting to Nomad</h4>

<h1 align="center">
  <img src="https://github.com/user-attachments/assets/77e291ef-11ae-4b12-94b1-3409f4356ceb" alt="nomad-futuristic-logo" style="width:200px;"/>
   <br/>
   MCP Nomad Go
</h1>

<p align="center">
  <a href="#quick-start">Quick start</a> ⚙
  <a href="#available-tools">Tools</a> ⚙
  <a href="#configuration">Configuration</a> ⚙
  <a href="#http-mode-optional">HTTP mode</a> ⚙
  <a href="https://github.com/kocierik/mcp-nomad/blob/main/DEVELOPMENT.md">Development ↗</a> ⚙
  <a href="https://modelcontextprotocol.io">About MCP ↗</a>
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/kocierik/mcp-nomad"><img src="https://goreportcard.com/badge/github.com/kocierik/mcp-nomad" alt="Go Report"></a>
  <a href="https://github.com/kocierik/mcp-nomad/releases/latest"><img src="https://img.shields.io/github/v/release/kocierik/mcp-nomad?logo=github&color=22ff22" alt="latest release badge"></a>
  <a href="https://github.com/kocierik/mcp-nomad/blob/main/LICENSE"><img src="https://img.shields.io/github/license/kocierik/mcp-nomad" alt="license badge"></a>
</p>

mcp-nomad lets AI assistants such as Claude inspect and operate a [HashiCorp Nomad](https://developer.hashicorp.com/nomad) cluster: jobs, allocations, logs, nodes, variables, ACLs and more.

https://github.com/user-attachments/assets/731621d7-0acf-4045-bacc-7b34a7d83648

## Quick start

Add the server to `claude_desktop_config.json` (requires Node.js):

```json
{
  "mcpServers": {
    "mcp_nomad": {
      "command": "npx",
      "args": ["-y", "@kocierik/mcp-nomad"],
      "env": {
        "NOMAD_ADDR": "http://127.0.0.1:4646",
        "NOMAD_TOKEN": "<your-nomad-acl-token>"
      }
    }
  }
}
```

Add `"-read-only"` to `args` to expose only tools that cannot change the cluster.

Other ways to install; then use `"command": "mcp-nomad"` with the same `env`:

- npm: `npm install -g @kocierik/mcp-nomad`
- Binary: download from [GitHub Releases](https://github.com/kocierik/mcp-nomad/releases/latest)
- Source: `go install github.com/kocierik/mcp-nomad@latest`

## Available tools

| Area | Read | Write |
| --- | --- | --- |
| Jobs | `list_jobs`, `get_job`, `get_job_summary`, `get_job_allocations`, `get_job_evaluations`, `get_evaluation`, `get_job_deployments`, `get_job_services` | `run_job`, `stop_job`, `scale_job` |
| Deployments | `list_deployments`, `get_deployment` | |
| Allocations | `list_allocations`, `get_allocation`, `get_allocation_logs` | `stop_allocation` |
| Nodes | `list_nodes`, `get_node` | `drain_node`, `eligibility_node` |
| Namespaces | `list_namespaces` | `create_namespace`, `delete_namespace` |
| Variables | `list_variables`, `get_variable` | `create_variable`, `delete_variable` |
| Volumes | `list_volumes`, `get_volume` | `delete_volume` |
| ACL | `list_acl_tokens`, `get_acl_token`, `list_acl_policies`, `get_acl_policy`, `list_acl_roles`, `get_acl_role` | `create_acl_token`, `delete_acl_token`, `create_acl_policy`, `delete_acl_policy`, `create_acl_role`, `delete_acl_role`, `bootstrap_acl_token` |
| Sentinel | `list_sentinel_policies`, `get_sentinel_policy` | `create_sentinel_policy`, `delete_sentinel_policy` |
| Cluster | `get_cluster_leader`, `list_cluster_peers`, `list_regions` | |

With `-read-only` the write column is not available at all. Whatever the mode, the Nomad ACLs attached to `NOMAD_TOKEN` decide what is actually allowed: use a token with only the permissions you want the assistant to have.

The server also provides resources (job specs and history, node status, allocation status and logs) and prompts.

## Configuration

| Variable | Purpose |
| --- | --- |
| `NOMAD_ADDR` | Nomad API address (default `http://127.0.0.1:4646`) |
| `NOMAD_TOKEN` | Nomad ACL token |
| `NOMAD_NAMESPACE` | Default namespace when a request does not specify one |
| `NOMAD_REGION` | Default region for multi-region clusters |
| `NOMAD_CACERT`, `NOMAD_TLS_SERVER_NAME`, `NOMAD_SKIP_VERIFY` | TLS to Nomad |
| `NOMAD_CLIENT_CERT`, `NOMAD_CLIENT_KEY` | Mutual TLS to Nomad |

| Flag | Default | Purpose |
| --- | --- | --- |
| `-transport` | `stdio` | `stdio`, `sse` or `streamable-http` |
| `-read-only` | off | Only register tools that cannot change the cluster |
| `-port` | `8080` | HTTP port |
| `-listen-address` | `127.0.0.1` | HTTP listen address |

## HTTP mode (optional)

Claude Desktop uses `stdio` and needs nothing from this section.

The `sse` and `streamable-http` transports are protected by a bearer token, separate from the Nomad token:

```bash
export MCP_TOKEN="<choose-a-long-random-token>"
mcp-nomad -transport=streamable-http
```

- Every request must send `Authorization: Bearer <MCP_TOKEN>`; the server refuses to start without it.
- The server listens on `127.0.0.1` only. Use `-listen-address=0.0.0.0` (for example in Docker) only behind a trusted network or HTTPS reverse proxy.
- Browser clients from other origins must be listed in `MCP_ALLOWED_ORIGINS` (comma-separated, e.g. `http://localhost:6274`).

Details: [SECURITY.md](SECURITY.md).

> **Upgrading from 0.3.x:** HTTP transports now require `MCP_TOKEN` and listen on `127.0.0.1` by default.

## Testing with MCP Inspector

```bash
npx @modelcontextprotocol/inspector npx @kocierik/mcp-nomad
```

To test HTTP mode, start the server as shown above with `MCP_ALLOWED_ORIGINS="http://localhost:6274"`, open `http://localhost:8080/mcp` in the Inspector and enter `MCP_TOKEN` as the bearer token.

## License

MIT, see [LICENSE](LICENSE).
