package prompts

import (
	"testing"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/stretchr/testify/assert"
)

func TestEffectiveNamespaceFromPrompt_explicitOverridesEnv(t *testing.T) {
	t.Setenv("NOMAD_NAMESPACE", "from-env")
	ns := effectiveNamespaceFromPrompt(map[string]string{"namespace": "explicit"})
	assert.Equal(t, "explicit", ns)
}

func TestEffectiveNamespaceFromPrompt_trimExplicit(t *testing.T) {
	t.Setenv("NOMAD_NAMESPACE", "")
	ns := effectiveNamespaceFromPrompt(map[string]string{"namespace": "  trimmed  "})
	assert.Equal(t, "trimmed", ns)
}

func TestEffectiveNamespaceFromPrompt_usesNOMAD_NAMESPACEWhenUnset(t *testing.T) {
	t.Setenv("NOMAD_NAMESPACE", "prod")

	// Only non-namespace keys — same shape as prompts passing action/job_id alone
	ns := effectiveNamespaceFromPrompt(map[string]string{"action": "list", "job_id": "web"})
	assert.Equal(t, "prod", ns)
}

func TestEffectiveNamespaceFromPrompt_defaultWithoutEnv(t *testing.T) {
	t.Setenv("NOMAD_NAMESPACE", "")
	ns := effectiveNamespaceFromPrompt(map[string]string{})
	assert.Equal(t, utils.NomadDefaultNamespace, ns)
}
