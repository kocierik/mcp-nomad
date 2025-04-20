// File: types/deployments.go
package types

// DeploymentSummary represents a summary of a deployment
type DeploymentSummary struct {
	ID        string `json:"id"`
	JobID     string `json:"job_id"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
}

// Deployment represents a detailed view of a deployment
type Deployment struct {
	ID         string                         `json:"id"`
	JobID      string                         `json:"job_id"`
	Namespace  string                         `json:"namespace"`
	Status     string                         `json:"status"`
	TaskGroups map[string]DeploymentTaskGroup `json:"task_groups"`
}

// DeploymentTaskGroup represents the deployment status of a task group
type DeploymentTaskGroup struct {
	DesiredTotal    int `json:"desired_total"`
	PlacedAllocs    int `json:"placed_allocs"`
	HealthyAllocs   int `json:"healthy_allocs"`
	UnhealthyAllocs int `json:"unhealthy_allocs"`
}
