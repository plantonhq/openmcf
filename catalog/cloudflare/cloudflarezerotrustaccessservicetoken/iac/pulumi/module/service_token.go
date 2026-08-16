package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// serviceToken creates the Access service token. The client secret is returned
// only at creation and rotation; the export is marked secret (the SDK also
// registers client_secret as an additional secret output).
func serviceToken(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustAccessServiceToken.Spec

	args := &cloudflare.ZeroTrustAccessServiceTokenArgs{
		Name: pulumi.String(spec.Name),
	}

	// Scope: exactly one of account_id or zone_id is set (enforced by the spec).
	if spec.AccountId != "" {
		args.AccountId = pulumi.StringPtr(spec.AccountId)
	}
	if spec.ZoneId.GetValue() != "" {
		args.ZoneId = pulumi.StringPtr(spec.ZoneId.GetValue())
	}

	// Empty means Cloudflare's default of one year (8760h).
	if spec.Duration != "" {
		args.Duration = pulumi.StringPtr(spec.Duration)
	}

	// Rotation pair (spec-enforced both-or-neither). The SDK types the version
	// as float64; the spec's integer converts losslessly.
	if spec.ClientSecretVersion != nil {
		args.ClientSecretVersion = pulumi.Float64Ptr(float64(*spec.ClientSecretVersion))
		args.PreviousClientSecretExpiresAt = pulumi.StringPtr(spec.PreviousClientSecretExpiresAt)
	}

	createdToken, err := cloudflare.NewZeroTrustAccessServiceToken(
		ctx,
		"service_token",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create access service token")
	}

	ctx.Export(OpServiceTokenId, createdToken.ID())
	ctx.Export(OpClientId, createdToken.ClientId)
	ctx.Export(OpClientSecret, pulumi.ToSecret(createdToken.ClientSecret))
	ctx.Export(OpExpiresAt, createdToken.ExpiresAt)

	return nil
}
