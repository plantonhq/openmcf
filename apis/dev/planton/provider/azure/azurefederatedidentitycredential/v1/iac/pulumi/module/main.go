package module

import (
	"github.com/pkg/errors"
	azurefederatedidentitycredentialv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefederatedidentitycredential/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/azure/pulumiazureprovider"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/armmsi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *azurefederatedidentitycredentialv1.AzureFederatedIdentityCredentialStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the Azure provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static client secret, keyless web identity, or ambient chain).
	azureProvider, err := pulumiazureprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create azure provider")
	}

	spec := locals.AzureFederatedIdentityCredential.Spec

	// The provider SDK requires the resource group as its own argument even
	// though the parent identity's ARM ID already carries it; derive it so
	// the two can never disagree.
	resourceGroupName, err := resourceGroupNameFromIdentityId(locals.UserAssignedIdentityId)
	if err != nil {
		return err
	}

	// Lifecycle notes worth knowing before operating this resource:
	// - issuer, subject, and audience update IN PLACE; name and the parent
	//   identity are the credential's ARM identity, so changing either
	//   replaces it (delete + create).
	// - The provider serializes writes per parent identity (ARM rejects
	//   concurrent credential writes on one identity), so several credentials
	//   on the same identity deploy sequentially -- expected, not a hang.
	//
	// The audience: ARM models it as a single-element list and this SDK
	// flattens it to one string -- the same one-audience contract the
	// Terraform engine expresses as a list capped at one element.
	credential, err := armmsi.NewFederatedIdentityCredential(ctx,
		locals.AzureFederatedIdentityCredential.Metadata.Name,
		&armmsi.FederatedIdentityCredentialArgs{
			Name:              pulumi.String(spec.Name),
			ParentId:          pulumi.String(locals.UserAssignedIdentityId),
			ResourceGroupName: pulumi.String(resourceGroupName),
			Issuer:            pulumi.String(spec.Issuer),
			Subject:           pulumi.String(spec.Subject),
			Audience:          pulumi.String(locals.Audience),
		},
		pulumi.Provider(azureProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create federated identity credential %s", spec.Name)
	}

	// Export stack outputs from the created resource (not the spec) so they
	// carry the values Azure resolved. The trust coordinates are exported so
	// downstream automation (CI config generators, cluster onboarding) can
	// wire the external side of the trust without re-reading the spec.
	ctx.Export(OpFederatedIdentityCredentialId, credential.ID())
	ctx.Export(OpName, credential.Name)
	ctx.Export(OpUserAssignedIdentityId, credential.ParentId)
	ctx.Export(OpIssuer, credential.Issuer)
	ctx.Export(OpSubject, credential.Subject)
	ctx.Export(OpAudience, credential.Audience)

	return nil
}
