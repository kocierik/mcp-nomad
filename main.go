// Package mcp-nomad provides a Model Context Protocol (MCP) server for interacting with HashiCorp Nomad.
// It implements stdio, SSE, and StreamableHTTP transports, allowing for easy integration with various clients.
//
// Features:
// - Job management (list, get, run, stop)
// - Deployment management
// - Namespace management
// - Node management
// - Allocation management
// - Variable management
// - Volume management
// - ACL management
//
// Example usage:
//
//	go install github.com/kocierik/mcp-nomad@latest
//
//	// Start the server with stdio transport [default]
//	go run .
//
//	// Legacy SSE transport (HTTP transports require MCP_TOKEN)
//	go run . -transport=sse -port=8080
//
//	// StreamableHTTP (Inspector / HTTP clients — path /mcp by default)
//	go run . -transport=streamable-http -port=8080
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kocierik/mcp-nomad/prompts"
	"github.com/kocierik/mcp-nomad/tools"
	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/server"
)

var (
	version   = "0.3.2"
	buildTime = "unknown"
	goVersion = "unknown"
)

type toolRegistrationDeps interface {
	utils.JobAPI
	utils.DeploymentAPI
	utils.NamespaceAPI
	utils.NodeAPI
	utils.AllocationAPI
	utils.VariableAPI
	utils.VolumeAPI
	utils.ACLToolsDeps
	utils.LogAPI
	utils.ClusterToolsAPI
	utils.SentinelAPI
	utils.DynamicResourcesNomad
}

func main() {
	// Define flags
	transport := flag.String("transport", "stdio", "Transport type (stdio, sse, or streamable-http)")
	port := flag.String("port", "8080", "Port for HTTP server")
	listenAddress := flag.String("listen-address", "127.0.0.1", "Address for HTTP server listener")
	readOnly := flag.Bool("read-only", false, "Register only read tools")
	// nomadAddr := flag.String("nomad-addr", "http://localhost:4646", "Nomad server address")
	flag.Parse()

	nomadAddr := os.Getenv("NOMAD_ADDR")
	if nomadAddr == "" {
		nomadAddr = "http://127.0.0.1:4646"
	}

	// Get token from environment
	token := os.Getenv("NOMAD_TOKEN")
	mcpToken := os.Getenv("MCP_TOKEN")
	mcpAllowedOrigins := os.Getenv("MCP_ALLOWED_ORIGINS")

	// Set up logging
	logger := log.New(os.Stderr, "[NomadMCP] ", log.LstdFlags)

	// Create MCP server
	s := server.NewMCPServer(
		"Nomad MCP",
		version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	var securityConfig httpSecurityConfig
	if *transport == "sse" || *transport == "streamable-http" {
		if mcpToken == "" {
			logger.Fatalf("MCP_TOKEN is required for HTTP transports")
		}
		allowedOrigins, err := buildAllowedOrigins(mcpAllowedOrigins, *port)
		if err != nil {
			logger.Fatalf("Invalid MCP_ALLOWED_ORIGINS: %v", err)
		}
		securityConfig = httpSecurityConfig{
			tokenHash:      sha256Hex(mcpToken),
			allowedOrigins: allowedOrigins,
		}
	}

	// Initialize Nomad client with token
	nomadClient, err := utils.NewNomadClient(nomadAddr, token)
	if err != nil {
		logger.Fatalf("Failed to create Nomad client: %v", err)
	}

	registerToolsWithPolicy(s, nomadClient, logger, *readOnly)

	// Register all prompts
	prompts.RegisterPrompts(s)

	// Start the MCP server based on transport type
	logger.Println("Starting Nomad MCP server...")

	switch *transport {
	case "stdio":
		logger.Println("Server started on stdio")
		if err := server.ServeStdio(s); err != nil {
			logger.Fatalf("Server error: %v", err)
		}
	case "sse":
		baseURL := fmt.Sprintf("http://%s", net.JoinHostPort(advertisedListenHost(*listenAddress), *port))

		// Create SSE server
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL(baseURL),
		)

		httpServer := &http.Server{
			Addr:              net.JoinHostPort(*listenAddress, *port),
			Handler:           newHTTPAuthOriginMiddleware(sseServer, securityConfig),
			ReadHeaderTimeout: 30 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      0,
			IdleTimeout:       60 * time.Second,
		}

		logger.Printf("SSE server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatalf("Server error: %v", err)
		}
	case "streamable-http":
		// Create StreamableHTTP server
		streamableServer := server.NewStreamableHTTPServer(s)

		httpServer := &http.Server{
			Addr:              net.JoinHostPort(*listenAddress, *port),
			Handler:           newHTTPAuthOriginMiddleware(streamableServer, securityConfig),
			ReadHeaderTimeout: 30 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      0,
			IdleTimeout:       60 * time.Second,
		}

		logger.Printf("StreamableHTTP server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatalf("Server error: %v", err)
		}
	default:
		logger.Fatalf("Invalid transport type: %s. Must be 'stdio', 'sse', or 'streamable-http'", *transport)
	}
}

// Register all tools with the MCP server
func registerToolsWithPolicy(s *server.MCPServer, nomadClient toolRegistrationDeps, logger *log.Logger, readOnly bool) {
	registerTools(s, nomadClient, logger)
	applyToolPolicy(s, readOnly)
}

func registerTools(s *server.MCPServer, nomadClient toolRegistrationDeps, logger *log.Logger) {
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

	// Register log tools
	tools.RegisterLogTools(s, nomadClient, logger)

	// Register resources
	tools.RegisterResourcesWithOptions(s, nomadClient, logger, tools.ResourceOptions{
		ReadmeContent:  readmeAsset(),
		LicenseContent: licenseAsset(),
		ServerVersion:  version,
		BuildTime:      buildTime,
		GoVersion:      goVersion,
		StartTime:      time.Now().UTC(),
	})

	// Register cluster tools
	tools.RegisterClusterTools(s, nomadClient, logger)

	// Register Sentinel tools
	tools.RegisterSentinelTools(s, nomadClient, logger)
}
