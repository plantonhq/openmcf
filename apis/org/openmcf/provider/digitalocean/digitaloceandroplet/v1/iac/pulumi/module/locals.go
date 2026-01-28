package module

import (
	"strconv"

	digitaloceanprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean"
	digitaloceandropletv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean/digitaloceandroplet/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/digitalocean/digitaloceanlabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles common pointers and label maps used across the module.
type Locals struct {
	DigitalOceanProviderConfig *digitaloceanprovider.DigitalOceanProviderConfig
	DigitalOceanDroplet        *digitaloceandropletv1.DigitalOceanDroplet
	DigitalOceanLabels         map[string]string
}

// initializeLocals mirrors the pattern established in the VPC module.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceandropletv1.DigitalOceanDropletStackInput) *Locals {
	locals := &Locals{}

	locals.DigitalOceanDroplet = stackInput.Target

	// Standard Planton labels for DigitalOcean resources.
	locals.DigitalOceanLabels = map[string]string{
		digitaloceanlabelkeys.Resource:     strconv.FormatBool(true),
		digitaloceanlabelkeys.ResourceName: locals.DigitalOceanDroplet.Metadata.Name,
		digitaloceanlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_DigitalOceanDroplet.String(),
	}

	if locals.DigitalOceanDroplet.Metadata.Org != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Organization] = locals.DigitalOceanDroplet.Metadata.Org
	}

	if locals.DigitalOceanDroplet.Metadata.Env != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Environment] = locals.DigitalOceanDroplet.Metadata.Env
	}

	if locals.DigitalOceanDroplet.Metadata.Id != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.ResourceId] = locals.DigitalOceanDroplet.Metadata.Id
	}

	locals.DigitalOceanProviderConfig = stackInput.ProviderConfig

	return locals
}
