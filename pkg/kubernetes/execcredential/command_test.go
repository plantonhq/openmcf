package execcredential

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRun_AwsEksEndToEnd drives the command exactly as a deploy engine would: provider
// selection and cluster identity from the env contract, static AWS credentials under
// the SDK's standard names, protocol JSON on stdout. Fully offline -- EKS token
// minting is pure local signing.
func TestRun_AwsEksEndToEnd(t *testing.T) {
	t.Setenv(ProviderEnvVar, ProviderAwsEks)
	t.Setenv(EksClusterNameEnvVar, "demo-cluster")
	t.Setenv(EksRegionEnvVar, "us-west-2")
	t.Setenv(AwsAccessKeyIDEnvVar, "AKIAIOSFODNN7EXAMPLE")
	t.Setenv(AwsSecretAccessKeyEnvVar, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	// Blank out ambient AWS state that could shadow the test credentials.
	t.Setenv(AwsSessionTokenEnvVar, "")
	t.Setenv("AWS_PROFILE", "")

	var out bytes.Buffer
	require.NoError(t, run(context.Background(), &out))

	var doc struct {
		APIVersion string `json:"apiVersion"`
		Status     struct {
			Token               string `json:"token"`
			ExpirationTimestamp string `json:"expirationTimestamp"`
		} `json:"status"`
	}
	require.NoError(t, json.Unmarshal(out.Bytes(), &doc))

	assert.Equal(t, "client.authentication.k8s.io/v1", doc.APIVersion)
	assert.True(t, strings.HasPrefix(doc.Status.Token, "k8s-aws-v1."))
	assert.NotEmpty(t, doc.Status.ExpirationTimestamp)
}

func TestRun_RejectsMissingProvider(t *testing.T) {
	t.Setenv(ProviderEnvVar, "")

	err := run(context.Background(), &bytes.Buffer{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), ProviderEnvVar,
		"the error must name the missing contract variable")
}

func TestRun_RejectsUnknownProvider(t *testing.T) {
	t.Setenv(ProviderEnvVar, "digital_ocean_doks") // DOKS never uses exec credentials

	err := run(context.Background(), &bytes.Buffer{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "digital_ocean_doks")
}

// TestRun_AzureAksRejectsIncompleteIdentity proves the AKS arm is wired: a client
// secret without its identity coordinates fails inside the minter, before any
// network activity -- and the error names the provider being minted for.
func TestRun_AzureAksRejectsIncompleteIdentity(t *testing.T) {
	t.Setenv(ProviderEnvVar, ProviderAzureAks)
	t.Setenv(AksClientSecretEnvVar, "test-secret")
	t.Setenv(AksTenantIdEnvVar, "")
	t.Setenv(AksClientIdEnvVar, "")

	err := run(context.Background(), &bytes.Buffer{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "azure_aks")
}
