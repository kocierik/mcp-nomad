package utils

import (
	"net/url"
	"os"
	"strings"
)

// NomadDefaultNamespace is the sentinel name for the Nomad cluster default namespace.
const NomadDefaultNamespace = "default"

// AddNomadNamespaceQuery sets the Nomad HTTP API `namespace` query parameter when targeting a non-default namespace.
// query must be non-nil.
func AddNomadNamespaceQuery(query map[string]string, namespace string) {
	if namespace == "" || namespace == NomadDefaultNamespace {
		return
	}
	query["namespace"] = namespace
}

// JobVersionsNamespaceQuerySuffix returns `?namespace=...` for job version endpoints when not default (URL-encoded).
func JobVersionsNamespaceQuerySuffix(namespace string) string {
	if namespace == "" || namespace == NomadDefaultNamespace {
		return ""
	}
	return "?namespace=" + url.QueryEscape(namespace)
}

// EffectiveToolNamespace resolves the namespace MCP tools send to Nomad REST calls:
// non-empty explicit "namespace" in arguments wins, then NOMAD_NAMESPACE if set, otherwise NomadDefaultNamespace.
func EffectiveToolNamespace(arguments map[string]interface{}) string {
	if arguments != nil {
		if ns, ok := arguments["namespace"].(string); ok {
			if trimmed := strings.TrimSpace(ns); trimmed != "" {
				return trimmed
			}
		}
	}
	if env := strings.TrimSpace(os.Getenv("NOMAD_NAMESPACE")); env != "" {
		return env
	}
	return NomadDefaultNamespace
}
