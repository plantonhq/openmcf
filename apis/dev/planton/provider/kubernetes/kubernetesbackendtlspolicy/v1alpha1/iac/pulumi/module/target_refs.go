package module

import (
	kubernetesbackendtlspolicyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesbackendtlspolicy/v1alpha1"
	gatewayv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/gatewayapis/kubernetes/gateway/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildTargetRefs maps the policy target references (the same-namespace
// Services this policy secures) onto the typed crd2pulumi targetRefs array.
//
// group is a presence-required optional in the spec (the CRD requires the KEY
// but allows the empty value — Services live in the core group, the empty
// string), so it is ALWAYS set from the resolved value, never dropped when
// empty. sectionName is genuinely optional and only set when present so the
// policy covers the entire resource otherwise.
//
// Each reference's name is a KubernetesService foreign key resolved to its
// literal value before the module runs, so GetValue() returns the final name.
func buildTargetRefs(targetRefs []*kubernetesbackendtlspolicyv1alpha1.KubernetesBackendTlsPolicyTargetReference) gatewayv1.BackendTLSPolicySpecTargetRefsArray {
	arr := gatewayv1.BackendTLSPolicySpecTargetRefsArray{}
	for _, ref := range targetRefs {
		args := gatewayv1.BackendTLSPolicySpecTargetRefsArgs{
			Group: pulumi.String(ref.GetGroup()),
			Kind:  pulumi.String(ref.GetKind()),
			Name:  pulumi.String(ref.GetName().GetValue()),
		}
		if ref.SectionName != nil {
			args.SectionName = pulumi.String(ref.GetSectionName())
		}
		arr = append(arr, args)
	}
	return arr
}
