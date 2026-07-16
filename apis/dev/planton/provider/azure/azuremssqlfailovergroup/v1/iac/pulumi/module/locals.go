package module

import (
	"strings"

	azuremssqlfailovergroupv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremssqlfailovergroup/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMssqlFailoverGroup *azuremssqlfailovergroupv1.AzureMssqlFailoverGroup

	// ServerId is the resolved ARM ID of the primary logical server;
	// PartnerServerIds and DatabaseIds are the resolved ARM IDs of the
	// partner servers and the databases to replicate.
	ServerId         string
	PartnerServerIds []string
	DatabaseIds      []string

	// FailoverMode is the ARM string for the read-write policy mode.
	FailoverMode string

	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremssqlfailovergroupv1.AzureMssqlFailoverGroupStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMssqlFailoverGroup = stackInput.Target
	target := stackInput.Target
	spec := target.Spec

	locals.ServerId = spec.ServerId.GetValue()

	for _, p := range spec.PartnerServers {
		locals.PartnerServerIds = append(locals.PartnerServerIds, p.ServerId.GetValue())
	}
	for _, d := range spec.DatabaseIds {
		locals.DatabaseIds = append(locals.DatabaseIds, d.GetValue())
	}

	if spec.ReadWriteEndpointFailoverPolicy != nil {
		switch spec.ReadWriteEndpointFailoverPolicy.Mode {
		case azuremssqlfailovergroupv1.AzureMssqlFailoverGroupFailoverMode_AUTOMATIC:
			locals.FailoverMode = "Automatic"
		case azuremssqlfailovergroupv1.AzureMssqlFailoverGroupFailoverMode_MANUAL:
			locals.FailoverMode = "Manual"
		}
	}

	// PARITY-EXCEPTION: resource_kind here is the lowered CloudResourceKind
	// enum string and resource_id is omitted when metadata.id is empty,
	// while the Terraform module emits the family-wide snake-case literal
	// and falls back to metadata.name. Output-neutral (tags never feed stack
	// outputs); aligning the two shapes is a family-wide convention change.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMssqlFailoverGroup.String()),
	}
	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}
	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}
	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}
	for k, v := range spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}
