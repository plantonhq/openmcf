package kubetoken

import (
	"context"
	"encoding/base64"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fake static credentials in valid AWS shape; presigning is pure local signing, so
// these never reach a network endpoint.
const (
	testAccessKeyID     = "AKIAIOSFODNN7EXAMPLE"
	testSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
)

func mintTestToken(t *testing.T, opts EksTokenOptions) Token {
	t.Helper()
	token, err := MintEksToken(context.Background(), opts)
	require.NoError(t, err)
	return token
}

// TestMintEksToken_TokenShape guards the exact encoding EKS mandates: the k8s-aws-v1.
// prefix followed by a base64url (unpadded) presigned STS GetCallerIdentity URL.
func TestMintEksToken_TokenShape(t *testing.T) {
	token := mintTestToken(t, EksTokenOptions{
		ClusterName:     "demo-cluster",
		Region:          "us-east-1",
		AccessKeyID:     testAccessKeyID,
		SecretAccessKey: testSecretAccessKey,
	})

	require.True(t, strings.HasPrefix(token.Value, "k8s-aws-v1."),
		"EKS tokens must carry the k8s-aws-v1. prefix")

	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(token.Value, "k8s-aws-v1."))
	require.NoError(t, err, "the token payload must be unpadded base64url")

	presigned, err := url.Parse(string(decoded))
	require.NoError(t, err)

	assert.Equal(t, "https", presigned.Scheme)
	assert.Equal(t, "sts.us-east-1.amazonaws.com", presigned.Host,
		"the token must be signed against the cluster region's STS endpoint")

	query := presigned.Query()
	assert.Equal(t, "GetCallerIdentity", query.Get("Action"))
	assert.Equal(t, "60", query.Get("X-Amz-Expires"))
	assert.Contains(t, query.Get("X-Amz-Credential"), testAccessKeyID)

	// The cluster-binding header must be part of the signature -- it is how the EKS
	// API server ties the token to itself and rejects tokens minted for other clusters.
	assert.Contains(t, strings.ToLower(query.Get("X-Amz-SignedHeaders")), "x-k8s-aws-id")
}

// TestMintEksToken_HonestExpiry guards the reported validity window: the EKS server
// honors tokens for ~15 minutes; reporting now+14m keeps a refresh margin.
func TestMintEksToken_HonestExpiry(t *testing.T) {
	before := time.Now()
	token := mintTestToken(t, EksTokenOptions{
		ClusterName:     "demo-cluster",
		Region:          "us-east-1",
		AccessKeyID:     testAccessKeyID,
		SecretAccessKey: testSecretAccessKey,
	})

	assert.WithinDuration(t, before.Add(14*time.Minute), token.ExpiresAt, 30*time.Second)
}

// TestMintEksToken_SecretNeverInToken: the secret access key signs the URL but must
// never appear in it (only the access key ID travels, inside X-Amz-Credential).
func TestMintEksToken_SecretNeverInToken(t *testing.T) {
	token := mintTestToken(t, EksTokenOptions{
		ClusterName:     "demo-cluster",
		Region:          "us-east-1",
		AccessKeyID:     testAccessKeyID,
		SecretAccessKey: testSecretAccessKey,
	})

	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(token.Value, "k8s-aws-v1."))
	require.NoError(t, err)
	assert.NotContains(t, string(decoded), testSecretAccessKey)
}

func TestMintEksToken_RequiresClusterIdentity(t *testing.T) {
	_, err := MintEksToken(context.Background(), EksTokenOptions{Region: "us-east-1"})
	assert.Error(t, err, "missing cluster name must be rejected")

	_, err = MintEksToken(context.Background(), EksTokenOptions{ClusterName: "demo"})
	assert.Error(t, err, "missing region must be rejected")
}

// TestMintEksToken_SessionTokenTravelsAsSecurityToken: temporary (ASIA) credentials
// include the session token in the signed URL as X-Amz-Security-Token.
func TestMintEksToken_SessionTokenTravelsAsSecurityToken(t *testing.T) {
	token := mintTestToken(t, EksTokenOptions{
		ClusterName:     "demo-cluster",
		Region:          "us-east-1",
		AccessKeyID:     "ASIAIOSFODNN7EXAMPLE",
		SecretAccessKey: testSecretAccessKey,
		SessionToken:    "test-session-token",
	})

	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(token.Value, "k8s-aws-v1."))
	require.NoError(t, err)
	presigned, err := url.Parse(string(decoded))
	require.NoError(t, err)
	assert.Equal(t, "test-session-token", presigned.Query().Get("X-Amz-Security-Token"))
}
