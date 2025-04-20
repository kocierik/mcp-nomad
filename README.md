# Nomad MCP 🚀

A Go implementation of the Model Context Protocol (MCP) for HashiCorp Nomad, enabling seamless integration between LLM applications and Nomad cluster management.

## Overview

Nomad MCP provides a standardized interface for AI models to interact with Nomad clusters. It exposes Nomad's core functionality as tools and resources through the MCP protocol, allowing Large Language Models to:

- List, get, create, and manage Nomad jobs
- Monitor deployments 
- Manage namespaces
- View and control cluster nodes
- Access job templates and examples

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/nomad-mcp.git
cd nomad-mcp

# Build the executable
go build -o nomad-mcp

# Run the MCP server
./nomad-mcp
```

## Configuration

Nomad MCP uses the following environment variables:

- `NOMAD_ADDR` - The address of your Nomad server (default: http://localhost:4646)
- `NOMAD_TOKEN` - Nomad ACL token (if ACLs are enabled)

## Available Tools

### Job Management

- `list_jobs` - List all jobs in a namespace
- `get_job` - Retrieve details for a specific job
- `run_job` - Run or update a job from a specification
- `stop_job` - Stop a running job

### Deployment Management

- `list_deployments` - List all deployments
- `get_deployment` - Retrieve details for a specific deployment

### Namespace Management

- `list_namespaces` - List all namespaces
- `create_namespace` - Create a new namespace
- `delete_namespace` - Delete an existing namespace

### Node Management

- `list_nodes` - List all nodes in the cluster
- `get_node` - Retrieve details for a specific node
- `drain_node` - Enable or disable drain mode for a node

## Available Resources

### Job Templates

- `nomad://templates` - List all available job templates
- `nomad://templates/{name}` - Retrieve a specific job template

## Usage Examples

### Listing Jobs

```
Call tool: list_jobs
```

### Getting Job Details

```
Call tool: get_job
Parameters:
  job_id: web-app
  namespace: default
```

### Running a Job

```
Call tool: run_job
Parameters:
  job_spec: |
    job "example" {
      datacenters = ["dc1"]
      type = "service"
      group "example" {
        count = 1
        task "server" {
          driver = "docker"
          config {
            image = "nginx:latest"
            ports = ["http"]
          }
        }
      }
    }
```

### Using Job Templates

```
// First, list available templates
Read resource: nomad://templates

// Then, get a specific template
Read resource: nomad://templates/service

// Modify as needed and run the job
Call tool: run_job
Parameters:
  job_spec: <modified template content>
```

## Project Structure

```
nomad-mcp/
  ├── main.go              # Main entry point
  ├── tools/               # Tool implementations
  │   ├── jobs.go          # Job management tools
  │   ├── deployments.go   # Deployment management tools
  │   ├── namespaces.go    # Namespace management tools
  │   └── nodes.go         # Node management tools
  ├── utils/               # Utility functions
  │   ├── client.go        # Nomad API client
  │   └── templates.go     # Job template utilities
  └── types/               # Type definitions
      ├── jobs.go          # Job-related types
      ├── deployments.go   # Deployment-related types
      ├── namespaces.go    # Namespace-related types
      └── nodes.go         # Node-related types
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.