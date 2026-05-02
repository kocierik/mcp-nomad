package utils

import (
	"context"
	"fmt"
	"strings"
)

// GetAllocationLogs retrieves logs from a specific task in an allocation
func (c *NomadClient) GetAllocationLogs(ctx context.Context, allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
	if allocID == "" {
		return "", fmt.Errorf("allocation ID is required")
	}
	if task == "" {
		return "", fmt.Errorf("task name is required")
	}

	// Set default log type if not specified
	if logType == "" {
		logType = "stdout"
	}

	// Build query parameters
	queryParams := map[string]string{
		"task":   task,
		"type":   logType,
		"follow": fmt.Sprintf("%v", follow),
		"plain":  "true",
	}

	// If tail is specified, we want to read from the end
	if tail > 0 {
		queryParams["origin"] = "end"
		// Estimate bytes needed for tail lines (assume average 200 bytes per line)
		estimatedBytes := tail * 200
		queryParams["offset"] = fmt.Sprintf("%d", estimatedBytes)
	} else if offset > 0 {
		queryParams["offset"] = fmt.Sprintf("%d", offset)
	}

	// Make request to Nomad API
	path := fmt.Sprintf("client/fs/logs/%s", allocID)
	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get allocation logs: %v", err)
	}

	// If tail was specified, we need to process the response to get the correct number of lines
	if tail > 0 {
		lines := strings.Split(string(respBody), "\n")
		if len(lines) > int(tail) {
			// Take only the last 'tail' lines
			lines = lines[len(lines)-int(tail):]
		}
		return strings.Join(lines, "\n"), nil
	}

	return string(respBody), nil
}
