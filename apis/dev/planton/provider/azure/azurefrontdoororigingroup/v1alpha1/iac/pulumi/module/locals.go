package module

import (
	azurefrontdoororigingroupv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefrontdoororigingroup/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFrontDoorOriginGroup *azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroup
	ProfileId                 string
}

// healthProbeProtocolStrings maps the probe-protocol enum to ARM's values.
var healthProbeProtocolStrings = map[azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroupHealthProbeProtocol]string{
	azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroupHealthProbeProtocol_HTTP:  "Http",
	azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroupHealthProbeProtocol_HTTPS: "Https",
}

// healthProbeRequestTypeStrings maps the probe-method enum to ARM's
// values (unspecified deploys HEAD, Azure's default).
var healthProbeRequestTypeStrings = map[azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroupHealthProbeRequestType]string{
	azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroupHealthProbeRequestType_HEAD: "HEAD",
	azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroupHealthProbeRequestType_GET:  "GET",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefrontdoororigingroupv1alpha1.AzureFrontDoorOriginGroupStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFrontDoorOriginGroup = stackInput.Target
	locals.ProfileId = stackInput.Target.Spec.ProfileId.GetValue()

	// No Azure tags: ARM does not support tags on Front Door origin
	// groups, so the platform's identity tags live on the profile.

	return locals
}
