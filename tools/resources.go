package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterResources registers all resources with the MCP server
func RegisterResources(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Register static resources
	registerStaticResources(s, logger)

	// Register dynamic resources
	registerDynamicResources(s, nomadClient, logger)
}

// registerStaticResources registers static resources
func registerStaticResources(s *server.MCPServer, logger *log.Logger) {
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
}

// registerDynamicResources registers dynamic resources
func registerDynamicResources(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
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
}

// extractIDFromURI extracts an ID from a URI using the given prefix and suffix
func extractIDFromURI(uri, prefix, suffix string) string {
	// Find the start of the ID
	start := len(prefix)
	if len(uri) <= start {
		return ""
	}

	// Find the end of the ID
	end := len(uri) - len(suffix)
	if end <= start {
		return ""
	}

	return uri[start:end]
}
