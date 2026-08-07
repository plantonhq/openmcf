package module

import (
	azurefrontdoorsecretv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurefrontdoorsecret/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFrontDoorSecret  *azurefrontdoorsecretv1alpha1.AzureFrontDoorSecret
	ProfileId             string
	KeyVaultCertificateId string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefrontdoorsecretv1alpha1.AzureFrontDoorSecretStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFrontDoorSecret = stackInput.Target
	locals.ProfileId = stackInput.Target.Spec.ProfileId.GetValue()
	locals.KeyVaultCertificateId = stackInput.Target.Spec.KeyVaultCertificateId.GetValue()

	// No Azure tags: ARM does not support tags on Front Door secrets,
	// so the platform's identity tags live on the profile.

	return locals
}
