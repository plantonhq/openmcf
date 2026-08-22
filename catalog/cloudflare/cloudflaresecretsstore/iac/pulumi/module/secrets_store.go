package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// secretsStore creates the account-level Secrets Store. Both arguments are
// create-only at the API (the provider's Update is an empty stub and every
// field forces replacement) -- a name change replaces the store AND every
// secret inside it. Cloudflare also allows only one store per account: if
// one already exists (e.g. dashboard-created), this create fails and the
// existing store should be adopted by import instead.
func secretsStore(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareSecretsStore.Spec

	createdStore, err := cloudflare.NewSecretsStore(
		ctx,
		"secrets_store",
		&cloudflare.SecretsStoreArgs{
			AccountId: pulumi.String(spec.AccountId),
			Name:      pulumi.String(spec.Name),
		},
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create secrets store")
	}

	ctx.Export(OpStoreId, createdStore.ID())

	return nil
}
