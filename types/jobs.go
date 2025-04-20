// File: types/jobs.go
package types

// JobSummary represents a summary of a Nomad job
type JobSummary struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

// Job represents a Nomad job
type Job struct {
	ID          string      `json:"id"`
	Namespace   string      `json:"namespace"`
	Status      string      `json:"status"`
	Type        string      `json:"type"`
	Datacenters []string    `json:"datacenters"`
	TaskGroups  []TaskGroup `json:"task_groups"`
}

// TaskGroup represents a task group within a job
type TaskGroup struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Tasks []Task `json:"tasks"`
}

// Task represents a task within a task group
type Task struct {
	Name      string                 `json:"name"`
	Driver    string                 `json:"driver"`
	Config    map[string]interface{} `json:"config"`
	Resources Resources              `json:"resources"`
}

// Resources represents the resources required by a task
type Resources struct {
	CPU      int `json:"cpu"`
	MemoryMB int `json:"memory_mb"`
	DiskMB   int `json:"disk_mb,omitempty"`
}
