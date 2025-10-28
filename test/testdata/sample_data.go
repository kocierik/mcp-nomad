package testdata

import (
	"time"

	"github.com/kocierik/mcp-nomad/types"
)

// SampleJobSpecs contains various job specifications for testing
var SampleJobSpecs = map[string]string{
	"simple": `job "test-job" {
  datacenters = ["dc1"]
  type = "service"

  group "web" {
    count = 2

    task "nginx" {
      driver = "docker"

      config {
        image = "nginx:latest"
        ports = ["http"]
      }

      resources {
        cpu    = 100
        memory = 128
      }

      service {
        name = "nginx"
        port = "http"

        check {
          type     = "http"
          path     = "/"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}`,
	"json": `{
  "Job": {
    "ID": "test-json-job",
    "Name": "test-json-job",
    "Type": "service",
    "Datacenters": ["dc1"],
    "TaskGroups": [
      {
        "Name": "web",
        "Count": 1,
        "Tasks": [
          {
            "Name": "nginx",
            "Driver": "docker",
            "Config": {
              "image": "nginx:latest"
            },
            "Resources": {
              "CPU": 100,
              "MemoryMB": 128
            }
          }
        ]
      }
    ]
  }
}`,
	"invalid": `invalid hcl syntax here`,
}

// SampleJobs contains sample job data for testing
var SampleJobs = []types.JobSummary{
	{
		ID:          "test-job-1",
		Summary:     map[string]types.TaskSummary{"web": {Running: 2, Complete: 0, Failed: 0}},
		CreateIndex: 1,
		ModifyIndex: 1,
	},
	{
		ID:          "test-job-2",
		Summary:     map[string]types.TaskSummary{"batch": {Complete: 1, Failed: 0, Running: 0}},
		CreateIndex: 2,
		ModifyIndex: 2,
	},
}

// SampleNodes contains sample node data for testing
var SampleNodes = []types.NodeSummary{
	{
		ID:         "node-1",
		Name:       "node-1",
		Datacenter: "dc1",
		NodeClass:  "default",
		Status:     "ready",
	},
	{
		ID:         "node-2",
		Name:       "node-2",
		Datacenter: "dc1",
		NodeClass:  "gpu",
		Status:     "down",
	},
}

// SampleNamespaces contains sample namespace data for testing
var SampleNamespaces = []types.Namespace{
	{
		Name:        "default",
		Description: "Default namespace",
	},
	{
		Name:        "production",
		Description: "Production namespace",
	},
}

// SampleAllocations contains sample allocation data for testing
var SampleAllocations = []types.Allocation{
	{
		ID:                 "alloc-1",
		EvalID:             "eval-1",
		Name:               "test-job-1.web[0]",
		NodeID:             "node-1",
		JobID:              "test-job-1",
		TaskGroup:          "web",
		DesiredStatus:      "run",
		DesiredDescription: "Allocation is running",
		ClientStatus:       "running",
		ClientDescription:  "Allocation is running",
		CreateIndex:        1,
		ModifyIndex:        1,
		CreateTime:         time.Now().Add(-1 * time.Hour).Unix(),
		ModifyTime:         time.Now().Add(-30 * time.Minute).Unix(),
	},
}

// SampleVariables contains sample variable data for testing
var SampleVariables = []types.Variable{
	{
		Path:      "app/config",
		Namespace: "default",
		Value:     `{"Items":{"database_url":"postgres://localhost:5432/test","api_key":"secret123"}}`,
	},
}

// SampleACLTokens contains sample ACL token data for testing
var SampleACLTokens = []types.ACLToken{
	{
		AccessorID:  "token-1",
		SecretID:    "secret-1",
		Name:        "test-token",
		Type:        "client",
		Policies:    []string{"read-only"},
		Global:      false,
		CreateIndex: 1,
		ModifyIndex: 1,
	},
}

// SampleLogs contains sample log data for testing
var SampleLogs = map[string]string{
	"nginx_stdout": `2024-01-01T10:00:00Z [INFO] Starting nginx
2024-01-01T10:00:01Z [INFO] Configuration loaded
2024-01-01T10:00:02Z [INFO] Server started on port 80
2024-01-01T10:00:03Z [INFO] Ready to serve requests`,
	"nginx_stderr": `2024-01-01T10:00:00Z [WARN] Using default configuration
2024-01-01T10:00:01Z [ERROR] Failed to load SSL certificate`,
}

// SampleClusterData contains sample cluster information
var SampleClusterData = map[string][]byte{
	"leader": []byte(`{
  "Servers": [
    {
      "ID": "server-1",
      "Node": "node-1",
      "Address": "127.0.0.1:4647",
      "Leader": true,
      "Voter": true
    }
  ]
}`),
	"peers": []byte(`{
  "Servers": [
    {
      "ID": "server-1",
      "Node": "node-1",
      "Address": "127.0.0.1:4647",
      "Leader": true,
      "Voter": true
    },
    {
      "ID": "server-2",
      "Node": "node-2",
      "Address": "127.0.0.1:4648",
      "Leader": false,
      "Voter": true
    }
  ]
}`),
	"regions": []byte(`["global", "us-east-1", "us-west-2"]`),
}
