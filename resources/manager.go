// Package resources provides implementations of MCP resources for Nomad
package resources

import (
	"log"
	"strings"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/server"
)

// ResourceManager manages all resources for the MCP server
type ResourceManager struct {
	server      *server.MCPServer
	nomadClient *utils.NomadClient
	logger      *log.Logger
}

// NewResourceManager creates a new resource manager
func NewResourceManager(server *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) *ResourceManager {
	return &ResourceManager{
		server:      server,
		nomadClient: nomadClient,
		logger:      logger,
	}
}

// RegisterAll registers all resources with the MCP server
func (rm *ResourceManager) RegisterAll() {
	// Register static resources
	RegisterStaticResources(rm.server, rm.logger)

	// Register resources by category
	RegisterResourcesByCategory(rm.server, rm.nomadClient, rm.logger)

	rm.logger.Printf("Registered all resources")
}

// GetURICategories returns a list of URI categories used in this implementation
func (rm *ResourceManager) GetURICategories() []string {
	return []string{
		"docs://",
		"system://",
		"nomad://jobs/",
		"nomad://nodes/",
		"nomad://allocations/",
		"nomad://cluster/",
		"nomad://evaluations/",
		"nomad://services/",
		"nomad://policies/",
	}
}

// GetResourceListByCategory returns a list of resources by category
func (rm *ResourceManager) GetResourceListByCategory() map[string][]string {
	return map[string][]string{
		"Documentation": {
			"docs://readme",
			"docs://license",
			"docs://help",
			"docs://api",
			"docs://nomad",
			"docs://job-spec",
			"docs://drivers",
			"docs://security",
			"docs://cli",
			"docs://architecture",
		},
		"System": {
			"system://info",
		},
		"Jobs": {
			"nomad://jobs/{job_id}/spec",
			"nomad://jobs/{job_id}/history",
			"nomad://jobs/{job_id}/allocations",
			"nomad://jobs/{job_id}/evaluations",
		},
		"Nodes": {
			"nomad://nodes/{node_id}/status",
			"nomad://nodes/{node_id}/resources",
			"nomad://nodes/{node_id}/allocations",
		},
		"Allocations": {
			"nomad://allocations/{alloc_id}/logs",
			"nomad://allocations/{alloc_id}/status",
			"nomad://allocations/{alloc_id}/tasks",
		},
		"Cluster": {
			"nomad://cluster/metrics",
			"nomad://cluster/leader",
			"nomad://policies/list",
		},
		"Miscellaneous": {
			"nomad://evaluations/{eval_id}",
			"nomad://services/{service_name}/health",
		},
	}
}

// IsValidCategory checks if a resource URI belongs to a valid category
func (rm *ResourceManager) IsValidCategory(uri string) bool {
	categories := rm.GetURICategories()
	for _, category := range categories {
		if strings.HasPrefix(uri, category) {
			return true
		}
	}
	return false
}

// IsValidURI checks if a resource URI is valid
func (rm *ResourceManager) IsValidURI(uri string) bool {
	if !rm.IsValidCategory(uri) {
		return false
	}

	// Additional validation for parameterized URIs
	if strings.Contains(uri, "{") || strings.Contains(uri, "}") {
		// Template URIs are not valid concrete URIs
		return false
	}

	// Check specific URI patterns
	parts := strings.Split(uri, "/")

	// Validate job URIs
	if strings.HasPrefix(uri, "nomad://jobs/") {
		if len(parts) < 4 {
			return false
		}
		// Check for valid job operation
		validOps := map[string]bool{
			"spec":        true,
			"history":     true,
			"allocations": true,
			"evaluations": true,
		}
		return validOps[parts[3]]
	}

	// Validate node URIs
	if strings.HasPrefix(uri, "nomad://nodes/") {
		if len(parts) < 4 {
			return false
		}
		// Check for valid node operation
		validOps := map[string]bool{
			"status":      true,
			"resources":   true,
			"allocations": true,
		}
		return validOps[parts[3]]
	}

	// Validate allocation URIs
	if strings.HasPrefix(uri, "nomad://allocations/") {
		if len(parts) < 4 {
			return false
		}
		// Check for valid allocation operation
		validOps := map[string]bool{
			"logs":   true,
			"status": true,
			"tasks":  true,
		}
		return validOps[parts[3]]
	}

	return true
}
