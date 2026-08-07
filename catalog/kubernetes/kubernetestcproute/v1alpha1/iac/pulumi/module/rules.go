package module

import (
	kubernetesapis "github.com/plantonhq/planton/catalog/kubernetes"
	kubernetestcproutev1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetestcproute/v1alpha1"
	gatewayv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/gatewayapis/kubernetes/gateway/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildRules maps the Planton TCP route rules onto the typed crd2pulumi rules
// array. A TCP route rule carries only an optional name and the backend refs
// (no matches, no filters).
func buildRules(rules []*kubernetestcproutev1alpha1.KubernetesTcpRouteRule) gatewayv1.TCPRouteSpecRulesArray {
	arr := gatewayv1.TCPRouteSpecRulesArray{}
	for _, r := range rules {
		args := gatewayv1.TCPRouteSpecRulesArgs{}
		if name := r.GetName(); name != "" {
			args.Name = pulumi.String(name)
		}
		if backendRefs := r.GetBackendRefs(); len(backendRefs) > 0 {
			args.BackendRefs = buildBackendRefs(backendRefs)
		}
		arr = append(arr, args)
	}
	return arr
}

// buildBackendRefs maps the shared KubernetesGatewayApiBackendRef (group / kind /
// name / namespace / port / weight) onto the typed crd2pulumi backendRefs array.
// Optional fields are only set when present so controller defaults flow through.
// TCP routes have no per-backend filters.
//
// Each backend's name is a KubernetesService foreign key resolved to its
// literal value before the module runs, so GetValue() returns the final name.
func buildBackendRefs(backendRefs []*kubernetesapis.KubernetesGatewayApiBackendRef) gatewayv1.TCPRouteSpecRulesBackendRefsArray {
	arr := gatewayv1.TCPRouteSpecRulesBackendRefsArray{}
	for _, b := range backendRefs {
		args := gatewayv1.TCPRouteSpecRulesBackendRefsArgs{
			Name: pulumi.String(b.GetName().GetValue()),
		}
		if group := b.GetGroup(); group != "" {
			args.Group = pulumi.String(group)
		}
		if kind := b.GetKind(); kind != "" {
			args.Kind = pulumi.String(kind)
		}
		if namespace := b.GetNamespace(); namespace != "" {
			args.Namespace = pulumi.String(namespace)
		}
		if b.Port != nil {
			args.Port = pulumi.Int(int(b.GetPort()))
		}
		if b.Weight != nil {
			args.Weight = pulumi.Int(int(b.GetWeight()))
		}
		arr = append(arr, args)
	}
	return arr
}
