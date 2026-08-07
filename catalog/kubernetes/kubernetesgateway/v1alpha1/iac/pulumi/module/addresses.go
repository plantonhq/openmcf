package module

import (
	kubernetesgatewayv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesgateway/v1alpha1"
	gatewayv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/gatewayapis/kubernetes/gateway/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func buildAddresses(addresses []*kubernetesgatewayv1alpha1.KubernetesGatewayAddress) gatewayv1.GatewaySpecAddressesArray {
	arr := gatewayv1.GatewaySpecAddressesArray{}
	for _, a := range addresses {
		args := gatewayv1.GatewaySpecAddressesArgs{}
		if t := a.GetType(); t != "" {
			args.Type = pulumi.String(t)
		}
		if v := a.GetValue(); v != "" {
			args.Value = pulumi.String(v)
		}
		arr = append(arr, args)
	}
	return arr
}
