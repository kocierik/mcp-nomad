package utils

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/kocierik/mcp-nomad/types"
)

// ListVolumes lists all host volumes
func (c *NomadClient) ListVolumes(ctx context.Context, nodeID string, pluginID string, nextToken string, perPage int, filter string) ([]types.Volume, error) {
	path := "volumes"
	query := url.Values{}
	if nodeID != "" {
		query.Set("node_id", nodeID)
	}
	if pluginID != "" {
		query.Set("plugin_id", pluginID)
	}
	if nextToken != "" {
		query.Set("next_token", nextToken)
	}
	if perPage > 0 {
		query.Set("per_page", strconv.Itoa(perPage))
	}
	if filter != "" {
		query.Set("filter", filter)
	}

	var volumes []types.Volume
	if err := c.get(ctx, path+"?"+query.Encode(), &volumes); err != nil {
		return nil, fmt.Errorf("error listing volumes: %v", err)
	}

	return volumes, nil
}

// GetVolume retrieves a specific host volume
func (c *NomadClient) GetVolume(ctx context.Context, volumeID string) (*types.Volume, error) {
	path := fmt.Sprintf("/v1/volume/host/%s", volumeID)
	var volume types.Volume
	if err := c.get(ctx, path, &volume); err != nil {
		return nil, fmt.Errorf("error getting volume: %v", err)
	}

	return &volume, nil
}

// DeleteVolume deletes a host volume
func (c *NomadClient) DeleteVolume(ctx context.Context, volumeID string) error {
	path := fmt.Sprintf("/v1/volume/host/%s/delete", volumeID)
	if err := c.delete(ctx, path); err != nil {
		return fmt.Errorf("error deleting volume: %v", err)
	}

	return nil
}

// ListVolumeClaims lists all volume claims
func (c *NomadClient) ListVolumeClaims(ctx context.Context, namespace string, claimID string, jobID string, taskGroup string, volumeName string, nextToken string, perPage int) ([]types.VolumeClaim, error) {
	path := "volumes/"
	query := url.Values{}
	query.Set("namespace", namespace)

	if claimID != "" {
		query.Set("claim_id", claimID)
	}
	if jobID != "" {
		query.Set("job_id", jobID)
	}
	if taskGroup != "" {
		query.Set("task_group", taskGroup)
	}
	if volumeName != "" {
		query.Set("volume_name", volumeName)
	}
	if nextToken != "" {
		query.Set("next_token", nextToken)
	}
	if perPage > 0 {
		query.Set("per_page", strconv.Itoa(perPage))
	}

	var claims []types.VolumeClaim
	if err := c.get(ctx, path+"?"+query.Encode(), &claims); err != nil {
		return nil, fmt.Errorf("error listing volume claims: %v", err)
	}

	return claims, nil
}

// DeleteVolumeClaim deletes a volume claim
func (c *NomadClient) DeleteVolumeClaim(ctx context.Context, claimID string) error {
	path := fmt.Sprintf("/v1/volumes/claim/%s", claimID)
	if err := c.delete(ctx, path); err != nil {
		return fmt.Errorf("error deleting volume claim: %v", err)
	}

	return nil
}
