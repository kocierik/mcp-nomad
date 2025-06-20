package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterACLTools registers all ACL-related tools
func RegisterACLTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// ACL Token tools
	listACLTokensTool := mcp.NewTool("list_acl_tokens",
		mcp.WithDescription("List all ACL tokens"),
	)
	s.AddTool(listACLTokensTool, ListACLTokensHandler(nomadClient, logger))

	getACLTokenTool := mcp.NewTool("get_acl_token",
		mcp.WithDescription("Get details of a specific ACL token"),
		mcp.WithString("accessor_id",
			mcp.Required(),
			mcp.Description("Accessor ID of the token to get"),
		),
	)
	s.AddTool(getACLTokenTool, GetACLTokenHandler(nomadClient, logger))

	createACLTokenTool := mcp.NewTool("create_acl_token",
		mcp.WithDescription("Create a new ACL token"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the token"),
		),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of the token (client or management)"),
			mcp.Enum("client", "management"),
		),
		mcp.WithArray("policies",
			mcp.Description("List of policy names to associate with the token"),
		),
		mcp.WithBoolean("global",
			mcp.Description("Whether the token is global (default: false)"),
		),
	)
	s.AddTool(createACLTokenTool, CreateACLTokenHandler(nomadClient, logger))

	deleteACLTokenTool := mcp.NewTool("delete_acl_token",
		mcp.WithDescription("Delete an ACL token"),
		mcp.WithString("accessor_id",
			mcp.Required(),
			mcp.Description("Accessor ID of the token to delete"),
		),
	)
	s.AddTool(deleteACLTokenTool, DeleteACLTokenHandler(nomadClient, logger))

	// ACL Policy tools
	listACLPoliciesTool := mcp.NewTool("list_acl_policies",
		mcp.WithDescription("List all ACL policies"),
	)
	s.AddTool(listACLPoliciesTool, ListACLPoliciesHandler(nomadClient, logger))

	getACLPolicyTool := mcp.NewTool("get_acl_policy",
		mcp.WithDescription("Get details of a specific ACL policy"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the policy to get"),
		),
	)
	s.AddTool(getACLPolicyTool, GetACLPolicyHandler(nomadClient, logger))

	createACLPolicyTool := mcp.NewTool("create_acl_policy",
		mcp.WithDescription("Create a new ACL policy"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the policy"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the policy"),
		),
		mcp.WithString("rules",
			mcp.Required(),
			mcp.Description("JSON rules for the policy"),
		),
	)
	s.AddTool(createACLPolicyTool, CreateACLPolicyHandler(nomadClient, logger))

	deleteACLPolicyTool := mcp.NewTool("delete_acl_policy",
		mcp.WithDescription("Delete an ACL policy"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the policy to delete"),
		),
	)
	s.AddTool(deleteACLPolicyTool, DeleteACLPolicyHandler(nomadClient, logger))

	// ACL Role tools
	listACLRolesTool := mcp.NewTool("list_acl_roles",
		mcp.WithDescription("List all ACL roles"),
	)
	s.AddTool(listACLRolesTool, ListACLRolesHandler(nomadClient, logger))

	getACLRoleTool := mcp.NewTool("get_acl_role",
		mcp.WithDescription("Get details of a specific ACL role"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the role to get"),
		),
	)
	s.AddTool(getACLRoleTool, GetACLRoleHandler(nomadClient, logger))

	createACLRoleTool := mcp.NewTool("create_acl_role",
		mcp.WithDescription("Create a new ACL role"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the role"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the role"),
		),
		mcp.WithString("policies",
			mcp.Required(),
			mcp.Description("List of policy names to associate with the role"),
		),
	)
	s.AddTool(createACLRoleTool, CreateACLRoleHandler(nomadClient, logger))

	deleteACLRoleTool := mcp.NewTool("delete_acl_role",
		mcp.WithDescription("Delete an ACL role"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the role to delete"),
		),
	)
	s.AddTool(deleteACLRoleTool, DeleteACLRoleHandler(nomadClient, logger))

	// Bootstrap ACL token tool
	bootstrapACLTokenTool := mcp.NewTool("bootstrap_acl_token",
		mcp.WithDescription("Bootstrap the ACL system and get the initial management token"),
	)
	s.AddTool(bootstrapACLTokenTool, BootstrapACLTokenHandler(nomadClient, logger))
}

// ListACLTokensHandler handles the list_acl_tokens tool request
func ListACLTokensHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tokens, err := nomadClient.ListACLTokens()
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

// GetACLTokenHandler handles the get_acl_token tool request
func GetACLTokenHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		accessorID, ok := arguments["accessor_id"].(string)
		if !ok || accessorID == "" {
			return mcp.NewToolResultError("accessor_id is required"), nil
		}

		token, err := nomadClient.GetACLToken(accessorID)
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

// CreateACLTokenHandler handles the create_acl_token tool request
func CreateACLTokenHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		tokenType, ok := arguments["type"].(string)
		if !ok || tokenType == "" {
			return mcp.NewToolResultError("type is required"), nil
		}

		var policies []string
		if policiesParam, ok := arguments["policies"].([]interface{}); ok {
			for _, p := range policiesParam {
				if policy, ok := p.(string); ok {
					policies = append(policies, policy)
				}
			}
		}

		global := false
		if globalParam, ok := arguments["global"].(bool); ok {
			global = globalParam
		}

		token := types.ACLToken{
			Name:     name,
			Type:     tokenType,
			Policies: policies,
			Global:   global,
		}

		createdToken, err := nomadClient.CreateACLToken(token)
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

// DeleteACLTokenHandler handles the delete_acl_token tool request
func DeleteACLTokenHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		accessorID, ok := arguments["accessor_id"].(string)
		if !ok || accessorID == "" {
			return mcp.NewToolResultError("accessor_id is required"), nil
		}

		err := nomadClient.DeleteACLToken(accessorID)
		if err != nil {
			logger.Printf("Error deleting ACL token: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete ACL token", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("ACL token %s deleted successfully", accessorID)), nil
	}
}

// ListACLPoliciesHandler handles the list_acl_policies tool request
func ListACLPoliciesHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policies, err := nomadClient.ListACLPolicies()
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

// GetACLPolicyHandler handles the get_acl_policy tool request
func GetACLPolicyHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		policy, err := nomadClient.GetACLPolicy(name)
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

// CreateACLPolicyHandler handles the create_acl_policy tool request
func CreateACLPolicyHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		rules, ok := arguments["rules"].(string)
		if !ok || rules == "" {
			return mcp.NewToolResultError("rules is required"), nil
		}

		description := ""
		if desc, ok := arguments["description"].(string); ok {
			description = desc
		}

		policy := types.ACLPolicy{
			Name:        name,
			Description: description,
			Rules:       rules,
		}

		err := nomadClient.CreateACLPolicy(policy)
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

// DeleteACLPolicyHandler handles the delete_acl_policy tool request
func DeleteACLPolicyHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		err := nomadClient.DeleteACLPolicy(name)
		if err != nil {
			logger.Printf("Error deleting ACL policy: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete ACL policy", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("ACL policy %s deleted successfully", name)), nil
	}
}

// ListACLRolesHandler handles the list_acl_roles tool request
func ListACLRolesHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		roles, err := nomadClient.ListACLRoles()
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

// GetACLRoleHandler handles the get_acl_role tool request
func GetACLRoleHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		id, ok := arguments["id"].(string)
		if !ok || id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}

		role, err := nomadClient.GetACLRole(id)
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

// CreateACLRoleHandler handles the create_acl_role tool request
func CreateACLRoleHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		description := ""
		if desc, ok := arguments["description"].(string); ok {
			description = desc
		}

		policiesParam, ok := arguments["policies"]
		if !ok || policiesParam == nil {
			return mcp.NewToolResultError("Specify at least one policy"), nil
		}

		var policyNames []string

		if policiesArray, ok := policiesParam.([]interface{}); ok {
			for _, p := range policiesArray {
				if policyStr, ok := p.(string); ok && policyStr != "" {
					policyNames = append(policyNames, policyStr)
				}
			}
		} else if policyStr, ok := policiesParam.(string); ok && policyStr != "" {
			policyNames = append(policyNames, policyStr)
		} else if policyMap, ok := policiesParam.(map[string]interface{}); ok {

			if policyArr, ok := policyMap["Policies"].([]interface{}); ok {
				for _, p := range policyArr {
					if pm, ok := p.(map[string]interface{}); ok {
						if pName, ok := pm["Name"].(string); ok && pName != "" {
							policyNames = append(policyNames, pName)
						}
					}
				}
			}
		}

		if len(policyNames) == 0 {
			return mcp.NewToolResultError("Specify at least one policy"), nil
		}

		policyLinks := make([]map[string]string, len(policyNames))
		for i, name := range policyNames {
			policyLinks[i] = map[string]string{"Name": name}
		}

		role := types.ACLRole{
			Name:        name,
			Description: description,
			Policies:    policyLinks,
		}

		role, err := nomadClient.CreateACLRole(role)
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

// DeleteACLRoleHandler handles the delete_acl_role tool request
func DeleteACLRoleHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		id, ok := arguments["id"].(string)
		if !ok || id == "" {
			return mcp.NewToolResultError("id is required"), nil
		}

		err := nomadClient.DeleteACLRole(id)
		if err != nil {
			logger.Printf("Error deleting ACL role: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete ACL role", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("ACL role %s deleted successfully", id)), nil
	}
}

// BootstrapACLTokenHandler handles the bootstrap_acl_token tool request
func BootstrapACLTokenHandler(nomadClient *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		token, err := nomadClient.BootstrapACLToken()
		if err != nil {
			logger.Printf("Error bootstrapping ACL token: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to bootstrap ACL token", err), nil
		}

		// Save the token in the client
		nomadClient.SetToken(token.SecretID)

		tokenJSON, err := json.MarshalIndent(token, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format token details", err), nil
		}

		return mcp.NewToolResultText(string(tokenJSON)), nil
	}
}
