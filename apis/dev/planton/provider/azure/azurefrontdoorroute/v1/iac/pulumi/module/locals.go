package module

import (
	azurefrontdoorroutev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefrontdoorroute/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFrontDoorRoute *azurefrontdoorroutev1.AzureFrontDoorRoute
	EndpointId          string
	OriginGroupId       string
	OriginIds           []string
}

// protocolStrings maps the client-facing protocol enum to ARM's values.
var protocolStrings = map[azurefrontdoorroutev1.AzureFrontDoorRouteProtocol]string{
	azurefrontdoorroutev1.AzureFrontDoorRouteProtocol_HTTP:  "Http",
	azurefrontdoorroutev1.AzureFrontDoorRouteProtocol_HTTPS: "Https",
}

// forwardingProtocolStrings maps the origin-leg protocol enum to ARM's
// values (unspecified deploys MatchRequest, Azure's default).
var forwardingProtocolStrings = map[azurefrontdoorroutev1.AzureFrontDoorRouteForwardingProtocol]string{
	azurefrontdoorroutev1.AzureFrontDoorRouteForwardingProtocol_MATCH_REQUEST: "MatchRequest",
	azurefrontdoorroutev1.AzureFrontDoorRouteForwardingProtocol_HTTP_ONLY:     "HttpOnly",
	azurefrontdoorroutev1.AzureFrontDoorRouteForwardingProtocol_HTTPS_ONLY:    "HttpsOnly",
}

// queryStringCachingBehaviorStrings maps the cache-key enum to ARM's
// values (unspecified deploys IgnoreQueryString, Azure's default).
var queryStringCachingBehaviorStrings = map[azurefrontdoorroutev1.AzureFrontDoorRouteQueryStringCachingBehavior]string{
	azurefrontdoorroutev1.AzureFrontDoorRouteQueryStringCachingBehavior_IGNORE_QUERY_STRING:             "IgnoreQueryString",
	azurefrontdoorroutev1.AzureFrontDoorRouteQueryStringCachingBehavior_USE_QUERY_STRING:                "UseQueryString",
	azurefrontdoorroutev1.AzureFrontDoorRouteQueryStringCachingBehavior_IGNORE_SPECIFIED_QUERY_STRINGS:  "IgnoreSpecifiedQueryStrings",
	azurefrontdoorroutev1.AzureFrontDoorRouteQueryStringCachingBehavior_INCLUDE_SPECIFIED_QUERY_STRINGS: "IncludeSpecifiedQueryStrings",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefrontdoorroutev1.AzureFrontDoorRouteStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFrontDoorRoute = stackInput.Target
	spec := stackInput.Target.Spec

	locals.EndpointId = spec.EndpointId.GetValue()
	locals.OriginGroupId = spec.OriginGroupId.GetValue()

	for _, originId := range spec.OriginIds {
		locals.OriginIds = append(locals.OriginIds, originId.GetValue())
	}

	// No Azure tags: ARM does not support tags on Front Door routes, so
	// the platform's identity tags live on the profile.

	return locals
}
