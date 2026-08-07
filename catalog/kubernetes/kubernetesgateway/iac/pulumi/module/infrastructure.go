package module

import (
	kubernetesgatewayv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesgateway/v1alpha1"
	gatewayv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/gatewayapis/kubernetes/gateway/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func buildInfrastructure(infra *kubernetesgatewayv1alpha1.KubernetesGatewayInfrastructure) gatewayv1.GatewaySpecInfrastructureArgs {
	args := gatewayv1.GatewaySpecInfrastructureArgs{}
	if labels := infra.GetLabels(); len(labels) > 0 {
		args.Labels = pulumi.ToStringMap(labels)
	}
	if annotations := infra.GetAnnotations(); len(annotations) > 0 {
		args.Annotations = pulumi.ToStringMap(annotations)
	}
	if paramsRef := infra.GetParametersRef(); paramsRef != nil {
		args.ParametersRef = gatewayv1.GatewaySpecInfrastructureParametersRefArgs{
			Group: pulumi.String(paramsRef.GetGroup()),
			Kind:  pulumi.String(paramsRef.GetKind()),
			Name:  pulumi.String(paramsRef.GetName()),
		}
	}
	return args
}

func buildAllowedListeners(allowedListeners *kubernetesgatewayv1alpha1.KubernetesGatewayAllowedListeners) gatewayv1.GatewaySpecAllowedListenersArgs {
	args := gatewayv1.GatewaySpecAllowedListenersArgs{}
	if namespaces := allowedListeners.GetNamespaces(); namespaces != nil {
		nsArgs := gatewayv1.GatewaySpecAllowedListenersNamespacesArgs{}
		if from := namespaces.GetFrom(); from != "" {
			nsArgs.From = pulumi.String(from)
		}
		if selector := namespaces.GetSelector(); selector != nil {
			nsArgs.Selector = buildAllowedListenersSelector(selector)
		}
		args.Namespaces = nsArgs
	}
	return args
}
