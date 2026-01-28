package module

import (
	"strconv"

	digitaloceanprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean"
	digitaloceancertificatev1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/digitalocean/digitaloceancertificate/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/digitalocean/digitaloceanlabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals aggregates handy references for the rest of the module.
type Locals struct {
	DigitalOceanProviderConfig *digitaloceanprovider.DigitalOceanProviderConfig
	DigitalOceanCertificate    *digitaloceancertificatev1.DigitalOceanCertificate
	DigitalOceanLabels         map[string]string
}

// initializeLocals copies stack‑input fields into Locals and builds a reusable label map.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceancertificatev1.DigitalOceanCertificateStackInput) *Locals {
	var locals Locals

	locals.DigitalOceanCertificate = stackInput.Target

	// Standard Planton labels for DigitalOcean resources.
	locals.DigitalOceanLabels = map[string]string{
		digitaloceanlabelkeys.Resource:     strconv.FormatBool(true),
		digitaloceanlabelkeys.ResourceName: locals.DigitalOceanCertificate.Metadata.Name,
		digitaloceanlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_DigitalOceanCertificate.String(),
	}

	if locals.DigitalOceanCertificate.Metadata.Org != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Organization] = locals.DigitalOceanCertificate.Metadata.Org
	}

	if locals.DigitalOceanCertificate.Metadata.Env != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Environment] = locals.DigitalOceanCertificate.Metadata.Env
	}

	if locals.DigitalOceanCertificate.Metadata.Id != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.ResourceId] = locals.DigitalOceanCertificate.Metadata.Id
	}

	locals.DigitalOceanProviderConfig = stackInput.ProviderConfig

	return &locals
}
