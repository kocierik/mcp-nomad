package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeSelfSignedClientCert generates a self-signed client certificate/key pair,
// writes them to files in dir, and returns the cert path, key path, and the
// parsed leaf certificate for assertions.
func writeSelfSignedClientCert(t *testing.T, dir string) (certPath, keyPath string, leaf *x509.Certificate) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	tmpl := x509.Certificate{
		SerialNumber:          big.NewInt(0xfeedface),
		Subject:               pkix.Name{CommonName: "mcp-nomad-test-client"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	require.NoError(t, err)

	certPath = filepath.Join(dir, "client.crt")
	require.NoError(t, os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600))

	keyDER, err := x509.MarshalECPrivateKey(priv)
	require.NoError(t, err)
	keyPath = filepath.Join(dir, "client.key")
	require.NoError(t, os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0o600))

	leaf, err = x509.ParseCertificate(der)
	require.NoError(t, err)
	return certPath, keyPath, leaf
}

func TestBuildTLSConfig_MutualTLS(t *testing.T) {
	t.Run("loads client cert and key into cfg.Certificates", func(t *testing.T) {
		certPath, keyPath, leaf := writeSelfSignedClientCert(t, t.TempDir())

		t.Setenv("NOMAD_CLIENT_CERT", certPath)
		t.Setenv("NOMAD_CLIENT_KEY", keyPath)

		cfg := buildTLSConfig()

		require.Len(t, cfg.Certificates, 1, "expected exactly one client certificate to be loaded")
		require.NotNil(t, cfg.Certificates[0].Leaf, "loaded certificate leaf should be populated")
		assert.Equal(t, leaf.SerialNumber, cfg.Certificates[0].Leaf.SerialNumber,
			"the certificate loaded into cfg.Certificates must be the one named by NOMAD_CLIENT_CERT")
		assert.Equal(t, "mcp-nomad-test-client", cfg.Certificates[0].Leaf.Subject.CommonName)
	})

	t.Run("cert without key leaves cfg.Certificates empty", func(t *testing.T) {
		certPath, _, _ := writeSelfSignedClientCert(t, t.TempDir())

		t.Setenv("NOMAD_CLIENT_CERT", certPath)
		t.Setenv("NOMAD_CLIENT_KEY", "")

		cfg := buildTLSConfig()
		assert.Empty(t, cfg.Certificates)
	})

	t.Run("key without cert leaves cfg.Certificates empty", func(t *testing.T) {
		_, keyPath, _ := writeSelfSignedClientCert(t, t.TempDir())

		t.Setenv("NOMAD_CLIENT_CERT", "")
		t.Setenv("NOMAD_CLIENT_KEY", keyPath)

		cfg := buildTLSConfig()
		assert.Empty(t, cfg.Certificates)
	})

	t.Run("neither set leaves cfg.Certificates empty", func(t *testing.T) {
		t.Setenv("NOMAD_CLIENT_CERT", "")
		t.Setenv("NOMAD_CLIENT_KEY", "")

		cfg := buildTLSConfig()
		assert.Empty(t, cfg.Certificates)
	})

	t.Run("unreadable cert file is ignored without loading a certificate", func(t *testing.T) {
		dir := t.TempDir()
		missing := filepath.Join(dir, "does-not-exist.crt")
		_, keyPath, _ := writeSelfSignedClientCert(t, dir)

		t.Setenv("NOMAD_CLIENT_CERT", missing)
		t.Setenv("NOMAD_CLIENT_KEY", keyPath)

		cfg := buildTLSConfig()
		assert.Empty(t, cfg.Certificates, "a load failure must be swallowed and leave no certificate set")
	})
}
