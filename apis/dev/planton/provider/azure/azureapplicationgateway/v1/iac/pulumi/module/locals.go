package module

import (
	"strings"

	azureapplicationgatewayv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureapplicationgateway/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureApplicationGateway *azureapplicationgatewayv1.AzureApplicationGateway
	ResourceGroupName       string
	AzureTags               map[string]string
	GatewayIpConfigName     string
}

// The spec's enums arrive as proto enum values; each map below carries the
// complete vocabulary for its enum, translated to ARM's strings. A missing
// entry would silently drop the setting, so the maps are exhaustive by
// construction.

// On the v2 platform (and Basic), SKU name and tier carry the same value.
var skuStrings = map[azureapplicationgatewayv1.AzureApplicationGatewaySku]string{
	azureapplicationgatewayv1.AzureApplicationGatewaySku_BASIC:       "Basic",
	azureapplicationgatewayv1.AzureApplicationGatewaySku_STANDARD_V2: "Standard_v2",
	azureapplicationgatewayv1.AzureApplicationGatewaySku_WAF_V2:      "WAF_v2",
}

var protocolStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayProtocol]string{
	azureapplicationgatewayv1.AzureApplicationGatewayProtocol_HTTP:  "Http",
	azureapplicationgatewayv1.AzureApplicationGatewayProtocol_HTTPS: "Https",
	azureapplicationgatewayv1.AzureApplicationGatewayProtocol_TCP:   "Tcp",
	azureapplicationgatewayv1.AzureApplicationGatewayProtocol_TLS:   "Tls",
}

var ruleTypeStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayRequestRoutingRuleType]string{
	azureapplicationgatewayv1.AzureApplicationGatewayRequestRoutingRuleType_BASIC_ROUTING:      "Basic",
	azureapplicationgatewayv1.AzureApplicationGatewayRequestRoutingRuleType_PATH_BASED_ROUTING: "PathBasedRouting",
}

var ipAllocationStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayIpAllocation]string{
	azureapplicationgatewayv1.AzureApplicationGatewayIpAllocation_DYNAMIC: "Dynamic",
	azureapplicationgatewayv1.AzureApplicationGatewayIpAllocation_STATIC:  "Static",
}

var identityTypeStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayIdentityType]string{
	azureapplicationgatewayv1.AzureApplicationGatewayIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azureapplicationgatewayv1.AzureApplicationGatewayIdentityType_USER_ASSIGNED:            "UserAssigned",
	azureapplicationgatewayv1.AzureApplicationGatewayIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

var sslPolicyTypeStrings = map[azureapplicationgatewayv1.AzureApplicationGatewaySslPolicyType]string{
	azureapplicationgatewayv1.AzureApplicationGatewaySslPolicyType_PREDEFINED: "Predefined",
	azureapplicationgatewayv1.AzureApplicationGatewaySslPolicyType_CUSTOM:     "Custom",
	azureapplicationgatewayv1.AzureApplicationGatewaySslPolicyType_CUSTOM_V2:  "CustomV2",
}

var tlsProtocolStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayTlsProtocol]string{
	azureapplicationgatewayv1.AzureApplicationGatewayTlsProtocol_TLS_V1_0: "TLSv1_0",
	azureapplicationgatewayv1.AzureApplicationGatewayTlsProtocol_TLS_V1_1: "TLSv1_1",
	azureapplicationgatewayv1.AzureApplicationGatewayTlsProtocol_TLS_V1_2: "TLSv1_2",
	azureapplicationgatewayv1.AzureApplicationGatewayTlsProtocol_TLS_V1_3: "TLSv1_3",
}

var redirectTypeStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayRedirectType]string{
	azureapplicationgatewayv1.AzureApplicationGatewayRedirectType_PERMANENT: "Permanent",
	azureapplicationgatewayv1.AzureApplicationGatewayRedirectType_FOUND:     "Found",
	azureapplicationgatewayv1.AzureApplicationGatewayRedirectType_SEE_OTHER: "SeeOther",
	azureapplicationgatewayv1.AzureApplicationGatewayRedirectType_TEMPORARY: "Temporary",
}

var urlComponentStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayRewriteRuleUrlComponent]string{
	azureapplicationgatewayv1.AzureApplicationGatewayRewriteRuleUrlComponent_PATH_ONLY:         "path_only",
	azureapplicationgatewayv1.AzureApplicationGatewayRewriteRuleUrlComponent_QUERY_STRING_ONLY: "query_string_only",
}

var statusCodeStrings = map[azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode]string{
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_400: "HttpStatus400",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_403: "HttpStatus403",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_404: "HttpStatus404",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_405: "HttpStatus405",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_408: "HttpStatus408",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_500: "HttpStatus500",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_502: "HttpStatus502",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_503: "HttpStatus503",
	azureapplicationgatewayv1.AzureApplicationGatewayCustomErrorStatusCode_HTTP_STATUS_504: "HttpStatus504",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureapplicationgatewayv1.AzureApplicationGatewayStackInput) *Locals {
	locals := &Locals{}

	locals.AzureApplicationGateway = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// The gateway's single IP configuration, derived from the dedicated
	// subnet -- pure ARM plumbing users never name. Matches the Terraform
	// module's derivation.
	locals.GatewayIpConfigName = target.Spec.Name + "-gateway-ip"

	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureApplicationGateway.String()),
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
