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

	// Register resource for job templates
	registerJobTemplates(s, logger)

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
