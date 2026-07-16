package module

import (
	"strings"

	azuremysqlflexibleserverv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremysqlflexibleserver/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMysqlFlexibleServer *azuremysqlflexibleserverv1.AzureMysqlFlexibleServer
	ResourceGroupName        string
	AzureTags                map[string]string
}

// createModeStrings maps the spec's create-mode enum to ARM's values. An
// unspecified mode means a fresh (DEFAULT) server and is not sent at all,
// mirroring the Terraform module's null -- so the two engines produce the
// same ARM payload.
var createModeStrings = map[azuremysqlflexibleserverv1.AzureMysqlFlexibleServerCreateMode]string{
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerCreateMode_DEFAULT:               "Default",
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerCreateMode_POINT_IN_TIME_RESTORE: "PointInTimeRestore",
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerCreateMode_REPLICA:               "Replica",
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerCreateMode_GEO_RESTORE:           "GeoRestore",
}

// haModeStrings maps the spec's high-availability mode enum to ARM's values.
var haModeStrings = map[azuremysqlflexibleserverv1.AzureMysqlFlexibleServerHighAvailabilityMode]string{
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerHighAvailabilityMode_ZONE_REDUNDANT: "ZoneRedundant",
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerHighAvailabilityMode_SAME_ZONE:      "SameZone",
}

// publicNetworkAccessStrings maps the spec's public-network-access enum to
// ARM's values. Unspecified is not sent at all, letting Azure derive the
// value (Enabled publicly, Disabled when VNet-injected) -- mirroring the
// Terraform module's null.
var publicNetworkAccessStrings = map[azuremysqlflexibleserverv1.AzureMysqlFlexibleServerPublicNetworkAccess]string{
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerPublicNetworkAccess_ENABLED:  "Enabled",
	azuremysqlflexibleserverv1.AzureMysqlFlexibleServerPublicNetworkAccess_DISABLED: "Disabled",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremysqlflexibleserverv1.AzureMysqlFlexibleServerStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMysqlFlexibleServer = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMysqlFlexibleServer.String()),
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

	// The user's spec tags merge over the metadata-derived tags -- user
	// tags deliberately win so an org's governance conventions can
	// override the derived values where they collide.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
