package module

import (
	"strconv"

	civoprovider "github.com/plantonhq/planton/catalog/civo"
	civodatabasev1alpha1 "github.com/plantonhq/planton/catalog/civo/civodatabase/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals groups frequently‑used values so other files stay concise.
type Locals struct {
	CivoProviderConfig *civoprovider.CivoProviderConfig
	CivoDatabase       *civodatabasev1alpha1.CivoDatabase
	CivoTags           pulumi.StringArray
	CivoLabels         map[string]string
}

// initializeLocals copies the stack‑input into Locals and derives tags / labels.
func initializeLocals(_ *pulumi.Context, stackInput *civodatabasev1alpha1.CivoDatabaseStackInput) *Locals {
	locals := &Locals{}
	locals.CivoDatabase = stackInput.Target
	locals.CivoProviderConfig = stackInput.ProviderConfig

	// Tags (directly from spec.tags).
	if len(locals.CivoDatabase.Spec.Tags) > 0 {
		for _, t := range locals.CivoDatabase.Spec.Tags {
			locals.CivoTags = append(locals.CivoTags, pulumi.String(t))
		}
	}

	// Basic Planton labels (handy if Civo adds tag‑key:value support later).
	locals.CivoLabels = map[string]string{
		"planton:resource":      strconv.FormatBool(true),
		"planton:resource_name": locals.CivoDatabase.Metadata.Name,
		"planton:resource_kind": cloudresourcekind.CloudResourceKind_CivoDatabase.String(),
	}
	if locals.CivoDatabase.Metadata.Org != "" {
		locals.CivoLabels["planton:organization"] = locals.CivoDatabase.Metadata.Org
	}
	if locals.CivoDatabase.Metadata.Env != "" {
		locals.CivoLabels["planton:environment"] = locals.CivoDatabase.Metadata.Env
	}
	if locals.CivoDatabase.Metadata.Id != "" {
		locals.CivoLabels["planton:resource_id"] = locals.CivoDatabase.Metadata.Id
	}

	return locals
}
