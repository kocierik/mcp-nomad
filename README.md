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
  <a href="https://goreportcard.com/report/github.com/kocierik/mcp-nomad"><img src="https://goreportcard.com/badge/github.com/kocierik/mcp-nomad" alt="Go Report"></a>
  <a href="https://github.com/kocierik/mcp-nomad/releases/latest"><img src="https://img.shields.io/github/v/release/kocierik/mcp-nomad?logo=github&color=22ff22" alt="latest release badge"></a>
  <a href="https://github.com/kocierik/mcp-nomad/blob/main/LICENSE"><img src="https://img.shields.io/github/license/kocierik/mcp-nomad" alt="license badge"></a>
</p>

## Features

<details>
<summary>Job Management</summary>

- List jobs
- Get job details
- Run jobs
- Stop jobs
- Get job versions
- Get job submission
- List job allocations
- List job evaluations
- List job deployments
- Get job deployment
- Get job summary
- Update jobs
- Dispatch parameterized jobs
- Revert jobs
- Set job stability
- Create job evaluations
- Create job plans
- Force new periodic instances
- Get job scale status
- Scale task groups
- List job services
</details>




<details>
<summary>Deployment Management</summary>

- List deployments
- Get deployment details
</details>

<details>
<summary>Namespace Management</summary>

- List namespaces
- Create namespaces
- Delete namespaces
</details>

<details>
<summary>Node Management</summary>

- List nodes
- Get node details
- Drain nodes
- Set node eligibility
</details>

<details>
<summary>Allocation Management</summary>

- List allocations
- Get allocation details
- Get allocation logs
- Get task logs
</details>

<details>
<summary>Variable Management</summary>

- List variables with filtering and pagination
- Get variable details
- Create variables with namespace support
- Delete variables with CAS support

Example variable operations:
```bash
# List variables in a namespace
list_variables namespace="my-namespace" prefix="my/path" per_page=10

# Get a specific variable
get_variable path="my/path" namespace="my-namespace"

# Create a variable
create_variable path="my/path" key="username" value="john" namespace="my-namespace"

# Delete a variable
delete_variable path="my/path" namespace="my-namespace"
```
</details>

<details>
<summary>Volume Management</summary>

- List volumes
- Get volume details
- Delete volumes
- List volume claims
- Delete volume claims
</details>

<details>
<summary>ACL Management</summary>

- List ACL tokens
- Get ACL token details
- Create ACL tokens
- Delete ACL tokens
- List ACL policies
- Get ACL policy details
- Create ACL policies
- Delete ACL policies
- List ACL roles
- Get ACL role details
- Create ACL roles
- Delete ACL roles
- Bootstrap ACL system
</details>

<details>
<summary>Sentinel Policy Management</summary>

- List Sentinel policies
- Get Sentinel policy details
- Create Sentinel policies
- Delete Sentinel policies
</details>

<details>
<summary>Cluster Management</summary>

- Get cluster leader
- List cluster peers
- List regions
</details>

## Browse With Inspector

To use the latest published version with Inspector:

```bash
npx @modelcontextprotocol/inspector npx @kocierik/mcp-nomad
```


### Options Available
```
  -nomad-addr string
    	Nomad server address (default "http://localhost:4646")
  -port string
    	Port for SSE server (default "8080")
  -transport string
    	Transport type (stdio or sse) (default "stdio")
```

### Environment Variables

- `NOMAD_ADDR`: Nomad HTTP API address (default: http://localhost:4646)
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
        "-e", "NOMAD_TOKEN=secret-token-acl-optional", 
        "kocierik/mcpnomad-server:latest"
      ]
    }
  }
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
