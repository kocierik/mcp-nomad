package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kocierik/mcp-nomad/types"
)

// ListSentinelPolicies lists all Sentinel policies
func (c *NomadClient) ListSentinelPolicies(ctx context.Context) ([]types.SentinelPolicy, error) {
	respBody, err := c.makeRequest(ctx, "GET", "sentinel/policies", nil, nil)
	if err != nil {
		return nil, err
	}

	var policies []types.SentinelPolicy
	if err := json.Unmarshal(respBody, &policies); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policies, nil
}

// GetSentinelPolicy retrieves a specific Sentinel policy by name
func (c *NomadClient) GetSentinelPolicy(ctx context.Context, name string) (types.SentinelPolicy, error) {
	path := fmt.Sprintf("sentinel/policy/%s", name)

	respBody, err := c.makeRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return types.SentinelPolicy{}, err
	}

	var policy types.SentinelPolicy
	if err := json.Unmarshal(respBody, &policy); err != nil {
		return types.SentinelPolicy{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policy, nil
}

// CreateSentinelPolicy creates a new Sentinel policy
func (c *NomadClient) CreateSentinelPolicy(ctx context.Context, policy types.SentinelPolicy) error {
	path := fmt.Sprintf("sentinel/policy/%s", policy.Name)
	_, err := c.makeRequest(ctx, "POST", path, nil, policy)
	return err
}

// DeleteSentinelPolicy deletes a Sentinel policy
func (c *NomadClient) DeleteSentinelPolicy(ctx context.Context, name string) error {
	path := fmt.Sprintf("sentinel/policy/%s", name)
	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}
