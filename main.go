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
//	go get github.com/kocierik/mcp-nomad@v0.1.0
//
//	import "github.com/kocierik/mcp-nomad"
//
//	// Start the server with stdio transport [default]
//	go run main.go -transport=stdio
//
//	// Start the server with SSE transport
//	go run main.go -transport=sse -port=8080
//
//	// Start the server with StreamableHTTP transport
//	go run main.go -transport=streamable-http -port=8080
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kocierik/mcp-nomad/prompts"
	"github.com/kocierik/mcp-nomad/tools"
	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/server"
)

// authKey is a custom context key for storing the auth token
type authKey struct{}

// withAuthKey adds an auth key to the context
func withAuthKey(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

// authFromRequest extracts the auth token from the request headers
func authFromRequest(ctx context.Context, r *http.Request) context.Context {
	// If no token is provided, return the context as is
	token := r.Header.Get("Authorization")
	if token == "" {
		return ctx
	}
	return withAuthKey(ctx, token)
}

// authFromEnv extracts the auth token from the environment
func authFromEnv(ctx context.Context) context.Context {
	// If no token is provided, return the context as is
	token := os.Getenv("NOMAD_TOKEN")
	if token == "" {
		return ctx
	}
	return withAuthKey(ctx, token)
}

// validateOrigin checks if the request origin is allowed
func validateOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true // Allow requests without Origin header (e.g., from curl)
	}

	// Allow localhost origins
	allowedOrigins := []string{
		"http://localhost",
		"http://127.0.0.1",
		os.Getenv("NOMAD_ADDR"),
	}

	for _, allowed := range allowedOrigins {
		if strings.HasPrefix(origin, allowed) {
			return true
		}
	}

	return false
}

// originValidationMiddleware validates the origin of incoming requests
func originValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !validateOrigin(r) {
			http.Error(w, "Invalid origin", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Define flags
	transport := flag.String("transport", "stdio", "Transport type (stdio, sse, or streamable-http)")
	port := flag.String("port", "8080", "Port for HTTP server")
	// nomadAddr := flag.String("nomad-addr", "http://localhost:4646", "Nomad server address")
	flag.Parse()

	nomadAddr := os.Getenv("NOMAD_ADDR")
	if nomadAddr == "" {
		nomadAddr = "http://0.0.0.0:4646"
	}

	// Get token from environment
	token := os.Getenv("NOMAD_TOKEN")

	// Set up logging
	logger := log.New(os.Stderr, "[NomadMCP] ", log.LstdFlags)

	// Create MCP server
	s := server.NewMCPServer(
		"Nomad MCP",
		"0.1.4",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Initialize Nomad client with token
	nomadClient, err := utils.NewNomadClient(nomadAddr, token)
	if err != nil {
		logger.Fatalf("Failed to create Nomad client: %v", err)
	}

	// Register all tools
	registerTools(s, nomadClient, logger)

	// Register all prompts
	prompts.RegisterPrompts(s)

	// Start the MCP server based on transport type
	logger.Println("Starting Nomad MCP server...")

	switch *transport {
	case "stdio":
		logger.Println("Server started on stdio")
		if err := server.ServeStdio(s, server.WithStdioContextFunc(authFromEnv)); err != nil {
			logger.Fatalf("Server error: %v", err)
		}
	case "sse":
		// Parse the Nomad address to get the host
		nomadURL, err := url.Parse(nomadAddr)
		if err != nil {
			logger.Fatalf("Invalid nomad-addr: %v", err)
		}
		logger.Printf("Nomad URL: %s", nomadURL.Hostname())

		// Create SSE server
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL(fmt.Sprintf("http://%s:%s", nomadURL.Hostname(), *port)),
			server.WithSSEContextFunc(authFromRequest),
		)

		// Create HTTP server with origin validation middleware
		httpServer := &http.Server{
			Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", *port),
			Handler: originValidationMiddleware(sseServer),
		}

		logger.Printf("SSE server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatalf("Server error: %v", err)
		}
	case "streamable-http":
		// Parse the Nomad address to get the host
		nomadURL, err := url.Parse(nomadAddr)
		if err != nil {
			logger.Fatalf("Invalid nomad-addr: %v", err)
		}
		logger.Printf("Nomad URL: %s", nomadURL.Hostname())

		// Create StreamableHTTP server
		streamableServer := server.NewStreamableHTTPServer(s,
			server.WithHTTPContextFunc(authFromRequest),
		)

		// Create HTTP server with origin validation middleware
		httpServer := &http.Server{
			Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", *port),
			Handler: originValidationMiddleware(streamableServer),
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

	// Register log tools
	tools.RegisterLogTools(s, nomadClient, logger)

	// Register resources
	tools.RegisterResources(s, nomadClient, logger)

	// Register cluster tools
	tools.RegisterClusterTools(s, nomadClient, logger)

	// Register Sentinel tools
	tools.RegisterSentinelTools(s, nomadClient, logger)
}
