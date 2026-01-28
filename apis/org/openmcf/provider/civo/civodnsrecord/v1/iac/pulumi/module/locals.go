package module

import (
	"strconv"

	civoprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/civo"
	civodnsrecordv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/civo/civodnsrecord/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/civo/civolabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles quick references that multiple files need.
type Locals struct {
	CivoProviderConfig *civoprovider.CivoProviderConfig
	CivoDnsRecord      *civodnsrecordv1.CivoDnsRecord
	CivoLabels         map[string]string
	ZoneId             string
}

// initializeLocals mirrors the pattern used by other Planton modules.
func initializeLocals(_ *pulumi.Context, stackInput *civodnsrecordv1.CivoDnsRecordStackInput) *Locals {
	locals := &Locals{}
	locals.CivoDnsRecord = stackInput.Target

	target := stackInput.Target

	// Extract zone ID from StringValueOrRef
	locals.ZoneId = target.Spec.ZoneId.GetValue()

	// Standard Planton labels for Civo resources.
	locals.CivoLabels = map[string]string{
		civolabelkeys.Resource:     strconv.FormatBool(true),
		civolabelkeys.ResourceName: locals.CivoDnsRecord.Metadata.Name,
		civolabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_CivoDnsRecord.String(),
	}

	if locals.CivoDnsRecord.Metadata.Org != "" {
		locals.CivoLabels[civolabelkeys.Organization] = locals.CivoDnsRecord.Metadata.Org
	}
	if locals.CivoDnsRecord.Metadata.Env != "" {
		locals.CivoLabels[civolabelkeys.Environment] = locals.CivoDnsRecord.Metadata.Env
	}
	if locals.CivoDnsRecord.Metadata.Id != "" {
		locals.CivoLabels[civolabelkeys.ResourceId] = locals.CivoDnsRecord.Metadata.Id
	}

	locals.CivoProviderConfig = stackInput.ProviderConfig

	return locals
}
