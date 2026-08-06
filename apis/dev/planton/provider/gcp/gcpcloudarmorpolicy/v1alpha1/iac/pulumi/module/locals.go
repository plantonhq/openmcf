package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpcloudarmorpolicyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudarmorpolicy/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig   *gcpprovider.GcpProviderConfig
	GcpCloudArmorPolicy *gcpcloudarmorpolicyv1alpha1.GcpCloudArmorPolicy

	// ProjectId is empty when the manifest omits it — the provider's default
	// project then applies (the same ambient contract the Terraform module
	// honors by passing null).
	ProjectId string

	// PolicyName falls back to metadata.name — explicit conditional, so both
	// engines derive the identical cloud-side name.
	PolicyName string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpcloudarmorpolicyv1alpha1.GcpCloudArmorPolicyStackInput) *Locals {
	locals := &Locals{}
	locals.GcpCloudArmorPolicy = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	locals.ProjectId = locals.GcpCloudArmorPolicy.Spec.ProjectId.GetValue()

	locals.PolicyName = locals.GcpCloudArmorPolicy.Spec.PolicyName
	if locals.PolicyName == "" {
		locals.PolicyName = locals.GcpCloudArmorPolicy.Metadata.Name
	}

	// Cloud Armor security policies carry no labels on the released provider
	// line — attribution is deliberately not attempted on either engine.

	return locals
}
