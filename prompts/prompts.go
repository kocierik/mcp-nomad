package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	guideJSONTools = "Tool results are JSON text. Summarize for humans: IDs, status fields, counts, and errors first. " +
		"When listing many items, highlight anomalies (failed, draining, dead). " +
		"If the user asked for a comparison or diagnosis, structure bullets with evidence from the payload."
)

// RegisterPrompts registers all prompts for the Nomad MCP server
func RegisterPrompts(s *server.MCPServer) {
	registerJobPrompts(s)
	registerNodePrompts(s)
	registerNamespacePrompts(s)
	registerVariablePrompts(s)
	registerACLPrompts(s)
}

func registerJobPrompts(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("job_management",
		mcp.WithPromptDescription("Nomad jobs: maps actions to MCP tools list_jobs, get_job, run_job, stop_job, scale_job, and related job tools"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("list | get | run | stop | scale"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("job_id",
			mcp.ArgumentDescription("Required for get, stop, scale (Nomad job ID / name)"),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Target namespace; omit to use NOMAD_NAMESPACE env or default"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		jobID := request.Params.Arguments["job_id"]
		namespace := effectiveNamespaceFromPrompt(request.Params.Arguments)

		sys := fmt.Sprintf("You are a Nomad job assistant. Effective namespace for tools is %q (prompt `namespace` argument, then NOMAD_NAMESPACE env, else default). "+
			"Prefer the smallest set of tool calls. Multi-region clusters: NOMAD_REGION is forwarded on API requests when set. "+
			"%s "+
			"Relevant tools: list_jobs, get_job, run_job, stop_job, scale_job, get_job_allocations, get_job_evaluations, get_job_deployments, get_job_summary, get_job_services.",
			namespace, guideJSONTools)

		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage("system", mcp.NewTextContent(sys)))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **list_jobs** with namespace %q and optional status (pending, running, dead). "+
					"If the user needs full job objects, follow with **get_job** per ID. Explain ListJobs may return enriched summaries when the server merges stub+summary+get.", namespace),
			)))
		case "get":
			if jobID == "" {
				return nil, fmt.Errorf("job_id is required for get action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **get_job** with job_id %q, namespace %q. Point out Type, Status, TaskGroups, allocations count if present.", jobID, namespace),
			)))
		case "run":
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				"Use **run_job** with job_spec (HCL or JSON) and optional detach. After success, mention EvalID / modify index if returned; suggest get_job or list_jobs to verify.",
			)))
		case "stop":
			if jobID == "" {
				return nil, fmt.Errorf("job_id is required for stop action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **stop_job** with job_id %q, namespace %q, purge if the user wants removal from state.", jobID, namespace),
			)))
		case "scale":
			if jobID == "" {
				return nil, fmt.Errorf("job_id is required for scale action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **scale_job** with job_id %q, namespace %q, plus task group name and integer count from the user.", jobID, namespace),
			)))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult("Nomad Job Management", messages), nil
	})
}

func registerNodePrompts(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("node_management",
		mcp.WithPromptDescription("Nomad clients: list_nodes, get_node, drain_node, eligibility_node"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("list | get | drain | eligibility"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("node_id",
			mcp.ArgumentDescription("Required for get, drain, eligibility (Nomad node ID)"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		nodeID := request.Params.Arguments["node_id"]

		sys := "You are a Nomad node assistant. " + guideJSONTools + " Tools: list_nodes, get_node, drain_node, eligibility_node."
		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage("system", mcp.NewTextContent(sys)))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				"Use **list_nodes** with optional status filter when the user cares about ready/down. Highlight scheduling health and drain flags from each summary.",
			)))
		case "get":
			if nodeID == "" {
				return nil, fmt.Errorf("node_id is required for get action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **get_node** for node_id %q. Explain Status, Drain, SchedulingEligibility and host/driver info.", nodeID),
			)))
		case "drain":
			if nodeID == "" {
				return nil, fmt.Errorf("node_id is required for drain action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **drain_node** for %q with enable true/false and deadline seconds if the user specifies maintenance windows.", nodeID),
			)))
		case "eligibility":
			if nodeID == "" {
				return nil, fmt.Errorf("node_id is required for eligibility action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **eligibility_node** for %q; eligible argument is typically `eligible` or `ineligible` per Nomad API.", nodeID),
			)))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult("Nomad Node Management", messages), nil
	})
}

func registerNamespacePrompts(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("namespace_management",
		mcp.WithPromptDescription("Namespaces: list_namespaces, create_namespace, delete_namespace"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("list | create | delete"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("name",
			mcp.ArgumentDescription("Namespace name (required for create and delete)"),
		),
		mcp.WithArgument("description",
			mcp.ArgumentDescription("Optional description when creating"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		name := request.Params.Arguments["name"]
		description := request.Params.Arguments["description"]

		sys := "You are a Nomad namespace assistant. " + guideJSONTools + " Tools: list_namespaces, create_namespace, delete_namespace."
		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage("system", mcp.NewTextContent(sys)))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				"Use **list_namespaces**. Summarize Name and Description columns; warn if production and default coexist.",
			)))
		case "create":
			if name == "" {
				return nil, fmt.Errorf("name is required for create action")
			}
			createHint := ""
			if description != "" {
				createHint = fmt.Sprintf(" Description: %q.", description)
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **create_namespace** with Name %q.%s Confirm with list_namespaces if needed.", name, createHint),
			)))
		case "delete":
			if name == "" {
				return nil, fmt.Errorf("name is required for delete action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **delete_namespace** for %q only after the user acknowledges impact on workloads in that namespace.", name),
			)))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult("Nomad Namespace Management", messages), nil
	})
}

func registerVariablePrompts(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("variable_management",
		mcp.WithPromptDescription("Nomad Variables: list_variables, get_variable, create_variable, delete_variable"),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("list | get | create | delete"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("path",
			mcp.ArgumentDescription("Variable path (required for get, create, delete)"),
		),
		mcp.WithArgument("namespace",
			mcp.ArgumentDescription("Variable namespace; omit to use NOMAD_NAMESPACE env or default"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		action := request.Params.Arguments["action"]
		path := request.Params.Arguments["path"]
		namespace := effectiveNamespaceFromPrompt(request.Params.Arguments)

		sys := "You are a Nomad Variables assistant. Namespace matches tools via prompt `namespace` or NOMAD_NAMESPACE. " + guideJSONTools + " Tools: list_variables, get_variable, create_variable, delete_variable."
		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage("system", mcp.NewTextContent(sys)))

		switch action {
		case "list":
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **list_variables** with namespace %q and optional prefix, next_token, per_page, filter from the user.", namespace),
			)))
		case "get":
			if path == "" {
				return nil, fmt.Errorf("path is required for get action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **get_variable** for path %q in namespace %q. Variables store Items as structured data; summarize keys, not secrets, unless the user owns the cluster.", path, namespace),
			)))
		case "create":
			if path == "" {
				return nil, fmt.Errorf("path is required for create action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **create_variable**: path %q, namespace %q, key/value from user; respect CAS / lock_operation if they mention concurrency.", path, namespace),
			)))
		case "delete":
			if path == "" {
				return nil, fmt.Errorf("path is required for delete action")
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(
				fmt.Sprintf("Use **delete_variable** for path %q, namespace %q; include CAS index if user provides optimistic locking.", path, namespace),
			)))
		default:
			return nil, fmt.Errorf("invalid action: %s", action)
		}

		return mcp.NewGetPromptResult("Nomad Variable Management", messages), nil
	})
}

func registerACLPrompts(s *server.MCPServer) {
	s.AddPrompt(mcp.NewPrompt("acl_management",
		mcp.WithPromptDescription("ACL tokens, policies, roles: map to list_acl_tokens, get_acl_token, create_acl_token, bootstrap_acl_token, *_policy*, *_role* tools"),
		mcp.WithArgument("resource",
			mcp.ArgumentDescription("token | policy | role"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("action",
			mcp.ArgumentDescription("list | get | create | delete"),
			mcp.RequiredArgument(),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		resource := request.Params.Arguments["resource"]
		action := request.Params.Arguments["action"]

		sys := "You are a Nomad ACL assistant. Treat tokens and policies as sensitive: never echo SecretID broadly; remind users about least privilege. " +
			"Initial cluster ACL setup uses **bootstrap_acl_token** only when the user explicitly intends to bootstrap. " + guideJSONTools
		var messages []mcp.PromptMessage
		messages = append(messages, mcp.NewPromptMessage("system", mcp.NewTextContent(sys)))

		switch resource {
		case "token":
			var extra string
			switch action {
			case "list":
				extra = "Use **list_acl_tokens**. Summarize AccessorID, Name, Global; do not dump secret material."
			case "get":
				extra = "Use **get_acl_token** with accessor_id. Describe linked policies/roles; avoid unnecessary SecretID exposure."
			case "create":
				extra = "Use **create_acl_token** with fields the user confirms (type, policies, roles). Explain client vs management token impact."
			case "delete":
				extra = "Use **delete_acl_token** only after the user confirms the accessor_id to revoke."
			default:
				return nil, fmt.Errorf("invalid action for token: %s", action)
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(extra)))

		case "policy":
			var extra string
			switch action {
			case "list":
				extra = "Use **list_acl_policies**."
			case "get":
				extra = "Use **get_acl_policy** with the policy name."
			case "create":
				extra = "Use **create_acl_policy** with name and rules body in the shape the tool expects."
			case "delete":
				extra = "Use **delete_acl_policy** with policy name."
			default:
				return nil, fmt.Errorf("invalid action for policy: %s", action)
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(extra)))

		case "role":
			var extra string
			switch action {
			case "list":
				extra = "Use **list_acl_roles**."
			case "get":
				extra = "Use **get_acl_role** with role ID."
			case "create":
				extra = "Use **create_acl_role**; gather policies to attach."
			case "delete":
				extra = "Use **delete_acl_role** after confirmation."
			default:
				return nil, fmt.Errorf("invalid action for role: %s", action)
			}
			messages = append(messages, mcp.NewPromptMessage("assistant", mcp.NewTextContent(extra)))

		default:
			return nil, fmt.Errorf("invalid resource type: %s", resource)
		}

		return mcp.NewGetPromptResult("Nomad ACL Management", messages), nil
	})
}
