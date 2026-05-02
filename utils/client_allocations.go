package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kocierik/mcp-nomad/types"
)

// GetAllocation returns the details of an allocation
func (c *NomadClient) GetAllocation(ctx context.Context, allocID string) (types.Allocation, error) {
	path := fmt.Sprintf("allocation/%s", allocID)

	var alloc types.Allocation
	err := c.get(ctx, path, &alloc)
	if err != nil {
		return types.Allocation{}, err
	}

	return alloc, nil
}

// ListAllocations lists allocations via GET /v1/allocations (namespace optional) when jobID is empty.
// When jobID is non-empty, it uses GET /v1/job/:job_id/allocations for that namespace (consistent with Nomad API).
func (c *NomadClient) ListAllocations(ctx context.Context, namespace, jobID string) ([]types.Allocation, error) {
	if strings.TrimSpace(jobID) != "" {
		return c.ListJobAllocations(ctx, jobID, namespace)
	}

	queryParams := make(map[string]string)
	AddNomadNamespaceQuery(queryParams, namespace)

	respBody, err := c.makeRequest(ctx, "GET", "allocations", queryParams, nil)
	if err != nil {
		return nil, err
	}

	var allocations []types.Allocation
	if err := json.Unmarshal(respBody, &allocations); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return allocations, nil
}

// StopAllocation stops a running allocation (POST /v1/allocation/:id/stop).
func (c *NomadClient) StopAllocation(ctx context.Context, allocationID string) error {
	allocationID = strings.TrimSpace(allocationID)
	if allocationID == "" {
		return fmt.Errorf("allocation ID is required")
	}
	path := fmt.Sprintf("allocation/%s/stop", allocationID)
	_, err := c.makeRequest(ctx, "POST", path, nil, nil)
	return err
}
