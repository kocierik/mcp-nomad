package utils

import (
	"context"
	"encoding/json"
	"fmt"

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

// ListAllocations lists all allocations in the cluster
func (c *NomadClient) ListAllocations(ctx context.Context) ([]types.Allocation, error) {
	respBody, err := c.makeRequest(ctx, "GET", "allocations", nil, nil)
	if err != nil {
		return nil, err
	}

	var allocations []types.Allocation
	if err := json.Unmarshal(respBody, &allocations); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return allocations, nil
}
