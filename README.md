# Nomad MCP Server

A server that provides a set of tools for managing Nomad clusters through the MCP (Model Control Protocol) interface.

## Features

- Job management (list, get, run, stop, restart, scale)
- Deployment management (list, get)
- Namespace management (list, create, delete)
- Node management (list, get, drain)
- Allocation management (list, get, stop)
- Variable management (list, get, create, delete)
- Volume management (list, get, create, delete)
- ACL management (tokens, policies, roles)
- Access to job templates

## Installation

1. Clone the repository:
```bash
git clone https://github.com/kocierik/nomad-mcp-server.git
cd nomad-mcp-server
```

2. Build the server:
```bash
go build
```

## Configuration

The server requires the following environment variables:

- `NOMAD_ADDR`: The address of the Nomad API server (default: http://localhost:4646)
- `NOMAD_TOKEN`: The Nomad API token (optional)

## Usage

Run the server:
```bash
./nomad-mcp-server
```

## Tools

### Job Tools

#### list_jobs
Lists all jobs in a namespace.

Parameters:
- `namespace` (string, optional): Namespace to list jobs from
- `status` (string, optional): Filter jobs by status

#### get_job
Gets details of a specific job.

Parameters:
- `job_id` (string, required): ID of the job to get
- `namespace` (string, optional): Namespace of the job

#### run_job
Runs a new job.

Parameters:
- `job_spec` (string, required): Job specification in HCL or JSON format
- `namespace` (string, optional): Namespace to run the job in
- `detach` (boolean, optional): Whether to detach from the job

Note: The job specification can be provided in either HCL or JSON format. If HCL is provided, it will be automatically converted to JSON before submission to Nomad.

#### stop_job
Stops a job.

Parameters:
- `job_id` (string, required): ID of the job to stop
- `namespace` (string, optional): Namespace of the job
- `purge` (boolean, optional): Whether to purge the job

#### restart_job
Restarts a job.

Parameters:
- `job_id` (string, required): ID of the job to restart
- `namespace` (string, optional): Namespace of the job

#### scale_job
Scales a job's task group.

Parameters:
- `job_id` (string, required): ID of the job to scale
- `group` (string, required): Name of the task group to scale
- `count` (integer, required): Desired count for the task group
- `namespace` (string, optional): Namespace of the job

#### get_job_allocations
Gets allocations for a job.

Parameters:
- `job_id` (string, required): ID of the job
- `namespace` (string, optional): Namespace of the job

#### get_job_evaluations
Gets evaluations for a job.

Parameters:
- `job_id` (string, required): ID of the job
- `namespace` (string, optional): Namespace of the job

### Deployment Tools

#### list_deployments
Lists all deployments.

Parameters:
- `namespace` (string, optional): Namespace to list deployments from

#### get_deployment
Gets details of a specific deployment.

Parameters:
- `deployment_id` (string, required): ID of the deployment to get

### Namespace Tools

#### list_namespaces
Lists all namespaces.

#### create_namespace
Creates a new namespace.

Parameters:
- `namespace_spec` (string, required): JSON specification of the namespace to create

#### delete_namespace
Deletes a namespace.

Parameters:
- `name` (string, required): Name of the namespace to delete

### Node Tools

#### list_nodes
Lists all nodes in the cluster.

Parameters:
- `status` (string, optional): Filter nodes by status

#### get_node
Gets details of a specific node.

Parameters:
- `node_id` (string, required): ID of the node to get

#### drain_node
Drains a node.

Parameters:
- `node_id` (string, required): ID of the node to drain
- `enable` (boolean, required): Whether to enable or disable drain mode
- `deadline` (integer, optional): Deadline in seconds for the drain operation

### Allocation Tools

#### list_allocations
Lists all allocations.

Parameters:
- `namespace` (string, optional): Namespace to list allocations from

#### get_allocation
Gets details of a specific allocation.

Parameters:
- `allocation_id` (string, required): ID of the allocation to get

#### stop_allocation
Stops an allocation.

Parameters:
- `allocation_id` (string, required): ID of the allocation to stop

### Variable Tools

#### list_variables
Lists all variables with an optional prefix filter.

Parameters:
- `prefix` (string, optional): Prefix to filter variables by

#### get_variable
Gets details of a specific variable.

Parameters:
- `path` (string, required): Path of the variable to get

#### create_variable
Creates or updates a variable.

Parameters:
- `path` (string, required): Path of the variable to create/update
- `items` (object, required): Key-value pairs for the variable

#### delete_variable
Deletes a variable.

Parameters:
- `path` (string, required): Path of the variable to delete

### Volume Tools

#### list_volumes
Lists all volumes in a namespace.

Parameters:
- `namespace` (string, optional): Namespace to list volumes from

#### get_volume
Gets details of a specific volume.

Parameters:
- `volume_id` (string, required): ID of the volume to get
- `namespace` (string, optional): Namespace of the volume

#### create_volume
Creates a new volume.

Parameters:
- `volume_spec` (string, required): JSON specification of the volume to create

#### delete_volume
Deletes a volume.

Parameters:
- `volume_id` (string, required): ID of the volume to delete
- `namespace` (string, optional): Namespace of the volume

### ACL Tools

#### Token Management

##### list_acl_tokens
Lists all ACL tokens.

##### get_acl_token
Gets details of a specific ACL token.

Parameters:
- `accessor_id` (string, required): Accessor ID of the token to get

##### create_acl_token
Creates a new ACL token.

Parameters:
- `token_spec` (string, required): JSON specification of the token to create

##### delete_acl_token
Deletes an ACL token.

Parameters:
- `accessor_id` (string, required): Accessor ID of the token to delete

#### Policy Management

##### list_acl_policies
Lists all ACL policies.

##### get_acl_policy
Gets details of a specific ACL policy.

Parameters:
- `name` (string, required): Name of the policy to get

##### create_acl_policy
Creates a new ACL policy.

Parameters:
- `policy_spec` (string, required): JSON specification of the policy to create

##### delete_acl_policy
Deletes an ACL policy.

Parameters:
- `name` (string, required): Name of the policy to delete

#### Role Management

##### list_acl_roles
Lists all ACL roles.

##### get_acl_role
Gets details of a specific ACL role.

Parameters:
- `id` (string, required): ID of the role to get

##### create_acl_role
Creates a new ACL role.

Parameters:
- `role_spec` (string, required): JSON specification of the role to create

##### delete_acl_role
Deletes an ACL role.

Parameters:
- `id` (string, required): ID of the role to delete

## Resources

### Job Templates

The server provides access to job templates that can be used as a starting point for creating new jobs. These templates are available through the `job_templates` resource.

## License

MIT