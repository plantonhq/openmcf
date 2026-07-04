package module

import (
	"strings"

	azurepostgresqlflexibleserverv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurepostgresqlflexibleserver/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzurePostgresqlFlexibleServer *azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServer
	ResourceGroupName             string
	AzureTags                     map[string]string
}

// createModeStrings maps the spec's create-mode enum to ARM's values. An
// unspecified mode means a fresh (DEFAULT) server and is not sent at all,
// mirroring the Terraform module's null -- so the two engines produce the
// same ARM payload.
var createModeStrings = map[azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerCreateMode]string{
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerCreateMode_DEFAULT:               "Default",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerCreateMode_POINT_IN_TIME_RESTORE: "PointInTimeRestore",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerCreateMode_REPLICA:               "Replica",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerCreateMode_GEO_RESTORE:           "GeoRestore",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerCreateMode_REVIVE_DROPPED:        "ReviveDropped",
}

// haModeStrings maps the spec's high-availability mode enum to ARM's values.
var haModeStrings = map[azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerHighAvailabilityMode]string{
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerHighAvailabilityMode_ZONE_REDUNDANT: "ZoneRedundant",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerHighAvailabilityMode_SAME_ZONE:      "SameZone",
}

// identityTypeStrings maps the spec's identity-type enum to ARM's values.
var identityTypeStrings = map[azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerIdentityType]string{
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerIdentityType_USER_ASSIGNED:            "UserAssigned",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

// principalTypeStrings maps the spec's Entra principal-type enum to ARM's
// values.
var principalTypeStrings = map[azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerAadPrincipalType]string{
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerAadPrincipalType_USER:              "User",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerAadPrincipalType_GROUP:             "Group",
	azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerAadPrincipalType_SERVICE_PRINCIPAL: "ServicePrincipal",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurepostgresqlflexibleserverv1.AzurePostgresqlFlexibleServerStackInput) *Locals {
	locals := &Locals{}

	locals.AzurePostgresqlFlexibleServer = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzurePostgresqlFlexibleServer.String()),
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
