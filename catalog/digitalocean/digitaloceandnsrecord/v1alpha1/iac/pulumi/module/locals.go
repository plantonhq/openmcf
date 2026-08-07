package module

import (
	"strconv"

	digitaloceanprovider "github.com/plantonhq/planton/catalog/digitalocean"
	digitaloceandnsrecordv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceandnsrecord/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/digitalocean/digitaloceanlabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds quick references used by other files.
type Locals struct {
	DigitalOceanProviderConfig *digitaloceanprovider.DigitalOceanProviderConfig
	DigitalOceanDnsRecord      *digitaloceandnsrecordv1alpha1.DigitalOceanDnsRecord
	DigitalOceanLabels         map[string]string
}

// initializeLocals sets up local values from stack input.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceandnsrecordv1alpha1.DigitalOceanDnsRecordStackInput) *Locals {
	locals := &Locals{}

	locals.DigitalOceanDnsRecord = stackInput.Target

	// Standard Planton labels for DigitalOcean resources.
	locals.DigitalOceanLabels = map[string]string{
		digitaloceanlabelkeys.Resource:     strconv.FormatBool(true),
		digitaloceanlabelkeys.ResourceName: locals.DigitalOceanDnsRecord.Metadata.Name,
		digitaloceanlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_DigitalOceanDnsRecord.String(),
	}

	if locals.DigitalOceanDnsRecord.Metadata.Org != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Organization] = locals.DigitalOceanDnsRecord.Metadata.Org
	}
	if locals.DigitalOceanDnsRecord.Metadata.Env != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Environment] = locals.DigitalOceanDnsRecord.Metadata.Env
	}
	if locals.DigitalOceanDnsRecord.Metadata.Id != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.ResourceId] = locals.DigitalOceanDnsRecord.Metadata.Id
	}

	locals.DigitalOceanProviderConfig = stackInput.ProviderConfig

	return locals
}
