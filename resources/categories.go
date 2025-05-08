// Package resources provides implementations of MCP resources for Nomad
package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerJobResources registers all job-related resources
func registerJobResources(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Job specification resource
	jobSpecTemplate := mcp.NewResourceTemplate(
		"nomad://jobs/{job_id}/spec",
		"Job Specification",
		mcp.WithTemplateDescription("Returns the specification for a specific job"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(jobSpecTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		jobID := extractIDFromURI(request.Params.URI, "jobs/", "/spec")
		if jobID == "" {
			return nil, fmt.Errorf("invalid job ID in URI")
		}

		job, err := nomadClient.GetJob(jobID, "default")
		if err != nil {
			logger.Printf("Error getting job spec: %v", err)
			return nil, err
		}

		jobJSON, err := json.MarshalIndent(job, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(jobJSON),
			},
		}, nil
	})

	// Job history resource
	jobHistoryTemplate := mcp.NewResourceTemplate(
		"nomad://jobs/{job_id}/history",
		"Job History",
		mcp.WithTemplateDescription("Returns the history of a specific job"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(jobHistoryTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		jobID := extractIDFromURI(request.Params.URI, "jobs/", "/history")
		if jobID == "" {
			return nil, fmt.Errorf("invalid job ID in URI")
		}

		// Get job versions
		versions, err := nomadClient.GetJobVersions(jobID, "default")
		if err != nil {
			logger.Printf("Error getting job versions: %v", err)
			return nil, err
		}

		versionsJSON, err := json.MarshalIndent(versions, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(versionsJSON),
			},
		}, nil
	})

	// Job allocations resource
	jobAllocationsTemplate := mcp.NewResourceTemplate(
		"nomad://jobs/{job_id}/allocations",
		"Job Allocations",
		mcp.WithTemplateDescription("Returns the allocations for a specific job"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(jobAllocationsTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		jobID := extractIDFromURI(request.Params.URI, "jobs/", "/allocations")
		if jobID == "" {
			return nil, fmt.Errorf("invalid job ID in URI")
		}

		allocations, err := nomadClient.ListJobAllocations(jobID, "default")
		if err != nil {
			logger.Printf("Error getting job allocations: %v", err)
			return nil, err
		}

		allocJSON, err := json.MarshalIndent(allocations, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(allocJSON),
			},
		}, nil
	})

	// Job evaluations resource
	jobEvaluationsTemplate := mcp.NewResourceTemplate(
		"nomad://jobs/{job_id}/evaluations",
		"Job Evaluations",
		mcp.WithTemplateDescription("Returns the evaluations for a specific job"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(jobEvaluationsTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		jobID := extractIDFromURI(request.Params.URI, "jobs/", "/evaluations")
		if jobID == "" {
			return nil, fmt.Errorf("invalid job ID in URI")
		}

		evaluations, err := nomadClient.ListJobEvaluations(jobID, "default")
		if err != nil {
			logger.Printf("Error getting job evaluations: %v", err)
			return nil, err
		}

		evalsJSON, err := json.MarshalIndent(evaluations, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(evalsJSON),
			},
		}, nil
	})
}

// registerNodeResources registers all node-related resources
func registerNodeResources(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Node status resource
	nodeStatusTemplate := mcp.NewResourceTemplate(
		"nomad://nodes/{node_id}/status",
		"Node Status",
		mcp.WithTemplateDescription("Returns the status information for a specific node"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(nodeStatusTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		nodeID := extractIDFromURI(request.Params.URI, "nodes/", "/status")
		if nodeID == "" {
			return nil, fmt.Errorf("invalid node ID in URI")
		}

		node, err := nomadClient.GetNode(nodeID)
		if err != nil {
			logger.Printf("Error getting node status: %v", err)
			return nil, err
		}

		nodeJSON, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(nodeJSON),
			},
		}, nil
	})

	// Node resources resource
	nodeResourcesTemplate := mcp.NewResourceTemplate(
		"nomad://nodes/{node_id}/resources",
		"Node Resources",
		mcp.WithTemplateDescription("Returns the resource information for a specific node"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(nodeResourcesTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		nodeID := extractIDFromURI(request.Params.URI, "nodes/", "/resources")
		if nodeID == "" {
			return nil, fmt.Errorf("invalid node ID in URI")
		}

		node, err := nomadClient.GetNode(nodeID)
		if err != nil {
			logger.Printf("Error getting node resources: %v", err)
			return nil, err
		}

		// Extract resource information
		resources := map[string]interface{}{
			"cpu": map[string]interface{}{
				"total":    node.Resources.CPU,
				"reserved": node.Reserved.CPU,
			},
			"memory": map[string]interface{}{
				"total":    node.Resources.MemoryMB,
				"reserved": node.Reserved.MemoryMB,
			},
			"disk": map[string]interface{}{
				"total":    node.Resources.DiskMB,
				"reserved": node.Reserved.DiskMB,
			},
		}

		resourcesJSON, err := json.MarshalIndent(resources, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(resourcesJSON),
			},
		}, nil
	})

	// Node allocations resource
	nodeAllocationsTemplate := mcp.NewResourceTemplate(
		"nomad://nodes/{node_id}/allocations",
		"Node Allocations",
		mcp.WithTemplateDescription("Returns allocations running on a specific node"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(nodeAllocationsTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		nodeID := extractIDFromURI(request.Params.URI, "nodes/", "/allocations")
		if nodeID == "" {
			return nil, fmt.Errorf("invalid node ID in URI")
		}

		// Use MakeRequest for direct API call since there isn't a specific client method
		path := fmt.Sprintf("node/%s/allocations", nodeID)
		allocsData, err := nomadClient.MakeRequest("GET", path, nil, nil)
		if err != nil {
			logger.Printf("Error getting node allocations: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(allocsData),
			},
		}, nil
	})
}

// registerAllocationResources registers all allocation-related resources
func registerAllocationResources(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Allocation logs resource
	allocationLogsTemplate := mcp.NewResourceTemplate(
		"nomad://allocations/{alloc_id}/logs",
		"Allocation Logs",
		mcp.WithTemplateDescription("Returns the logs for a specific allocation"),
		mcp.WithTemplateMIMEType("text/plain"),
	)

	s.AddResourceTemplate(allocationLogsTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		allocID := extractIDFromURI(request.Params.URI, "allocations/", "/logs")
		if allocID == "" {
			return nil, fmt.Errorf("invalid allocation ID in URI")
		}

		allocLogs, err := nomadClient.GetAllocationLogs(allocID)
		if err != nil {
			logger.Printf("Error getting allocation logs: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     allocLogs,
			},
		}, nil
	})

	// Allocation status resource
	allocationStatusTemplate := mcp.NewResourceTemplate(
		"nomad://allocations/{alloc_id}/status",
		"Allocation Status",
		mcp.WithTemplateDescription("Returns the status information for a specific allocation"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(allocationStatusTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		allocID := extractIDFromURI(request.Params.URI, "allocations/", "/status")
		if allocID == "" {
			return nil, fmt.Errorf("invalid allocation ID in URI")
		}

		alloc, err := nomadClient.GetAllocation(allocID)
		if err != nil {
			logger.Printf("Error getting allocation status: %v", err)
			return nil, err
		}

		allocJSON, err := json.MarshalIndent(alloc, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(allocJSON),
			},
		}, nil
	})

	// Allocation tasks resource
	allocationTasksTemplate := mcp.NewResourceTemplate(
		"nomad://allocations/{alloc_id}/tasks",
		"Allocation Tasks",
		mcp.WithTemplateDescription("Returns the tasks for a specific allocation"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(allocationTasksTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		allocID := extractIDFromURI(request.Params.URI, "allocations/", "/tasks")
		if allocID == "" {
			return nil, fmt.Errorf("invalid allocation ID in URI")
		}

		alloc, err := nomadClient.GetAllocation(allocID)
		if err != nil {
			logger.Printf("Error getting allocation tasks: %v", err)
			return nil, err
		}

		// Extract only the tasks information
		tasksJSON, err := json.MarshalIndent(alloc.TaskStates, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(tasksJSON),
			},
		}, nil
	})
}

// registerClusterResources registers all cluster-level resources
func registerClusterResources(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Cluster metrics resource
	clusterMetricsResource := mcp.NewResource(
		"nomad://cluster/metrics",
		"Cluster Metrics",
		mcp.WithResourceDescription("Returns metrics for the entire Nomad cluster"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(clusterMetricsResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		metricsData, err := nomadClient.ListClusterPeers()
		if err != nil {
			logger.Printf("Error getting cluster metrics: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "nomad://cluster/metrics",
				MIMEType: "application/json",
				Text:     string(metricsData),
			},
		}, nil
	})

	// Cluster policies resource
	clusterPoliciesResource := mcp.NewResource(
		"nomad://policies/list",
		"Cluster Policies",
		mcp.WithResourceDescription("Returns a list of all policies in the cluster"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(clusterPoliciesResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		policies, err := nomadClient.ListACLPolicies()
		if err != nil {
			logger.Printf("Error getting cluster policies: %v", err)
			return nil, err
		}

		policiesJSON, err := json.MarshalIndent(policies, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "nomad://policies/list",
				MIMEType: "application/json",
				Text:     string(policiesJSON),
			},
		}, nil
	})

	// Cluster leader resource
	clusterLeaderResource := mcp.NewResource(
		"nomad://cluster/leader",
		"Cluster Leader",
		mcp.WithResourceDescription("Returns information about the cluster leader"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(clusterLeaderResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		leaderData, err := nomadClient.GetClusterLeader()
		if err != nil {
			logger.Printf("Error getting cluster leader: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "nomad://cluster/leader",
				MIMEType: "application/json",
				Text:     string(leaderData),
			},
		}, nil
	})
}

// registerMiscResources registers miscellaneous resources
func registerMiscResources(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Evaluation resource
	evaluationTemplate := mcp.NewResourceTemplate(
		"nomad://evaluations/{eval_id}",
		"Evaluation Details",
		mcp.WithTemplateDescription("Returns details about a specific evaluation"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(evaluationTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		evalID := extractIDFromURI(request.Params.URI, "evaluations/", "")
		if evalID == "" {
			return nil, fmt.Errorf("invalid evaluation ID in URI")
		}

		path := fmt.Sprintf("evaluation/%s", evalID)
		evalData, err := nomadClient.MakeRequest("GET", path, nil, nil)
		if err != nil {
			logger.Printf("Error getting evaluation: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(evalData),
			},
		}, nil
	})

	// Service health resource
	serviceHealthTemplate := mcp.NewResourceTemplate(
		"nomad://services/{service_name}/health",
		"Service Health Status",
		mcp.WithTemplateDescription("Returns health information for a specific service"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(serviceHealthTemplate, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		serviceName := extractIDFromURI(request.Params.URI, "services/", "/health")
		if serviceName == "" {
			return nil, fmt.Errorf("invalid service name in URI")
		}

		path := fmt.Sprintf("service/%s", serviceName)
		serviceData, err := nomadClient.MakeRequest("GET", path, nil, nil)
		if err != nil {
			logger.Printf("Error getting service health: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(serviceData),
			},
		}, nil
	})
}

// RegisterResourcesByCategory organizes and registers resources by category
func RegisterResourcesByCategory(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Register resources by category
	registerJobResources(s, nomadClient, logger)
	registerNodeResources(s, nomadClient, logger)
	registerAllocationResources(s, nomadClient, logger)
	registerClusterResources(s, nomadClient, logger)
	registerMiscResources(s, nomadClient, logger)
}
