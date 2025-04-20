package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/nomad-mcp-server/types"
	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

// ListACLTokensHandler returns a handler for listing ACL tokens
func ListACLTokensHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tokens, err := client.ListACLTokens()
		if err != nil {
			logger.Printf("Error listing ACL tokens: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list ACL tokens", err), nil
		}

		tokensJSON, err := json.MarshalIndent(tokens, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format token list", err), nil
		}

		return mcp.NewToolResultText(string(tokensJSON)), nil
	}
}

// GetACLTokenHandler returns a handler for getting a specific ACL token
func GetACLTokenHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessorID, ok := request.Params.Arguments["accessor_id"].(string)
		if !ok || accessorID == "" {
			return mcp.NewToolResultError("accessor_id is required"), nil
		}

		token, err := client.GetACLToken(accessorID)
		if err != nil {
			logger.Printf("Error getting ACL token: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get ACL token", err), nil
		}

		tokenJSON, err := json.MarshalIndent(token, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format token details", err), nil
		}

		return mcp.NewToolResultText(string(tokenJSON)), nil
	}
}

// CreateACLTokenHandler returns a handler for creating a new ACL token
func CreateACLTokenHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		tokenType, ok := request.Params.Arguments["type"].(string)
		if !ok || tokenType == "" {
			return mcp.NewToolResultError("type is required"), nil
		}

		var policies []string
		if policiesParam, ok := request.Params.Arguments["policies"].([]interface{}); ok {
			for _, p := range policiesParam {
				if policy, ok := p.(string); ok {
					policies = append(policies, policy)
				}
			}
		}

		global := false
		if globalParam, ok := request.Params.Arguments["global"].(bool); ok {
			global = globalParam
		}

		token := types.ACLToken{
			Name:     name,
			Type:     tokenType,
			Policies: policies,
			Global:   global,
		}

		createdToken, err := client.CreateACLToken(token)
		if err != nil {
			logger.Printf("Error creating ACL token: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create ACL token", err), nil
		}

		tokenJSON, err := json.MarshalIndent(createdToken, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format token details", err), nil
		}

		return mcp.NewToolResultText(string(tokenJSON)), nil
	}
}

// DeleteACLTokenHandler returns a handler for deleting an ACL token
func DeleteACLTokenHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accessorID, ok := request.Params.Arguments["accessor_id"].(string)
		if !ok || accessorID == "" {
			return mcp.NewToolResultError("accessor_id is required"), nil
		}

		err := client.DeleteACLToken(accessorID)
		if err != nil {
			logger.Printf("Error deleting ACL token: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete ACL token", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("ACL token %s deleted successfully", accessorID)), nil
	}
}

// ListACLPoliciesHandler returns a handler for listing ACL policies
func ListACLPoliciesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policies, err := client.ListACLPolicies()
		if err != nil {
			logger.Printf("Error listing ACL policies: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list ACL policies", err), nil
		}

		policiesJSON, err := json.MarshalIndent(policies, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format policy list", err), nil
		}

		return mcp.NewToolResultText(string(policiesJSON)), nil
	}
}

// GetACLPolicyHandler returns a handler for getting a specific ACL policy
func GetACLPolicyHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		policy, err := client.GetACLPolicy(name)
		if err != nil {
			logger.Printf("Error getting ACL policy: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get ACL policy", err), nil
		}

		policyJSON, err := json.MarshalIndent(policy, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format policy details", err), nil
		}

		return mcp.NewToolResultText(string(policyJSON)), nil
	}
}

// CreateACLPolicyHandler returns a handler for creating a new ACL policy
func CreateACLPolicyHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		rules, ok := request.Params.Arguments["rules"].(string)
		if !ok || rules == "" {
			return mcp.NewToolResultError("rules is required"), nil
		}

		description := ""
		if desc, ok := request.Params.Arguments["description"].(string); ok {
			description = desc
		}

		policy := types.ACLPolicy{
			Name:        name,
			Description: description,
			Rules:       rules,
		}

		err := client.CreateACLPolicy(policy)
		if err != nil {
			logger.Printf("Error creating ACL policy: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create ACL policy", err), nil
		}

		policyJSON, err := json.MarshalIndent(policy, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format policy details", err), nil
		}

		return mcp.NewToolResultText(string(policyJSON)), nil
	}
}

// DeleteACLPolicyHandler returns a handler for deleting an ACL policy
func DeleteACLPolicyHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		err := client.DeleteACLPolicy(name)
		if err != nil {
			logger.Printf("Error deleting ACL policy: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete ACL policy", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("ACL policy %s deleted successfully", name)), nil
	}
}

// ListACLRolesHandler returns a handler for listing ACL roles
func ListACLRolesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		roles, err := client.ListACLRoles()
		if err != nil {
			logger.Printf("Error listing ACL roles: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list ACL roles", err), nil
		}

		rolesJSON, err := json.MarshalIndent(roles, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format role list", err), nil
		}

		return mcp.NewToolResultText(string(rolesJSON)), nil
	}
}

// GetACLRoleHandler returns a handler for getting a specific ACL role
func GetACLRoleHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, ok := request.Params.Arguments["id"].(string)
		if !ok || id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}

		role, err := client.GetACLRole(id)
		if err != nil {
			logger.Printf("Error getting ACL role: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get ACL role", err), nil
		}

		roleJSON, err := json.MarshalIndent(role, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format role details", err), nil
		}

		return mcp.NewToolResultText(string(roleJSON)), nil
	}
}

// CreateACLRoleHandler returns a handler for creating a new ACL role
func CreateACLRoleHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		description := ""
		if desc, ok := request.Params.Arguments["description"].(string); ok {
			description = desc
		}

		var policies []string
		if policiesParam, ok := request.Params.Arguments["policies"].([]interface{}); ok {
			for _, p := range policiesParam {
				if policy, ok := p.(string); ok {
					policies = append(policies, policy)
				}
			}
		}

		role := types.ACLRole{
			Name:        name,
			Description: description,
			Policies:    policies,
		}

		err := client.CreateACLRole(role)
		if err != nil {
			logger.Printf("Error creating ACL role: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create ACL role", err), nil
		}

		roleJSON, err := json.MarshalIndent(role, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format role details", err), nil
		}

		return mcp.NewToolResultText(string(roleJSON)), nil
	}
}

// DeleteACLRoleHandler returns a handler for deleting an ACL role
func DeleteACLRoleHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, ok := request.Params.Arguments["id"].(string)
		if !ok || id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}

		err := client.DeleteACLRole(id)
		if err != nil {
			logger.Printf("Error deleting ACL role: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete ACL role", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("ACL role %s deleted successfully", id)), nil
	}
}

// BootstrapACLTokenHandler returns a handler for bootstrapping the ACL system
func BootstrapACLTokenHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		token, err := client.BootstrapACLToken()
		if err != nil {
			logger.Printf("Error bootstrapping ACL token: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to bootstrap ACL token", err), nil
		}

		// Save the token in the client
		client.SetToken(token.SecretID)

		tokenJSON, err := json.MarshalIndent(token, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format token details", err), nil
		}

		return mcp.NewToolResultText(string(tokenJSON)), nil
	}
}
