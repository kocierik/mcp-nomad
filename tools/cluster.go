package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterClusterTools registers all cluster-related tools
func RegisterClusterTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Get cluster leader tool
	getClusterLeaderTool := mcp.NewTool("get_cluster_leader",
		mcp.WithDescription("Get the current leader of the Nomad cluster"),
	)
	s.AddTool(getClusterLeaderTool, GetClusterLeaderHandler(nomadClient, logger))

	// List cluster peers tool
	listClusterPeersTool := mcp.NewTool("list_cluster_peers",
		mcp.WithDescription("List the peers in the Nomad cluster"),
	)
	s.AddTool(listClusterPeersTool, ListClusterPeersHandler(nomadClient, logger))

	// List regions tool
	listRegionsTool := mcp.NewTool("list_regions",
		mcp.WithDescription("List all available regions in the Nomad cluster"),
	)
	s.AddTool(listRegionsTool, ListRegionsHandler(nomadClient, logger))
}

// GetClusterLeaderHandler returns a handler for getting the cluster leader
func GetClusterLeaderHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		body, err := client.GetClusterLeader()

		if err != nil {
			logger.Printf("Error getting cluster configuration: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get cluster configuration", err), nil
		}

		// Parse the response to find the leader
		var config map[string]interface{}
		if err := json.Unmarshal(body, &config); err != nil {
			logger.Printf("Error parsing cluster configuration: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to parse cluster configuration", err), nil
		}

		servers, ok := config["Servers"].([]interface{})
		if !ok {
			return mcp.NewToolResultError("Could not find servers in configuration"), nil
		}

		leaderAddr := ""
		for _, srv := range servers {
			serverMap, ok := srv.(map[string]interface{})
			if !ok {
				continue
			}
			isLeader, ok := serverMap["Leader"].(bool)
			if ok && isLeader {
				leaderAddr, _ = serverMap["Address"].(string)
				break
			}
		}

		if leaderAddr == "" {
			return mcp.NewToolResultError("Could not determine cluster leader"), nil
		}

		return mcp.NewToolResultText(leaderAddr), nil
	}
}

// ListClusterPeersHandler returns a handler for listing cluster peers
func ListClusterPeersHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		body, err := client.ListClusterPeers()

		if err != nil {
			logger.Printf("Error getting cluster configuration: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get cluster configuration", err), nil
		}

		// Parse the response to find peers
		var config map[string]interface{}
		if err := json.Unmarshal(body, &config); err != nil {
			logger.Printf("Error parsing cluster configuration: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to parse cluster configuration", err), nil
		}

		servers, ok := config["Servers"].([]interface{})
		if !ok {
			return mcp.NewToolResultError("Could not find servers in configuration"), nil
		}

		var peers []string
		for _, srv := range servers {
			serverMap, ok := srv.(map[string]interface{})
			if !ok {
				continue
			}
			if addr, ok := serverMap["Address"].(string); ok {
				peers = append(peers, addr)
			}
		}

		peersJSON, err := json.MarshalIndent(peers, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format peer list", err), nil
		}

		return mcp.NewToolResultText(string(peersJSON)), nil
	}
}

// ListRegionsHandler returns a handler for listing regions
func ListRegionsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		body, err := client.MakeRequest("GET", "regions", nil, nil)
		if err != nil {
			logger.Printf("Error listing regions: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list regions", err), nil
		}

		// The response body is expected to be a JSON array of strings
		return mcp.NewToolResultText(string(body)), nil
	}
}
