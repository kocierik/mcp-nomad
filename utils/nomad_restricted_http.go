package utils

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// ErrMCPOutboundHTTPForbidden indicates MakeRequest refused a path/method pairing.
// Exported MakeRequest from NomadClient is restricted to MCP tool needs (see validateMCPOutboundNomadHTTP).
var ErrMCPOutboundHTTPForbidden = errors.New("nomad: HTTP method/path not permitted for MCP MakeRequest")

var postAllocationStopPath = regexp.MustCompile(`^allocation/[^/]+/stop$`)

var mcpOutboundNomadExactGETPaths = map[string]struct{}{
	"operator/raft/configuration": {},
	"regions":                     {},
}

func validateMCPOutboundNomadHTTP(method, normalizedRelPath string) error {
	m := strings.TrimSpace(strings.ToUpper(method))
	rel := strings.TrimPrefix(strings.TrimSpace(normalizedRelPath), "/")
	switch m {
	case http.MethodGet:
		if _, ok := mcpOutboundNomadExactGETPaths[rel]; ok {
			return nil
		}
		return fmt.Errorf("%w: unsupported GET path %q", ErrMCPOutboundHTTPForbidden, rel)
	case http.MethodPost:
		if postAllocationStopPath.MatchString(rel) {
			return nil
		}
		return fmt.Errorf("%w: unsupported POST path %q", ErrMCPOutboundHTTPForbidden, rel)
	default:
		return fmt.Errorf("%w: unsupported method %q", ErrMCPOutboundHTTPForbidden, m)
	}
}
