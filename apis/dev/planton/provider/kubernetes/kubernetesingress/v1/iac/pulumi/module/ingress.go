package module

import (
	"github.com/pkg/errors"
	kubernetesingressv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesingress/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	kubernetesnetworkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createIngress creates the networking/v1 Ingress resource.
//
// Creation deliberately does NOT block on a controller claiming the Ingress
// (the skipAwait annotation below; the Terraform module's
// wait_for_load_balancer=false is the exact same choice). An Ingress object is
// valid without a controller — infra charts routinely deploy the workload and
// its exposure before the ingress controller wave — and blocking every deploy
// until a controller populates the load-balancer status would couple this kind
// to cluster addon ordering. The load-balancer address handles surface through
// the outputs as soon as a controller reconciles the object.
func createIngress(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetesnetworkingv1.Ingress, error) {
	spec := locals.Spec

	ingressSpecArgs := &kubernetesnetworkingv1.IngressSpecArgs{}

	if spec.GetIngressClassName() != "" {
		ingressSpecArgs.IngressClassName = pulumi.String(spec.GetIngressClassName())
	}

	if spec.GetDefaultBackend() != nil {
		ingressSpecArgs.DefaultBackend = buildBackend(spec.GetDefaultBackend())
	}

	if len(spec.GetTls()) > 0 {
		tlsArray := kubernetesnetworkingv1.IngressTLSArray{}
		for _, t := range spec.GetTls() {
			tlsArgs := &kubernetesnetworkingv1.IngressTLSArgs{}
			if len(t.GetHosts()) > 0 {
				tlsArgs.Hosts = pulumi.ToStringArray(t.GetHosts())
			}
			// secret_name is a StringValueOrRef foreign key, resolved to the
			// literal Secret name before the module runs. With cert-manager the
			// Secret need not exist yet — the issuer annotation instructs
			// cert-manager to create it under exactly this name.
			if t.GetSecretName().GetValue() != "" {
				tlsArgs.SecretName = pulumi.String(t.GetSecretName().GetValue())
			}
			tlsArray = append(tlsArray, tlsArgs)
		}
		ingressSpecArgs.Tls = tlsArray
	}

	if len(spec.GetRules()) > 0 {
		ruleArray := kubernetesnetworkingv1.IngressRuleArray{}
		for _, rule := range spec.GetRules() {
			ruleArgs := &kubernetesnetworkingv1.IngressRuleArgs{}
			if rule.GetHost() != "" {
				ruleArgs.Host = pulumi.String(rule.GetHost())
			}
			pathArray := kubernetesnetworkingv1.HTTPIngressPathArray{}
			for _, p := range rule.GetPaths() {
				pathArgs := &kubernetesnetworkingv1.HTTPIngressPathArgs{
					PathType: pulumi.String(resolvePathType(p.GetPathType())),
					Backend:  buildBackend(p.GetBackend()),
				}
				if p.GetPath() != "" {
					pathArgs.Path = pulumi.String(p.GetPath())
				}
				pathArray = append(pathArray, pathArgs)
			}
			ruleArgs.Http = &kubernetesnetworkingv1.HTTPIngressRuleValueArgs{
				Paths: pathArray,
			}
			ruleArray = append(ruleArray, ruleArgs)
		}
		ingressSpecArgs.Rules = ruleArray
	}

	// skipAwait rides the object annotations but is Pulumi engine metadata, not
	// user configuration — it is added here rather than in locals so the
	// user-facing annotation set stays exactly what the spec declared.
	annotations := make(map[string]string, len(locals.Annotations)+1)
	for k, v := range locals.Annotations {
		annotations[k] = v
	}
	annotations["pulumi.com/skipAwait"] = "true"

	ingress, err := kubernetesnetworkingv1.NewIngress(
		ctx,
		locals.Name,
		&kubernetesnetworkingv1.IngressArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.Name),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(annotations),
			},
			Spec: ingressSpecArgs,
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create ingress %s/%s", locals.Namespace, locals.Name)
	}

	return ingress, nil
}

// buildBackend converts a proto backend into Pulumi Ingress backend args. The
// Kubernetes API requires exactly one of port name/number — the spec's CEL rule
// guarantees it, so the two branches here are exhaustive.
func buildBackend(b *kubernetesingressv1.KubernetesIngressBackend) *kubernetesnetworkingv1.IngressBackendArgs {
	portArgs := &kubernetesnetworkingv1.ServiceBackendPortArgs{}
	if b.GetPortNumber() != 0 {
		portArgs.Number = pulumi.Int(int(b.GetPortNumber()))
	} else {
		portArgs.Name = pulumi.String(b.GetPortName())
	}

	return &kubernetesnetworkingv1.IngressBackendArgs{
		Service: &kubernetesnetworkingv1.IngressServiceBackendArgs{
			// service_name is a StringValueOrRef foreign key (default kind
			// KubernetesService), resolved to the literal Service name before
			// the module runs — this is where a workload's exported `service`
			// output lands when charts wire exposure to a workload.
			Name: pulumi.String(b.GetServiceName().GetValue()),
			Port: portArgs,
		},
	}
}

// resolvePathType maps the protobuf enum to the Kubernetes API path type string.
func resolvePathType(t kubernetesingressv1.KubernetesIngressHttpPath_KubernetesIngressPathType) string {
	switch t {
	case kubernetesingressv1.KubernetesIngressHttpPath_exact:
		return "Exact"
	case kubernetesingressv1.KubernetesIngressHttpPath_implementation_specific:
		return "ImplementationSpecific"
	default:
		return "Prefix"
	}
}
