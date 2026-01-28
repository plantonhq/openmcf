package module

import (
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// exportOutputs exports the stack outputs for the Gateway API CRDs installation.
func exportOutputs(ctx *pulumi.Context, locals *Locals, crds *pulumiyaml.ConfigFile) error {
	// Export installed version
	ctx.Export("installed_version", pulumi.String(locals.Version))

	// Export installed channel
	ctx.Export("installed_channel", pulumi.String(locals.ChannelName))

	// Export list of installed CRDs
	var installedCrds []string
	if locals.IsExperimental {
		installedCrds = []string{
			"gatewayclasses.gateway.networking.k8s.io",
			"gateways.gateway.networking.k8s.io",
			"httproutes.gateway.networking.k8s.io",
			"referencegrants.gateway.networking.k8s.io",
			"tcproutes.gateway.networking.k8s.io",
			"udproutes.gateway.networking.k8s.io",
			"tlsroutes.gateway.networking.k8s.io",
			"grpcroutes.gateway.networking.k8s.io",
		}
	} else {
		installedCrds = []string{
			"gatewayclasses.gateway.networking.k8s.io",
			"gateways.gateway.networking.k8s.io",
			"httproutes.gateway.networking.k8s.io",
			"referencegrants.gateway.networking.k8s.io",
		}
	}

	crdArray := make(pulumi.StringArray, len(installedCrds))
	for i, crd := range installedCrds {
		crdArray[i] = pulumi.String(crd)
	}
	ctx.Export("installed_crds", crdArray)

	return nil
}
