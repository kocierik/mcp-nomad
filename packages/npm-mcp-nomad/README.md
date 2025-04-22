# MCP Nomad

This is a distribution of MCP server connecting to Nomad written in Golang.

Currently available:

| 💬 prompt | 🗂️ resource | 🤖 tool |

## Example usage with Claude Desktop

To use this MCP server with Claude Desktop you would firstly need to install it.

You have two options at the moment - use pre-built binaries or build it from source.

### Building from source

```bash
# Clone the repository
git clone https://github.com/kocierik/nomad-mcp-server.git
cd nomad-mcp-server

# install dependencies
go mod tidy
```

Then check if the server is working by running:

```bash
go run main.go -transport=stdio
```

Once built, you can proceed to add configuration to `claude_desktop_config.json` file:

```json
{
    "mcpServers": {
        "nomad_mcp": {
            "command": "/path/to/nomad-mcp-server",
            "args": [
                "-transport=stdio",
                "-nomad-addr=http://localhost:4646"
            ],
            "env": {
                "NOMAD_TOKEN": "${NOMAD_TOKEN}" // token for ACL
            }
        }
    }
}
```

### Using from Claude Desktop

Now you should be able to run Claude Desktop and:
- See Nomad clusters available to attach to conversation as a resource
- Ask Claude to list jobs and their status
- Ask Claude to deploy and manage jobs
- Ask Claude to manage namespaces and ACLs
- Ask Claude to monitor allocations and deployments
- Ask Claude to manage nodes and their status
- Ask Claude to handle variables and volumes

### Demo Operations

Here are some example operations you can perform:

1. **Job Management**:
   ```
   List all jobs in the default namespace
   Show me the status of job "my-service"
   Scale the "web" task group in job "my-service" to 3 instances
   ```

2. **Node Management**:
   ```
   List all nodes in the cluster
   Show me the status of node "node-1234"
   Enable drain mode for node "node-1234"
   ```

3. **Resource Management**:
   ```
   Create a new namespace called "development"
   List all variables in the system
   Show me all volumes in namespace "production"
   ```

4. **ACL Management**:
   ```
   Create a new ACL policy for read-only access
   List all ACL tokens
   Create a new role with specific policies
   ```

