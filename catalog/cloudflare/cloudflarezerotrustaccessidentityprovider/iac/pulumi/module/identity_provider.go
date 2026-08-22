package module

import (
	"github.com/pkg/errors"
	cloudflarezerotrustaccessidentityproviderv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustaccessidentityprovider/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// identityProvider creates the Access identity provider. The TYPE is immutable
// at Cloudflare -- changing it replaces the provider (new ID), invalidating
// policy rules referencing the old one. The SCIM secret is minted once when
// SCIM is first enabled and redacted on every later read.
func identityProvider(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustAccessIdentityProvider.Spec

	args := &cloudflare.ZeroTrustAccessIdentityProviderArgs{
		Name: pulumi.String(spec.Name),
		Type: pulumi.String(spec.Type),
		// The provider REQUIRES config even for types that need no parameters
		// (onetimepin) -- the empty object is the correct payload then.
		Config: buildConfigArgs(spec.Config),
	}

	// Scope: exactly one of account_id or zone_id is set (enforced by the spec).
	if spec.AccountId != "" {
		args.AccountId = pulumi.StringPtr(spec.AccountId)
	}
	if spec.ZoneId.GetValue() != "" {
		args.ZoneId = pulumi.StringPtr(spec.ZoneId.GetValue())
	}

	if spec.SamlCertificateSetId != "" {
		args.SamlCertificateSetId = pulumi.StringPtr(spec.SamlCertificateSetId)
	}

	if spec.ScimConfig != nil {
		scimArgs := &cloudflare.ZeroTrustAccessIdentityProviderScimConfigArgs{
			Enabled:         pulumi.BoolPtr(spec.ScimConfig.Enabled),
			SeatDeprovision: pulumi.BoolPtr(spec.ScimConfig.SeatDeprovision),
			UserDeprovision: pulumi.BoolPtr(spec.ScimConfig.UserDeprovision),
		}
		if spec.ScimConfig.IdentityUpdateBehavior != "" {
			scimArgs.IdentityUpdateBehavior = pulumi.StringPtr(spec.ScimConfig.IdentityUpdateBehavior)
		}
		args.ScimConfig = scimArgs
	}

	// A safety latch: while true, Cloudflare refuses API updates and deletes.
	if spec.ReadOnly {
		args.ReadOnly = pulumi.BoolPtr(true)
	}

	createdIdp, err := cloudflare.NewZeroTrustAccessIdentityProvider(
		ctx,
		"identity_provider",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create access identity provider")
	}

	ctx.Export(OpIdentityProviderId, createdIdp.ID())
	ctx.Export(OpScimBaseUrl, createdIdp.ScimConfig.ScimBaseUrl())
	ctx.Export(OpScimSecret, pulumi.ToSecret(createdIdp.ScimConfig.Secret()))

	return nil
}

// buildConfigArgs maps the spec's config onto the SDK's args, sending only the
// fields the manifest set -- the provider's per-type validators trigger on any
// non-null value, so an unset field must never travel.
func buildConfigArgs(
	config *cloudflarezerotrustaccessidentityproviderv1alpha1.CloudflareZeroTrustAccessIdentityProviderConfig,
) *cloudflare.ZeroTrustAccessIdentityProviderConfigArgs {
	configArgs := &cloudflare.ZeroTrustAccessIdentityProviderConfigArgs{}
	if config == nil {
		return configArgs
	}

	if len(config.Claims) > 0 {
		configArgs.Claims = pulumi.ToStringArray(config.Claims)
	}
	if config.ClientId != "" {
		configArgs.ClientId = pulumi.StringPtr(config.ClientId)
	}
	if config.ClientSecret.GetValue() != "" {
		configArgs.ClientSecret = pulumi.StringPtr(config.ClientSecret.GetValue())
	}
	if config.EmailClaimName != "" {
		configArgs.EmailClaimName = pulumi.StringPtr(config.EmailClaimName)
	}
	if config.PkceEnabled != nil {
		configArgs.PkceEnabled = pulumi.BoolPtr(*config.PkceEnabled)
	}

	// azureAD family.
	if config.ConditionalAccessEnabled != nil {
		configArgs.ConditionalAccessEnabled = pulumi.BoolPtr(*config.ConditionalAccessEnabled)
	}
	if config.DirectoryId != "" {
		configArgs.DirectoryId = pulumi.StringPtr(config.DirectoryId)
	}
	if config.Prompt != "" {
		configArgs.Prompt = pulumi.StringPtr(config.Prompt)
	}
	if config.SupportGroups != nil {
		configArgs.SupportGroups = pulumi.BoolPtr(*config.SupportGroups)
	}

	// centrify family.
	if config.CentrifyAccount != "" {
		configArgs.CentrifyAccount = pulumi.StringPtr(config.CentrifyAccount)
	}
	if config.CentrifyAppId != "" {
		configArgs.CentrifyAppId = pulumi.StringPtr(config.CentrifyAppId)
	}

	// google-apps family.
	if config.AppsDomain != "" {
		configArgs.AppsDomain = pulumi.StringPtr(config.AppsDomain)
	}

	// oidc family.
	if config.AuthUrl != "" {
		configArgs.AuthUrl = pulumi.StringPtr(config.AuthUrl)
	}
	if config.CertsUrl != "" {
		configArgs.CertsUrl = pulumi.StringPtr(config.CertsUrl)
	}
	if len(config.Scopes) > 0 {
		configArgs.Scopes = pulumi.ToStringArray(config.Scopes)
	}
	if config.TokenUrl != "" {
		configArgs.TokenUrl = pulumi.StringPtr(config.TokenUrl)
	}

	// okta / onelogin / pingone families.
	if config.AuthorizationServerId != "" {
		configArgs.AuthorizationServerId = pulumi.StringPtr(config.AuthorizationServerId)
	}
	if config.OktaAccount != "" {
		configArgs.OktaAccount = pulumi.StringPtr(config.OktaAccount)
	}
	if config.OneloginAccount != "" {
		configArgs.OneloginAccount = pulumi.StringPtr(config.OneloginAccount)
	}
	if config.PingEnvId != "" {
		configArgs.PingEnvId = pulumi.StringPtr(config.PingEnvId)
	}

	// saml family.
	if len(config.Attributes) > 0 {
		configArgs.Attributes = pulumi.ToStringArray(config.Attributes)
	}
	if config.EmailAttributeName != "" {
		configArgs.EmailAttributeName = pulumi.StringPtr(config.EmailAttributeName)
	}
	if config.EnableEncryption != nil {
		configArgs.EnableEncryption = pulumi.BoolPtr(*config.EnableEncryption)
	}
	if len(config.HeaderAttributes) > 0 {
		headerAttributes := make(cloudflare.ZeroTrustAccessIdentityProviderConfigHeaderAttributeArray, 0, len(config.HeaderAttributes))
		for _, headerAttribute := range config.HeaderAttributes {
			attributeArgs := &cloudflare.ZeroTrustAccessIdentityProviderConfigHeaderAttributeArgs{}
			if headerAttribute.AttributeName != "" {
				attributeArgs.AttributeName = pulumi.StringPtr(headerAttribute.AttributeName)
			}
			if headerAttribute.HeaderName != "" {
				attributeArgs.HeaderName = pulumi.StringPtr(headerAttribute.HeaderName)
			}
			headerAttributes = append(headerAttributes, attributeArgs)
		}
		configArgs.HeaderAttributes = headerAttributes
	}
	if len(config.IdpPublicCerts) > 0 {
		configArgs.IdpPublicCerts = pulumi.ToStringArray(config.IdpPublicCerts)
	}
	if config.IssuerUrl != "" {
		configArgs.IssuerUrl = pulumi.StringPtr(config.IssuerUrl)
	}
	if config.SignRequest != nil {
		configArgs.SignRequest = pulumi.BoolPtr(*config.SignRequest)
	}
	if config.SsoTargetUrl != "" {
		configArgs.SsoTargetUrl = pulumi.StringPtr(config.SsoTargetUrl)
	}

	if config.RestrictToAccountMembers != nil {
		configArgs.RestrictToAccountMembers = pulumi.BoolPtr(*config.RestrictToAccountMembers)
	}

	return configArgs
}
