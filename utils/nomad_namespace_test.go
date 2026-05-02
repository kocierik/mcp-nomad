package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNomadNamespaceQuery(t *testing.T) {
	t.Parallel()
	q := make(map[string]string)
	AddNomadNamespaceQuery(q, "")
	AddNomadNamespaceQuery(q, NomadDefaultNamespace)
	assert.Empty(t, q)

	AddNomadNamespaceQuery(q, "staging")
	assert.Equal(t, "staging", q["namespace"])
}

func TestJobVersionsNamespaceQuerySuffix(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", JobVersionsNamespaceQuerySuffix(""))
	assert.Equal(t, "", JobVersionsNamespaceQuerySuffix(NomadDefaultNamespace))
	assert.Equal(t, "?namespace=staging", JobVersionsNamespaceQuerySuffix("staging"))
	assert.Equal(t, "?namespace=my%2Fns", JobVersionsNamespaceQuerySuffix("my/ns"))
}

func TestEffectiveToolNamespace(t *testing.T) {
	t.Run("explicit beats env", func(t *testing.T) {
		t.Setenv("NOMAD_NAMESPACE", "from-env")
		ns := EffectiveToolNamespace(map[string]interface{}{"namespace": "explicit-ns"})
		assert.Equal(t, "explicit-ns", ns)
	})

	t.Run("trim spaces on explicit", func(t *testing.T) {
		t.Setenv("NOMAD_NAMESPACE", "")
		ns := EffectiveToolNamespace(map[string]interface{}{"namespace": "  web  "})
		assert.Equal(t, "web", ns)
	})

	t.Run("env when argument missing", func(t *testing.T) {
		t.Setenv("NOMAD_NAMESPACE", "prod")
		ns := EffectiveToolNamespace(map[string]interface{}{})
		assert.Equal(t, "prod", ns)
	})

	t.Run("default without env", func(t *testing.T) {
		t.Setenv("NOMAD_NAMESPACE", "")
		ns := EffectiveToolNamespace(map[string]interface{}{})
		assert.Equal(t, NomadDefaultNamespace, ns)
	})

	t.Run("explicit default skips env per-key", func(t *testing.T) {
		t.Setenv("NOMAD_NAMESPACE", "other")
		ns := EffectiveToolNamespace(map[string]interface{}{"namespace": NomadDefaultNamespace})
		assert.Equal(t, NomadDefaultNamespace, ns)
	})
}
