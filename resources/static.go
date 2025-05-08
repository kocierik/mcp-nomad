// Package resources provides implementations of MCP resources for Nomad
package resources

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterStaticResources registers all static resources with the MCP server
func RegisterStaticResources(s *server.MCPServer, logger *log.Logger) {
	// README resource
	readmeResource := mcp.NewResource(
		"docs://readme",
		"Project README",
		mcp.WithResourceDescription("The project's README file"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(readmeResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile("README.md")
		if err != nil {
			logger.Printf("Error reading README: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	})

	// License resource
	licenseResource := mcp.NewResource(
		"docs://license",
		"Project License",
		mcp.WithResourceDescription("The project's license file"),
		mcp.WithMIMEType("text/plain"),
	)

	s.AddResource(licenseResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile("LICENSE")
		if err != nil {
			logger.Printf("Error reading LICENSE: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://license",
				MIMEType: "text/plain",
				Text:     string(content),
			},
		}, nil
	})

	// System Info resource
	systemInfoResource := mcp.NewResource(
		"system://info",
		"System Information",
		mcp.WithResourceDescription("Information about the Nomad cluster and MCP server"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(systemInfoResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		info := map[string]interface{}{
			"server_name":    "Nomad MCP Server",
			"server_version": "1.0.0",
			"start_time":     time.Now().Format(time.RFC3339),
			"capabilities": []string{
				"resources",
				"tools",
				"prompts",
			},
		}

		infoJSON, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "system://info",
				MIMEType: "application/json",
				Text:     string(infoJSON),
			},
		}, nil
	})

	// Help documentation resource
	helpResource := mcp.NewResource(
		"docs://help",
		"Help Documentation",
		mcp.WithResourceDescription("Documentation on how to use the MCP Nomad integration"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(helpResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		helpText := `# MCP Nomad Integration Help

This integration allows you to interact with your Nomad cluster using the Model Context Protocol.

## Available Resources

- Job specifications: nomad://jobs/{job_id}/spec
- Node status: nomad://nodes/{node_id}/status
- Allocation logs: nomad://allocations/{alloc_id}/logs
- Job history: nomad://jobs/{job_id}/history
- Node resources: nomad://nodes/{node_id}/resources
- Allocation status: nomad://allocations/{alloc_id}/status
- Cluster metrics: nomad://cluster/metrics
- Evaluations: nomad://evaluations/{eval_id}
- Service health: nomad://services/{service_name}/health
- Cluster policies: nomad://policies/list
- Cluster leader: nomad://cluster/leader
- Job allocations: nomad://jobs/{job_id}/allocations
- Job evaluations: nomad://jobs/{job_id}/evaluations
- Node allocations: nomad://nodes/{node_id}/allocations
- Allocation tasks: nomad://allocations/{alloc_id}/tasks

## Available Tools

Various tools are available for managing Nomad jobs, nodes, allocations, and other cluster resources.

### Job Tools
- Run a job
- Stop a job
- Get job information
- Scale a job

### Node Tools
- Drain a node
- Set node eligibility
- Get node information

### Cluster Tools
- Check cluster health
- List regions
- List namespaces

### Allocation Tools
- Get allocation logs
- Get allocation information
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://help",
				MIMEType: "text/markdown",
				Text:     helpText,
			},
		}, nil
	})

	// Nomad API documentation resource
	apiDocsResource := mcp.NewResource(
		"docs://api",
		"Nomad API Documentation",
		mcp.WithResourceDescription("Documentation on the Nomad HTTP API endpoints"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(apiDocsResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		apiDocsText := `# Nomad API Documentation

## Overview
This resource provides a high-level overview of the Nomad HTTP API endpoints used by this MCP integration.

## Common Endpoints

### Jobs
- GET /v1/jobs - List all jobs
- GET /v1/job/{job_id} - Get job details
- POST /v1/jobs - Create/update a job
- DELETE /v1/job/{job_id} - Stop a job

### Nodes
- GET /v1/nodes - List all nodes
- GET /v1/node/{node_id} - Get node details
- POST /v1/node/{node_id}/drain - Drain a node
- POST /v1/node/{node_id}/eligibility - Set node eligibility

### Allocations
- GET /v1/allocations - List all allocations
- GET /v1/allocation/{alloc_id} - Get allocation details
- GET /v1/client/fs/logs/{alloc_id} - Get allocation logs

### Evaluations
- GET /v1/evaluations - List all evaluations
- GET /v1/evaluation/{eval_id} - Get evaluation details

### Deployments
- GET /v1/deployments - List all deployments
- GET /v1/deployment/{deployment_id} - Get deployment details

### Namespaces
- GET /v1/namespaces - List all namespaces
- GET /v1/namespace/{name} - Get namespace details

## API Authentication
Most requests to the Nomad API require an ACL token, which should be provided in the X-Nomad-Token header.
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://api",
				MIMEType: "text/markdown",
				Text:     apiDocsText,
			},
		}, nil
	})

	// Complete Nomad documentation resource
	completeDocsResource := mcp.NewResource(
		"docs://nomad",
		"Complete Nomad Documentation",
		mcp.WithResourceDescription("Complete documentation for HashiCorp Nomad"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(completeDocsResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		completeDocsText := `# HashiCorp Nomad Documentation

## Introduction to Nomad

Nomad is a flexible workload orchestrator that enables an organization to easily deploy and manage any containerized or legacy application using a single, unified workflow. Nomad can run a diverse workload of Docker, non-containerized, microservice, and batch applications.

Nomad enables developers to use declarative infrastructure-as-code for deploying applications. Nomad uses bin packing to efficiently schedule jobs and optimize for resource utilization. Nomad is supported on macOS, Windows, and Linux.

### Key Features

- **Simple and Lightweight**: Nomad is a single binary that integrates into existing infrastructure. It's easy to operate and maintain.
- **Flexible Workload Support**: Nomad can run diverse workloads including Docker, non-containerized, microservices, and batch applications.
- **High Performance**: Nomad can schedule thousands of containers per second, providing low latency for production workloads.
- **Multi-Datacenter and Multi-Region Federation**: Nomad can span multiple datacenters in a single region or across multiple regions.
- **Bin Packing**: Nomad automatically places jobs to maximize resource utilization.
- **Declarative Jobs**: Jobs in Nomad are defined using a declarative job specification format.
- **Service Discovery Integration**: Nomad integrates with Consul for service discovery.

## Architecture

Nomad operates as a single binary with a client-server architecture. The server component provides the scheduling capabilities, while the client component runs on every machine that will host applications and registers with the servers.

### Server

The server component is responsible for accepting jobs from users, managing clients, and computing task placements. Servers form a consensus group using the Raft protocol for leader election and state replication.

### Client

The client component registers with the servers, watches for work, and executes tasks. Clients communicate with the servers using a Remote Procedure Call (RPC) protocol over a TLS connection.

### Regions and Datacenters

Nomad models infrastructure as regions and datacenters. A region is a logical boundary that may contain multiple datacenters. Servers manage state and make scheduling decisions for the region. Clients are assigned to a specific datacenter.

## Job Specification

Jobs are the primary configuration unit in Nomad. A job is a declarative specification that defines a set of tasks to be run. Jobs can be updated, scaled, or stopped.

### Job Structure

A job consists of one or more groups, each containing one or more tasks. A group is scheduled on a single client and acts as a unit of scaling and migration.

` + "```hcl" + `
job "example" {
  datacenters = ["dc1"]
  
  group "cache" {
    count = 1
    
    task "redis" {
      driver = "docker"
      
      config {
        image = "redis:6.0"
        port_map {
          db = 6379
        }
      }
      
      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
` + "```" + `

### Task Drivers

Nomad supports various task drivers that determine how a task is executed:

- **Docker**: For running Docker containers
- **Exec**: For executing commands directly on the client
- **Java**: For running Java applications
- **QEMU**: For running VMs via QEMU/KVM
- **Raw Exec**: For executing commands with elevated privileges
- **Container**: For running containers (LXC)

## Scheduling

Nomad uses a modular scheduler design with three primary scheduler types:

- **Service Scheduler**: For long-lived services
- **Batch Scheduler**: For batch jobs with a finite completion time
- **System Scheduler**: For running jobs on every client

The scheduling process involves three steps:

1. **Feasibility Checking**: Filtering out clients that cannot run the job
2. **Ranking**: Ranking feasible clients based on a scoring criteria
3. **Allocation**: Selecting the highest-ranking clients for placement

## Deployments

Nomad provides deployment capabilities for service jobs. A deployment is created when a job is created or updated. The deployment tracks the status of the job update and can perform automatic health checking and rollback.

### Update Strategies

- **Rolling Updates**: Incrementally update a service by replacing instances one at a time
- **Blue/Green Deployments**: Deploy a new version alongside the old one and switch traffic over
- **Canary Deployments**: Deploy a new version to a small subset of instances for testing

## Service Discovery

Nomad integrates with Consul for service discovery. When a task is started, it can register itself with Consul, allowing other services to discover and connect to it.

## Security

Nomad provides multiple security features:

- **TLS Communication**: All RPC communication can be encrypted with TLS
- **ACL System**: For access control and authorization
- **Sentinel Policies**: For governance and advanced policy enforcement

### Access Control Lists (ACLs)

The ACL system in Nomad controls access to data and APIs using a set of rules. It operates in a allow-by-default or deny-by-default policy.

## Monitoring

Nomad provides multiple ways to monitor the cluster:

- **Metrics**: Exposes metrics in formats compatible with various monitoring systems
- **Telemetry**: Integration with metrics collection systems
- **Logs**: Detailed logs for debugging and auditing
- **Web UI**: Visual dashboard for monitoring and management

## API

Nomad provides a comprehensive HTTP API for interacting with all aspects of the system. The API supports:

- **Job Management**: Create, update, and stop jobs
- **Cluster Management**: Monitor and manage the cluster
- **Resource Allocation**: Track and manage resource allocations
- **Evaluation Management**: Track and manage evaluations
- **ACL Management**: Manage ACL tokens and policies

## CLI

Nomad provides a command-line interface for interacting with the cluster:

` + "```" + `
# Start a Nomad agent
nomad agent -dev

# Submit a job
nomad job run example.nomad

# List running jobs
nomad job status

# Stop a job
nomad job stop example
` + "```" + `

## Enterprise Features

Nomad Enterprise offers additional features for organizations:

- **Namespaces**: For logically partitioning jobs and resources
- **Resource Quotas**: For limiting resource consumption
- **Sentinel Policies**: For governance and policy enforcement
- **Advanced Federation**: Enhanced multi-region capabilities

## Best Practices

- **Job Sizing**: Properly size jobs to avoid resource waste or contention
- **Task Groups**: Group related tasks that should run together
- **Constraints**: Use constraints to control placement
- **Resources**: Explicitly specify resource requirements
- **Service Discovery**: Use service discovery for communication between services
- **Scaling**: Plan for scaling based on application needs
- **Monitoring**: Set up comprehensive monitoring for the cluster

## Common Use Cases

- **Microservices**: Running containerized microservices
- **Batch Processing**: Running batch jobs and data processing
- **CI/CD Pipelines**: Running build and deployment pipelines
- **Legacy Applications**: Running non-containerized applications
- **High-Performance Computing**: Running compute-intensive workloads

## Comparison with Other Orchestrators

Nomad differentiates itself from other orchestrators like Kubernetes in several ways:

- **Simplicity**: Single binary, simple architecture
- **Flexibility**: Support for diverse workloads, not just containers
- **Lightweight**: Minimal resource overhead
- **Ease of Operation**: Easier to set up, run, and maintain
- **Performance**: Higher scheduling throughput
- **Integration**: Seamless integration with the HashiCorp stack

## Common Troubleshooting

- **Server Not Starting**: Check configuration, logs, and network connectivity
- **Client Not Connecting**: Check network connectivity and client configuration
- **Job Not Being Scheduled**: Check constraints, resources, and driver compatibility
- **Task Failing**: Check task logs, resource constraints, and driver issues

## Community and Support

- **Documentation**: [Nomad Documentation](https://www.nomadproject.io/docs)
- **Guides**: [Nomad Guides](https://learn.hashicorp.com/nomad)
- **Community Forum**: [Nomad Community Forum](https://discuss.hashicorp.com/c/nomad)
- **GitHub**: [Nomad GitHub Repository](https://github.com/hashicorp/nomad)
- **Professional Support**: Available for enterprise customers

## Reference

- **Job Specification**: [Job Specification Documentation](https://www.nomadproject.io/docs/job-specification)
- **Drivers**: [Task Drivers Documentation](https://www.nomadproject.io/docs/drivers)
- **API**: [API Documentation](https://www.nomadproject.io/api-docs)
- **CLI**: [CLI Documentation](https://www.nomadproject.io/docs/commands)
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://nomad",
				MIMEType: "text/markdown",
				Text:     completeDocsText,
			},
		}, nil
	})

	// Nomad job specification documentation
	jobSpecResource := mcp.NewResource(
		"docs://job-spec",
		"Nomad Job Specification Documentation",
		mcp.WithResourceDescription("Detailed documentation for Nomad job specifications"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(jobSpecResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		jobSpecText := `# Nomad Job Specification

## Overview

The Nomad job specification (jobspec) is a declarative configuration file that defines how Nomad should run your workloads. The jobspec is written in HashiCorp Configuration Language (HCL) or JSON and describes a job which is the primary configuration unit in Nomad.

## Job Stanza

The job stanza is the top-level stanza in the jobspec. It represents a single application or service deployed across the cluster.

` + "```hcl" + `
job "example" {
  datacenters = ["dc1"]
  type = "service"
  
  // Job configuration continues...
}
` + "```" + `

### Parameters

- **id** (string: job-name): A unique identifier for the job.
- **name** (string: job-id): A human-friendly name for the job.
- **datacenters** (array<string>): A list of datacenters where this job should be run.
- **type** (string: "service"): The type of job. Can be "service", "batch", or "system".
- **namespace** (string: "default"): The namespace to run the job in.
- **region** (string: "global"): The region to run the job in.
- **priority** (integer: 50): The job's priority relative to other jobs.
- **all_at_once** (boolean: false): Whether all task groups should be updated at once.
- **constraints** (array<Constraint>): Constraints restricting job placement.
- **affinity** (array<Affinity>): Affinities affecting job placement.
- **spread** (array<Spread>): Spread criteria for job placement.
- **meta** (map<string|string>): Metadata for the job.
- **parameterized** (Parameterized): Parameterized job configuration.
- **periodic** (Periodic): Periodic job configuration.
- **reschedule** (Reschedule): Reschedule policy for the job.
- **migrate** (Migrate): Migration strategy for the job.
- **update** (Update): Update strategy for the job.
- **vault** (Vault): Vault configuration for the job.
- **consul** (Consul): Consul configuration for the job.

## Group Stanza

The group stanza defines a series of tasks that should be co-located on the same client. Each job must have at least one group.

` + "```hcl" + `
job "example" {
  // ...
  
  group "cache" {
    count = 1
    
    // Group configuration continues...
  }
}
` + "```" + `

### Parameters

- **name** (string: group-name): A unique name for the group.
- **count** (integer: 1): The number of instances of this group to run.
- **constraints** (array<Constraint>): Constraints restricting group placement.
- **affinity** (array<Affinity>): Affinities affecting group placement.
- **spread** (array<Spread>): Spread criteria for group placement.
- **restart** (Restart): Restart policy for tasks in this group.
- **reschedule** (Reschedule): Reschedule policy for the group.
- **meta** (map<string|string>): Metadata for the group.
- **ephemeral_disk** (EphemeralDisk): Configuration for the ephemeral disk.
- **update** (Update): Update strategy for the group.
- **migrate** (Migrate): Migration strategy for the group.
- **network** (Network): Network configuration for the group.
- **services** (array<Service>): Services provided by the group.
- **shutdown_delay** (string: "0s"): Delay between stopping tasks and killing them.
- **stop_after_client_disconnect** (string: ""): Time to wait after a client disconnect before stopping the allocation.
- **vault** (Vault): Vault configuration for the group.
- **consul** (Consul): Consul configuration for the group.

## Task Stanza

The task stanza creates an individual unit of work, such as a Docker container, web application, or batch process. Each group must have at least one task.

` + "```hcl" + `
job "example" {
  // ...
  
  group "cache" {
    // ...
    
    task "redis" {
      driver = "docker"
      
      config {
        image = "redis:6.0"
      }
      
      // Task configuration continues...
    }
  }
}
` + "```" + `

### Parameters

- **name** (string: task-name): A unique name for the task.
- **driver** (string): The driver to use for the task. Examples: "docker", "exec", "java".
- **config** (map<string|string>): Driver-specific configuration.
- **constraints** (array<Constraint>): Constraints restricting task placement.
- **affinity** (array<Affinity>): Affinities affecting task placement.
- **env** (map<string|string>): Environment variables for the task.
- **resources** (Resources): Resources required by the task.
- **restart** (Restart): Restart policy for the task.
- **meta** (map<string|string>): Metadata for the task.
- **kill_timeout** (string: "5s"): Time to wait for a task to gracefully stop.
- **logs** (Logs): Log rotation configuration.
- **artifact** (array<Artifact>): Artifacts to download for the task.
- **template** (array<Template>): Templates to render for the task.
- **dispatch_payload** (DispatchPayload): Configuration for handling dispatch payloads.
- **volume_mount** (array<VolumeMount>): Volume mounts for the task.
- **leader** (boolean: false): Whether this task is the leader for the group.
- **shutdown_delay** (string: "0s"): Delay between stopping the task and killing it.
- **kill_signal** (string: ""): Signal to use for killing the task.
- **vault** (Vault): Vault configuration for the task.
- **consul** (Consul): Consul configuration for the task.

## Resources Stanza

The resources stanza describes the resources required by a task.

` + "```hcl" + `
task "redis" {
  // ...
  
  resources {
    cpu    = 500 # 500 MHz
    memory = 256 # 256 MB
    
    network {
      port "db" {
        static = 6379
      }
    }
  }
}
` + "```" + `

### Parameters

- **cpu** (integer: 100): The CPU required in MHz.
- **memory** (integer: 300): The memory required in MB.
- **disk** (integer: 0): The disk space required in MB.
- **network** (Network): Network requirements for the task.

## Network Stanza

The network stanza specifies the network requirements for a task or group.

` + "```hcl" + `
network {
  mode = "bridge"
  
  port "http" {
    static = 8080
    to     = 80
  }
}
` + "```" + `

### Parameters

- **mode** (string: "host"): The networking mode to use. Can be "host", "bridge", or "none".
- **hostname** (string: ""): The hostname to assign to the container.
- **port** (map<Port>): Port requirements and mappings.

## Service Stanza

The service stanza defines a service that should be registered with a service discovery provider, such as Consul.

` + "```hcl" + `
service {
  name = "redis"
  port = "db"
  
  check {
    type     = "tcp"
    port     = "db"
    interval = "10s"
    timeout  = "2s"
  }
}
` + "```" + `

### Parameters

- **name** (string: task-name): The service name to register.
- **port** (string: ""): The port to advertise for the service.
- **tags** (array<string>): Tags to associate with the service.
- **canary_tags** (array<string>): Tags to associate with canary instances.
- **enable_tag_override** (boolean: false): Whether to allow tag updates from Consul.
- **address_mode** (string: "auto"): How to determine the service address.
- **check** (array<Check>): Health checks for the service.
- **check_restart** (CheckRestart): When to restart on health check failures.
- **connect** (Connect): Configuration for Consul Connect.
- **meta** (map<string|string>): Metadata for the service.

## Constraint Stanza

The constraint stanza defines placement constraints for a job, group, or task.

` + "```hcl" + `
constraint {
  attribute = "${attr.kernel.name}"
  value     = "linux"
  operator  = "="
}
` + "```" + `

### Parameters

- **attribute** (string): The attribute to examine for the constraint.
- **value** (string): The value to compare against.
- **operator** (string): The comparison operator. Can be "=", "!=", ">", ">=", "<", "<=", "distinct_hosts", "distinct_property".

## Template Stanza

The template stanza instructs Nomad to manage a dynamic configuration file, rendered from a template.

` + "```hcl" + `
template {
  data        = "---\\nkey: {{ key \"service/redis/config\" }}\\n"
  destination = "local/redis.yml"
}
` + "```" + `

### Parameters

- **source** (string: ""): The path to the template source.
- **destination** (string: required): The destination path to write the rendered template.
- **data** (string: ""): The raw template to execute.
- **change_mode** (string: "restart"): What to do when the template changes.
- **change_signal** (string: ""): The signal to send when the template changes.
- **splay** (string: "5s"): How long to wait before restarting tasks.
- **perms** (string: "644"): The file permissions to assign.
- **left_delimiter** (string: "{{"): The left template delimiter.
- **right_delimiter** (string: "}}"): The right template delimiter.
- **env** (boolean: false): Whether to process environment variables.
- **vault** (Vault): Vault configuration for the template.
- **consul** (Consul): Consul configuration for the template.

## Variables Stanza

The variables stanza defines variables that can be accessed within the job.

` + "```hcl" + `
variables {
  region = "us-west-1"
  environment = "staging"
}
` + "```" + `

## Examples

### Complete Job Example

` + "```hcl" + `
job "redis" {
  datacenters = ["dc1"]
  type = "service"
  
  group "cache" {
    count = 1
    
    network {
      port "db" {
        static = 6379
      }
    }
    
    service {
      name = "redis"
      port = "db"
      
      check {
        type     = "tcp"
        port     = "db"
        interval = "10s"
        timeout  = "2s"
      }
    }
    
    task "redis" {
      driver = "docker"
      
      config {
        image = "redis:6.0"
        ports = ["db"]
      }
      
      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
` + "```" + `

### Batch Job Example

` + "```hcl" + `
job "batch-job" {
  datacenters = ["dc1"]
  type = "batch"
  
  group "example" {
    count = 1
    
    task "command" {
      driver = "exec"
      
      config {
        command = "/bin/bash"
        args    = ["-c", "echo 'Hello, World!' > /tmp/output.txt"]
      }
      
      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}
` + "```" + `

### System Job Example

` + "```hcl" + `
job "system-job" {
  datacenters = ["dc1"]
  type = "system"
  
  group "example" {
    task "server" {
      driver = "docker"
      
      config {
        image = "nginx:latest"
        ports = ["http"]
      }
      
      resources {
        cpu    = 100
        memory = 128
        
        network {
          port "http" {
            static = 80
          }
        }
      }
    }
  }
}
` + "```" + `
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://job-spec",
				MIMEType: "text/markdown",
				Text:     jobSpecText,
			},
		}, nil
	})

	// Nomad drivers documentation
	driversResource := mcp.NewResource(
		"docs://drivers",
		"Nomad Task Drivers Documentation",
		mcp.WithResourceDescription("Documentation on Nomad task drivers"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(driversResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		driversText := `# Nomad Task Drivers

## Overview

Task drivers in Nomad are used to execute tasks and provide isolation. They control how a task is run and determine what kind of isolation is needed.

## Available Drivers

### Docker Driver

The Docker driver is used to run Docker containers. It leverages the Docker engine for running containers.

#### Configuration

` + "```hcl" + `
task "redis" {
  driver = "docker"
  
  config {
    image = "redis:6.0"
    port_map {
      db = 6379
    }
    volumes = [
      "local/data:/data"
    ]
    auth {
      username = "username"
      password = "password"
    }
  }
}
` + "```" + `

#### Options

- **image** (string: required) - The Docker image to run.
- **command** (string: "") - The command to run when starting the container.
- **args** (array<string>: []) - Additional arguments to pass to the command.
- **auth** (block: optional) - Authentication information for private registries.
- **ports** (array<string>: []) - Ports to expose from the container.
- **port_map** (block: optional) - Mapping of port labels to ports within the container.
- **volumes** (array<string>: []) - Volumes to mount into the container.
- **network_mode** (string: "bridge") - Network mode for the container.

### Exec Driver

The exec driver is used to run command-line applications directly on the client. It provides minimal isolation.

#### Configuration

` + "```hcl" + `
task "script" {
  driver = "exec"
  
  config {
    command = "/bin/bash"
    args = ["-c", "echo 'Hello, World!' > /tmp/output.txt"]
  }
}
` + "```" + `

#### Options

- **command** (string: required) - The command to run.
- **args** (array<string>: []) - Arguments to the command.

### Java Driver

The Java driver is used to run Java applications. It manages the JVM and application classpath.

#### Configuration

` + "```hcl" + `
task "app" {
  driver = "java"
  
  config {
    jar_path = "local/app.jar"
    jvm_options = ["-Xmx2048m", "-Xms256m"]
    args = ["arg1", "arg2"]
  }
}
` + "```" + `

#### Options

- **jar_path** (string: required) - Path to the JAR file.
- **class** (string: "") - Name of the class to run.
- **jvm_options** (array<string>: []) - JVM options.
- **args** (array<string>: []) - Arguments to the Java application.

### QEMU Driver

The QEMU driver is used to run VMs with QEMU/KVM. It provides full VM isolation.

#### Configuration

` + "```hcl" + `
task "vm" {
  driver = "qemu"
  
  config {
    image_path = "local/image.qcow2"
    accelerator = "kvm"
    memory = 1024
    args = ["-device", "virtio-net,netdev=user.0"]
  }
}
` + "```" + `

#### Options

- **image_path** (string: required) - Path to the VM image.
- **accelerator** (string: "kvm") - The accelerator to use.
- **memory** (integer: 512) - The amount of memory to provide to the VM in MB.
- **cpu** (integer: 1) - The number of CPU cores to provide to the VM.

### Raw Exec Driver

The raw exec driver is used to run command-line applications with elevated privileges. It provides no isolation and should be used with caution.

#### Configuration

` + "```hcl" + `
task "installer" {
  driver = "raw_exec"
  
  config {
    command = "/bin/bash"
    args = ["-c", "apt-get update && apt-get install -y nginx"]
  }
}
` + "```" + `

#### Options

- **command** (string: required) - The command to run.
- **args** (array<string>: []) - Arguments to the command.

## Driver Compatibility

| Driver   | Linux | Windows | macOS |
|----------|-------|---------|-------|
| Docker   | Yes   | Yes     | Yes   |
| Exec     | Yes   | Yes     | Yes   |
| Java     | Yes   | Yes     | Yes   |
| QEMU     | Yes   | No      | No    |
| Raw Exec | Yes   | Yes     | Yes   |

## Best Practices

- Use the appropriate driver for your workload
- Prefer drivers with higher isolation in production environments
- Limit access to raw_exec driver
- Configure resource limits for each task
- Use volume mounts to persist data when needed
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://drivers",
				MIMEType: "text/markdown",
				Text:     driversText,
			},
		}, nil
	})

	// Nomad security documentation
	securityResource := mcp.NewResource(
		"docs://security",
		"Nomad Security Documentation",
		mcp.WithResourceDescription("Documentation on Nomad security features"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(securityResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		securityText := `# Nomad Security

## Overview

Nomad is designed with security as a core consideration. It provides multiple security features to protect your infrastructure and workloads.

## TLS Encryption

Nomad uses TLS (Transport Layer Security) to encrypt communications between clients and servers. This helps protect sensitive data and prevents man-in-the-middle attacks.

### Configuring TLS

TLS should be configured in the server and client configuration files:

` + "```hcl" + `
# Server configuration
server {
  enabled = true
  encrypt = "your-gossip-encryption-key"
  
  tls {
    enabled = true
    cert_file = "/path/to/cert.pem"
    key_file = "/path/to/key.pem"
    ca_file = "/path/to/ca.pem"
    verify_server_hostname = true
  }
}

# Client configuration
client {
  enabled = true
  
  tls {
    enabled = true
    cert_file = "/path/to/cert.pem"
    key_file = "/path/to/key.pem"
    ca_file = "/path/to/ca.pem"
    verify_server_hostname = true
  }
}
` + "```" + `

## Access Control Lists (ACLs)

Nomad's ACL system controls access to data and APIs. It operates in a deny-by-default model where all requests are denied unless explicitly allowed by a policy.

### ACL System Components

1. **Tokens**: Credentials used to authenticate API requests.
2. **Policies**: Rules that grant or deny access to resources.
3. **Capabilities**: Specific privileges that can be granted (e.g., read, write, deny).

### Configuring ACLs

To enable ACLs, add the following to your server configuration:

` + "```hcl" + `
acl {
  enabled = true
}
` + "```" + `

### ACL Policies

ACL policies define access rules. Here's an example policy:

` + "```hcl" + `
# Allow read-only access to all jobs
job {
  policy = "read"
}

# Allow full access to a specific namespace
namespace "dev" {
  policy = "write"
}

# Allow job submission, but not stopping jobs
job {
  policy = "write"
  capabilities = ["submit-job"]
}
` + "```" + `

## Vault Integration

Nomad integrates with HashiCorp Vault to provide secrets to tasks. This allows secure access to sensitive information like database credentials or API keys.

### Configuring Vault Integration

` + "```hcl" + `
vault {
  enabled = true
  address = "https://vault.example.com:8200"
  token = "VAULT_TOKEN"
  create_from_role = "nomad-cluster"
}
` + "```" + `

### Using Vault in Jobs

` + "```hcl" + `
job "app" {
  // ...
  
  group "web" {
    // ...
    
    task "server" {
      // ...
      
      vault {
        policies = ["app-policy"]
      }
      
      template {
        data = <<EOH
          DB_USER="{{ with secret "database/creds/app" }}{{ .Data.username }}{{ end }}"
          DB_PASS="{{ with secret "database/creds/app" }}{{ .Data.password }}{{ end }}"
        EOH
        destination = "secrets/config.env"
        env = true
      }
    }
  }
}
` + "```" + `

## Sentinel Policies (Enterprise)

Nomad Enterprise includes support for Sentinel, a policy-as-code framework for enforcing governance rules.

### Example Sentinel Policy

` + "```hcl" + `
# Require Docker driver to use specific registries
import "strings"

driver_is_docker = rule {
    keys_contain(job.task_groups[0].tasks[0], "driver") and
    job.task_groups[0].tasks[0].driver is "docker"
}

docker_image_is_allowed = rule {
    driver_is_docker and
    strings.has_prefix(job.task_groups[0].tasks[0].config.image, "registry.example.com/")
}

main = rule {
    docker_image_is_allowed
}
` + "```" + `

## Best Practices

1. **Enable TLS**: Always enable TLS in production environments.
2. **Rotate Certificates**: Regularly rotate TLS certificates and encryption keys.
3. **Enable ACLs**: Use ACLs to restrict access based on the principle of least privilege.
4. **Namespaces**: Use namespaces (Enterprise) to isolate different teams or environments.
5. **Audit Logging**: Enable audit logging to track access and changes.
6. **Validate Job Submissions**: Use Sentinel policies to validate job submissions.
7. **Secure Vault Integration**: Properly secure the Vault token used by Nomad.
8. **Network Segmentation**: Use network segmentation to isolate Nomad clusters.
9. **Regular Updates**: Keep Nomad updated to the latest version to get security patches.
10. **Restrict Raw Exec**: Limit access to the raw_exec driver due to its lack of isolation.
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://security",
				MIMEType: "text/markdown",
				Text:     securityText,
			},
		}, nil
	})

	// Nomad CLI documentation
	cliResource := mcp.NewResource(
		"docs://cli",
		"Nomad CLI Documentation",
		mcp.WithResourceDescription("Documentation on Nomad command-line interface"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(cliResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		cliText := `# Nomad Command-Line Interface

## Overview

The Nomad CLI provides a command-line interface for interacting with Nomad clusters. It allows you to manage jobs, nodes, allocations, and other resources.

## Installation

The Nomad CLI is distributed as a single binary. You can download it from the [Nomad website](https://www.nomadproject.io/downloads).

## Configuration

The Nomad CLI can be configured using environment variables or a configuration file.

### Environment Variables

- **NOMAD_ADDR**: The address of the Nomad server (default: http://127.0.0.1:4646)
- **NOMAD_REGION**: The region of the Nomad server
- **NOMAD_NAMESPACE**: The namespace to use
- **NOMAD_TOKEN**: The ACL token to use
- **NOMAD_CACERT**: The path to a CA certificate file
- **NOMAD_CLIENT_CERT**: The path to a client certificate file
- **NOMAD_CLIENT_KEY**: The path to a client key file

### Configuration File

The Nomad CLI can also be configured using a configuration file. The default location is $HOME/.nomadrc. Example:

` + "```hcl" + `
address = "https://nomad.example.com:4646"
token = "your-token"
ca_cert = "/path/to/ca.pem"
client_cert = "/path/to/cert.pem"
client_key = "/path/to/key.pem"
` + "```" + `

## Common Commands

### Server Management

` + "```bash" + `
# Start a Nomad agent in server mode
nomad agent -config=server.hcl

# List the members of the Nomad cluster
nomad server members

# Get information about the leader
nomad server leader

# Force a new leader election
nomad server force-leader
` + "```" + `

### Job Management

` + "```bash" + `
# List all jobs
nomad job status

# Get information about a specific job
nomad job status <job_id>

# Submit a job
nomad job run job.nomad

# Stop a job
nomad job stop <job_id>

# Get job allocations
nomad job allocs <job_id>

# Get job evaluations
nomad job evaluations <job_id>

# Get job history
nomad job history <job_id>

# Generate a specification file for a job
nomad job init

# Validate a job specification
nomad job validate job.nomad

# Plan a job update
nomad job plan job.nomad

# Scale a job
nomad job scale <job_id> <group> <count>

# Dispatch a parameterized job
nomad job dispatch <job_id> <input>
` + "```" + `

### Node Management

` + "```bash" + `
# List all nodes
nomad node status

# Get information about a specific node
nomad node status <node_id>

# Drain a node
nomad node drain -enable <node_id>

# Disable draining on a node
nomad node drain -disable <node_id>

# Set node eligibility
nomad node eligibility -enable <node_id>
nomad node eligibility -disable <node_id>
` + "```" + `

### Allocation Management

` + "```bash" + `
# List all allocations
nomad alloc status

# Get information about a specific allocation
nomad alloc status <alloc_id>

# Get logs from an allocation
nomad alloc logs <alloc_id>

# Stream logs from an allocation
nomad alloc logs -f <alloc_id>

# Stop an allocation
nomad alloc stop <alloc_id>

# Restart a task in an allocation
nomad alloc restart <alloc_id>
` + "```" + `

### Evaluation Management

` + "```bash" + `
# List all evaluations
nomad eval status

# Get information about a specific evaluation
nomad eval status <eval_id>
` + "```" + `

### Namespace Management (Enterprise)

` + "```bash" + `
# List all namespaces
nomad namespace list

# Create a namespace
nomad namespace create -description "Development environment" dev

# Delete a namespace
nomad namespace delete dev
` + "```" + `

### ACL Management

` + "```bash" + `
# Bootstrap the ACL system
nomad acl bootstrap

# List ACL policies
nomad acl policy list

# Create an ACL policy
nomad acl policy create -description "Developer policy" -rules @developer-policy.hcl developer

# Create an ACL token
nomad acl token create -name "Developer token" -policy developer

# List ACL tokens
nomad acl token list
` + "```" + `

## Advanced Usage

### Using the Nomad CLI in Scripts

The Nomad CLI can be used in scripts to automate tasks. You can use the -json flag to get machine-readable output:

` + "```bash" + `
# Get job information in JSON format
nomad job status -json <job_id> > job.json

# Parse JSON output with jq
nomad job status -json <job_id> | jq '.JobID'
` + "```" + `

### Using the Nomad CLI with Different Environments

You can use different environment variables or configuration files for different environments:

` + "```bash" + `
# Use a specific environment
NOMAD_ADDR=https://nomad-prod.example.com:4646 NOMAD_TOKEN=prod-token nomad job status

# Use a specific configuration file
NOMAD_CONFIG_PATH=/path/to/prod.nomadrc nomad job status
` + "```" + `

## Best Practices

1. **Use ACL Tokens**: Always use ACL tokens with appropriate permissions.
2. **Validate Jobs**: Always validate jobs before submitting them.
3. **Use Job Plans**: Use job plans to see the impact of job changes before applying them.
4. **Scripting**: Use JSON output and tools like jq for scripting.
5. **Configuration Files**: Use configuration files for different environments.
6. **Monitoring**: Use the CLI to monitor the health of your cluster.
7. **Automation**: Automate common tasks using the CLI in scripts.
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://cli",
				MIMEType: "text/markdown",
				Text:     cliText,
			},
		}, nil
	})

	// Nomad architecture documentation
	architectureResource := mcp.NewResource(
		"docs://architecture",
		"Nomad Architecture Documentation",
		mcp.WithResourceDescription("Documentation on Nomad's architecture and components"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(architectureResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		architectureText := `# Nomad Architecture

## Overview

Nomad has a flexible and scalable architecture that allows it to run workloads across multiple datacenters and regions. This document provides an overview of Nomad's architecture and its components.

## System Components

Nomad has two main components:

1. **Server**: Handles the cluster management, scheduling, and coordination tasks.
2. **Client**: Runs the tasks and communicates with the server.

### Server

The server component is responsible for:

- **Job registration**: Accepting and storing job definitions
- **Scheduling**: Placing tasks on appropriate clients
- **Evaluations**: Processing changes to the system and making scheduling decisions
- **Allocation**: Tracking the lifecycle of tasks
- **State management**: Maintaining the system state

Servers form a consensus group for leader election and state replication using the Raft protocol. Only the leader server can make scheduling decisions.

### Client

The client component is responsible for:

- **Task execution**: Running tasks as specified by the server
- **Resource isolation**: Ensuring tasks have the resources they need
- **Task monitoring**: Monitoring the health and status of tasks
- **Resource usage reporting**: Reporting resource usage back to the server

Clients register with the servers and communicate with them using a Remote Procedure Call (RPC) protocol over a TLS connection.

## Architectural Concepts

### Regions

A region is a logical boundary that may contain multiple datacenters. Each region operates independently and has its own servers.

### Datacenters

A datacenter is a physical or logical grouping of resources within a region. Clients are assigned to a specific datacenter.

### Jobs

A job is the primary configuration unit in Nomad. It defines a set of tasks and their requirements.

### Task Groups

A task group is a set of tasks that must be run together on the same client.

### Tasks

A task is the smallest unit of work in Nomad. It could be a Docker container, a binary, or another type of workload.

### Allocations

An allocation is the mapping of a task group to a client. It represents the act of Nomad placing a task group on a node.

### Evaluations

An evaluation is the process by which Nomad makes scheduling decisions. Evaluations are created in response to changes in the system, such as registering a job, updating a job, or a node failing.

### Deployments

A deployment is the process of updating a job. It tracks the status of the job update and can perform automatic health checking and rollback.

## Data Flow

1. A user submits a job to a Nomad server.
2. The server creates an evaluation to place the job.
3. The scheduler processes the evaluation and creates allocations.
4. The server sends the allocations to the appropriate clients.
5. The clients execute the tasks as specified in the allocations.
6. The clients report the status of the tasks back to the server.

## High Availability

Nomad is designed to be highly available. Multiple servers can be run in a single region, forming a consensus group. If the leader fails, another server is elected as the leader.

## Federation

Nomad supports federation across multiple regions. Each region operates independently but can communicate with other regions. This allows for global job submission and cross-region queries.

## Network Architecture

### Server-to-Server Communication

Servers within a region communicate with each other using a gossip protocol for membership and Raft for consensus. The gossip protocol uses UDP on port 4648, and Raft uses TCP on port 4647.

### Server-to-Client Communication

Servers communicate with clients using RPC over TLS. This communication uses TCP on port 4647.

### Client-to-Service Communication

Clients communicate with the services they run using various protocols and ports as specified in the job configuration.

### API Communication

The Nomad API is used for external communication with the cluster. It uses HTTP/HTTPS on port 4646.

## Security Architecture

Nomad has several security features:

1. **TLS**: All communication can be encrypted using TLS.
2. **ACLs**: Access Control Lists for authorization.
3. **Vault Integration**: Secure handling of secrets.
4. **mTLS**: Mutual TLS for client authentication.

## Scheduling Architecture

Nomad uses a modular scheduler with three types:

1. **Service Scheduler**: For long-lived services.
2. **Batch Scheduler**: For batch jobs.
3. **System Scheduler**: For system jobs that run on every client.

The scheduling process involves:

1. **Feasibility Checking**: Filtering out ineligible clients.
2. **Ranking**: Scoring eligible clients based on a scoring algorithm.
3. **Selection**: Selecting the highest-scoring client for each allocation.

## Deployment Architecture

Nomad supports several deployment strategies:

1. **Rolling Updates**: Incrementally updating tasks.
2. **Blue/Green Deployments**: Deploying a new version alongside the old one.
3. **Canary Deployments**: Testing a new version with a small subset of traffic.

## Scalability

Nomad is designed to be highly scalable. A single region can support thousands of clients and millions of allocations. Federation allows for global deployments across multiple regions.

## Integration Architecture

Nomad integrates with several other HashiCorp tools:

1. **Consul**: For service discovery and service mesh.
2. **Vault**: For secrets management.
3. **Terraform**: For infrastructure as code.

## Best Practices

1. **Server Sizing**: Recommended 3-5 servers per region for high availability.
2. **Client Sizing**: Size clients based on the workload requirements.
3. **Network Topology**: Place servers in the same datacenter for low latency.
4. **Federation**: Use federation for multi-region deployments.
5. **Monitoring**: Set up monitoring and alerting for the Nomad cluster.
6. **Backup**: Regularly backup the Nomad state.
7. **Upgrade Strategy**: Plan and test upgrades before applying them to production.
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://architecture",
				MIMEType: "text/markdown",
				Text:     architectureText,
			},
		}, nil
	})
}
