package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobVersionMethods_DecodeVersionsEnvelopeAndUseNamespaceQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(*NomadClient) (int, string, error)
	}{
		{
			name: "GetJobVersions",
			call: func(c *NomadClient) (int, string, error) {
				versions, err := c.GetJobVersions(context.Background(), "app", "apps")
				if err != nil {
					return 0, "", err
				}
				if len(versions) == 0 {
					return 0, "", nil
				}
				return len(versions), versions[0].ID, nil
			},
		},
		{
			name: "ListJobVersions",
			call: func(c *NomadClient) (int, string, error) {
				versions, err := c.ListJobVersions(context.Background(), "app", "apps")
				if err != nil {
					return 0, "", err
				}
				if len(versions) == 0 {
					return 0, "", nil
				}
				return len(versions), versions[0].ID, nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/v1/job/app/versions", r.URL.Path)
				require.Equal(t, "apps", r.URL.Query().Get("namespace"))
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"Versions":[{"ID":"app"}],"Diffs":null}`))
			}))
			defer server.Close()

			client := &NomadClient{address: server.URL, httpClient: server.Client()}
			count, id, err := tt.call(client)
			require.NoError(t, err)
			assert.Equal(t, 1, count)
			assert.Equal(t, "app", id)
		})
	}
}
