package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/kocierik/nomad-mcp-server/tools"
	"github.com/kocierik/nomad-mcp-server/utils"
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

	// Register job-related tools
	registerJobTools(s, nomadClient, logger)

	// Register deployment tools
	registerDeploymentTools(s, nomadClient, logger)

	// Register namespace tools
	registerNamespaceTools(s, nomadClient, logger)

	// Register node tools
	registerNodeTools(s, nomadClient, logger)

	// Register allocation tools
	registerAllocationTools(s, nomadClient, logger)

	// Register variable tools
	registerVariableTools(s, nomadClient, logger)

	// Register additional job operation tools
	registerJobOperationTools(s, nomadClient, logger)

	// Register resource for job templates
	registerJobTemplates(s, logger)

	// Register volume tools
	registerVolumeTools(s, nomadClient, logger)

	// Register ACL tools
	registerACLTools(s, nomadClient, logger)

	// Start the MCP server using stdio
	logger.Println("Starting Nomad MCP server...")
	if err := server.ServeStdio(s); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}

// Register job-related tools
func registerJobTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List jobs tool
	listJobsTool := mcp.NewTool("list_jobs",
		mcp.WithDescription("List all jobs in Nomad"),
		mcp.WithString("namespace",
			mcp.Description("The namespace to list jobs from (default: default)"),
		),
		mcp.WithString("status",
			mcp.Description("Filter jobs by status (pending, running, dead)"),
			mcp.Enum("pending", "running", "dead", ""),
		),
	)
	s.AddTool(listJobsTool, tools.ListJobsHandler(nomadClient, logger))

	// Get job tool
	getJobTool := mcp.NewTool("get_job",
		mcp.WithDescription("Get job details by ID"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to retrieve"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobTool, tools.GetJobHandler(nomadClient, logger))

	// Run job tool
	runJobTool := mcp.NewTool("run_job",
		mcp.WithDescription("Run a new job or update an existing job"),
		mcp.WithString("job_spec",
			mcp.Required(),
			mcp.Description("The job specification in HCL or JSON format"),
		),
		mcp.WithBoolean("detach",
			mcp.Description("Return immediately instead of monitoring deployment"),
		),
	)
	s.AddTool(runJobTool, tools.RunJobHandler(nomadClient, logger))

	// Stop job tool
	stopJobTool := mcp.NewTool("stop_job",
		mcp.WithDescription("Stop a running job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to stop"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
		mcp.WithBoolean("purge",
			mcp.Description("Purge the job from Nomad instead of just stopping it"),
		),
	)
	s.AddTool(stopJobTool, tools.StopJobHandler(nomadClient, logger))
}

// Register deployment tools
func registerDeploymentTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List deployments tool
	listDeploymentsTool := mcp.NewTool("list_deployments",
		mcp.WithDescription("List all deployments"),
		mcp.WithString("namespace",
			mcp.Description("The namespace to list deployments from (default: default)"),
		),
	)
	s.AddTool(listDeploymentsTool, tools.ListDeploymentsHandler(nomadClient, logger))

	// Get deployment tool
	getDeploymentTool := mcp.NewTool("get_deployment",
		mcp.WithDescription("Get deployment details by ID"),
		mcp.WithString("deployment_id",
			mcp.Required(),
			mcp.Description("The ID of the deployment to retrieve"),
		),
	)
	s.AddTool(getDeploymentTool, tools.GetDeploymentHandler(nomadClient, logger))
}

// Register namespace tools
func registerNamespaceTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List namespaces tool
	listNamespacesTool := mcp.NewTool("list_namespaces",
		mcp.WithDescription("List all namespaces in Nomad"),
	)
	s.AddTool(listNamespacesTool, tools.ListNamespacesHandler(nomadClient, logger))

	// Create namespace tool
	createNamespaceTool := mcp.NewTool("create_namespace",
		mcp.WithDescription("Create a new namespace"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the namespace to create"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the namespace"),
		),
	)
	s.AddTool(createNamespaceTool, tools.CreateNamespaceHandler(nomadClient, logger))

	// Delete namespace tool
	deleteNamespaceTool := mcp.NewTool("delete_namespace",
		mcp.WithDescription("Delete a namespace"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the namespace to delete"),
		),
	)
	s.AddTool(deleteNamespaceTool, tools.DeleteNamespaceHandler(nomadClient, logger))
}

// Register node tools
func registerNodeTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List nodes tool
	listNodesTool := mcp.NewTool("list_nodes",
		mcp.WithDescription("List all nodes in the Nomad cluster"),
		mcp.WithString("status",
			mcp.Description("Filter nodes by status"),
			mcp.Enum("ready", "down", ""),
		),
	)
	s.AddTool(listNodesTool, tools.ListNodesHandler(nomadClient, logger))

	// Get node tool
	getNodeTool := mcp.NewTool("get_node",
		mcp.WithDescription("Get details for a specific node"),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("The ID of the node to retrieve"),
		),
	)
	s.AddTool(getNodeTool, tools.GetNodeHandler(nomadClient, logger))

	// Drain node tool
	drainNodeTool := mcp.NewTool("drain_node",
		mcp.WithDescription("Enable or disable drain mode for a node"),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("The ID of the node to drain"),
		),
		mcp.WithBoolean("enable",
			mcp.Required(),
			mcp.Description("Enable or disable drain mode"),
		),
		mcp.WithNumber("deadline",
			mcp.Description("Deadline in seconds for the drain operation (default: -1, no deadline)"),
		),
	)
	s.AddTool(drainNodeTool, tools.DrainNodeHandler(nomadClient, logger))

	// Eligibility node tool
	eligibilityNodeTool := mcp.NewTool("eligibility_node",
		mcp.WithDescription("Set eligibility for a node"),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("The ID of the node to set eligibility for"),
		),
		mcp.WithString("eligible",
			mcp.Required(),
			mcp.Description("The eligibility status to set (eligible or ineligible)"),
		),
	)
	s.AddTool(eligibilityNodeTool, tools.EligibilityNodeHandler(nomadClient, logger))
}

// Register allocation tools
func registerAllocationTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List allocations tool
	listAllocationsTool := mcp.NewTool("list_allocations",
		mcp.WithDescription("List all allocations in Nomad"),
		mcp.WithString("namespace",
			mcp.Description("The namespace to list allocations from (default: default)"),
		),
		mcp.WithString("job_id",
			mcp.Description("Filter allocations by job ID"),
		),
	)
	s.AddTool(listAllocationsTool, tools.ListAllocationsHandler(nomadClient, logger))

	// Get allocation tool
	getAllocationTool := mcp.NewTool("get_allocation",
		mcp.WithDescription("Get allocation details by ID"),
		mcp.WithString("allocation_id",
			mcp.Required(),
			mcp.Description("The ID of the allocation to retrieve"),
		),
	)
	s.AddTool(getAllocationTool, tools.GetAllocationHandler(nomadClient, logger))

	// Stop allocation tool
	stopAllocationTool := mcp.NewTool("stop_allocation",
		mcp.WithDescription("Stop a running allocation"),
		mcp.WithString("allocation_id",
			mcp.Required(),
			mcp.Description("The ID of the allocation to stop"),
		),
	)
	s.AddTool(stopAllocationTool, tools.StopAllocationHandler(nomadClient, logger))
}

// Register variable tools
func registerVariableTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List variables tool
	listVariablesTool := mcp.NewTool("list_variables",
		mcp.WithDescription("List all variables in Nomad"),
		mcp.WithString("prefix",
			mcp.Description("Optional prefix to filter variables"),
		),
	)
	s.AddTool(listVariablesTool, tools.ListVariablesHandler(nomadClient, logger))

	// Get variable tool
	getVariableTool := mcp.NewTool("get_variable",
		mcp.WithDescription("Get variable details by path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path of the variable to retrieve"),
		),
	)
	s.AddTool(getVariableTool, tools.GetVariableHandler(nomadClient, logger))

	// Create variable tool
	createVariableTool := mcp.NewTool("create_variable",
		mcp.WithDescription("Create or update a variable"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path where to create the variable"),
		),
		mcp.WithObject("items",
			mcp.Required(),
			mcp.Description("The key-value pairs to store in the variable"),
		),
	)
	s.AddTool(createVariableTool, tools.CreateVariableHandler(nomadClient, logger))

	// Delete variable tool
	deleteVariableTool := mcp.NewTool("delete_variable",
		mcp.WithDescription("Delete a variable"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path of the variable to delete"),
		),
	)
	s.AddTool(deleteVariableTool, tools.DeleteVariableHandler(nomadClient, logger))
}

// Register additional job operation tools
func registerJobOperationTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {

	// Scale job tool
	scaleJobTool := mcp.NewTool("scale_job",
		mcp.WithDescription("Scale a job's task group"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to scale"),
		),
		mcp.WithString("group",
			mcp.Required(),
			mcp.Description("The task group to scale"),
		),
		mcp.WithNumber("count",
			mcp.Required(),
			mcp.Description("The desired count of the task group"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(scaleJobTool, tools.ScaleJobHandler(nomadClient, logger))

	// Get job allocations tool
	getJobAllocationsTool := mcp.NewTool("get_job_allocations",
		mcp.WithDescription("Get allocations for a job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobAllocationsTool, tools.GetJobAllocationsHandler(nomadClient, logger))

	// Get job evaluations tool
	getJobEvaluationsTool := mcp.NewTool("get_job_evaluations",
		mcp.WithDescription("Get evaluations for a job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobEvaluationsTool, tools.GetJobEvaluationsHandler(nomadClient, logger))
}

// Register job templates as resources
func registerJobTemplates(s *server.MCPServer, logger *log.Logger) {
	// Job templates resource
	jobTemplatesResource := mcp.NewResource(
		"nomad://templates",
		"Nomad Job Templates",
		mcp.WithResourceDescription("Available Nomad job templates"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(jobTemplatesResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		templates, err := utils.GetJobTemplates()
		if err != nil {
			logger.Printf("Error getting job templates: %v", err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "nomad://templates",
				MIMEType: "application/json",
				Text:     templates,
			},
		}, nil
	})

	// Template by name resource
	templateResource := mcp.NewResourceTemplate(
		"nomad://templates/{name}",
		"Nomad Job Template",
		mcp.WithTemplateDescription("Specific Nomad job template"),
		mcp.WithTemplateMIMEType("text/x-nomad-hcl"),
	)

	s.AddResourceTemplate(templateResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract template name from URI
		name := utils.ExtractTemplateNameFromURI(request.Params.URI)
		if name == "" {
			return nil, fmt.Errorf("invalid template URI: %s", request.Params.URI)
		}

		template, err := utils.GetJobTemplate(name)
		if err != nil {
			logger.Printf("Error getting job template %s: %v", name, err)
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/x-nomad-hcl",
				Text:     template,
			},
		}, nil
	})
}

func registerVolumeTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List volumes tool
	listVolumesTool := mcp.NewTool("list_volumes",
		mcp.WithDescription("List all volumes in a namespace"),
		mcp.WithString("namespace",
			mcp.Description("Namespace to list volumes from (optional)"),
		),
	)
	s.AddTool(listVolumesTool, tools.ListVolumesHandler(nomadClient, logger))

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
	s.AddTool(getVolumeTool, tools.GetVolumeHandler(nomadClient, logger))

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
	s.AddTool(deleteVolumeTool, tools.DeleteVolumeHandler(nomadClient, logger))
}

func registerACLTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// ACL Token tools
	listACLTokensTool := mcp.NewTool("list_acl_tokens",
		mcp.WithDescription("List all ACL tokens"),
	)
	s.AddTool(listACLTokensTool, tools.ListACLTokensHandler(nomadClient, logger))

	getACLTokenTool := mcp.NewTool("get_acl_token",
		mcp.WithDescription("Get details of a specific ACL token"),
		mcp.WithString("accessor_id",
			mcp.Required(),
			mcp.Description("Accessor ID of the token to get"),
		),
	)
	s.AddTool(getACLTokenTool, tools.GetACLTokenHandler(nomadClient, logger))

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
	s.AddTool(createACLTokenTool, tools.CreateACLTokenHandler(nomadClient, logger))

	deleteACLTokenTool := mcp.NewTool("delete_acl_token",
		mcp.WithDescription("Delete an ACL token"),
		mcp.WithString("accessor_id",
			mcp.Required(),
			mcp.Description("Accessor ID of the token to delete"),
		),
	)
	s.AddTool(deleteACLTokenTool, tools.DeleteACLTokenHandler(nomadClient, logger))

	// ACL Policy tools
	listACLPoliciesTool := mcp.NewTool("list_acl_policies",
		mcp.WithDescription("List all ACL policies"),
	)
	s.AddTool(listACLPoliciesTool, tools.ListACLPoliciesHandler(nomadClient, logger))

	getACLPolicyTool := mcp.NewTool("get_acl_policy",
		mcp.WithDescription("Get details of a specific ACL policy"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the policy to get"),
		),
	)
	s.AddTool(getACLPolicyTool, tools.GetACLPolicyHandler(nomadClient, logger))

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
			mcp.Description("HCL rules for the policy"),
		),
	)
	s.AddTool(createACLPolicyTool, tools.CreateACLPolicyHandler(nomadClient, logger))

	deleteACLPolicyTool := mcp.NewTool("delete_acl_policy",
		mcp.WithDescription("Delete an ACL policy"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the policy to delete"),
		),
	)
	s.AddTool(deleteACLPolicyTool, tools.DeleteACLPolicyHandler(nomadClient, logger))

	// ACL Role tools
	listACLRolesTool := mcp.NewTool("list_acl_roles",
		mcp.WithDescription("List all ACL roles"),
	)
	s.AddTool(listACLRolesTool, tools.ListACLRolesHandler(nomadClient, logger))

	getACLRoleTool := mcp.NewTool("get_acl_role",
		mcp.WithDescription("Get details of a specific ACL role"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the role to get"),
		),
	)
	s.AddTool(getACLRoleTool, tools.GetACLRoleHandler(nomadClient, logger))

	createACLRoleTool := mcp.NewTool("create_acl_role",
		mcp.WithDescription("Create a new ACL role"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the role"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the role"),
		),
		mcp.WithArray("policies",
			mcp.Description("List of policy names to associate with the role"),
		),
	)
	s.AddTool(createACLRoleTool, tools.CreateACLRoleHandler(nomadClient, logger))

	deleteACLRoleTool := mcp.NewTool("delete_acl_role",
		mcp.WithDescription("Delete an ACL role"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("ID of the role to delete"),
		),
	)
	s.AddTool(deleteACLRoleTool, tools.DeleteACLRoleHandler(nomadClient, logger))

	// Bootstrap ACL token tool
	bootstrapACLTokenTool := mcp.NewTool("bootstrap_acl_token",
		mcp.WithDescription("Bootstrap the ACL system and get the initial management token"),
	)
	s.AddTool(bootstrapACLTokenTool, tools.BootstrapACLTokenHandler(nomadClient, logger))
}
