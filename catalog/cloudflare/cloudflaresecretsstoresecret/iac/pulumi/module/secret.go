package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// secret creates one secret inside the account Secrets Store. The value is
// write-only at Cloudflare (never returned, never drift-detected) and is
// marked secret in Pulumi state. account_id, store_id, and name are
// create-only; value, scopes, and comment update in place (a merge-patch).
// The spec's CEL wall already guarantees scopes arrive in Cloudflare's
// canonical alphabetical order -- the API returns them sorted, and an
// unsorted config would drift forever against the provider's ordered list.
func secret(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareSecretsStoreSecret.Spec

	args := &cloudflare.SecretsStoreSecretArgs{
		AccountId: pulumi.String(spec.AccountId),
		StoreId:   pulumi.String(spec.StoreId.GetValue()),
		Name:      pulumi.String(spec.Name),
		Value:     pulumi.ToSecret(pulumi.String(spec.Value.GetValue())).(pulumi.StringOutput),
		Scopes:    pulumi.ToStringArray(spec.Scopes),
	}

	if spec.Comment != "" {
		args.Comment = pulumi.StringPtr(spec.Comment)
	}

	createdSecret, err := cloudflare.NewSecretsStoreSecret(
		ctx,
		"secrets_store_secret",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create secrets store secret")
	}

	ctx.Export(OpSecretId, createdSecret.ID())
	ctx.Export(OpStoreId, createdSecret.StoreId)

	return nil
}
