package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterVolumeTools registers all volume-related tools
func RegisterVolumeTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List volumes tool
	listVolumesTool := mcp.NewTool("list_volumes",
		mcp.WithDescription("List all volumes in a namespace"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to list volumes from (optional)"),
		),
	)
	s.AddTool(listVolumesTool, ListVolumesHandler(nomadClient, logger))

	// Get volume tool
	getVolumeTool := mcp.NewTool("get_volume",
		mcp.WithDescription("Get details of a specific volume"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to get"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace of the volume (optional)"),
		),
	)
	s.AddTool(getVolumeTool, GetVolumeHandler(nomadClient, logger))

	// Delete volume tool
	deleteVolumeTool := mcp.NewTool("delete_volume",
		mcp.WithDescription("Delete a volume"),
		mcp.WithString("volume_id",
			mcp.Required(),
			mcp.Description("ID of the volume to delete"),
		),
		mcp.WithString("namespace",
			mcp.Description("Namespace of the volume (optional)"),
		),
	)
	s.AddTool(deleteVolumeTool, DeleteVolumeHandler(nomadClient, logger))
}

// ListVolumesHandler returns a handler for listing volumes
func ListVolumesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get optional parameters
		nodeID, _ := request.Params.Arguments["node_id"].(string)
		pluginID, _ := request.Params.Arguments["plugin_id"].(string)
		nextToken, _ := request.Params.Arguments["next_token"].(string)
		perPage, _ := request.Params.Arguments["per_page"].(int)
		filter, _ := request.Params.Arguments["filter"].(string)

		// Validate node_id and plugin_id if provided
		if nodeID != "" && len(nodeID)%2 != 0 {
			return mcp.NewToolResultError("node_id must have an even number of hexadecimal characters"), nil
		}
		if pluginID != "" && len(pluginID)%2 != 0 {
			return mcp.NewToolResultError("plugin_id must have an even number of hexadecimal characters"), nil
		}

		// List volumes with the specified parameters
		volumes, err := client.ListVolumes(nodeID, pluginID, nextToken, perPage, filter)
		if err != nil {
			logger.Printf("Error listing volumes: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list volumes", err), nil
		}

		// Format the response
		volumesJSON, err := json.MarshalIndent(volumes, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format volume list", err), nil
		}

		return mcp.NewToolResultText(string(volumesJSON)), nil
	}
}

// GetVolumeHandler returns a handler for getting volume details
func GetVolumeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok || volumeID == "" {
			return mcp.NewToolResultError("volume_id is required"), nil
		}

		volume, err := client.GetVolume(volumeID)

		if err != nil {
			logger.Printf("Error getting volume: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get volume", err), nil
		}

		volumeJSON, err := json.MarshalIndent(volume, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format volume details", err), nil
		}

		return mcp.NewToolResultText(string(volumeJSON)), nil
	}
}

// DeleteVolumeHandler returns a handler for deleting a volume
func DeleteVolumeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required parameters
		volumeID, ok := request.Params.Arguments["volume_id"].(string)
		if !ok || volumeID == "" {
			return mcp.NewToolResultError("volume_id is required"), nil
		}

		err := client.DeleteVolume(volumeID)

		if err != nil {
			logger.Printf("Error deleting volume: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete volume", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Volume %s deleted successfully", volumeID)), nil
	}
}
