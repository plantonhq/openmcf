package module

import (
	"strings"

	azurerediscachev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurerediscache/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureRedisCache   *azurerediscachev1.AzureRedisCache
	ResourceGroupName string
	AzureTags         map[string]string
	// SkuName is ARM's tier value, materialized from the spec enum with
	// the documented STANDARD default (stack inputs never carry proto
	// defaults).
	SkuName string
	// Family is Azure's size-family letter, fully determined by the tier:
	// "C" for Basic/Standard, "P" for Premium -- never spelled twice.
	Family string
}

// skuStrings maps the spec's sku enum to ARM's tier values.
var skuStrings = map[azurerediscachev1.AzureRedisCacheSku]string{
	azurerediscachev1.AzureRedisCacheSku_BASIC:    "Basic",
	azurerediscachev1.AzureRedisCacheSku_STANDARD: "Standard",
	azurerediscachev1.AzureRedisCacheSku_PREMIUM:  "Premium",
}

// dayOfWeekStrings maps the patch-schedule day enum to ARM's capitalized
// English day names.
var dayOfWeekStrings = map[azurerediscachev1.AzureRedisCachePatchScheduleDay]string{
	azurerediscachev1.AzureRedisCachePatchScheduleDay_MONDAY:    "Monday",
	azurerediscachev1.AzureRedisCachePatchScheduleDay_TUESDAY:   "Tuesday",
	azurerediscachev1.AzureRedisCachePatchScheduleDay_WEDNESDAY: "Wednesday",
	azurerediscachev1.AzureRedisCachePatchScheduleDay_THURSDAY:  "Thursday",
	azurerediscachev1.AzureRedisCachePatchScheduleDay_FRIDAY:    "Friday",
	azurerediscachev1.AzureRedisCachePatchScheduleDay_SATURDAY:  "Saturday",
	azurerediscachev1.AzureRedisCachePatchScheduleDay_SUNDAY:    "Sunday",
}

// persistenceAuthStrings maps the persistence auth enum to ARM's values.
var persistenceAuthStrings = map[azurerediscachev1.AzureRedisCachePersistenceAuthMethod]string{
	azurerediscachev1.AzureRedisCachePersistenceAuthMethod_SAS:              "SAS",
	azurerediscachev1.AzureRedisCachePersistenceAuthMethod_MANAGED_IDENTITY: "ManagedIdentity",
}

// identityTypeStrings maps the identity-type enum to ARM's values.
var identityTypeStrings = map[azurerediscachev1.AzureRedisCacheIdentityType]string{
	azurerediscachev1.AzureRedisCacheIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azurerediscachev1.AzureRedisCacheIdentityType_USER_ASSIGNED:            "UserAssigned",
	azurerediscachev1.AzureRedisCacheIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurerediscachev1.AzureRedisCacheStackInput) *Locals {
	locals := &Locals{}

	locals.AzureRedisCache = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Materialize the tier default: unspecified deploys STANDARD (the
	// spec's documented default -- stack inputs never carry proto
	// defaults), then derive the size-family letter from the tier.
	locals.SkuName = skuStrings[target.Spec.SkuName]
	if locals.SkuName == "" {
		locals.SkuName = "Standard"
	}
	locals.Family = "C"
	if locals.SkuName == "Premium" {
		locals.Family = "P"
	}

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureRedisCache.String()),
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

	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
