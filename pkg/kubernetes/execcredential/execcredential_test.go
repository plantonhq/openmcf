package execcredential

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/plantonhq/planton/pkg/kubernetes/kubetoken"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEmit_ProtocolConformance guards the exact wire shape client-go and the tofu/pulumi
// kubernetes providers parse: apiVersion v1, kind ExecCredential, and an RFC3339
// status.expirationTimestamp (the honest expiry that drives re-invocation).
func TestEmit_ProtocolConformance(t *testing.T) {
	expiresAt := time.Date(2026, 7, 9, 12, 30, 0, 0, time.UTC)

	var out bytes.Buffer
	require.NoError(t, Emit(&out, kubetoken.Token{Value: "k8s-aws-v1.abc", ExpiresAt: expiresAt}))

	var doc struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Status     struct {
			Token               string `json:"token"`
			ExpirationTimestamp string `json:"expirationTimestamp"`
		} `json:"status"`
	}
	require.NoError(t, json.Unmarshal(out.Bytes(), &doc))

	assert.Equal(t, "client.authentication.k8s.io/v1", doc.APIVersion)
	assert.Equal(t, "ExecCredential", doc.Kind)
	assert.Equal(t, "k8s-aws-v1.abc", doc.Status.Token)

	parsed, err := time.Parse(time.RFC3339, doc.Status.ExpirationTimestamp)
	require.NoError(t, err, "expirationTimestamp must be RFC3339")
	assert.True(t, parsed.Equal(expiresAt))
}

// TestEmit_NonUtcExpiryNormalized: client-go accepts any RFC3339 offset, but emitting
// UTC keeps the output canonical regardless of the host's local timezone.
func TestEmit_NonUtcExpiryNormalized(t *testing.T) {
	ist := time.FixedZone("IST", 5*3600+1800)
	expiresAt := time.Date(2026, 7, 9, 18, 0, 0, 0, ist)

	var out bytes.Buffer
	require.NoError(t, Emit(&out, kubetoken.Token{Value: "tok", ExpiresAt: expiresAt}))

	var doc struct {
		Status struct {
			ExpirationTimestamp string `json:"expirationTimestamp"`
		} `json:"status"`
	}
	require.NoError(t, json.Unmarshal(out.Bytes(), &doc))

	parsed, err := time.Parse(time.RFC3339, doc.Status.ExpirationTimestamp)
	require.NoError(t, err)
	assert.True(t, parsed.Equal(expiresAt))
	assert.Equal(t, time.UTC, parsed.Location())
}
