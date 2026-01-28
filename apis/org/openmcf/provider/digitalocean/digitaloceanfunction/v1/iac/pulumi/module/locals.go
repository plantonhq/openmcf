package module

import (
	"strconv"

	digitaloceanprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean"
	digitaloceanfunctionv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean/digitaloceanfunction/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/digitalocean/digitaloceanlabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	DigitalOceanProviderConfig *digitaloceanprovider.DigitalOceanProviderConfig
	DigitalOceanFunction       *digitaloceanfunctionv1.DigitalOceanFunction
	DigitalOceanLabels         map[string]string
}

// initializeLocals copies stack‑input fields into the Locals struct and builds a reusable label map.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceanfunctionv1.DigitalOceanFunctionStackInput) *Locals {
	locals := &Locals{}

	locals.DigitalOceanFunction = stackInput.Target

	// Standard Planton labels for DigitalOcean resources.
	locals.DigitalOceanLabels = map[string]string{
		digitaloceanlabelkeys.Resource:     strconv.FormatBool(true),
		digitaloceanlabelkeys.ResourceName: locals.DigitalOceanFunction.Metadata.Name,
		digitaloceanlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_DigitalOceanFunction.String(),
	}

	if locals.DigitalOceanFunction.Metadata.Org != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Organization] = locals.DigitalOceanFunction.Metadata.Org
	}

	if locals.DigitalOceanFunction.Metadata.Env != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Environment] = locals.DigitalOceanFunction.Metadata.Env
	}

	if locals.DigitalOceanFunction.Metadata.Id != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.ResourceId] = locals.DigitalOceanFunction.Metadata.Id
	}

	locals.DigitalOceanProviderConfig = stackInput.ProviderConfig

	return locals
}
