package module

import (
	kubernetesapis "github.com/plantonhq/planton/catalog/kubernetes"
	kuberneteslistenersetv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kuberneteslistenerset/v1alpha1"
	gatewayv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/gatewayapis/kubernetes/gateway/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildListeners maps the Planton listener entries onto the typed crd2pulumi
// listener args array. Optional fields are only set when present so
// upstream/controller defaults flow through unchanged. The TLS and
// allowed-routes sub-shapes come from the shared Gateway API messages (the
// same ones KubernetesGateway listeners use), so the two kinds cannot drift.
func buildListeners(listeners []*kuberneteslistenersetv1alpha1.KubernetesListenerSetListener) gatewayv1.ListenerSetSpecListenersArray {
	arr := gatewayv1.ListenerSetSpecListenersArray{}
	for _, l := range listeners {
		args := gatewayv1.ListenerSetSpecListenersArgs{
			Name:     pulumi.String(l.GetName()),
			Port:     pulumi.Int(int(l.GetPort())),
			Protocol: pulumi.String(l.GetProtocol()),
		}
		if hostname := l.GetHostname(); hostname != "" {
			args.Hostname = pulumi.String(hostname)
		}
		if tls := l.GetTls(); tls != nil {
			args.Tls = buildListenerTls(tls)
		}
		if allowedRoutes := l.GetAllowedRoutes(); allowedRoutes != nil {
			args.AllowedRoutes = buildAllowedRoutes(allowedRoutes)
		}
		arr = append(arr, args)
	}
	return arr
}

func buildListenerTls(tls *kubernetesapis.KubernetesGatewayApiListenerTlsConfig) gatewayv1.ListenerSetSpecListenersTlsArgs {
	args := gatewayv1.ListenerSetSpecListenersTlsArgs{}
	if mode := tls.GetMode(); mode != "" {
		args.Mode = pulumi.String(mode)
	}
	if refs := tls.GetCertificateRefs(); len(refs) > 0 {
		certRefs := gatewayv1.ListenerSetSpecListenersTlsCertificateRefsArray{}
		for _, ref := range refs {
			certRefs = append(certRefs, buildListenerCertificateRef(ref))
		}
		args.CertificateRefs = certRefs
	}
	if options := tls.GetOptions(); len(options) > 0 {
		args.Options = pulumi.ToStringMap(options)
	}
	return args
}

// buildListenerCertificateRef maps a TLS certificate Secret reference. The
// reference's name is a KubernetesSecret foreign key resolved to its literal
// value before the module runs, so GetValue() returns the final Secret name.
func buildListenerCertificateRef(ref *kubernetesapis.KubernetesGatewayApiSecretObjectReference) gatewayv1.ListenerSetSpecListenersTlsCertificateRefsArgs {
	args := gatewayv1.ListenerSetSpecListenersTlsCertificateRefsArgs{
		Name: pulumi.String(ref.GetName().GetValue()),
	}
	if group := ref.GetGroup(); group != "" {
		args.Group = pulumi.String(group)
	}
	if kind := ref.GetKind(); kind != "" {
		args.Kind = pulumi.String(kind)
	}
	if namespace := ref.GetNamespace(); namespace != "" {
		args.Namespace = pulumi.String(namespace)
	}
	return args
}

func buildAllowedRoutes(allowedRoutes *kubernetesapis.KubernetesGatewayApiAllowedRoutes) gatewayv1.ListenerSetSpecListenersAllowedRoutesArgs {
	args := gatewayv1.ListenerSetSpecListenersAllowedRoutesArgs{}
	if namespaces := allowedRoutes.GetNamespaces(); namespaces != nil {
		nsArgs := gatewayv1.ListenerSetSpecListenersAllowedRoutesNamespacesArgs{}
		if from := namespaces.GetFrom(); from != "" {
			nsArgs.From = pulumi.String(from)
		}
		if selector := namespaces.GetSelector(); selector != nil {
			nsArgs.Selector = buildAllowedRoutesSelector(selector)
		}
		args.Namespaces = nsArgs
	}
	if kinds := allowedRoutes.GetKinds(); len(kinds) > 0 {
		kindArr := gatewayv1.ListenerSetSpecListenersAllowedRoutesKindsArray{}
		for _, k := range kinds {
			kindArgs := gatewayv1.ListenerSetSpecListenersAllowedRoutesKindsArgs{
				Kind: pulumi.String(k.GetKind()),
			}
			if group := k.GetGroup(); group != "" {
				kindArgs.Group = pulumi.String(group)
			}
			kindArr = append(kindArr, kindArgs)
		}
		args.Kinds = kindArr
	}
	return args
}

func buildAllowedRoutesSelector(selector *kubernetesapis.KubernetesGatewayApiLabelSelector) gatewayv1.ListenerSetSpecListenersAllowedRoutesNamespacesSelectorArgs {
	args := gatewayv1.ListenerSetSpecListenersAllowedRoutesNamespacesSelectorArgs{}
	if matchLabels := selector.GetMatchLabels(); len(matchLabels) > 0 {
		args.MatchLabels = pulumi.ToStringMap(matchLabels)
	}
	if expressions := selector.GetMatchExpressions(); len(expressions) > 0 {
		exprArr := gatewayv1.ListenerSetSpecListenersAllowedRoutesNamespacesSelectorMatchExpressionsArray{}
		for _, e := range expressions {
			exprArgs := gatewayv1.ListenerSetSpecListenersAllowedRoutesNamespacesSelectorMatchExpressionsArgs{
				Key:      pulumi.String(e.GetKey()),
				Operator: pulumi.String(e.GetOperator()),
			}
			if values := e.GetValues(); len(values) > 0 {
				exprArgs.Values = pulumi.ToStringArray(values)
			}
			exprArr = append(exprArr, exprArgs)
		}
		args.MatchExpressions = exprArr
	}
	return args
}
