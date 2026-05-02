package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kocierik/mcp-nomad/types"
)

// ListNamespaces lists all namespaces
func (c *NomadClient) ListNamespaces(ctx context.Context) ([]types.Namespace, error) {
	respBody, err := c.makeRequest(ctx, "GET", "namespaces", nil, nil)
	if err != nil {
		return nil, err
	}

	var namespaces []types.Namespace
	if err := json.Unmarshal(respBody, &namespaces); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return namespaces, nil
}

// CreateNamespace creates a new namespace
func (c *NomadClient) CreateNamespace(ctx context.Context, namespace types.Namespace) error {
	_, err := c.makeRequest(ctx, "POST", "namespace", nil, namespace)
	return err
}

// DeleteNamespace deletes a namespace
func (c *NomadClient) DeleteNamespace(ctx context.Context, name string) error {
	path := fmt.Sprintf("namespace/%s", name)
	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}
