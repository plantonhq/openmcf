package module

import (
	"strings"

	azurelinuxwebappv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurelinuxwebapp/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureLinuxWebApp  *azurelinuxwebappv1.AzureLinuxWebApp
	ResourceGroupName string
	AzureTags         map[string]string
}

// Spec enums -> Azure's wire values. Spelled out row by row so a
// vocabulary drift fails loudly at preview time instead of deploying a
// wrong value.

var clientCertificateModeStrings = map[azurelinuxwebappv1.AzureLinuxWebAppClientCertificateMode]string{
	azurelinuxwebappv1.AzureLinuxWebAppClientCertificateMode_REQUIRED:                  "Required",
	azurelinuxwebappv1.AzureLinuxWebAppClientCertificateMode_OPTIONAL:                  "Optional",
	azurelinuxwebappv1.AzureLinuxWebAppClientCertificateMode_OPTIONAL_INTERACTIVE_USER: "OptionalInteractiveUser",
}

var tlsVersionStrings = map[azurelinuxwebappv1.AzureLinuxWebAppTlsVersion]string{
	azurelinuxwebappv1.AzureLinuxWebAppTlsVersion_TLS_1_0: "1.0",
	azurelinuxwebappv1.AzureLinuxWebAppTlsVersion_TLS_1_1: "1.1",
	azurelinuxwebappv1.AzureLinuxWebAppTlsVersion_TLS_1_2: "1.2",
	azurelinuxwebappv1.AzureLinuxWebAppTlsVersion_TLS_1_3: "1.3",
}

var ftpsStateStrings = map[azurelinuxwebappv1.AzureLinuxWebAppFtpsState]string{
	azurelinuxwebappv1.AzureLinuxWebAppFtpsState_ALL_ALLOWED: "AllAllowed",
	azurelinuxwebappv1.AzureLinuxWebAppFtpsState_FTPS_ONLY:   "FtpsOnly",
	azurelinuxwebappv1.AzureLinuxWebAppFtpsState_DISABLED:    "Disabled",
}

var loadBalancingModeStrings = map[azurelinuxwebappv1.AzureLinuxWebAppLoadBalancingMode]string{
	azurelinuxwebappv1.AzureLinuxWebAppLoadBalancingMode_LEAST_REQUESTS:         "LeastRequests",
	azurelinuxwebappv1.AzureLinuxWebAppLoadBalancingMode_WEIGHTED_ROUND_ROBIN:   "WeightedRoundRobin",
	azurelinuxwebappv1.AzureLinuxWebAppLoadBalancingMode_LEAST_RESPONSE_TIME:    "LeastResponseTime",
	azurelinuxwebappv1.AzureLinuxWebAppLoadBalancingMode_WEIGHTED_TOTAL_TRAFFIC: "WeightedTotalTraffic",
	azurelinuxwebappv1.AzureLinuxWebAppLoadBalancingMode_REQUEST_HASH:           "RequestHash",
	azurelinuxwebappv1.AzureLinuxWebAppLoadBalancingMode_PER_SITE_ROUND_ROBIN:   "PerSiteRoundRobin",
}

var managedPipelineModeStrings = map[azurelinuxwebappv1.AzureLinuxWebAppManagedPipelineMode]string{
	azurelinuxwebappv1.AzureLinuxWebAppManagedPipelineMode_INTEGRATED: "Integrated",
	azurelinuxwebappv1.AzureLinuxWebAppManagedPipelineMode_CLASSIC:    "Classic",
}

var ipRestrictionActionStrings = map[azurelinuxwebappv1.AzureLinuxWebAppIpRestrictionAction]string{
	azurelinuxwebappv1.AzureLinuxWebAppIpRestrictionAction_ALLOW: "Allow",
	azurelinuxwebappv1.AzureLinuxWebAppIpRestrictionAction_DENY:  "Deny",
}

var javaServerStrings = map[azurelinuxwebappv1.AzureLinuxWebAppJavaServer]string{
	azurelinuxwebappv1.AzureLinuxWebAppJavaServer_JAVA_SE:  "JAVA",
	azurelinuxwebappv1.AzureLinuxWebAppJavaServer_TOMCAT:   "TOMCAT",
	azurelinuxwebappv1.AzureLinuxWebAppJavaServer_JBOSSEAP: "JBOSSEAP",
}

var connectionStringTypeStrings = map[azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType]string{
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_MYSQL:            "MySQL",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_SQL_SERVER:       "SQLServer",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_SQL_AZURE:        "SQLAzure",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_CUSTOM:           "Custom",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_NOTIFICATION_HUB: "NotificationHub",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_SERVICE_BUS:      "ServiceBus",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_EVENT_HUB:        "EventHub",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_API_HUB:          "APIHub",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_DOC_DB:           "DocDb",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_REDIS_CACHE:      "RedisCache",
	azurelinuxwebappv1.AzureLinuxWebAppConnectionStringType_POSTGRESQL:       "PostgreSQL",
}

var storageMountTypeStrings = map[azurelinuxwebappv1.AzureLinuxWebAppStorageMountType]string{
	azurelinuxwebappv1.AzureLinuxWebAppStorageMountType_AZURE_FILES: "AzureFiles",
	azurelinuxwebappv1.AzureLinuxWebAppStorageMountType_AZURE_BLOB:  "AzureBlob",
}

var logLevelStrings = map[azurelinuxwebappv1.AzureLinuxWebAppLogLevel]string{
	azurelinuxwebappv1.AzureLinuxWebAppLogLevel_OFF:         "Off",
	azurelinuxwebappv1.AzureLinuxWebAppLogLevel_ERROR:       "Error",
	azurelinuxwebappv1.AzureLinuxWebAppLogLevel_WARNING:     "Warning",
	azurelinuxwebappv1.AzureLinuxWebAppLogLevel_INFORMATION: "Information",
	azurelinuxwebappv1.AzureLinuxWebAppLogLevel_VERBOSE:     "Verbose",
}

var identityTypeStrings = map[azurelinuxwebappv1.AzureLinuxWebAppIdentityType]string{
	azurelinuxwebappv1.AzureLinuxWebAppIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azurelinuxwebappv1.AzureLinuxWebAppIdentityType_USER_ASSIGNED:            "UserAssigned",
	azurelinuxwebappv1.AzureLinuxWebAppIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

var backupFrequencyUnitStrings = map[azurelinuxwebappv1.AzureLinuxWebAppBackupFrequencyUnit]string{
	azurelinuxwebappv1.AzureLinuxWebAppBackupFrequencyUnit_DAY:  "Day",
	azurelinuxwebappv1.AzureLinuxWebAppBackupFrequencyUnit_HOUR: "Hour",
}

var unauthenticatedActionStrings = map[azurelinuxwebappv1.AzureLinuxWebAppUnauthenticatedAction]string{
	azurelinuxwebappv1.AzureLinuxWebAppUnauthenticatedAction_REDIRECT_TO_LOGIN_PAGE: "RedirectToLoginPage",
	azurelinuxwebappv1.AzureLinuxWebAppUnauthenticatedAction_ALLOW_ANONYMOUS:        "AllowAnonymous",
	azurelinuxwebappv1.AzureLinuxWebAppUnauthenticatedAction_RETURN_401:             "Return401",
	azurelinuxwebappv1.AzureLinuxWebAppUnauthenticatedAction_RETURN_403:             "Return403",
}

var forwardProxyConventionStrings = map[azurelinuxwebappv1.AzureLinuxWebAppForwardProxyConvention]string{
	azurelinuxwebappv1.AzureLinuxWebAppForwardProxyConvention_FORWARD_PROXY_NO_PROXY: "NoProxy",
	azurelinuxwebappv1.AzureLinuxWebAppForwardProxyConvention_FORWARD_PROXY_STANDARD: "Standard",
	azurelinuxwebappv1.AzureLinuxWebAppForwardProxyConvention_FORWARD_PROXY_CUSTOM:   "Custom",
}

var cookieExpirationConventionStrings = map[azurelinuxwebappv1.AzureLinuxWebAppCookieExpirationConvention]string{
	azurelinuxwebappv1.AzureLinuxWebAppCookieExpirationConvention_FIXED_TIME:                "FixedTime",
	azurelinuxwebappv1.AzureLinuxWebAppCookieExpirationConvention_IDENTITY_PROVIDER_DERIVED: "IdentityProviderDerived",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurelinuxwebappv1.AzureLinuxWebAppStackInput) *Locals {
	locals := &Locals{}

	locals.AzureLinuxWebApp = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureLinuxWebApp.String()),
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
