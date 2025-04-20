// File: types/nodes.go
package types

// NodeSummary represents a summary of a Nomad node
type NodeSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Datacenter string `json:"datacenter"`
	NodeClass  string `json:"node_class"`
}

// Node represents a detailed view of a Nomad node
type Node struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Status     string            `json:"status"`
	Datacenter string            `json:"datacenter"`
	Drain      bool              `json:"drain"`
	Drivers    map[string]bool   `json:"drivers"`
	Resources  NodeResources     `json:"resources"`
	Reserved   NodeResources     `json:"reserved"`
	NodeClass  string            `json:"node_class"`
	Meta       map[string]string `json:"meta"`
}

// NodeResources represents the resources of a node
type NodeResources struct {
	CPU      int `json:"cpu"`
	MemoryMB int `json:"memory_mb"`
	DiskMB   int `json:"disk_mb"`
}
