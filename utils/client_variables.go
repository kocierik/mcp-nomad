package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kocierik/mcp-nomad/types"
)

// ListVariables lists variables in the specified namespace
func (c *NomadClient) ListVariables(ctx context.Context, namespace, prefix string, nextToken string, perPage int, filter string) ([]types.Variable, error) {
	path := "vars"

	queryParams := make(map[string]string)
	AddNomadNamespaceQuery(queryParams, namespace)
	if prefix != "" {
		queryParams["prefix"] = prefix
	}
	if nextToken != "" {
		queryParams["next_token"] = nextToken
	}
	if perPage > 0 {
		queryParams["per_page"] = strconv.Itoa(perPage)
	}
	if filter != "" {
		queryParams["filter"] = filter
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var variables []types.Variable
	if err := json.Unmarshal(respBody, &variables); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return variables, nil
}

// GetVariable retrieves a specific variable by path
func (c *NomadClient) GetVariable(ctx context.Context, path, namespace string) (types.Variable, error) {
	apiPath := fmt.Sprintf("var/%s", path)

	queryParams := make(map[string]string)
	AddNomadNamespaceQuery(queryParams, namespace)

	respBody, err := c.makeRequest(ctx, "GET", apiPath, queryParams, nil)
	if err != nil {
		return types.Variable{}, err
	}

	var variable types.Variable
	if err := json.Unmarshal(respBody, &variable); err != nil {
		return types.Variable{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return variable, nil
}

// CreateVariable creates a new variable
func (c *NomadClient) CreateVariable(ctx context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error {
	apiPath := fmt.Sprintf("var/%s", variable.Path)

	// Parse the Value string into a map to use as request body
	var requestBody map[string]interface{}
	if err := json.Unmarshal([]byte(variable.Value), &requestBody); err != nil {
		return fmt.Errorf("failed to parse variable value: %v", err)
	}

	// Add CAS if provided
	if cas > 0 {
		requestBody["CAS"] = cas
	}

	// Add lock operation if provided
	if lockOperation != "" {
		requestBody["LockOperation"] = lockOperation
	}

	// Add namespace as query parameter if provided
	queryParams := make(map[string]string)
	AddNomadNamespaceQuery(queryParams, namespace)

	_, err := c.makeRequest(ctx, "PUT", apiPath, queryParams, requestBody)
	return err
}

// DeleteVariable deletes a variable by path
func (c *NomadClient) DeleteVariable(ctx context.Context, path, namespace string, cas int) error {
	apiPath := fmt.Sprintf("var/%s", path)

	queryParams := make(map[string]string)
	AddNomadNamespaceQuery(queryParams, namespace)
	if cas > 0 {
		queryParams["cas"] = strconv.Itoa(cas)
	}

	_, err := c.makeRequest(ctx, "DELETE", apiPath, queryParams, nil)
	return err
}
