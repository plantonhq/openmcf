package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// serviceAccount provisions the account and (optionally) a key.
// Returns the created account and key pointers (key may be nil).
func serviceAccount(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
) (*serviceaccount.Account, *serviceaccount.Key, error) {
	spec := locals.GcpServiceAccount.Spec

	// Display name falls back to the resource's metadata name so every account
	// is human-identifiable in the console even when the field is omitted.
	displayName := spec.DisplayName
	if displayName == "" {
		displayName = locals.GcpServiceAccount.Metadata.Name
	}

	accountArgs := &serviceaccount.AccountArgs{
		// account_id and project are immutable (ForceNew in the provider):
		// changing either destroys and recreates the account, invalidating every
		// IAM binding and workload identity referencing the old email.
		AccountId:   pulumi.String(spec.ServiceAccountId),
		DisplayName: pulumi.String(displayName),
		// GCP flips disabled state via separate Enable/Disable API calls (not the
		// regular update mask); the provider handles that internally, so toggling
		// this field never recreates the account.
		Disabled: pulumi.Bool(spec.GetDisabled()),
	}

	// Omitted description stays unset (matching the Terraform module's null)
	// rather than being sent as an empty string.
	if spec.Description != "" {
		accountArgs.Description = pulumi.String(spec.Description)
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		accountArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdServiceAccount, err := serviceaccount.NewAccount(
		ctx,
		locals.GcpServiceAccount.Metadata.Name,
		accountArgs,
		pulumi.Provider(gcpProvider),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create service account")
	}

	// Optional user-managed JSON key. Created only when spec.create_key is true —
	// keyless patterns (Workload Identity, impersonation, federation) are the
	// recommended default. The private key is returned once at creation and is
	// marked secret in state by the engine.
	var createdKey *serviceaccount.Key
	if spec.GetCreateKey() {
		createdKey, err = serviceaccount.NewKey(
			ctx,
			fmt.Sprintf("%s-key", locals.GcpServiceAccount.Metadata.Name),
			&serviceaccount.KeyArgs{
				ServiceAccountId: createdServiceAccount.Name,
			},
			pulumi.Provider(gcpProvider),
			pulumi.DependsOn([]pulumi.Resource{createdServiceAccount}),
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create service account key")
		}
	}

	return createdServiceAccount, createdKey, nil
}
