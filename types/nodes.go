// File: types/nodes.go
package types

// NodeSummary represents a summary of a Nomad node
type NodeSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Datacenter string `json:"datacenter"`
	NodeClass  string `json:"NodeClass"`
	NodePool   string `json:"NodePool"`
}

// Node represents a detailed view of a Nomad node
type Node struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	Status     string                `json:"status"`
	Datacenter string                `json:"datacenter"`
	Drain      bool                  `json:"drain"`
	Drivers    map[string]DriverInfo `json:"drivers"`
	Resources  NodeResources         `json:"resources"`
	Reserved   NodeResources         `json:"reserved"`
	NodeClass  string                `json:"NodeClass"`
	NodePool   string                `json:"NodePool"`
	Meta       map[string]string     `json:"meta"`
}

// DriverInfo represents a node's fingerprinted task-driver status.
// Nomad's /v1/node/:id returns each driver as an object, not a bool — the
// previous map[string]bool value type made get_node fail to decode.
type DriverInfo struct {
	Detected          bool   `json:"Detected"`
	Healthy           bool   `json:"Healthy"`
	HealthDescription string `json:"HealthDescription"`
}

// NodeResources represents the resources of a node
type NodeResources struct {
	CPU      int `json:"CPU"`
	MemoryMB int `json:"MemoryMB"`
	DiskMB   int `json:"DiskMB"`
}
