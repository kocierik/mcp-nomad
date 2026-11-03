# Security Policy

## Supported Versions

The following versions of **mcp-nomad** are currently supported with security updates:

| Version | Supported |
| ------- | --------- |
| 0.3.x   | ✅        |
| < 0.3   | ❌        |

> Only the latest minor release line receives security updates. Users are strongly encouraged to upgrade to the most recent version.

## HTTP Transport Security Model

- `stdio` remains the default transport and does not use HTTP ingress auth.
- `sse` and `streamable-http` require a separate MCP ingress token: `MCP_TOKEN`.
- `NOMAD_TOKEN` is never used as MCP ingress auth.
- Clients must send `Authorization: Bearer <MCP_TOKEN>` on every HTTP request.
- Auth failures return 401 with `WWW-Authenticate: Bearer`.
- Origin validation failures return 403.

### Origin Policy

- `MCP_ALLOWED_ORIGINS` is optional and comma-separated.
- Values must be exact HTTP(S) origins only (`scheme://host[:port]`).
- Path/query/fragment/userinfo are rejected.
- Empty list entries are rejected.
- Canonical default ports are normalized (`http` => `:80`, `https` => `:443`).
- Requests with no `Origin` header are allowed after auth.
- Duplicate or ambiguous `Origin` headers are rejected.
- Default loopback origins are always allowed for the configured port: `localhost`, `127.0.0.1`, `[::1]`.

### Deployment Guidance

- HTTP transport binds to `127.0.0.1` by default.
- Use `-listen-address=0.0.0.0` only when remote/container access is explicitly intended.
- Terminate HTTPS at a trusted reverse proxy and restrict ingress to trusted networks.
- Static bearer token auth is shared-identity ingress control; this project does not claim OAuth delegation compliance.

---

## Reporting a Vulnerability

If you discover a security vulnerability in **mcp-nomad**, please report it responsibly.

### How to Report
- Open a **private security advisory** on GitHub (preferred), or  
- Contact the maintainer via GitHub profile or email

### What to Include
Please provide as much detail as possible:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if available)

### Response Time
- Initial response: **within 48 hours**
- Status update: **within 5 days**

### Disclosure Policy
- If the vulnerability is confirmed, a fix will be prepared and released as soon as possible
- You may be credited for the discovery (if desired)
- Public disclosure should be coordinated and only happen **after a patch is available**

### Out of Scope
- Issues in unsupported versions
- Non-security bugs or misconfigurations without security impact
