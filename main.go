package main

import (
	"log"
	"os"

	"github.com/kocierik/nomad-mcp-server/prompts"
	"github.com/kocierik/nomad-mcp-server/tools"
	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Set up logging
	logger := log.New(os.Stderr, "[NomadMCP] ", log.LstdFlags)

	// Create MCP server
	s := server.NewMCPServer(
		"Nomad MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Initialize Nomad client
	nomadClient, err := utils.NewNomadClient()
	if err != nil {
		logger.Fatalf("Failed to create Nomad client: %v", err)
	}

	// Register all tools
	registerTools(s, nomadClient, logger)

	// Register all prompts
	prompts.RegisterPrompts(s)

	// Start the MCP server using stdio
	logger.Println("Starting Nomad MCP server...")
	if err := server.ServeStdio(s); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}

// Register all tools with the MCP server
func registerTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Register job-related tools
	tools.RegisterJobTools(s, nomadClient, logger)

	// Register deployment tools
	tools.RegisterDeploymentTools(s, nomadClient, logger)

	// Register namespace tools
	tools.RegisterNamespaceTools(s, nomadClient, logger)

	// Register node tools
	tools.RegisterNodeTools(s, nomadClient, logger)

	// Register allocation tools
	tools.RegisterAllocationTools(s, nomadClient, logger)

	// Register variable tools
	tools.RegisterVariableTools(s, nomadClient, logger)

	// Register volume tools
	tools.RegisterVolumeTools(s, nomadClient, logger)

	// Register ACL tools
	tools.RegisterACLTools(s, nomadClient, logger)
}
