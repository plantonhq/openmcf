package module

import (
	"strings"

	azurefunctionappv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefunctionapp/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFunctionApp  *azurefunctionappv1.AzureFunctionApp
	ResourceGroupName string
	AzureTags         map[string]string
}

// Spec enums -> Azure's wire values. Spelled out row by row so a
// vocabulary drift fails loudly at preview time instead of deploying a
// wrong value.

var clientCertificateModeStrings = map[azurefunctionappv1.AzureFunctionAppClientCertificateMode]string{
	azurefunctionappv1.AzureFunctionAppClientCertificateMode_REQUIRED:                  "Required",
	azurefunctionappv1.AzureFunctionAppClientCertificateMode_OPTIONAL:                  "Optional",
	azurefunctionappv1.AzureFunctionAppClientCertificateMode_OPTIONAL_INTERACTIVE_USER: "OptionalInteractiveUser",
}

var tlsVersionStrings = map[azurefunctionappv1.AzureFunctionAppTlsVersion]string{
	azurefunctionappv1.AzureFunctionAppTlsVersion_TLS_1_0: "1.0",
	azurefunctionappv1.AzureFunctionAppTlsVersion_TLS_1_1: "1.1",
	azurefunctionappv1.AzureFunctionAppTlsVersion_TLS_1_2: "1.2",
	azurefunctionappv1.AzureFunctionAppTlsVersion_TLS_1_3: "1.3",
}

var ftpsStateStrings = map[azurefunctionappv1.AzureFunctionAppFtpsState]string{
	azurefunctionappv1.AzureFunctionAppFtpsState_ALL_ALLOWED: "AllAllowed",
	azurefunctionappv1.AzureFunctionAppFtpsState_FTPS_ONLY:   "FtpsOnly",
	azurefunctionappv1.AzureFunctionAppFtpsState_DISABLED:    "Disabled",
}

var loadBalancingModeStrings = map[azurefunctionappv1.AzureFunctionAppLoadBalancingMode]string{
	azurefunctionappv1.AzureFunctionAppLoadBalancingMode_LEAST_REQUESTS:         "LeastRequests",
	azurefunctionappv1.AzureFunctionAppLoadBalancingMode_WEIGHTED_ROUND_ROBIN:   "WeightedRoundRobin",
	azurefunctionappv1.AzureFunctionAppLoadBalancingMode_LEAST_RESPONSE_TIME:    "LeastResponseTime",
	azurefunctionappv1.AzureFunctionAppLoadBalancingMode_WEIGHTED_TOTAL_TRAFFIC: "WeightedTotalTraffic",
	azurefunctionappv1.AzureFunctionAppLoadBalancingMode_REQUEST_HASH:           "RequestHash",
	azurefunctionappv1.AzureFunctionAppLoadBalancingMode_PER_SITE_ROUND_ROBIN:   "PerSiteRoundRobin",
}

var managedPipelineModeStrings = map[azurefunctionappv1.AzureFunctionAppManagedPipelineMode]string{
	azurefunctionappv1.AzureFunctionAppManagedPipelineMode_INTEGRATED: "Integrated",
	azurefunctionappv1.AzureFunctionAppManagedPipelineMode_CLASSIC:    "Classic",
}

var ipRestrictionActionStrings = map[azurefunctionappv1.AzureFunctionAppIpRestrictionAction]string{
	azurefunctionappv1.AzureFunctionAppIpRestrictionAction_ALLOW: "Allow",
	azurefunctionappv1.AzureFunctionAppIpRestrictionAction_DENY:  "Deny",
}

var connectionStringTypeStrings = map[azurefunctionappv1.AzureFunctionAppConnectionStringType]string{
	azurefunctionappv1.AzureFunctionAppConnectionStringType_MYSQL:            "MySQL",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_SQL_SERVER:       "SQLServer",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_SQL_AZURE:        "SQLAzure",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_CUSTOM:           "Custom",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_NOTIFICATION_HUB: "NotificationHub",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_SERVICE_BUS:      "ServiceBus",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_EVENT_HUB:        "EventHub",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_API_HUB:          "APIHub",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_DOC_DB:           "DocDb",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_REDIS_CACHE:      "RedisCache",
	azurefunctionappv1.AzureFunctionAppConnectionStringType_POSTGRESQL:       "PostgreSQL",
}

var storageMountTypeStrings = map[azurefunctionappv1.AzureFunctionAppStorageMountType]string{
	azurefunctionappv1.AzureFunctionAppStorageMountType_AZURE_FILES: "AzureFiles",
	azurefunctionappv1.AzureFunctionAppStorageMountType_AZURE_BLOB:  "AzureBlob",
}

var identityTypeStrings = map[azurefunctionappv1.AzureFunctionAppIdentityType]string{
	azurefunctionappv1.AzureFunctionAppIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azurefunctionappv1.AzureFunctionAppIdentityType_USER_ASSIGNED:            "UserAssigned",
	azurefunctionappv1.AzureFunctionAppIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

var backupFrequencyUnitStrings = map[azurefunctionappv1.AzureFunctionAppBackupFrequencyUnit]string{
	azurefunctionappv1.AzureFunctionAppBackupFrequencyUnit_DAY:  "Day",
	azurefunctionappv1.AzureFunctionAppBackupFrequencyUnit_HOUR: "Hour",
}

var unauthenticatedActionStrings = map[azurefunctionappv1.AzureFunctionAppUnauthenticatedAction]string{
	azurefunctionappv1.AzureFunctionAppUnauthenticatedAction_REDIRECT_TO_LOGIN_PAGE: "RedirectToLoginPage",
	azurefunctionappv1.AzureFunctionAppUnauthenticatedAction_ALLOW_ANONYMOUS:        "AllowAnonymous",
	azurefunctionappv1.AzureFunctionAppUnauthenticatedAction_RETURN_401:             "Return401",
	azurefunctionappv1.AzureFunctionAppUnauthenticatedAction_RETURN_403:             "Return403",
}

var forwardProxyConventionStrings = map[azurefunctionappv1.AzureFunctionAppForwardProxyConvention]string{
	azurefunctionappv1.AzureFunctionAppForwardProxyConvention_FORWARD_PROXY_NO_PROXY: "NoProxy",
	azurefunctionappv1.AzureFunctionAppForwardProxyConvention_FORWARD_PROXY_STANDARD: "Standard",
	azurefunctionappv1.AzureFunctionAppForwardProxyConvention_FORWARD_PROXY_CUSTOM:   "Custom",
}

var cookieExpirationConventionStrings = map[azurefunctionappv1.AzureFunctionAppCookieExpirationConvention]string{
	azurefunctionappv1.AzureFunctionAppCookieExpirationConvention_FIXED_TIME:                "FixedTime",
	azurefunctionappv1.AzureFunctionAppCookieExpirationConvention_IDENTITY_PROVIDER_DERIVED: "IdentityProviderDerived",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefunctionappv1.AzureFunctionAppStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFunctionApp = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureFunctionApp.String()),
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
