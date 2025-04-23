
<h4 align="center">Golang-based MCP server connecting to Nomad</h4>

<h1 align="center">
  <img src="https://github.com/user-attachments/assets/77e291ef-11ae-4b12-94b1-3409f4356ceb" alt="nomad-futuristic-logo" style="width:200px;"/>
   <br/>
   MCP Nomad Go
</h1>

<p align="center">
  <a href="#features">Features</a> ⚙
  <a href="#browse-with-inspector">Browse With Inspector</a> ⚙
  <a href="#use-with-claude">Use With Claude</a> ⚙
  <a href="https://github.com/kocierik/mcp-nomad/blob/main/CONTRIBUTING.md">Contributing ↗</a> ⚙
  <a href="https://modelcontextprotocol.io">About MCP ↗</a>
</p>

<p align="center">
  <a href="https://github.com/kocierik/mcp-nomad/actions/workflows/test.yaml"><img src="https://github.com/kocierik/mcp-nomad/actions/workflows/test.yaml/badge.svg"></a>
  <a href="https://goreportcard.com/report/github.com/kocierik/mcp-nomad"><img src="https://goreportcard.com/badge/github.com/kocierik/mcp-nomad" alt="Go Report"></a>
  <a href="https://github.com/kocierik/mcp-nomad/releases/latest"><img src="https://img.shields.io/github/v/release/kocierik/mcp-nomad?logo=github&color=22ff22" alt="latest release badge"></a>
  <a href="https://github.com/kocierik/mcp-nomad/blob/main/LICENSE"><img src="https://img.shields.io/github/license/kocierik/mcp-nomad" alt="license badge"></a>
</p>

## Features

MCP 💬 prompt 🗂️ resource 🤖 tool 

- 🗂️🤖 List Nomad jobs
- 💬🤖 List Nomad nodes
- 🤖 Get Nomad job status
- 🤖 Get Nomad allocation logs
- 🤖 Restart a Nomad job
- 🤖 List deployments
- 🤖 View allocation info
- 💬 Get node metrics

## Browse With Inspector

To use the latest published version with Inspector:

```bash
npx @modelcontextprotocol/inspector npx @kocierik/mcp-nomad
```


### Environment Variables

- `NOMAD_ADDR`: Nomad HTTP API address (e.g. http://localhost:4646)
- `NOMAD_TOKEN`: Nomad ACL token (optional)
 


## Use With Claude

https://github.com/user-attachments/assets/731621d7-0acf-4045-bacc-7b34a7d83648


### Installation Options

|              | <a href="#using-smithery">Smithery</a> | <a href="#using-mcp-get">mcp-get</a> | <a href="#prebuilt-from-npm">Pre-built NPM</a> | <a href="#from-github-releases">Pre-built in Github</a> | <a href="#building-from-source">From sources</a> | <a href="#using-docker">Using Docker</a> |
| ------------ | -------------------------------------- | ------------------------------------ | ---------------------------------------------- | ------------------------------------------------------- | ------------------------------------------------ | ---------------------------------------- |
| Claude Setup | Auto                                   | Auto                                 | Manual                                         | Manual                                                  | Manual                                           | Manual                                   |
| Prerequisite | Node.js                                | Node.js                              | Node.js                                        | None                                                    | Golang                                           | Docker                                   |

### Using Smithery

```bash
npx -y @smithery/cli install @kocierik/mcp-nomad --client claude
```

### Using mcp-get

```bash
npx @michaellatman/mcp-get@latest install @kocierik/mcp-nomad
```

### Prebuilt from npm

```bash
npm install -g @kocierik/mcp-nomad
```

Update your `claude_desktop_config.json`:

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
      "command": "mcp-nomad-go",
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

### Using Docker linux

```bash
docker run -i --rm --network=host kocierik/mcpnomad-server:latest
```

### Using Docker macos/windows

```bash
docker run -i --rm \
  -e NOMAD_ADDR=http://host.docker.internal:4646 \
  kocierik/mcpnomad-server:latest
```

### For Claude macos/windows:

```json
{
  "mcpServers": {
    "mcp_nomad": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e", "NOMAD_ADDR=http://localhost:4646",
        "-e", "NOMAD_TOKEN=secret-token-acl-optional", 
        "-e", "NOMAD_ADDR=http://host.docker.internal:4646",
        "mcpnomad/server:latest"
      ]
    }
  }
}
```

### For Claude linux:

```json
{
  "mcpServers": {
    "mcp_nomad": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "NOMAD_ADDR=http://172.17.0.1:4646",
        "kocierik/mcpnomad-server:latest"
      ]
    }
  }
}
```
