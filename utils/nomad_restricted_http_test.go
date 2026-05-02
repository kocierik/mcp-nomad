package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateMCPOutboundNomadHTTP_Allowed(t *testing.T) {
	t.Parallel()
	require.NoError(t, validateMCPOutboundNomadHTTP("GET", "regions"))
	require.NoError(t, validateMCPOutboundNomadHTTP("get", "operator/raft/configuration"))
	require.NoError(t, validateMCPOutboundNomadHTTP("POST", "allocation/abc-123/stop"))
}

func TestValidateMCPOutboundNomadHTTP_Rejects(t *testing.T) {
	t.Parallel()
	require.ErrorIs(t, validateMCPOutboundNomadHTTP("GET", "jobs"), ErrMCPOutboundHTTPForbidden)
	require.ErrorIs(t, validateMCPOutboundNomadHTTP("DELETE", "job/foo"), ErrMCPOutboundHTTPForbidden)
	require.ErrorIs(t, validateMCPOutboundNomadHTTP("POST", "job/foo/evaluate"), ErrMCPOutboundHTTPForbidden)
	require.ErrorIs(t, validateMCPOutboundNomadHTTP("POST", "allocation/foo/extra/stop"), ErrMCPOutboundHTTPForbidden)
}
