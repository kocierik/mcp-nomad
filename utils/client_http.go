package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// normalizeAPIPath strips a leading "/" and redundant "v1/" — makeRequest always adds /v1/.
func normalizeAPIPath(p string) string {
	p = strings.TrimPrefix(strings.TrimSpace(p), "/")
	p = strings.TrimPrefix(p, "v1/")
	return p
}

func applyRegionFromEnvironment(query url.Values, queryKeys map[string]bool) {
	// Nomad forwards cross-region RPC when "region" is set (REST query param).
	// Mirrors NOMAD_REGION used by Nomad CLI; see Nomad HTTP API docs.
	if query.Has("region") || queryKeys["region"] {
		return
	}
	if reg := strings.TrimSpace(os.Getenv("NOMAD_REGION")); reg != "" {
		query.Set("region", reg)
	}
}

// makeRequest is a helper function to make HTTP requests to the Nomad API.
func (c *NomadClient) makeRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	rel := normalizeAPIPath(path)
	base := strings.TrimSuffix(c.address, "/")
	baseURL := fmt.Sprintf("%s/v1/%s", base, rel)

	query := url.Values{}
	queryKeysPresent := map[string]bool{}
	for key, value := range queryParams {
		queryKeysPresent[key] = true
		query.Set(key, value)
	}
	applyRegionFromEnvironment(query, queryKeysPresent)

	if encoded := query.Encode(); encoded != "" {
		baseURL = fmt.Sprintf("%s?%s", baseURL, encoded)
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add ACL token to headers if available
	if c.token != "" {
		req.Header.Set("X-Nomad-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, NewNomadHTTPError(resp.StatusCode, method, rel, respBody)
	}

	return respBody, nil
}

// MakeRequest is a helper function to make HTTP requests to the Nomad API
func (c *NomadClient) MakeRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	return c.makeRequest(ctx, method, path, queryParams, body)
}

// Helper methods for HTTP requests
func (c *NomadClient) get(ctx context.Context, path string, result interface{}) error {
	respBody, err := c.makeRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, result)
}

func (c *NomadClient) delete(ctx context.Context, path string) error {
	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}
