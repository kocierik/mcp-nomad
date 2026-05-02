package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kocierik/mcp-nomad/types"
)

// ListNodes lists all nodes in the cluster
func (c *NomadClient) ListNodes(ctx context.Context, status string) ([]types.NodeSummary, error) {
	queryParams := make(map[string]string)
	if status != "" {
		queryParams["status"] = status
	}

	respBody, err := c.makeRequest(ctx, "GET", "nodes", queryParams, nil)
	if err != nil {
		return nil, err
	}

	var nodes []types.NodeSummary
	if err := json.Unmarshal(respBody, &nodes); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return nodes, nil
}

// GetNode retrieves a specific node by ID
func (c *NomadClient) GetNode(ctx context.Context, nodeID string) (types.Node, error) {
	path := fmt.Sprintf("node/%s", nodeID)

	respBody, err := c.makeRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return types.Node{}, err
	}

	var node types.Node
	if err := json.Unmarshal(respBody, &node); err != nil {
		return types.Node{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return node, nil
}

// DrainNode enables or disables drain mode for a node
func (c *NomadClient) DrainNode(ctx context.Context, nodeID string, enable bool, deadline int64) (string, error) {
	path := fmt.Sprintf("node/%s/drain", nodeID)

	drainSpec := map[string]interface{}{
		"DrainSpec": map[string]interface{}{
			"Deadline":         deadline,
			"IgnoreSystemJobs": false,
		},
		"Meta": map[string]string{
			"reason": "Initiated via API",
		},
	}

	if !enable {
		drainSpec = map[string]interface{}{
			"DrainSpec": nil,
			"Meta": map[string]string{
				"reason": "Drain disabled via API",
			},
		}
	}

	respBody, err := c.makeRequest(ctx, "POST", path, nil, drainSpec)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if enable {
		if deadline > 0 {
			return fmt.Sprintf("Node drain enabled with deadline %d seconds", deadline), nil
		}
		return "Node drain enabled with no deadline", nil
	}
	return "Node drain disabled", nil
}

// EligibilityNode sets scheduling eligibility on a node
func (c *NomadClient) EligibilityNode(ctx context.Context, nodeID string, eligible string) (types.NodeSummary, error) {
	path := fmt.Sprintf("node/%s/eligibility", nodeID)

	eligibilitySpec := map[string]interface{}{
		"Eligibility": eligible,
	}

	respBody, err := c.makeRequest(ctx, "POST", path, nil, eligibilitySpec)
	if err != nil {
		return types.NodeSummary{}, err
	}

	var nodes types.NodeSummary
	if err := json.Unmarshal(respBody, &nodes); err != nil {
		return types.NodeSummary{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return nodes, nil
}
