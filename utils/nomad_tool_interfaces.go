package utils

import (
	"context"

	"github.com/kocierik/mcp-nomad/types"
)

// Narrow interfaces consumed by MCP tool handlers enable safer refactors and tests via mocks.

// RawNomadCaller is exposed for endpoints that tools address with arbitrary paths/methods.
type RawNomadCaller interface {
	MakeRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}) ([]byte, error)
}

var _ RawNomadCaller = (*NomadClient)(nil)

// NomadACLTokenBootstrapper resets the in-memory client token after ACL bootstrap (tools/acl.go).
type NomadACLTokenBootstrapper interface {
	SetToken(token string)
}

var _ NomadACLTokenBootstrapper = (*NomadClient)(nil)

// JobAPI is implemented by NomadClient and used by job-related MCP tools plus dynamic resources.
type JobAPI interface {
	ListJobs(ctx context.Context, namespace, status string) ([]types.JobSummary, error)
	GetJob(ctx context.Context, jobID, namespace string) (types.Job, error)
	RunJob(ctx context.Context, jobSpec string, detach bool) (map[string]interface{}, error)
	StopJob(ctx context.Context, jobID, namespace string, purge bool) (map[string]interface{}, error)
	ScaleTaskGroup(ctx context.Context, jobID, group string, count int, namespace string) error
	ListJobAllocations(ctx context.Context, jobID, namespace string) ([]types.Allocation, error)
	ListJobEvaluations(ctx context.Context, jobID, namespace string) ([]types.Evaluation, error)
	ListJobDeployments(ctx context.Context, jobID, namespace string) ([]types.JobDeployment, error)
	GetJobSummary(ctx context.Context, jobID, namespace string) (types.JobSummary, error)
	ListJobServices(ctx context.Context, jobID, namespace string) ([]types.Service, error)
	GetJobVersions(ctx context.Context, jobID, namespace string) ([]types.Job, error)
}

var _ JobAPI = (*NomadClient)(nil)

// NodeAPI backs node MCP tools (and helpers that inspect node payloads in resources).
type NodeAPI interface {
	ListNodes(ctx context.Context, status string) ([]types.NodeSummary, error)
	GetNode(ctx context.Context, nodeID string) (types.Node, error)
	DrainNode(ctx context.Context, nodeID string, enable bool, deadline int64) (string, error)
	EligibilityNode(ctx context.Context, nodeID string, eligible string) (types.NodeSummary, error)
}

var _ NodeAPI = (*NomadClient)(nil)

// NamespaceAPI backs namespace tools.
type NamespaceAPI interface {
	ListNamespaces(ctx context.Context) ([]types.Namespace, error)
	CreateNamespace(ctx context.Context, namespace types.Namespace) error
	DeleteNamespace(ctx context.Context, name string) error
}

var _ NamespaceAPI = (*NomadClient)(nil)

// DeploymentAPI backs deployment MCP tools (global deployments listing).
type DeploymentAPI interface {
	ListDeployments(ctx context.Context, namespace string) ([]types.DeploymentSummary, error)
	GetDeployment(ctx context.Context, deploymentID string) (types.Deployment, error)
}

var _ DeploymentAPI = (*NomadClient)(nil)

// VolumeAPI backs CSI/host volume MCP tools currently exposed via MCP.
type VolumeAPI interface {
	ListVolumes(ctx context.Context, nodeID string, pluginID string, nextToken string, perPage int, filter string) ([]types.Volume, error)
	GetVolume(ctx context.Context, volumeID string) (*types.Volume, error)
	DeleteVolume(ctx context.Context, volumeID string) error
}

var _ VolumeAPI = (*NomadClient)(nil)

// VariableAPI backs Nomad Variables tools.
type VariableAPI interface {
	ListVariables(ctx context.Context, namespace, prefix string, nextToken string, perPage int, filter string) ([]types.Variable, error)
	GetVariable(ctx context.Context, path, namespace string) (types.Variable, error)
	CreateVariable(ctx context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error
	DeleteVariable(ctx context.Context, path, namespace string, cas int) error
}

var _ VariableAPI = (*NomadClient)(nil)

// AllocationAPI backs allocation MCP tools (no arbitrary HTTP; cluster tools use ClusterToolsAPI).
type AllocationAPI interface {
	ListAllocations(ctx context.Context, namespace, jobID string) ([]types.Allocation, error)
	GetAllocation(ctx context.Context, allocID string) (types.Allocation, error)
	StopAllocation(ctx context.Context, allocID string) error
}

var _ AllocationAPI = (*NomadClient)(nil)

// LogAPI backs allocation log tools.
type LogAPI interface {
	GetAllocationLogs(ctx context.Context, allocID, task, logType string, follow bool, tail, offset int64) (string, error)
}

var _ LogAPI = (*NomadClient)(nil)

// ACLAPI backs ACL MCP tools except SetToken refresh after bootstrap.
type ACLAPI interface {
	ListACLTokens(ctx context.Context) ([]types.ACLToken, error)
	GetACLToken(ctx context.Context, accessorID string) (types.ACLToken, error)
	CreateACLToken(ctx context.Context, token types.ACLToken) (types.ACLToken, error)
	DeleteACLToken(ctx context.Context, accessorID string) error
	ListACLPolicies(ctx context.Context) ([]types.ACLPolicy, error)
	GetACLPolicy(ctx context.Context, name string) (types.ACLPolicy, error)
	CreateACLPolicy(ctx context.Context, policy types.ACLPolicy) error
	DeleteACLPolicy(ctx context.Context, name string) error
	ListACLRoles(ctx context.Context) ([]types.ACLRole, error)
	GetACLRole(ctx context.Context, id string) (types.ACLRole, error)
	CreateACLRole(ctx context.Context, role types.ACLRole) (types.ACLRole, error)
	DeleteACLRole(ctx context.Context, id string) error
	BootstrapACLToken(ctx context.Context) (types.ACLToken, error)
}

var _ ACLAPI = (*NomadClient)(nil)

// ACLToolsDeps composes ACL API access with bootstrap token propagation.
type ACLToolsDeps interface {
	ACLAPI
	NomadACLTokenBootstrapper
}

var _ ACLToolsDeps = (*NomadClient)(nil)

// SentinelAPI backs Sentinel MCP tools where enabled.
type SentinelAPI interface {
	ListSentinelPolicies(ctx context.Context) ([]types.SentinelPolicy, error)
	GetSentinelPolicy(ctx context.Context, name string) (types.SentinelPolicy, error)
	CreateSentinelPolicy(ctx context.Context, policy types.SentinelPolicy) error
	DeleteSentinelPolicy(ctx context.Context, name string) error
}

var _ SentinelAPI = (*NomadClient)(nil)

// ClusterToolsAPI backs cluster/regions MCP tools.
type ClusterToolsAPI interface {
	RawNomadCaller
	ListClusterPeers(ctx context.Context) ([]byte, error)
}

var _ ClusterToolsAPI = (*NomadClient)(nil)

// DynamicResourcesNomad is the subset of NomadClient used when publishing MCP dynamic resources.
type DynamicResourcesNomad interface {
	GetJob(ctx context.Context, jobID, namespace string) (types.Job, error)
	GetJobVersions(ctx context.Context, jobID, namespace string) ([]types.Job, error)
	GetNode(ctx context.Context, nodeID string) (types.Node, error)
	GetAllocation(ctx context.Context, allocID string) (types.Allocation, error)
	GetAllocationLogs(ctx context.Context, allocID, task, logType string, follow bool, tail, offset int64) (string, error)
}

var _ DynamicResourcesNomad = (*NomadClient)(nil)
