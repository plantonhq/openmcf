package module

import (
	azurefrontdoororigingroupv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefrontdoororigingroup/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFrontDoorOriginGroup *azurefrontdoororigingroupv1.AzureFrontDoorOriginGroup
	ProfileId                 string
}

// healthProbeProtocolStrings maps the probe-protocol enum to ARM's values.
var healthProbeProtocolStrings = map[azurefrontdoororigingroupv1.AzureFrontDoorOriginGroupHealthProbeProtocol]string{
	azurefrontdoororigingroupv1.AzureFrontDoorOriginGroupHealthProbeProtocol_HTTP:  "Http",
	azurefrontdoororigingroupv1.AzureFrontDoorOriginGroupHealthProbeProtocol_HTTPS: "Https",
}

// healthProbeRequestTypeStrings maps the probe-method enum to ARM's
// values (unspecified deploys HEAD, Azure's default).
var healthProbeRequestTypeStrings = map[azurefrontdoororigingroupv1.AzureFrontDoorOriginGroupHealthProbeRequestType]string{
	azurefrontdoororigingroupv1.AzureFrontDoorOriginGroupHealthProbeRequestType_HEAD: "HEAD",
	azurefrontdoororigingroupv1.AzureFrontDoorOriginGroupHealthProbeRequestType_GET:  "GET",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefrontdoororigingroupv1.AzureFrontDoorOriginGroupStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFrontDoorOriginGroup = stackInput.Target
	locals.ProfileId = stackInput.Target.Spec.ProfileId.GetValue()

	// No Azure tags: ARM does not support tags on Front Door origin
	// groups, so the platform's identity tags live on the profile.

	return locals
}
