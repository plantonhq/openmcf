package pulumigoogleprovider

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validServiceAccountKey carries the four fields validateServiceAccountKey requires plus a
// PEM-shaped private key. Not a real credential.
const validServiceAccountKey = `{
	"type": "service_account",
	"project_id": "test-project",
	"private_key": "-----BEGIN PRIVATE KEY-----\nfake\n-----END PRIVATE KEY-----\n",
	"client_email": "test@test-project.iam.gserviceaccount.com"
}`

const testWifProvider = "//iam.googleapis.com/projects/123456/locations/global/workloadIdentityPools/test-pool/providers/test-provider"

const testIdentityToken = "eyJhbGciOiJSUzI1NiJ9.payload.sig"

const testServiceAccountEmail = "provisioner@test-project.iam.gserviceaccount.com"

func webIdentityConfig() *gcpprovider.GcpProviderConfig {
	return &gcpprovider.GcpProviderConfig{
		WebIdentity: &gcpprovider.GcpWebIdentityProviderConfig{
			WebIdentityToken:    testIdentityToken,
			Audience:            testWifProvider,
			ServiceAccountEmail: testServiceAccountEmail,
		},
	}
}

// requireWebIdentityCredentialsEntry unwraps the raw keyless property map down to the single
// external_credentials entry, asserting the wire shape on the way: the field must be a
// single-element ARRAY (the raw terraform list form the provider-config encoder accepts --
// see the package doc), never the typed object form.
func requireWebIdentityCredentialsEntry(t *testing.T, inputs *providerInputs) pulumi.Map {
	t.Helper()

	require.NotNil(t, inputs.webIdentityProps)
	// The keyless form never populates typed args: exactly one form per dispatch.
	require.Nil(t, inputs.args)

	require.Len(t, inputs.webIdentityProps, 1, "keyless props must carry only externalCredentials")
	credentials, ok := inputs.webIdentityProps["externalCredentials"].(pulumi.Array)
	require.True(t, ok, "externalCredentials must be a pulumi.Array (raw list wire shape)")
	require.Len(t, credentials, 1, "external_credentials is a max-items-one block: exactly one entry")

	entry, ok := credentials[0].(pulumi.Map)
	require.True(t, ok, "the external credentials entry must be a pulumi.Map")
	return entry
}

func TestBuildProviderInputs_NilConfig_Ambient(t *testing.T) {
	inputs, err := buildProviderInputs(nil)
	require.NoError(t, err)
	require.NotNil(t, inputs.args)

	assert.Nil(t, inputs.args.Credentials)
	assert.Nil(t, inputs.args.ExternalCredentials)
	assert.Nil(t, inputs.webIdentityProps)
}

func TestBuildProviderInputs_EmptyConfig_Ambient(t *testing.T) {
	// No service-account key and no web identity -> ambient ADC chain.
	inputs, err := buildProviderInputs(&gcpprovider.GcpProviderConfig{})
	require.NoError(t, err)
	require.NotNil(t, inputs.args)

	assert.Nil(t, inputs.args.Credentials)
	assert.Nil(t, inputs.args.ExternalCredentials)
	assert.Nil(t, inputs.webIdentityProps)
}

func TestBuildProviderInputs_ServiceAccountKey(t *testing.T) {
	cfg := &gcpprovider.GcpProviderConfig{
		ServiceAccountKey: validServiceAccountKey,
	}

	inputs, err := buildProviderInputs(cfg)
	require.NoError(t, err)
	require.NotNil(t, inputs.args)

	assert.Equal(t, pulumi.String(validServiceAccountKey), inputs.args.Credentials)
	// Static and keyless are mutually exclusive.
	assert.Nil(t, inputs.args.ExternalCredentials)
	assert.Nil(t, inputs.webIdentityProps)
}

func TestBuildProviderInputs_ServiceAccountKey_InvalidJson_Errors(t *testing.T) {
	_, err := buildProviderInputs(&gcpprovider.GcpProviderConfig{
		ServiceAccountKey: "not-json",
	})
	assert.Error(t, err)
}

func TestBuildProviderInputs_ServiceAccountKey_MissingField_Errors(t *testing.T) {
	_, err := buildProviderInputs(&gcpprovider.GcpProviderConfig{
		ServiceAccountKey: `{"type": "service_account", "project_id": "p"}`,
	})
	assert.Error(t, err)
}

func TestBuildProviderInputs_ServiceAccountKey_NonPemPrivateKey_Errors(t *testing.T) {
	_, err := buildProviderInputs(&gcpprovider.GcpProviderConfig{
		ServiceAccountKey: `{
			"type": "service_account",
			"project_id": "test-project",
			"private_key": "MIIEvQIBADANBgkqhkiG9w0BAQEFAASC",
			"client_email": "test@test-project.iam.gserviceaccount.com"
		}`,
	})
	assert.Error(t, err)
}

func TestBuildProviderInputs_WebIdentity_RawArrayWireShape(t *testing.T) {
	inputs, err := buildProviderInputs(webIdentityConfig())
	require.NoError(t, err)

	entry := requireWebIdentityCredentialsEntry(t, inputs)

	// Exactly the three provider-schema keys, nothing else.
	require.Len(t, entry, 3)

	// The audience must be passed through verbatim (byte-identity with the token's `aud`).
	assert.Equal(t, pulumi.String(testWifProvider), entry["audience"])
	assert.Equal(t, pulumi.String(testServiceAccountEmail), entry["serviceAccountEmail"])
}

func TestBuildProviderInputs_WebIdentity_IdentityTokenIsSecretWrapped(t *testing.T) {
	inputs, err := buildProviderInputs(webIdentityConfig())
	require.NoError(t, err)

	entry := requireWebIdentityCredentialsEntry(t, inputs)

	// The SDK does NOT auto-secret-wrap identity_token, so the builder must: the wrapped
	// value is a secret Output, no longer the plain pulumi.String.
	require.NotNil(t, entry["identityToken"])
	assert.NotEqual(t, pulumi.String(testIdentityToken), entry["identityToken"])
	_, isPlainString := entry["identityToken"].(pulumi.String)
	assert.False(t, isPlainString, "identityToken must be a ToSecret-wrapped Output, not a plain string")
}

func TestBuildProviderInputs_WebIdentity_TakesPrecedenceOverStaleKey(t *testing.T) {
	// A config carrying both dispatches to keyless: web identity is the deliberate mode
	// switch, a lingering service_account_key must not win.
	cfg := webIdentityConfig()
	cfg.ServiceAccountKey = validServiceAccountKey

	inputs, err := buildProviderInputs(cfg)
	require.NoError(t, err)

	requireWebIdentityCredentialsEntry(t, inputs)
}

func TestBuildProviderInputs_WebIdentity_MissingToken_Errors(t *testing.T) {
	cfg := webIdentityConfig()
	cfg.WebIdentity.WebIdentityToken = ""

	_, err := buildProviderInputs(cfg)
	assert.Error(t, err)
}

func TestBuildProviderInputs_WebIdentity_MissingAudience_Errors(t *testing.T) {
	cfg := webIdentityConfig()
	cfg.WebIdentity.Audience = ""

	_, err := buildProviderInputs(cfg)
	assert.Error(t, err)
}

func TestBuildProviderInputs_WebIdentity_MissingServiceAccountEmail_Errors(t *testing.T) {
	cfg := webIdentityConfig()
	cfg.WebIdentity.ServiceAccountEmail = ""

	_, err := buildProviderInputs(cfg)
	assert.Error(t, err)
}

func TestGcpPluginVersion_NonEmptySemver(t *testing.T) {
	version := gcpPluginVersion()
	require.NotEmpty(t, version)
	assert.NotEqual(t, "v", version[:1], "plugin version must be a bare semver, no v prefix")
}

// TestGcpPluginVersion_MatchesSdkPin guards fallbackGcpPluginVersion against go.mod drift: the
// repo's pulumi-gcp pin must equal the fallback constant, so an SDK bump fails here until the
// constant is updated to match. (Test binaries embed no dependency build info, so the guard
// reads go.mod directly; in sandboxed build modes without go.mod on disk it skips.)
func TestGcpPluginVersion_MatchesSdkPin(t *testing.T) {
	goMod, err := findRepoGoMod()
	if err != nil {
		t.Skipf("go.mod not reachable in this build mode; drift guard not applicable: %v", err)
	}

	content, err := os.ReadFile(goMod)
	require.NoError(t, err)

	pinPattern := regexp.MustCompile(regexp.QuoteMeta(gcpSdkModulePath) + `\s+v(\S+)`)
	match := pinPattern.FindSubmatch(content)
	require.NotNil(t, match, "pulumi-gcp SDK pin not found in %s", goMod)

	assert.Equal(t, fallbackGcpPluginVersion, string(match[1]),
		"fallbackGcpPluginVersion must match the pulumi-gcp pin in go.mod")
}

// findRepoGoMod ascends from the test's working directory to the module root's go.mod.
func findRepoGoMod() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no go.mod found walking up from the test working directory")
		}
		dir = parent
	}
}

func TestProviderResourceName(t *testing.T) {
	// State continuity: the base name must stay "google".
	assert.Equal(t, "google", ProviderResourceName(nil))
	assert.Equal(t, "google-replica", ProviderResourceName([]string{"replica"}))
	assert.Equal(t, "google-dns-zone", ProviderResourceName([]string{"dns", "zone"}))
}
