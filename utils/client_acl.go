package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kocierik/mcp-nomad/types"
)

// ListACLTokens lists all ACL tokens
func (c *NomadClient) ListACLTokens(ctx context.Context) ([]types.ACLToken, error) {
	respBody, err := c.makeRequest(ctx, "GET", "acl/tokens", nil, nil)
	if err != nil {
		return nil, err
	}

	var tokens []types.ACLToken
	if err := json.Unmarshal(respBody, &tokens); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return tokens, nil
}

// GetACLToken retrieves a specific ACL token by accessor ID
func (c *NomadClient) GetACLToken(ctx context.Context, accessorID string) (types.ACLToken, error) {
	path := fmt.Sprintf("acl/token/%s", accessorID)

	respBody, err := c.makeRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return types.ACLToken{}, err
	}

	var token types.ACLToken
	if err := json.Unmarshal(respBody, &token); err != nil {
		return types.ACLToken{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return token, nil
}

// CreateACLToken creates a new ACL token
func (c *NomadClient) CreateACLToken(ctx context.Context, token types.ACLToken) (types.ACLToken, error) {
	respBody, err := c.makeRequest(ctx, "POST", "acl/token", nil, token)
	if err != nil {
		return types.ACLToken{}, err
	}

	var newToken types.ACLToken
	if err := json.Unmarshal(respBody, &newToken); err != nil {
		return types.ACLToken{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return newToken, nil
}

// DeleteACLToken deletes an ACL token
func (c *NomadClient) DeleteACLToken(ctx context.Context, accessorID string) error {
	path := fmt.Sprintf("acl/token/%s", accessorID)
	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}

// ListACLPolicies lists all ACL policies
func (c *NomadClient) ListACLPolicies(ctx context.Context) ([]types.ACLPolicy, error) {
	respBody, err := c.makeRequest(ctx, "GET", "acl/policies", nil, nil)
	if err != nil {
		return nil, err
	}

	var policies []types.ACLPolicy
	if err := json.Unmarshal(respBody, &policies); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policies, nil
}

// GetACLPolicy retrieves a specific ACL policy by name
func (c *NomadClient) GetACLPolicy(ctx context.Context, name string) (types.ACLPolicy, error) {
	path := fmt.Sprintf("acl/policy/%s", name)

	respBody, err := c.makeRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return types.ACLPolicy{}, err
	}

	var policy types.ACLPolicy
	if err := json.Unmarshal(respBody, &policy); err != nil {
		return types.ACLPolicy{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policy, nil
}

// CreateACLPolicy creates a new ACL policy
func (c *NomadClient) CreateACLPolicy(ctx context.Context, policy types.ACLPolicy) error {
	path := fmt.Sprintf("acl/policy/%s", policy.Name)

	_, err := c.makeRequest(ctx, "POST", path, nil, policy)
	return err
}

// DeleteACLPolicy deletes an ACL policy
func (c *NomadClient) DeleteACLPolicy(ctx context.Context, name string) error {
	path := fmt.Sprintf("acl/policy/%s", name)
	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}

// ListACLRoles lists all ACL roles
func (c *NomadClient) ListACLRoles(ctx context.Context) ([]types.ACLRole, error) {
	respBody, err := c.makeRequest(ctx, "GET", "acl/roles", nil, nil)
	if err != nil {
		return nil, err
	}

	var roles []types.ACLRole
	if err := json.Unmarshal(respBody, &roles); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return roles, nil
}

// GetACLRole retrieves a specific ACL role by ID
func (c *NomadClient) GetACLRole(ctx context.Context, id string) (types.ACLRole, error) {
	path := fmt.Sprintf("acl/role/%s", id)

	respBody, err := c.makeRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return types.ACLRole{}, err
	}

	var role types.ACLRole
	if err := json.Unmarshal(respBody, &role); err != nil {
		return types.ACLRole{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return role, nil
}

// CreateACLRole creates a new ACL role
func (c *NomadClient) CreateACLRole(ctx context.Context, role types.ACLRole) (types.ACLRole, error) {
	respBody, err := c.makeRequest(ctx, "POST", "acl/role", nil, role)
	if err != nil {
		return types.ACLRole{}, err
	}

	var newRole types.ACLRole
	if err := json.Unmarshal(respBody, &newRole); err != nil {
		return types.ACLRole{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return newRole, nil
}

// DeleteACLRole deletes an ACL role
func (c *NomadClient) DeleteACLRole(ctx context.Context, id string) error {
	path := fmt.Sprintf("acl/role/%s", id)
	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}

// BootstrapACLToken bootstraps the ACL system and returns the initial management token
func (c *NomadClient) BootstrapACLToken(ctx context.Context) (types.ACLToken, error) {
	respBody, err := c.makeRequest(ctx, "POST", "acl/bootstrap", nil, nil)
	if err != nil {
		return types.ACLToken{}, err
	}

	var token types.ACLToken
	if err := json.Unmarshal(respBody, &token); err != nil {
		return types.ACLToken{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return token, nil
}
