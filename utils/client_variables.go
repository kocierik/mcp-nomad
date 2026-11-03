package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/kocierik/mcp-nomad/types"
)

type variableRequestPayload struct {
	Path      string            `json:"Path,omitempty"`
	Namespace string            `json:"Namespace,omitempty"`
	Items     map[string]string `json:"Items,omitempty"`
}

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
	for i := range variables {
		hydrateLegacyVariableValue(&variables[i])
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
	hydrateLegacyVariableValue(&variable)

	return variable, nil
}

// CreateVariable creates a new variable
func (c *NomadClient) CreateVariable(ctx context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error {
	if lockOperation != "" {
		return fmt.Errorf("unsupported lock operation %q: lease payload is not supported", lockOperation)
	}

	apiPath := fmt.Sprintf("var/%s", variable.Path)

	requestBody := variableRequestPayload{}
	if variable.Value != "" {
		if err := json.Unmarshal([]byte(variable.Value), &requestBody); err != nil {
			return fmt.Errorf("failed to parse variable value: %v", err)
		}
	}

	if variable.Path != "" {
		requestBody.Path = variable.Path
	}
	if namespace != "" {
		requestBody.Namespace = namespace
	} else if variable.Namespace != "" {
		requestBody.Namespace = variable.Namespace
	}

	if len(variable.Items) > 0 {
		requestBody.Items = variable.Items
	}
	if len(requestBody.Items) == 0 {
		return fmt.Errorf("variable value is required")
	}

	queryParams := make(map[string]string)
	AddNomadNamespaceQuery(queryParams, namespace)

	if variable.CAS != nil {
		if *variable.CAS > uint64(math.MaxInt) {
			return fmt.Errorf("cas exceeds supported range")
		}
		queryParams["cas"] = strconv.FormatUint(*variable.CAS, 10)
	} else if cas > 0 {
		queryParams["cas"] = strconv.Itoa(cas)
	}

	_, err := c.makeRequest(ctx, "PUT", apiPath, queryParams, requestBody)
	return err
}

func hydrateLegacyVariableValue(variable *types.Variable) {
	if variable == nil || variable.Value != "" || len(variable.Items) == 0 {
		return
	}
	payload, err := json.Marshal(map[string]map[string]string{"Items": variable.Items})
	if err != nil {
		return
	}
	variable.Value = string(payload)
}

// DeleteVariable deletes a variable by path
func (c *NomadClient) DeleteVariable(ctx context.Context, path, namespace string, cas int) error {
	if cas > 0 {
		casValue := uint64(cas)
		return c.DeleteVariableCAS(ctx, path, namespace, &casValue)
	}
	return c.DeleteVariableCAS(ctx, path, namespace, nil)
}

// DeleteVariableCAS deletes a variable by path with an optional exact CAS query parameter.
func (c *NomadClient) DeleteVariableCAS(ctx context.Context, path, namespace string, cas *uint64) error {
	apiPath := fmt.Sprintf("var/%s", path)

	queryParams := make(map[string]string)
	AddNomadNamespaceQuery(queryParams, namespace)
	if cas != nil {
		queryParams["cas"] = strconv.FormatUint(*cas, 10)
	}

	_, err := c.makeRequest(ctx, "DELETE", apiPath, queryParams, nil)
	return err
}
