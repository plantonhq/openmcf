package pulumiazureprovider

import (
	"testing"

	azureprovider "github.com/plantonhq/planton/catalog/azure"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildProviderArgs_NilConfig_Ambient(t *testing.T) {
	args, err := buildProviderArgs(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, args)

	assert.Nil(t, args.ClientId)
	assert.Nil(t, args.ClientSecret)
	assert.Nil(t, args.TenantId)
	assert.Nil(t, args.SubscriptionId)
	assert.Nil(t, args.UseOidc)
	assert.Nil(t, args.OidcToken)
}

func TestBuildProviderArgs_RunnerMode_IdentityCoordinatesOnly(t *testing.T) {
	// No client secret and no web identity -> identity coordinates only (ambient chain).
	cfg := &azureprovider.AzureProviderConfig{
		ClientId:       "11111111-1111-1111-1111-111111111111",
		TenantId:       "22222222-2222-2222-2222-222222222222",
		SubscriptionId: "33333333-3333-3333-3333-333333333333",
	}

	args, err := buildProviderArgs(cfg, nil)
	require.NoError(t, err)

	assert.Equal(t, pulumi.String("11111111-1111-1111-1111-111111111111"), args.ClientId)
	assert.Equal(t, pulumi.String("22222222-2222-2222-2222-222222222222"), args.TenantId)
	assert.Equal(t, pulumi.String("33333333-3333-3333-3333-333333333333"), args.SubscriptionId)
	assert.Nil(t, args.ClientSecret)
	assert.Nil(t, args.UseOidc)
	assert.Nil(t, args.OidcToken)
}

func TestBuildProviderArgs_StaticCredentials(t *testing.T) {
	cfg := &azureprovider.AzureProviderConfig{
		ClientId:       "11111111-1111-1111-1111-111111111111",
		ClientSecret:   "static-client-secret",
		TenantId:       "22222222-2222-2222-2222-222222222222",
		SubscriptionId: "33333333-3333-3333-3333-333333333333",
	}

	args, err := buildProviderArgs(cfg, nil)
	require.NoError(t, err)

	assert.Equal(t, pulumi.String("11111111-1111-1111-1111-111111111111"), args.ClientId)
	assert.Equal(t, pulumi.String("static-client-secret"), args.ClientSecret)
	assert.Equal(t, pulumi.String("22222222-2222-2222-2222-222222222222"), args.TenantId)
	assert.Equal(t, pulumi.String("33333333-3333-3333-3333-333333333333"), args.SubscriptionId)
	// Static and keyless are mutually exclusive.
	assert.Nil(t, args.UseOidc)
	assert.Nil(t, args.OidcToken)
}

func TestBuildProviderArgs_WebIdentity_SetsInlineOidcToken(t *testing.T) {
	cfg := &azureprovider.AzureProviderConfig{
		ClientId:       "11111111-1111-1111-1111-111111111111",
		TenantId:       "22222222-2222-2222-2222-222222222222",
		SubscriptionId: "33333333-3333-3333-3333-333333333333",
		WebIdentity: &azureprovider.AzureWebIdentityProviderConfig{
			WebIdentityToken: "eyJhbGciOiJSUzI1NiJ9.payload.sig",
		},
	}

	args, err := buildProviderArgs(cfg, nil)
	require.NoError(t, err)

	assert.Equal(t, pulumi.Bool(true), args.UseOidc)
	// This classic SDK's NewProvider auto-secret-wraps OidcToken, so the builder passes the
	// plain string and must not double-wrap it.
	assert.Equal(t, pulumi.String("eyJhbGciOiJSUzI1NiJ9.payload.sig"), args.OidcToken)
	assert.Equal(t, pulumi.String("11111111-1111-1111-1111-111111111111"), args.ClientId)
	assert.Equal(t, pulumi.String("22222222-2222-2222-2222-222222222222"), args.TenantId)
	assert.Equal(t, pulumi.String("33333333-3333-3333-3333-333333333333"), args.SubscriptionId)
	// Keyless must never carry a client secret.
	assert.Nil(t, args.ClientSecret)
}

func TestBuildProviderArgs_WebIdentity_TakesPrecedenceOverStaleSecret(t *testing.T) {
	// A config carrying both dispatches to keyless: web identity is the deliberate mode
	// switch, a lingering client_secret must not win.
	cfg := &azureprovider.AzureProviderConfig{
		ClientId:       "11111111-1111-1111-1111-111111111111",
		ClientSecret:   "stale-client-secret",
		TenantId:       "22222222-2222-2222-2222-222222222222",
		SubscriptionId: "33333333-3333-3333-3333-333333333333",
		WebIdentity: &azureprovider.AzureWebIdentityProviderConfig{
			WebIdentityToken: "eyJhbGciOiJSUzI1NiJ9.payload.sig",
		},
	}

	args, err := buildProviderArgs(cfg, nil)
	require.NoError(t, err)

	assert.Equal(t, pulumi.Bool(true), args.UseOidc)
	assert.NotNil(t, args.OidcToken)
	assert.Nil(t, args.ClientSecret)
}

func TestBuildProviderArgs_WebIdentity_MissingToken_Errors(t *testing.T) {
	cfg := &azureprovider.AzureProviderConfig{
		ClientId:       "11111111-1111-1111-1111-111111111111",
		TenantId:       "22222222-2222-2222-2222-222222222222",
		SubscriptionId: "33333333-3333-3333-3333-333333333333",
		WebIdentity:    &azureprovider.AzureWebIdentityProviderConfig{},
	}

	_, err := buildProviderArgs(cfg, nil)
	assert.Error(t, err)
}

func TestBuildProviderArgs_NilFeatures_LeftUnset(t *testing.T) {
	args, err := buildProviderArgs(nil, nil)
	require.NoError(t, err)

	assert.Nil(t, args.Features)
}

func TestBuildProviderArgs_Features_PassedThrough(t *testing.T) {
	// The features block is a module-design input, orthogonal to credential dispatch:
	// it must arrive on the args verbatim and leave the credential fields untouched.
	features := azure.ProviderFeaturesArgs{
		MachineLearning: azure.ProviderFeaturesMachineLearningArgs{
			PurgeSoftDeletedWorkspaceOnDestroy: pulumi.Bool(true),
		},
	}
	cfg := &azureprovider.AzureProviderConfig{
		SubscriptionId: "33333333-3333-3333-3333-333333333333",
	}

	args, err := buildProviderArgs(cfg, features)
	require.NoError(t, err)

	assert.Equal(t, features, args.Features)
	assert.Equal(t, pulumi.String("33333333-3333-3333-3333-333333333333"), args.SubscriptionId)
	assert.Nil(t, args.ClientSecret)
	assert.Nil(t, args.UseOidc)
}

func TestProviderResourceName(t *testing.T) {
	// State continuity: the base name must stay "azure".
	assert.Equal(t, "azure", ProviderResourceName(nil))
	assert.Equal(t, "azure-replica", ProviderResourceName([]string{"replica"}))
	assert.Equal(t, "azure-dns-zone", ProviderResourceName([]string{"dns", "zone"}))
}
