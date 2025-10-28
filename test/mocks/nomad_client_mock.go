package mocks

import (
	"github.com/kocierik/mcp-nomad/types"
)

// MockNomadClient implements a mock version of NomadClient for testing
type MockNomadClient struct {
	// Job methods
	ListJobsFunc       func(namespace, status string) ([]types.JobSummary, error)
	GetJobFunc         func(jobID, namespace string) (types.Job, error)
	RunJobFunc         func(jobSpec string, detach bool) (map[string]interface{}, error)
	StopJobFunc        func(jobID, namespace string, purge bool) (map[string]interface{}, error)
	ScaleTaskGroupFunc func(jobID, group string, count int, namespace string) error

	// Node methods
	ListNodesFunc       func(status string) ([]types.NodeSummary, error)
	GetNodeFunc         func(nodeID string) (types.Node, error)
	DrainNodeFunc       func(nodeID string, enable bool, deadline int64) (string, error)
	EligibilityNodeFunc func(nodeID string, eligible string) (types.NodeSummary, error)

	// Namespace methods
	ListNamespacesFunc  func() ([]types.Namespace, error)
	CreateNamespaceFunc func(namespace types.Namespace) error
	DeleteNamespaceFunc func(name string) error

	// Allocation methods
	ListAllocationsFunc   func() ([]types.Allocation, error)
	GetAllocationFunc     func(allocID string) (types.Allocation, error)
	GetAllocationLogsFunc func(allocID, task, logType string, follow bool, tail, offset int64) (string, error)

	// Variable methods
	ListVariablesFunc  func(namespace, prefix string, nextToken string, perPage int, filter string) ([]types.Variable, error)
	GetVariableFunc    func(path, namespace string) (types.Variable, error)
	CreateVariableFunc func(variable types.Variable, namespace string, cas int, lockOperation string) error
	DeleteVariableFunc func(path, namespace string, cas int) error

	// ACL methods
	ListACLTokensFunc  func() ([]types.ACLToken, error)
	GetACLTokenFunc    func(accessorID string) (types.ACLToken, error)
	CreateACLTokenFunc func(token types.ACLToken) (types.ACLToken, error)
	DeleteACLTokenFunc func(accessorID string) error

	// Cluster methods
	GetClusterLeaderFunc func() ([]byte, error)
	ListClusterPeersFunc func() ([]byte, error)
	ListRegionsFunc      func() ([]byte, error)
}

// Job methods
func (m *MockNomadClient) ListJobs(namespace, status string) ([]types.JobSummary, error) {
	if m.ListJobsFunc != nil {
		return m.ListJobsFunc(namespace, status)
	}
	return []types.JobSummary{}, nil
}

func (m *MockNomadClient) GetJob(jobID, namespace string) (types.Job, error) {
	if m.GetJobFunc != nil {
		return m.GetJobFunc(jobID, namespace)
	}
	return types.Job{}, nil
}

func (m *MockNomadClient) RunJob(jobSpec string, detach bool) (map[string]interface{}, error) {
	if m.RunJobFunc != nil {
		return m.RunJobFunc(jobSpec, detach)
	}
	return map[string]interface{}{}, nil
}

func (m *MockNomadClient) StopJob(jobID, namespace string, purge bool) (map[string]interface{}, error) {
	if m.StopJobFunc != nil {
		return m.StopJobFunc(jobID, namespace, purge)
	}
	return map[string]interface{}{}, nil
}

func (m *MockNomadClient) ScaleTaskGroup(jobID, group string, count int, namespace string) error {
	if m.ScaleTaskGroupFunc != nil {
		return m.ScaleTaskGroupFunc(jobID, group, count, namespace)
	}
	return nil
}

// Node methods
func (m *MockNomadClient) ListNodes(status string) ([]types.NodeSummary, error) {
	if m.ListNodesFunc != nil {
		return m.ListNodesFunc(status)
	}
	return []types.NodeSummary{}, nil
}

func (m *MockNomadClient) GetNode(nodeID string) (types.Node, error) {
	if m.GetNodeFunc != nil {
		return m.GetNodeFunc(nodeID)
	}
	return types.Node{}, nil
}

func (m *MockNomadClient) DrainNode(nodeID string, enable bool, deadline int64) (string, error) {
	if m.DrainNodeFunc != nil {
		return m.DrainNodeFunc(nodeID, enable, deadline)
	}
	return "", nil
}

func (m *MockNomadClient) EligibilityNode(nodeID string, eligible string) (types.NodeSummary, error) {
	if m.EligibilityNodeFunc != nil {
		return m.EligibilityNodeFunc(nodeID, eligible)
	}
	return types.NodeSummary{}, nil
}

// Namespace methods
func (m *MockNomadClient) ListNamespaces() ([]types.Namespace, error) {
	if m.ListNamespacesFunc != nil {
		return m.ListNamespacesFunc()
	}
	return []types.Namespace{}, nil
}

func (m *MockNomadClient) CreateNamespace(namespace types.Namespace) error {
	if m.CreateNamespaceFunc != nil {
		return m.CreateNamespaceFunc(namespace)
	}
	return nil
}

func (m *MockNomadClient) DeleteNamespace(name string) error {
	if m.DeleteNamespaceFunc != nil {
		return m.DeleteNamespaceFunc(name)
	}
	return nil
}

// Allocation methods
func (m *MockNomadClient) ListAllocations() ([]types.Allocation, error) {
	if m.ListAllocationsFunc != nil {
		return m.ListAllocationsFunc()
	}
	return []types.Allocation{}, nil
}

func (m *MockNomadClient) GetAllocation(allocID string) (types.Allocation, error) {
	if m.GetAllocationFunc != nil {
		return m.GetAllocationFunc(allocID)
	}
	return types.Allocation{}, nil
}

func (m *MockNomadClient) GetAllocationLogs(allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
	if m.GetAllocationLogsFunc != nil {
		return m.GetAllocationLogsFunc(allocID, task, logType, follow, tail, offset)
	}
	return "", nil
}

// Variable methods
func (m *MockNomadClient) ListVariables(namespace, prefix string, nextToken string, perPage int, filter string) ([]types.Variable, error) {
	if m.ListVariablesFunc != nil {
		return m.ListVariablesFunc(namespace, prefix, nextToken, perPage, filter)
	}
	return []types.Variable{}, nil
}

func (m *MockNomadClient) GetVariable(path, namespace string) (types.Variable, error) {
	if m.GetVariableFunc != nil {
		return m.GetVariableFunc(path, namespace)
	}
	return types.Variable{}, nil
}

func (m *MockNomadClient) CreateVariable(variable types.Variable, namespace string, cas int, lockOperation string) error {
	if m.CreateVariableFunc != nil {
		return m.CreateVariableFunc(variable, namespace, cas, lockOperation)
	}
	return nil
}

func (m *MockNomadClient) DeleteVariable(path, namespace string, cas int) error {
	if m.DeleteVariableFunc != nil {
		return m.DeleteVariableFunc(path, namespace, cas)
	}
	return nil
}

// ACL methods
func (m *MockNomadClient) ListACLTokens() ([]types.ACLToken, error) {
	if m.ListACLTokensFunc != nil {
		return m.ListACLTokensFunc()
	}
	return []types.ACLToken{}, nil
}

func (m *MockNomadClient) GetACLToken(accessorID string) (types.ACLToken, error) {
	if m.GetACLTokenFunc != nil {
		return m.GetACLTokenFunc(accessorID)
	}
	return types.ACLToken{}, nil
}

func (m *MockNomadClient) CreateACLToken(token types.ACLToken) (types.ACLToken, error) {
	if m.CreateACLTokenFunc != nil {
		return m.CreateACLTokenFunc(token)
	}
	return types.ACLToken{}, nil
}

func (m *MockNomadClient) DeleteACLToken(accessorID string) error {
	if m.DeleteACLTokenFunc != nil {
		return m.DeleteACLTokenFunc(accessorID)
	}
	return nil
}

// Cluster methods
func (m *MockNomadClient) GetClusterLeader() ([]byte, error) {
	if m.GetClusterLeaderFunc != nil {
		return m.GetClusterLeaderFunc()
	}
	return []byte{}, nil
}

func (m *MockNomadClient) ListClusterPeers() ([]byte, error) {
	if m.ListClusterPeersFunc != nil {
		return m.ListClusterPeersFunc()
	}
	return []byte{}, nil
}

func (m *MockNomadClient) ListRegions() ([]byte, error) {
	if m.ListRegionsFunc != nil {
		return m.ListRegionsFunc()
	}
	return []byte{}, nil
}

// Utility methods that might be needed
func (m *MockNomadClient) SetToken(token string)               {}
func (m *MockNomadClient) GetToken() string                    { return "" }
func (m *MockNomadClient) SetDefaultTailLines(lines int) error { return nil }
func (m *MockNomadClient) GetDefaultTailLines() int            { return 100 }
func (m *MockNomadClient) MakeRequest(method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	return []byte{}, nil
}
