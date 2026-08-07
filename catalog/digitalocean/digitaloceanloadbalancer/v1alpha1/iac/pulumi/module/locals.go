package module

import (
	"strconv"

	digitaloceanprovider "github.com/plantonhq/planton/catalog/digitalocean"
	digitaloceanloadbalancerv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceanloadbalancer/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/digitalocean/digitaloceanlabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	DigitalOceanProviderConfig *digitaloceanprovider.DigitalOceanProviderConfig
	DigitalOceanLoadBalancer   *digitaloceanloadbalancerv1alpha1.DigitalOceanLoadBalancer
	DigitalOceanLabels         map[string]string
}

// initializeLocals copies stack‑input fields into the Locals struct and builds
// a reusable label map—mirrors the VPC module pattern.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceanloadbalancerv1alpha1.DigitalOceanLoadBalancerStackInput) *Locals {
	locals := &Locals{}

	locals.DigitalOceanLoadBalancer = stackInput.Target

	// Standard Planton labels for DigitalOcean resources.
	locals.DigitalOceanLabels = map[string]string{
		digitaloceanlabelkeys.Resource:     strconv.FormatBool(true),
		digitaloceanlabelkeys.ResourceName: locals.DigitalOceanLoadBalancer.Spec.LoadBalancerName,
		digitaloceanlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_DigitalOceanLoadBalancer.String(),
	}

	if locals.DigitalOceanLoadBalancer.Metadata.Org != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Organization] = locals.DigitalOceanLoadBalancer.Metadata.Org
	}

	if locals.DigitalOceanLoadBalancer.Metadata.Env != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Environment] = locals.DigitalOceanLoadBalancer.Metadata.Env
	}

	if locals.DigitalOceanLoadBalancer.Metadata.Id != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.ResourceId] = locals.DigitalOceanLoadBalancer.Metadata.Id
	}

	locals.DigitalOceanProviderConfig = stackInput.ProviderConfig

	return locals
}
