// Package resources provides implementations of MCP resources for Nomad.
//
// The package organizes resources by category (jobs, nodes, allocations, etc.)
// and provides a ResourceManager to manage them. MCP resources provide access
// to read-only data from the Nomad cluster, which can be used as context by LLMs.
//
// Resource categories:
// - Static resources: readme, license, help documentation
// - Job resources: specification, history, allocations, evaluations
// - Node resources: status, resources, allocations
// - Allocation resources: logs, status, tasks
// - Cluster resources: metrics, leader, policies
// - Miscellaneous resources: evaluations, service health
//
// Each resource has a unique URI that follows the format:
// [protocol]://[path]
//
// Examples:
// - docs://readme
// - nomad://jobs/my-job/spec
// - nomad://nodes/abc123/status
// - nomad://allocations/xyz789/logs
// - nomad://cluster/metrics
package resources
