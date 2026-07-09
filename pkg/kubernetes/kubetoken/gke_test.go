package kubetoken

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeServiceAccountKeyJSON builds a syntactically valid Google service-account key
// with a freshly generated RSA key, so JWT signing works without any real GCP secret.
func fakeServiceAccountKeyJSON(t *testing.T) string {
	t.Helper()

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKey),
	})

	key, err := json.Marshal(map[string]string{
		"type":         "service_account",
		"project_id":   "test-project",
		"client_email": "test-sa@test-project.iam.gserviceaccount.com",
		"private_key":  string(keyPem),
	})
	require.NoError(t, err)
	return string(key)
}

// TestMintGkeToken_ExchangesKeyForAccessToken proves the full offline loop: parse the
// key, sign the JWT, exchange it at the (fake) token endpoint, and surface the access
// token with the endpoint-reported expiry.
func TestMintGkeToken_ExchangesKeyForAccessToken(t *testing.T) {
	tokenEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		// The oauth2 jwt flow presents the signed service-account JWT as an assertion.
		assert.NotEmpty(t, r.Form.Get("assertion"))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"test-access-token","token_type":"Bearer","expires_in":3600}`)
	}))
	defer tokenEndpoint.Close()

	before := time.Now()
	token, err := MintGkeToken(context.Background(), GkeTokenOptions{
		ServiceAccountKeyJSON: fakeServiceAccountKeyJSON(t),
		TokenURL:              tokenEndpoint.URL,
	})
	require.NoError(t, err)

	assert.Equal(t, "test-access-token", token.Value)
	assert.WithinDuration(t, before.Add(time.Hour), token.ExpiresAt, 30*time.Second,
		"expiry must come from the token response, never be assumed")
}

func TestMintGkeToken_RejectsMissingOrMalformedKey(t *testing.T) {
	_, err := MintGkeToken(context.Background(), GkeTokenOptions{})
	assert.Error(t, err, "an empty key must be rejected")

	_, err = MintGkeToken(context.Background(), GkeTokenOptions{ServiceAccountKeyJSON: "{not json"})
	assert.Error(t, err, "a malformed key must be rejected at parse time, before any exchange")
}
