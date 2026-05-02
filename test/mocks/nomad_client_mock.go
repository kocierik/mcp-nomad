package mocks

import (
	"context"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/kocierik/mcp-nomad/utils"
)

// Compile-time: mock stays aligned with the narrow MCP tool-facing interfaces on NomadClient.
var (
	_ utils.JobAPI                = (*MockNomadClient)(nil)
	_ utils.NodeAPI               = (*MockNomadClient)(nil)
	_ utils.NamespaceAPI          = (*MockNomadClient)(nil)
	_ utils.DeploymentAPI         = (*MockNomadClient)(nil)
	_ utils.VolumeAPI             = (*MockNomadClient)(nil)
	_ utils.VariableAPI           = (*MockNomadClient)(nil)
	_ utils.AllocationAPI         = (*MockNomadClient)(nil)
	_ utils.LogAPI                = (*MockNomadClient)(nil)
	_ utils.ACLToolsDeps          = (*MockNomadClient)(nil)
	_ utils.SentinelAPI           = (*MockNomadClient)(nil)
	_ utils.ClusterToolsAPI       = (*MockNomadClient)(nil)
	_ utils.DynamicResourcesNomad = (*MockNomadClient)(nil)
)

// MockNomadClient implements the tool-facing subsets of NomadClient for testing.
type MockNomadClient struct {
	// Job methods
	ListJobsFunc             func(context.Context, string, string) ([]types.JobSummary, error)
	GetJobFunc               func(context.Context, string, string) (types.Job, error)
	RunJobFunc               func(context.Context, string, bool) (map[string]interface{}, error)
	StopJobFunc              func(context.Context, string, string, bool) (map[string]interface{}, error)
	ScaleTaskGroupFunc       func(context.Context, string, string, int, string) error
	ListJobAllocationsFunc   func(context.Context, string, string) ([]types.Allocation, error)
	ListJobEvaluationsFunc   func(context.Context, string, string) ([]types.Evaluation, error)
	ListJobDeploymentsFunc   func(context.Context, string, string) ([]types.JobDeployment, error)
	GetJobSummaryFunc        func(context.Context, string, string) (types.JobSummary, error)
	ListJobServicesFunc      func(context.Context, string, string) ([]types.Service, error)
	GetJobVersionsFunc       func(context.Context, string, string) ([]types.Job, error)
	ListDeploymentsFunc      func(context.Context, string) ([]types.DeploymentSummary, error)
	GetDeploymentFunc        func(context.Context, string) (types.Deployment, error)
	ListVolumesFunc          func(context.Context, string, string, string, int, string) ([]types.Volume, error)
	GetVolumeFunc            func(context.Context, string) (*types.Volume, error)
	DeleteVolumeFunc         func(context.Context, string) error
	ListNodesFunc            func(context.Context, string) ([]types.NodeSummary, error)
	GetNodeFunc              func(context.Context, string) (types.Node, error)
	DrainNodeFunc            func(context.Context, string, bool, int64) (string, error)
	EligibilityNodeFunc      func(context.Context, string, string) (types.NodeSummary, error)
	ListNamespacesFunc       func(context.Context) ([]types.Namespace, error)
	CreateNamespaceFunc      func(context.Context, types.Namespace) error
	DeleteNamespaceFunc      func(context.Context, string) error
	ListAllocationsFunc      func(context.Context, string, string) ([]types.Allocation, error)
	GetAllocationFunc        func(context.Context, string) (types.Allocation, error)
	StopAllocationFunc       func(context.Context, string) error
	GetAllocationLogsFunc    func(context.Context, string, string, string, bool, int64, int64) (string, error)
	ListVariablesFunc        func(context.Context, string, string, string, int, string) ([]types.Variable, error)
	GetVariableFunc          func(context.Context, string, string) (types.Variable, error)
	CreateVariableFunc       func(context.Context, types.Variable, string, int, string) error
	DeleteVariableFunc       func(context.Context, string, string, int) error
	ListACLTokensFunc        func(context.Context) ([]types.ACLToken, error)
	GetACLTokenFunc          func(context.Context, string) (types.ACLToken, error)
	CreateACLTokenFunc       func(context.Context, types.ACLToken) (types.ACLToken, error)
	DeleteACLTokenFunc       func(context.Context, string) error
	ListACLPoliciesFunc      func(context.Context) ([]types.ACLPolicy, error)
	GetACLPolicyFunc         func(context.Context, string) (types.ACLPolicy, error)
	CreateACLPolicyFunc      func(context.Context, types.ACLPolicy) error
	DeleteACLPolicyFunc      func(context.Context, string) error
	ListACLRolesFunc         func(context.Context) ([]types.ACLRole, error)
	GetACLRoleFunc           func(context.Context, string) (types.ACLRole, error)
	CreateACLRoleFunc        func(context.Context, types.ACLRole) (types.ACLRole, error)
	DeleteACLRoleFunc        func(context.Context, string) error
	BootstrapACLTokenFunc    func(context.Context) (types.ACLToken, error)
	ListSentinelPoliciesFunc func(context.Context) ([]types.SentinelPolicy, error)
	GetSentinelPolicyFunc    func(context.Context, string) (types.SentinelPolicy, error)
	CreateSentinelPolicyFunc func(context.Context, types.SentinelPolicy) error
	DeleteSentinelPolicyFunc func(context.Context, string) error
	ListClusterPeersFunc     func(context.Context) ([]byte, error)
	MakeRequestFunc          func(context.Context, string, string, map[string]string, interface{}) ([]byte, error)

	token string // SetToken persists here for assertions in tests
}

func (m *MockNomadClient) ListJobs(ctx context.Context, namespace, status string) ([]types.JobSummary, error) {
	if m.ListJobsFunc != nil {
		return m.ListJobsFunc(ctx, namespace, status)
	}
	return []types.JobSummary{}, nil
}

func (m *MockNomadClient) GetJob(ctx context.Context, jobID, namespace string) (types.Job, error) {
	if m.GetJobFunc != nil {
		return m.GetJobFunc(ctx, jobID, namespace)
	}
	return types.Job{}, nil
}

func (m *MockNomadClient) RunJob(ctx context.Context, jobSpec string, detach bool) (map[string]interface{}, error) {
	if m.RunJobFunc != nil {
		return m.RunJobFunc(ctx, jobSpec, detach)
	}
	return map[string]interface{}{}, nil
}

func (m *MockNomadClient) StopJob(ctx context.Context, jobID, namespace string, purge bool) (map[string]interface{}, error) {
	if m.StopJobFunc != nil {
		return m.StopJobFunc(ctx, jobID, namespace, purge)
	}
	return map[string]interface{}{}, nil
}

func (m *MockNomadClient) ScaleTaskGroup(ctx context.Context, jobID, group string, count int, namespace string) error {
	if m.ScaleTaskGroupFunc != nil {
		return m.ScaleTaskGroupFunc(ctx, jobID, group, count, namespace)
	}
	return nil
}

func (m *MockNomadClient) ListJobAllocations(ctx context.Context, jobID, namespace string) ([]types.Allocation, error) {
	if m.ListJobAllocationsFunc != nil {
		return m.ListJobAllocationsFunc(ctx, jobID, namespace)
	}
	return nil, nil
}

func (m *MockNomadClient) ListJobEvaluations(ctx context.Context, jobID, namespace string) ([]types.Evaluation, error) {
	if m.ListJobEvaluationsFunc != nil {
		return m.ListJobEvaluationsFunc(ctx, jobID, namespace)
	}
	return nil, nil
}

func (m *MockNomadClient) ListJobDeployments(ctx context.Context, jobID, namespace string) ([]types.JobDeployment, error) {
	if m.ListJobDeploymentsFunc != nil {
		return m.ListJobDeploymentsFunc(ctx, jobID, namespace)
	}
	return nil, nil
}

func (m *MockNomadClient) GetJobSummary(ctx context.Context, jobID, namespace string) (types.JobSummary, error) {
	if m.GetJobSummaryFunc != nil {
		return m.GetJobSummaryFunc(ctx, jobID, namespace)
	}
	return types.JobSummary{}, nil
}

func (m *MockNomadClient) ListJobServices(ctx context.Context, jobID, namespace string) ([]types.Service, error) {
	if m.ListJobServicesFunc != nil {
		return m.ListJobServicesFunc(ctx, jobID, namespace)
	}
	return nil, nil
}

func (m *MockNomadClient) GetJobVersions(ctx context.Context, jobID, namespace string) ([]types.Job, error) {
	if m.GetJobVersionsFunc != nil {
		return m.GetJobVersionsFunc(ctx, jobID, namespace)
	}
	return nil, nil
}

func (m *MockNomadClient) ListDeployments(ctx context.Context, namespace string) ([]types.DeploymentSummary, error) {
	if m.ListDeploymentsFunc != nil {
		return m.ListDeploymentsFunc(ctx, namespace)
	}
	return []types.DeploymentSummary{}, nil
}

func (m *MockNomadClient) GetDeployment(ctx context.Context, deploymentID string) (types.Deployment, error) {
	if m.GetDeploymentFunc != nil {
		return m.GetDeploymentFunc(ctx, deploymentID)
	}
	return types.Deployment{}, nil
}

func (m *MockNomadClient) ListVolumes(ctx context.Context, nodeID string, pluginID string, nextToken string, perPage int, filter string) ([]types.Volume, error) {
	if m.ListVolumesFunc != nil {
		return m.ListVolumesFunc(ctx, nodeID, pluginID, nextToken, perPage, filter)
	}
	return []types.Volume{}, nil
}

func (m *MockNomadClient) GetVolume(ctx context.Context, volumeID string) (*types.Volume, error) {
	if m.GetVolumeFunc != nil {
		return m.GetVolumeFunc(ctx, volumeID)
	}
	return nil, nil
}

func (m *MockNomadClient) DeleteVolume(ctx context.Context, volumeID string) error {
	if m.DeleteVolumeFunc != nil {
		return m.DeleteVolumeFunc(ctx, volumeID)
	}
	return nil
}

func (m *MockNomadClient) ListNodes(ctx context.Context, status string) ([]types.NodeSummary, error) {
	if m.ListNodesFunc != nil {
		return m.ListNodesFunc(ctx, status)
	}
	return []types.NodeSummary{}, nil
}

func (m *MockNomadClient) GetNode(ctx context.Context, nodeID string) (types.Node, error) {
	if m.GetNodeFunc != nil {
		return m.GetNodeFunc(ctx, nodeID)
	}
	return types.Node{}, nil
}

func (m *MockNomadClient) DrainNode(ctx context.Context, nodeID string, enable bool, deadline int64) (string, error) {
	if m.DrainNodeFunc != nil {
		return m.DrainNodeFunc(ctx, nodeID, enable, deadline)
	}
	return "", nil
}

func (m *MockNomadClient) EligibilityNode(ctx context.Context, nodeID string, eligible string) (types.NodeSummary, error) {
	if m.EligibilityNodeFunc != nil {
		return m.EligibilityNodeFunc(ctx, nodeID, eligible)
	}
	return types.NodeSummary{}, nil
}

func (m *MockNomadClient) ListNamespaces(ctx context.Context) ([]types.Namespace, error) {
	if m.ListNamespacesFunc != nil {
		return m.ListNamespacesFunc(ctx)
	}
	return []types.Namespace{}, nil
}

func (m *MockNomadClient) CreateNamespace(ctx context.Context, namespace types.Namespace) error {
	if m.CreateNamespaceFunc != nil {
		return m.CreateNamespaceFunc(ctx, namespace)
	}
	return nil
}

func (m *MockNomadClient) DeleteNamespace(ctx context.Context, name string) error {
	if m.DeleteNamespaceFunc != nil {
		return m.DeleteNamespaceFunc(ctx, name)
	}
	return nil
}

func (m *MockNomadClient) ListAllocations(ctx context.Context, namespace, jobID string) ([]types.Allocation, error) {
	if m.ListAllocationsFunc != nil {
		return m.ListAllocationsFunc(ctx, namespace, jobID)
	}
	return []types.Allocation{}, nil
}

func (m *MockNomadClient) GetAllocation(ctx context.Context, allocID string) (types.Allocation, error) {
	if m.GetAllocationFunc != nil {
		return m.GetAllocationFunc(ctx, allocID)
	}
	return types.Allocation{}, nil
}

func (m *MockNomadClient) StopAllocation(ctx context.Context, allocID string) error {
	if m.StopAllocationFunc != nil {
		return m.StopAllocationFunc(ctx, allocID)
	}
	return nil
}

func (m *MockNomadClient) MakeRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	if m.MakeRequestFunc != nil {
		return m.MakeRequestFunc(ctx, method, path, queryParams, body)
	}
	return []byte{}, nil
}

func (m *MockNomadClient) GetAllocationLogs(ctx context.Context, allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
	if m.GetAllocationLogsFunc != nil {
		return m.GetAllocationLogsFunc(ctx, allocID, task, logType, follow, tail, offset)
	}
	return "", nil
}

func (m *MockNomadClient) ListVariables(ctx context.Context, namespace, prefix string, nextToken string, perPage int, filter string) ([]types.Variable, error) {
	if m.ListVariablesFunc != nil {
		return m.ListVariablesFunc(ctx, namespace, prefix, nextToken, perPage, filter)
	}
	return []types.Variable{}, nil
}

func (m *MockNomadClient) GetVariable(ctx context.Context, path, namespace string) (types.Variable, error) {
	if m.GetVariableFunc != nil {
		return m.GetVariableFunc(ctx, path, namespace)
	}
	return types.Variable{}, nil
}

func (m *MockNomadClient) CreateVariable(ctx context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error {
	if m.CreateVariableFunc != nil {
		return m.CreateVariableFunc(ctx, variable, namespace, cas, lockOperation)
	}
	return nil
}

func (m *MockNomadClient) DeleteVariable(ctx context.Context, path, namespace string, cas int) error {
	if m.DeleteVariableFunc != nil {
		return m.DeleteVariableFunc(ctx, path, namespace, cas)
	}
	return nil
}

func (m *MockNomadClient) ListACLTokens(ctx context.Context) ([]types.ACLToken, error) {
	if m.ListACLTokensFunc != nil {
		return m.ListACLTokensFunc(ctx)
	}
	return []types.ACLToken{}, nil
}

func (m *MockNomadClient) GetACLToken(ctx context.Context, accessorID string) (types.ACLToken, error) {
	if m.GetACLTokenFunc != nil {
		return m.GetACLTokenFunc(ctx, accessorID)
	}
	return types.ACLToken{}, nil
}

func (m *MockNomadClient) CreateACLToken(ctx context.Context, token types.ACLToken) (types.ACLToken, error) {
	if m.CreateACLTokenFunc != nil {
		return m.CreateACLTokenFunc(ctx, token)
	}
	return types.ACLToken{}, nil
}

func (m *MockNomadClient) DeleteACLToken(ctx context.Context, accessorID string) error {
	if m.DeleteACLTokenFunc != nil {
		return m.DeleteACLTokenFunc(ctx, accessorID)
	}
	return nil
}

func (m *MockNomadClient) ListACLPolicies(ctx context.Context) ([]types.ACLPolicy, error) {
	if m.ListACLPoliciesFunc != nil {
		return m.ListACLPoliciesFunc(ctx)
	}
	return []types.ACLPolicy{}, nil
}

func (m *MockNomadClient) GetACLPolicy(ctx context.Context, name string) (types.ACLPolicy, error) {
	if m.GetACLPolicyFunc != nil {
		return m.GetACLPolicyFunc(ctx, name)
	}
	return types.ACLPolicy{}, nil
}

func (m *MockNomadClient) CreateACLPolicy(ctx context.Context, policy types.ACLPolicy) error {
	if m.CreateACLPolicyFunc != nil {
		return m.CreateACLPolicyFunc(ctx, policy)
	}
	return nil
}

func (m *MockNomadClient) DeleteACLPolicy(ctx context.Context, name string) error {
	if m.DeleteACLPolicyFunc != nil {
		return m.DeleteACLPolicyFunc(ctx, name)
	}
	return nil
}

func (m *MockNomadClient) ListACLRoles(ctx context.Context) ([]types.ACLRole, error) {
	if m.ListACLRolesFunc != nil {
		return m.ListACLRolesFunc(ctx)
	}
	return []types.ACLRole{}, nil
}

func (m *MockNomadClient) GetACLRole(ctx context.Context, id string) (types.ACLRole, error) {
	if m.GetACLRoleFunc != nil {
		return m.GetACLRoleFunc(ctx, id)
	}
	return types.ACLRole{}, nil
}

func (m *MockNomadClient) CreateACLRole(ctx context.Context, role types.ACLRole) (types.ACLRole, error) {
	if m.CreateACLRoleFunc != nil {
		return m.CreateACLRoleFunc(ctx, role)
	}
	return types.ACLRole{}, nil
}

func (m *MockNomadClient) DeleteACLRole(ctx context.Context, id string) error {
	if m.DeleteACLRoleFunc != nil {
		return m.DeleteACLRoleFunc(ctx, id)
	}
	return nil
}

func (m *MockNomadClient) BootstrapACLToken(ctx context.Context) (types.ACLToken, error) {
	if m.BootstrapACLTokenFunc != nil {
		return m.BootstrapACLTokenFunc(ctx)
	}
	return types.ACLToken{}, nil
}

func (m *MockNomadClient) ListSentinelPolicies(ctx context.Context) ([]types.SentinelPolicy, error) {
	if m.ListSentinelPoliciesFunc != nil {
		return m.ListSentinelPoliciesFunc(ctx)
	}
	return []types.SentinelPolicy{}, nil
}

func (m *MockNomadClient) GetSentinelPolicy(ctx context.Context, name string) (types.SentinelPolicy, error) {
	if m.GetSentinelPolicyFunc != nil {
		return m.GetSentinelPolicyFunc(ctx, name)
	}
	return types.SentinelPolicy{}, nil
}

func (m *MockNomadClient) CreateSentinelPolicy(ctx context.Context, policy types.SentinelPolicy) error {
	if m.CreateSentinelPolicyFunc != nil {
		return m.CreateSentinelPolicyFunc(ctx, policy)
	}
	return nil
}

func (m *MockNomadClient) DeleteSentinelPolicy(ctx context.Context, name string) error {
	if m.DeleteSentinelPolicyFunc != nil {
		return m.DeleteSentinelPolicyFunc(ctx, name)
	}
	return nil
}

func (m *MockNomadClient) ListClusterPeers(ctx context.Context) ([]byte, error) {
	if m.ListClusterPeersFunc != nil {
		return m.ListClusterPeersFunc(ctx)
	}
	return []byte{}, nil
}

func (m *MockNomadClient) SetToken(token string) {
	m.token = token
}

func (m *MockNomadClient) GetToken() string {
	return m.token
}

func (m *MockNomadClient) SetDefaultTailLines(lines int) error { return nil }
func (m *MockNomadClient) GetDefaultTailLines() int            { return 100 }
