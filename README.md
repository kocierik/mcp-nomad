# Nomad MCP Server

A server that implements the MCP protocol to interact with HashiCorp Nomad. This server provides tools for managing Nomad jobs, deployments, namespaces, nodes, allocations, variables, and more.

## Features

- Job Management
  - List jobs
  - Get job details
  - Run jobs
  - Stop jobs
  - Restart jobs
  - Scale jobs
  - Get job allocations
  - Get job evaluations
- Deployment Management
  - List deployments
  - Get deployment details
- Namespace Management
  - List namespaces
  - Create namespaces
  - Delete namespaces
- Node Management
  - List nodes
  - Get node details
  - Drain nodes
- Allocation Management
  - List allocations
  - Get allocation details
  - Stop allocations
- Variable Management
  - List variables
  - Get variable details
  - Create/update variables
  - Delete variables
- Job Templates
  - Access to predefined job templates

## Development

1. Clone the repository:
```bash
# Clone the repository
git clone https://github.com/kocierik/nomad-mcp-server.git
cd nomad-mcp-server

# install dependencies
go mod tidy

# Run the MCP inspector
npx @modelcontextprotocol/inspector go run main.go
```

## Configuration

The server requires a Nomad API endpoint to be configured. By default, it will use `http://localhost:4646`. You can set the following environment variables to configure the connection:

- `NOMAD_ADDR`: The address of the Nomad server (default: `http://localhost:4646`)
- `NOMAD_TOKEN`: The authentication token for the Nomad server (optional)

## Usage

Run the server:
```bash
./nomad-mcp-server
```

The server implements the MCP protocol and can be used with any MCP-compatible client.

## Tools

### Job Tools

- `list_jobs`: List all jobs in Nomad
  - Parameters:
    - `namespace` (optional): The namespace to list jobs from
    - `status` (optional): Filter jobs by status (pending, running, dead)

- `get_job`: Get job details by ID
  - Parameters:
    - `job_id` (required): The ID of the job to retrieve
    - `namespace` (optional): The namespace of the job

- `run_job`: Run a new job or update an existing job
  - Parameters:
    - `job_spec` (required): The job specification in HCL or JSON format
    - `detach` (optional): Return immediately instead of monitoring deployment

- `stop_job`: Stop a running job
  - Parameters:
    - `job_id` (required): The ID of the job to stop
    - `namespace` (optional): The namespace of the job
    - `purge` (optional): Purge the job from Nomad instead of just stopping it

- `restart_job`: Restart a job
  - Parameters:
    - `job_id` (required): The ID of the job to restart
    - `namespace` (optional): The namespace of the job

- `scale_job`: Scale a job's task group
  - Parameters:
    - `job_id` (required): The ID of the job to scale
    - `group` (required): The task group to scale
    - `count` (required): The desired count of the task group
    - `namespace` (optional): The namespace of the job

- `get_job_allocations`: Get allocations for a job
  - Parameters:
    - `job_id` (required): The ID of the job
    - `namespace` (optional): The namespace of the job

- `get_job_evaluations`: Get evaluations for a job
  - Parameters:
    - `job_id` (required): The ID of the job
    - `namespace` (optional): The namespace of the job

### Deployment Tools

- `list_deployments`: List all deployments
  - Parameters:
    - `namespace` (optional): The namespace to list deployments from

- `get_deployment`: Get deployment details by ID
  - Parameters:
    - `deployment_id` (required): The ID of the deployment to retrieve

### Namespace Tools

- `list_namespaces`: List all namespaces in Nomad
- `create_namespace`: Create a new namespace
  - Parameters:
    - `name` (required): The name of the namespace to create
    - `description` (optional): Description of the namespace
- `delete_namespace`: Delete a namespace
  - Parameters:
    - `name` (required): The name of the namespace to delete

### Node Tools

- `list_nodes`: List all nodes in the Nomad cluster
  - Parameters:
    - `status` (optional): Filter nodes by status (ready, down)
- `get_node`: Get details for a specific node
  - Parameters:
    - `node_id` (required): The ID of the node to retrieve
- `drain_node`: Enable or disable drain mode for a node
  - Parameters:
    - `node_id` (required): The ID of the node to drain
    - `enable` (required): Enable or disable drain mode
    - `deadline` (optional): Deadline in seconds for the drain operation

### Allocation Tools

- `list_allocations`: List all allocations in Nomad
  - Parameters:
    - `namespace` (optional): The namespace to list allocations from
    - `job_id` (optional): Filter allocations by job ID
- `get_allocation`: Get allocation details by ID
  - Parameters:
    - `allocation_id` (required): The ID of the allocation to retrieve
- `stop_allocation`: Stop a running allocation
  - Parameters:
    - `allocation_id` (required): The ID of the allocation to stop

### Variable Tools

- `list_variables`: List all variables in Nomad
  - Parameters:
    - `prefix` (optional): Optional prefix to filter variables
- `get_variable`: Get variable details by path
  - Parameters:
    - `path` (required): The path of the variable to retrieve
- `create_variable`: Create or update a variable
  - Parameters:
    - `path` (required): The path where to create the variable
    - `items` (required): The key-value pairs to store in the variable
- `delete_variable`: Delete a variable
  - Parameters:
    - `path` (required): The path of the variable to delete

## Resources

### Job Templates

The server provides access to predefined job templates through the following resources:

- `nomad://templates`: List of all available job templates
- `nomad://templates/{name}`: Specific job template by name

## License

This project is licensed under the MIT License - see the LICENSE file for details.