package utils

import "context"

// GetClusterLeader return the info of the cluster leader
func (c *NomadClient) GetClusterLeader(ctx context.Context) ([]byte, error) {
	respBody, err := c.makeRequest(ctx, "GET", "operator/raft/configuration", nil, nil)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// ListClusterPeers return the list of the cluster nodes
func (c *NomadClient) ListClusterPeers(ctx context.Context) ([]byte, error) {
	respBody, err := c.makeRequest(ctx, "GET", "operator/raft/configuration", nil, nil)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// ListRegions return the regions listed
func (c *NomadClient) ListRegions(ctx context.Context) ([]byte, error) {
	return c.MakeRequest(ctx, "GET", "regions", nil, nil)
}
