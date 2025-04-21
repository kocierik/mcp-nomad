package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterPrompts registers all prompts for the Nomad MCP server
func RegisterPrompts(s *server.MCPServer) {
	// Job management prompt
	s.AddPrompt(mcp.NewPrompt("job_management",
		mcp.WithPromptDescription("Assist with Nomad job management tasks"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("The action to perform (list, get, run, stop, scale)"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("job_id",
			mcp.ArgumentDescription("The ID of the job (required for get, stop, scale)"),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("The namespace to operate in (default: default)"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		jobID := request.Params.Arguments["job_id"]
		namespace := request.Params.Arguments["namespace"]
		if namespace == "" {
			namespace = "default"
		}

		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage(
			"system",
			mcp.NewTextContent("You are a Nomad job management assistant. Help users manage their Nomad jobs effectively."),
		))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you list jobs in the %s namespace. You can filter by status (pending, running, dead) if needed.", namespace)),
			))
		case "get":
			if jobID == "" {
				return nil, fmt.Errorf("job_id is required for get action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you get details for job %s in the %s namespace.", jobID, namespace)),
			))
		case "run":
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent("I'll help you run a new job. Please provide the job specification in HCL or JSON format."),
			))
		case "stop":
			if jobID == "" {
				return nil, fmt.Errorf("job_id is required for stop action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you stop job %s in the %s namespace. You can also purge the job if needed.", jobID, namespace)),
			))
		case "scale":
			if jobID == "" {
				return nil, fmt.Errorf("job_id is required for scale action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you scale job %s in the %s namespace. Please specify the task group and desired count.", jobID, namespace)),
			))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult(
			"Nomad Job Management",
			messages,
		), nil
	})

	// Node management prompt
	s.AddPrompt(mcp.NewPrompt("node_management",
		mcp.WithPromptDescription("Assist with Nomad node management tasks"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("The action to perform (list, get, drain, eligibility)"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("node_id",
			mcp.ArgumentDescription("The ID of the node (required for get, drain, eligibility)"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		nodeID := request.Params.Arguments["node_id"]

		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage(
			"system",
			mcp.NewTextContent("You are a Nomad node management assistant. Help users manage their Nomad nodes effectively."),
		))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent("I'll help you list nodes in the cluster. You can filter by status (ready, down) if needed."),
			))
		case "get":
			if nodeID == "" {
				return nil, fmt.Errorf("node_id is required for get action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you get details for node %s.", nodeID)),
			))
		case "drain":
			if nodeID == "" {
				return nil, fmt.Errorf("node_id is required for drain action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you manage drain mode for node %s. Please specify whether to enable or disable drain mode.", nodeID)),
			))
		case "eligibility":
			if nodeID == "" {
				return nil, fmt.Errorf("node_id is required for eligibility action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you set eligibility for node %s. Please specify whether to make it eligible or ineligible.", nodeID)),
			))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult(
			"Nomad Node Management",
			messages,
		), nil
	})

	// Namespace management prompt
	s.AddPrompt(mcp.NewPrompt("namespace_management",
		mcp.WithPromptDescription("Assist with Nomad namespace management tasks"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("The action to perform (list, create, delete)"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("name",
			mcp.ArgumentDescription("The name of the namespace (required for create, delete)"),
		),
		mcp.WithArgument("description",
			mcp.ArgumentDescription("Description of the namespace (optional for create)"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		name := request.Params.Arguments["name"]
		description := request.Params.Arguments["description"]

		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage(
			"system",
			mcp.NewTextContent("You are a Nomad namespace management assistant. Help users manage their Nomad namespaces effectively."),
		))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent("I'll help you list all namespaces in Nomad."),
			))
		case "create":
			if name == "" {
				return nil, fmt.Errorf("name is required for create action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you create namespace %s%s.", name, func() string {
					if description != "" {
						return fmt.Sprintf(" with description: %s", description)
					}
					return ""
				}())),
			))
		case "delete":
			if name == "" {
				return nil, fmt.Errorf("name is required for delete action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you delete namespace %s.", name)),
			))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult(
			"Nomad Namespace Management",
			messages,
		), nil
	})

	// Variable management prompt
	s.AddPrompt(mcp.NewPrompt("variable_management",
		mcp.WithPromptDescription("Assist with Nomad variable management tasks"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("The action to perform (list, get, create, delete)"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("path",
			mcp.ArgumentDescription("The path of the variable (required for get, create, delete)"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		path := request.Params.Arguments["path"]

		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage(
			"system",
			mcp.NewTextContent("You are a Nomad variable management assistant. Help users manage their Nomad variables effectively."),
		))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent("I'll help you list variables in Nomad. You can optionally filter by prefix."),
			))
		case "get":
			if path == "" {
				return nil, fmt.Errorf("path is required for get action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you get the variable at path %s.", path)),
			))
		case "create":
			if path == "" {
				return nil, fmt.Errorf("path is required for create action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you create a variable at path %s. Please provide the key-value pairs to store.", path)),
			))
		case "delete":
			if path == "" {
				return nil, fmt.Errorf("path is required for delete action")
			}
			messages = append(messages, mcp.NewPromptMessage(
				"assistant",
				mcp.NewTextContent(fmt.Sprintf("I'll help you delete the variable at path %s.", path)),
			))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult(
			"Nomad Variable Management",
			messages,
		), nil
	})

	// ACL management prompt
	s.AddPrompt(mcp.NewPrompt("acl_management",
		mcp.WithPromptDescription("Assist with Nomad ACL management tasks"),
		mcp.WithArgument("resource",
			mcp.ArgumentDescription("The ACL resource type (token, policy, role)"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("The action to perform (list, get, create, delete)"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		resource := request.Params.Arguments["resource"]
		action := request.Params.Arguments["action"]

		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage(
			"system",
			mcp.NewTextContent("You are a Nomad ACL management assistant. Help users manage their Nomad ACL resources effectively."),
		))

		switch resource {
		case "token":
			switch action {
			case "list":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you list all ACL tokens."),
				))
			case "get":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you get details for a specific ACL token. Please provide the accessor ID."),
				))
			case "create":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you create a new ACL token. Please provide the name, type (client/management), and optional policies."),
				))
			case "delete":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you delete an ACL token. Please provide the accessor ID."),
				))
			default:
				return nil, fmt.Errorf("invalid action for token: %s", action)
			}
		case "policy":
			switch action {
			case "list":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you list all ACL policies."),
				))
			case "get":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you get details for a specific ACL policy. Please provide the policy name."),
				))
			case "create":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you create a new ACL policy. Please provide the name, description, and HCL rules."),
				))
			case "delete":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you delete an ACL policy. Please provide the policy name."),
				))
			default:
				return nil, fmt.Errorf("invalid action for policy: %s", action)
			}
		case "role":
			switch action {
			case "list":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you list all ACL roles."),
				))
			case "get":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you get details for a specific ACL role. Please provide the role ID."),
				))
			case "create":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you create a new ACL role. Please provide the name, description, and associated policies."),
				))
			case "delete":
				messages = append(messages, mcp.NewPromptMessage(
					"assistant",
					mcp.NewTextContent("I'll help you delete an ACL role. Please provide the role ID."),
				))
			default:
				return nil, fmt.Errorf("invalid action for role: %s", action)
			}
		default:
			return nil, fmt.Errorf("invalid resource type: %s", resource)
		}

		return mcp.NewGetPromptResult(
			"Nomad ACL Management",
			messages,
		), nil
	})
}
