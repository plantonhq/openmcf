package module

import (
	kubernetesbackendtlspolicyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesbackendtlspolicy/v1alpha1"
	gatewayv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/gatewayapis/kubernetes/gateway/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildValidation maps the policy's validation block (the trust anchor, the
// SNI/authentication hostname, and optional Subject Alternative Names) onto
// the typed crd2pulumi validation args. The trust-anchor arms are
// exactly-one-of (enforced by protovalidate CEL): caCertificateRefs and
// wellKnownCACertificates are each set only when present, so exactly one
// appears in the rendered CR. hostname is required and always set.
func buildValidation(validation *kubernetesbackendtlspolicyv1alpha1.KubernetesBackendTlsPolicyValidation) gatewayv1.BackendTLSPolicySpecValidationArgs {
	args := gatewayv1.BackendTLSPolicySpecValidationArgs{
		Hostname: pulumi.String(validation.GetHostname()),
	}

	if caCertificateRefs := validation.GetCaCertificateRefs(); len(caCertificateRefs) > 0 {
		args.CaCertificateRefs = buildCaCertificateRefs(caCertificateRefs)
	}

	if validation.WellKnownCaCertificates != nil {
		args.WellKnownCACertificates = pulumi.String(validation.GetWellKnownCaCertificates())
	}

	if subjectAltNames := validation.GetSubjectAltNames(); len(subjectAltNames) > 0 {
		args.SubjectAltNames = buildSubjectAltNames(subjectAltNames)
	}

	return args
}

// buildCaCertificateRefs maps the CA-bundle references (Core: one ConfigMap
// with the PEM bundle in a key named `ca.crt`) onto the typed crd2pulumi
// caCertificateRefs array.
//
// group is a presence-required optional in the spec (the CRD requires the KEY
// but allows the empty value — ConfigMaps live in the core group, the empty
// string), so it is ALWAYS set from the resolved value, never dropped when
// empty.
//
// Each reference's name is a KubernetesConfigMap foreign key resolved to its
// literal value before the module runs, so GetValue() returns the final name.
func buildCaCertificateRefs(caCertificateRefs []*kubernetesbackendtlspolicyv1alpha1.KubernetesBackendTlsPolicyCaCertificateReference) gatewayv1.BackendTLSPolicySpecValidationCaCertificateRefsArray {
	arr := gatewayv1.BackendTLSPolicySpecValidationCaCertificateRefsArray{}
	for _, ref := range caCertificateRefs {
		arr = append(arr, gatewayv1.BackendTLSPolicySpecValidationCaCertificateRefsArgs{
			Group: pulumi.String(ref.GetGroup()),
			Kind:  pulumi.String(ref.GetKind()),
			Name:  pulumi.String(ref.GetName().GetValue()),
		})
	}
	return arr
}

// buildSubjectAltNames maps the Subject Alternative Names the backend
// certificate must prove onto the typed crd2pulumi subjectAltNames array.
// type is a closed enum ("Hostname" | "URI") and always set; hostname and
// uri are each set only when present — protovalidate CEL enforces the
// type/value pairing, so exactly the matching value field appears.
func buildSubjectAltNames(subjectAltNames []*kubernetesbackendtlspolicyv1alpha1.KubernetesBackendTlsPolicySubjectAltName) gatewayv1.BackendTLSPolicySpecValidationSubjectAltNamesArray {
	arr := gatewayv1.BackendTLSPolicySpecValidationSubjectAltNamesArray{}
	for _, san := range subjectAltNames {
		args := gatewayv1.BackendTLSPolicySpecValidationSubjectAltNamesArgs{
			Type: pulumi.String(san.GetType()),
		}
		if san.Hostname != nil {
			args.Hostname = pulumi.String(san.GetHostname())
		}
		if san.Uri != nil {
			args.Uri = pulumi.String(san.GetUri())
		}
		arr = append(arr, args)
	}
	return arr
}
